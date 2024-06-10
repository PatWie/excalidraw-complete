package stores

import (
	"excalidraw-complete/core"
	"excalidraw-complete/stores/aws"
	"excalidraw-complete/stores/filesystem"
	"excalidraw-complete/stores/memory"
	"excalidraw-complete/stores/sqlite"
	"os"

	"github.com/sirupsen/logrus"
)

func GetStore() core.DocumentStore {
	storageType := os.Getenv("STORAGE_TYPE")
	var store core.DocumentStore

	storageField := logrus.Fields{
		"storageType": storageType,
	}

	switch storageType {
	case "filesystem":
		basePath := os.Getenv("LOCAL_STORAGE_PATH")
		storageField["basePath"] = basePath
		store = filesystem.NewDocumentStore(basePath)
	case "sqlite":
		dataSourceName := os.Getenv("DATA_SOURCE_NAME")
		storageField["dataSourceName"] = dataSourceName
		store = sqlite.NewDocumentStore(dataSourceName)
	case "s3":
		bucketName := os.Getenv("S3_BUCKET_NAME")
		storageField["bucketName"] = bucketName
		store = aws.NewDocumentStore(bucketName)
	default:
		store = memory.NewDocumentStore()
		storageField["storageType"] = "in-memory"
	}
	logrus.WithFields(storageField).Info("Use storage")
	return store
}
