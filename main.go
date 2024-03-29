package main

import (
	"bytes"
	"embed"
	_ "embed"
	"excalidraw-complete/core"
	"excalidraw-complete/documents/memory"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/zishang520/engine.io/v2/types"
	socketio "github.com/zishang520/socket.io/v2/socket"
)

type (
	DocumentCreateResponse struct {
		ID string `json:"id"`
	}

	UserToFollow struct {
		SocketId string
		Username string
	}
	OnUserFollowedPayload struct {
		UserToFollow UserToFollow
		action       string
	}
)

//go:embed all:frontend
var assets embed.FS

func Assets() (fs.FS, error) {
	return fs.Sub(assets, "frontend")
}

type FrontEndHandler struct{}

func (h FrontEndHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	frontendHandler(w, r)
}

func frontendHandler(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	upath = path.Clean(upath)

	sub, err := fs.Sub(assets, "frontend")
	if err != nil {
		panic(err)
	}
	http.FileServer(http.FS(sub)).ServeHTTP(w, r)

}

func main() {
	opts := socketio.DefaultServerOptions()
	opts.SetMaxHttpBufferSize(5000000)
	opts.SetPath("/socket.io")
	opts.SetAllowEIO3(true)
	opts.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})
	ioo := socketio.NewServer(nil, opts)
	ioo.On("connection", func(clients ...any) {
		socket := clients[0].(*socketio.Socket)
		ioo.To(socketio.Room(socket.Id())).Emit("init-room")
		me := socket.Id()
		socket.On("join-room", func(datas ...any) {
			room := socketio.Room(datas[0].(string))
			fmt.Printf("Socket %v has joined %v\n", me, room)
			socket.Join(room)
			ioo.In(room).FetchSockets()(func(usersInRoom []*socketio.RemoteSocket, _ error) {
				if len(usersInRoom) <= 1 {
					ioo.To(socketio.Room(me)).Emit("first-in-room")
				} else {
					fmt.Printf("emit new user %v in room %v\n", me, room)
					socket.Broadcast().To(room).Emit("new-user", me)
				}

				// Inform all clients by new users.
				newRoomUsers := []socketio.SocketId{}
				for _, user := range usersInRoom {
					newRoomUsers = append(newRoomUsers, user.Id())
				}
				fmt.Printf(" room %v has users %v\n", room, newRoomUsers)
				ioo.In(room).Emit(
					"room-user-change",
					newRoomUsers,
				)

			})
		})
		socket.On("server-broadcast", func(datas ...any) {
			roomID := datas[0].(string)
			fmt.Printf(" user %v sends update to room %v\n", me, roomID)
			socket.Broadcast().To(socketio.Room(roomID)).Emit("client-broadcast", datas[1], datas[2])
		})
		socket.On("server-volatile-broadcast", func(datas ...any) {
			roomID := datas[0].(string)
			fmt.Printf(" user %v sends volatile update to room %v\n", me, roomID)
			socket.Volatile().Broadcast().To(socketio.Room(roomID)).Emit("client-broadcast", datas[1], datas[2])
		})

		socket.On("user-follow", func(datas ...any) {
			// TODO()

		})
		socket.On("disconnecting", func(datas ...any) {
			for _, currentRoom := range socket.Rooms().Keys() {
				ioo.In(currentRoom).FetchSockets()(func(usersInRoom []*socketio.RemoteSocket, _ error) {
					allUsers := []socketio.SocketId{}
					remainingUsers := []socketio.SocketId{}
					fmt.Printf("disconnecting %v from room %v\n", me, currentRoom)
					for _, userInRoom := range usersInRoom {
						allUsers = append(allUsers, userInRoom.Id())
						if userInRoom.Id() != me {
							remainingUsers = append(remainingUsers, userInRoom.Id())
						}
					}
					if len(remainingUsers) > 0 {
						fmt.Printf("leaving user, room %v has users %v -> %v\n", currentRoom, allUsers, remainingUsers)
						ioo.In(currentRoom).Emit(
							"room-user-change",
							remainingUsers,
						)

					}

				})

			}

		})
		socket.On("disconnect", func(datas ...any) {
			socket.RemoveAllListeners("")
			socket.Disconnect(true)
		})
	})

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Content-Length", "X-CSRF-Token", "Token", "session", "Origin", "Host", "Connection", "Accept-Encoding", "Accept-Language", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	documentStore := memory.NewDocumentStore()

	r.Mount("/", FrontEndHandler{})

	r.Route("/api/v2", func(r chi.Router) {
		r.Post("/post/", func(w http.ResponseWriter, r *http.Request) {
			data := new(bytes.Buffer)
			_, err := io.Copy(data, r.Body)
			if err != nil {
				http.Error(w, "Failed to copy", http.StatusInternalServerError)
				return
			}
			id, err := documentStore.Create(r.Context(), &core.Document{Data: *data})
			if err != nil {
				http.Error(w, "Failed to save", http.StatusInternalServerError)
				return
			}

			render.JSON(w, r, DocumentCreateResponse{ID: id})
			render.Status(r, http.StatusOK)
		})
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				id := chi.URLParam(r, "id")
				document, err := documentStore.FindID(r.Context(), id)
				if err != nil {
					http.Error(w, "not found", http.StatusNotFound)
					return
				}
				w.Write(document.Data.Bytes())
			})
		})
	})

	r.Handle("/socket.io/", ioo.ServeHandler(nil))
	go http.ListenAndServe(":3002", r)
	fmt.Println("listen on 3002")

	exit := make(chan struct{})
	SignalC := make(chan os.Signal)

	signal.Notify(SignalC, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range SignalC {
			switch s {
			case os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				close(exit)
				return
			}
		}
	}()

	<-exit
	ioo.Close(nil)
	os.Exit(0)
}
