package utils

import (
	"context"
	"io"

	"github.com/99designs/gqlgen/graphql"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
	"github.com/zicops/zicops-course-creator/lib/googleprojectlib"
)

type UploadRequest struct {
	File       *graphql.Upload
	BucketPath string
	LspId      string
}

var UploaderQueue = make(chan UploadRequest, 10)
var ErrorQueue = make(chan error)

func init() {
	go func() {
		ctx := context.Background()
		for {
			req := <-UploaderQueue
			storageC := bucket.NewStorageHandler()
			gproject := googleprojectlib.GetGoogleProjectID()
			err := storageC.InitializeStorageClient(ctx, gproject, req.LspId)
			if err != nil {
				log.Errorf("Failed to upload video to course topic: %v", err.Error())
				ErrorQueue <- err
			}
			writer, err := storageC.UploadToGCS(ctx, req.BucketPath, map[string]string{})
			if err != nil {
				log.Errorf("Failed to upload video to course topic: %v", err.Error())
				ErrorQueue <- err
			}

			// read the file in chunks and upload incrementally
			// create chunks of 10mb
			buf := make([]byte, 10*1024*1024)
			for {
				n, err := req.File.File.Read(buf)
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Errorf("Failed to read file: %v", err.Error())
				}
				_, err = writer.Write(buf[:n])
				if err != nil {
					log.Errorf("Failed to upload file: %v", err.Error())
				}
			}
			err = writer.Close()
			if err != nil {
				log.Errorf("Failed to close writer: %v", err.Error())
				ErrorQueue <- err
			}
		}
	}()
}

func SendUploadRequestToUploaderQueue(ctx context.Context, file *graphql.Upload, bucketPath string, lspId string) error {
	// send message to uploader queue
	uploadRequest := UploadRequest{
		BucketPath: bucketPath,
		File:       file,
		LspId:      lspId,
	}

	// send message to uploader queue in utils
	UploaderQueue <- uploadRequest
	return nil
}
