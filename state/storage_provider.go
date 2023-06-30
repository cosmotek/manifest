package state

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const basepath = "state"

type StorageProvider interface {
	Put(key string, val any) error
	Get(key string, dst any) error
}

type S3BucketConfig struct {
	BucketName string `split_words:"true" required:"true"`
}

type S3Bucket struct {
	svc        *s3.S3
	bucketName string
}

func NewS3Bucket(conf S3BucketConfig) (*S3Bucket, error) {
	session, err := session.NewSession(&aws.Config{})
	if err != nil {
		return nil, err
	}

	return &S3Bucket{
		svc: s3.New(session),
	}, nil
}

func (s *S3Bucket) Put(key string, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	putObjInput := &s3.PutObjectInput{
		Body:        aws.ReadSeekCloser(bytes.NewReader(data)),
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(filepath.Join(basepath, key)),
		ContentType: aws.String("application/json"),
	}

	_, err = s.svc.PutObjectWithContext(context.TODO(), putObjInput)
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Bucket) Get(key string, dst any) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(filepath.Join(basepath, key)),
	}

	result, err := s.svc.GetObjectWithContext(context.TODO(), input)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	decoder := json.NewDecoder(result.Body)
	err = decoder.Decode(&dst)
	if err != nil {
		return err
	}

	return nil
}
