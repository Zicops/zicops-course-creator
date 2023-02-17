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
	File       graphql.Upload
	BucketPath string
	LspId      string
}

var (
	UploaderQueue = make(chan *UploadRequest, 5)
)

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
				panic(err.Error())
			}
			writer, err := storageC.UploadToGCS(ctx, req.BucketPath, map[string]string{})
			if err != nil {
				log.Errorf("Failed to upload video to course topic: %v", err.Error())
				panic(err.Error())
			}
			// read the file in chunks and upload incrementally
			buf := make([]byte, 1024)
			for {
				n, err := req.File.File.Read(buf)
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Errorf("Failed to read file: %v", err.Error())
					panic(err.Error())
				}
				_, err = writer.Write(buf[:n])
				if err != nil {
					log.Errorf("Failed to upload file: %v", err.Error())
					panic(err.Error())
				}
			}
			writer.Close()
		}
	}()
}
