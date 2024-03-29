# Exalidraw Complete

Frustrated on how difficult it is to setup excalidraw self-hosted but with data
storage and collaboration function this represents and attempt to run the
necessary function with a single binary implemented in go. This includes:

- the frontend UI
- a in-memory data layer
- socket.io implementation for collaboration

Apply the patch to the frontend and build excalidraw into `frontend`. Run
```bash
go run main.go
```

Everything will be served under `localhost:3002`

