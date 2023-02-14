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
	UploaderQueue = make(chan *UploadRequest, 1)
)

func init() {
	go func() {
		storageC := bucket.NewStorageHandler()
		ctx := context.Background()
		gproject := googleprojectlib.GetGoogleProjectID()
		err := storageC.InitializeStorageClient(ctx, gproject, "lspId")
		if err != nil {
			log.Errorf("Failed to upload video to course topic: %v", err.Error())
		}
		for {
			req := <-UploaderQueue
			writer, err := storageC.UploadToGCS(ctx, req.BucketPath, map[string]string{})
			if err != nil {
				log.Errorf("Failed to upload video to course topic: %v", err.Error())
				return
			}
			defer writer.Close()
			// Upload the file to GCS using writer without buffering and routines
			// This is because we want to upload the file as it is being read from the request
			// and not buffer it in memory
			_, err = io.Copy(writer, req.File.File)
			if err != nil {
				log.Errorf("Failed to upload video to course topic: %v", err.Error())
				return
			}
		}
	}()
}
