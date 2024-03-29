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

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("document with id %s not found", id)
		}
		return nil, err
	}

	document := core.Document{
		Data: *bytes.NewBuffer(data),
	}

	return &document, nil
}

func (s *documentStore) Create(ctx context.Context, document *core.Document) (string, error) {
	id := ulid.Make().String()
	filePath := filepath.Join(s.basePath, id)

	if err := os.WriteFile(filePath, document.Data.Bytes(), 0644); err != nil {
		return "", err
	}

	return id, nil
}
