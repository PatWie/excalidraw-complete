package stores

import (
	"excalidraw-complete/core"
	"excalidraw-complete/stores/aws"
	"excalidraw-complete/stores/filesystem"
	"excalidraw-complete/stores/memory"
	"excalidraw-complete/stores/sqlite"
	"os"
)

func GetStore() core.DocumentStore {
	storageType := os.Getenv("STORAGE_TYPE")
	var store core.DocumentStore

	switch storageType {
	case "filesystem":
		basePath := os.Getenv("LOCAL_STORAGE_PATH")
		store = filesystem.NewDocumentStore(basePath)
	case "sqlite":
		dataSourceName := os.Getenv("DATA_SOURCE_NAME")
		store = sqlite.NewDocumentStore(dataSourceName)
	case "s3":
		bucketName := os.Getenv("S3_BUCKET_NAME")
		store = aws.NewDocumentStore(bucketName)
	default:
		store = memory.NewDocumentStore()
	}
	return store
}
