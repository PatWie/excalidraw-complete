package filesystem

import (
	"bytes"
	"context"
	"excalidraw-complete/core"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type documentStore struct {
	basePath string // Directory where documents are stored.
}

func NewDocumentStore(basePath string) core.DocumentStore {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		log.Fatalf("failed to create base directory: %v", err)
	}

	return &documentStore{basePath: basePath}
}

func (s *documentStore) FindID(ctx context.Context, id string) (*core.Document, error) {
	filePath := filepath.Join(s.basePath, id)
	log := logrus.WithField("document_id", id)

	log.WithField("file_path", filePath).Info("Retrieving document by ID")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
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
	filePath := filepath.Join(s.basePath, id)
	log := logrus.WithFields(logrus.Fields{
		"document_id": id,
		"file_path":   filePath,
	})
	log.Info("Creating new document")

	if err := os.WriteFile(filePath, document.Data.Bytes(), 0644); err != nil {
		log.WithField("error", err).Error("Failed to create document")
		return "", err
	}

	log.Info("Document created successfully")
	return id, nil
}
