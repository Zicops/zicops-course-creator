package utils

import (
	"bytes"
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

var (
	UploaderQueue = make(chan *UploadRequest, 10)
)

func init() {
	go func() {
		for {
			req := <-UploaderQueue
			UploadFileToGCP(req.File, req.BucketPath, req.LspId)
		}
	}()
}

// UploadFileToGCP uploads file to GCP bucket
func UploadFileToGCP(file *graphql.Upload, bucketPath string, lspId string) {
	storageC := bucket.NewStorageHandler()
	ctx := context.Background()
	const chunkSize = 1024 * 1024 // 1 MB
	gproject := googleprojectlib.GetGoogleProjectID()
	err := storageC.InitializeStorageClient(ctx, gproject, lspId)
	if err != nil {
		log.Errorf("Failed to upload video to course topic: %v", err.Error())
	}
	writer, err := storageC.UploadToGCS(ctx, bucketPath, map[string]string{})
	if err != nil {
		log.Errorf("Failed to upload video to course topic: %v", err.Error())
		return
	}
	defer writer.Close()

	chunkChan := make(chan []byte, 100)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for chunk := range chunkChan {
			_, err := io.Copy(writer, bytes.NewReader(chunk))
			if err != nil {
				log.Errorf("Failed to upload video to course topic: %v", err.Error())
				return
			}
		}
	}()

	buf := make([]byte, chunkSize)
	for {
		if file == nil {
			break
		}
		n, err := file.File.Read(buf)
		if err != nil && err != io.EOF {
			log.Errorf("Failed to upload video to course topic: %v", err.Error())
			return
		}
		if n == 0 {
			break
		}
		chunkChan <- buf[:n]
	}
	close(chunkChan)
	wg.Wait()
}
