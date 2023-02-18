package utils

import (
	"context"
	"io"
	"sync"

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

var UploaderQueue = make(chan UploadRequest)
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
			var wg sync.WaitGroup
			for {
				n, err := req.File.File.Read(buf)
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Errorf("Failed to read file: %v", err.Error())
					ErrorQueue <- err
				}

				// increment the WaitGroup counter for each chunk
				wg.Add(1)
				go func(chunk []byte) {
					_, err = writer.Write(chunk)
					if err != nil {
						log.Errorf("Failed to upload file: %v", err.Error())
						ErrorQueue <- err
					}
					wg.Done()
				}(buf[:n])
			}

			// wait for all the chunks to be uploaded before closing the writer
			wg.Wait()
			err = writer.Close()
			if err != nil {
				log.Errorf("Failed to close writer: %v", err.Error())
				ErrorQueue <- err
			}
		}
	}()
}
