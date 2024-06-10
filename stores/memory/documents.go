package memory

import (
	"context"
	"excalidraw-complete/core"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

var savedDocuments = make(map[string]core.Document)

type documentStore struct {
}

func NewDocumentStore() core.DocumentStore {
	return &documentStore{}
}

func (s *documentStore) FindID(ctx context.Context, id string) (*core.Document, error) {
	log := logrus.WithField("document_id", id)
	if val, ok := savedDocuments[id]; ok {
		log.Info("Document retrieved successfully")
		return &val, nil
	}
	log.WithField("error", "document not found").Warn("Document with specified ID not found")
	return nil, fmt.Errorf("document with id %s not found", id)
}

func (s *documentStore) Create(ctx context.Context, document *core.Document) (string, error) {
	id := ulid.Make().String()
	savedDocuments[id] = *document
	log := logrus.WithFields(logrus.Fields{
		"document_id": id,
		"data_length": len(document.Data.Bytes()),
	})
	log.Info("Document created successfully")

	return id, nil
}
