package aws

import (
	"bytes"
	"context"
	"excalidraw-complete/core"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/oklog/ulid/v2"
)

type documentStore struct {
	s3Client *s3.Client
	bucket   string // Name of the S3 bucket
}

func NewDocumentStore(bucketName string) core.DocumentStore {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	return &documentStore{
		s3Client: s3Client,
		bucket:   bucketName,
	}
}

func (s *documentStore) FindID(ctx context.Context, id string) (*core.Document, error) {
	resp, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get document with id %s: %v", id, err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read document data: %v", err)
	}

	document := core.Document{
		Data: *bytes.NewBuffer(data),
	}

	return &document, nil
}

func (s *documentStore) Create(ctx context.Context, document *core.Document) (string, error) {
	id := ulid.Make().String()

	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(id),
		Body:   bytes.NewReader(document.Data.Bytes()),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload document: %v", err)
	}

	return id, nil
}
