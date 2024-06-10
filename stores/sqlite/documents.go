package sqlite

import (
	"bytes"
	"context"
	"excalidraw-complete/core"
	"fmt"

	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

var savedDocuments = make(map[string]core.Document)

type documentStore struct {
	db *sql.DB
}

func NewDocumentStore(dataSourceName string) core.DocumentStore {
	// db, err := sql.Open("sqlite3", ":memory:")
	db, err := sql.Open("sqlite3", dataSourceName)

	if err != nil {
		log.Fatal(err)
	}
	sts := `CREATE TABLE IF NOT EXISTS documents (id TEXT PRIMARY KEY, data BLOB);`
	_, err = db.Exec(sts)
	if err != nil {
		log.Fatal(err)
	}
	return &documentStore{db}
}

func (s *documentStore) FindID(ctx context.Context, id string) (*core.Document, error) {
	log := logrus.WithField("document_id", id)
	log.Debug("Retrieving document by ID")
	var data []byte
	err := s.db.QueryRowContext(ctx, "SELECT data FROM documents WHERE id = ?", id).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			log.WithField("error", "document not found").Warn("Document with specified ID not found")
			return nil, fmt.Errorf("document with id %s not found", id)
		}
		log.WithField("error", err).Error("Failed to retrieve document")
		return nil, err
	}
	document := core.Document{
		Data: *bytes.NewBuffer(data),
	}
	log.Info("Document retrieved successfully")
	return &document, nil
}

func (s *documentStore) Create(ctx context.Context, document *core.Document) (string, error) {
	id := ulid.Make().String()
	data := document.Data.Bytes()
	log := logrus.WithFields(logrus.Fields{
		"document_id": id,
		"data_length": len(data),
	})

	_, err := s.db.ExecContext(ctx, "INSERT INTO documents (id, data) VALUES (?, ?)", id, data)
	if err != nil {
		log.WithField("error", err).Error("Failed to create document")
		return "", err
	}
	log.Info("Document created successfully")
	return id, nil
}
