package memory

import (
	"context"
	"excalidraw-complete/core"
	"fmt"

	"github.com/oklog/ulid/v2"
)

var savedDocuments = make(map[string]core.Document)

type documentStore struct {
}

func NewDocumentStore() core.DocumentStore {
	return &documentStore{}
}

func (s *documentStore) FindID(ctx context.Context, id string) (*core.Document, error) {
	if val, ok := savedDocuments[id]; ok {
		return &val, nil
	}
	return nil, fmt.Errorf("document with id %s not found", id)
}

func (s *documentStore) Create(ctx context.Context, document *core.Document) (string, error) {
	id := ulid.Make().String()
	savedDocuments[id] = *document
	return id, nil
}
