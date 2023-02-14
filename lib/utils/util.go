package utils

import (
	"bytes"
	"context"
	"io"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
	"github.com/zicops/zicops-course-creator/lib/googleprojectlib"
)

func UploadFileToGCP(file *io.ReadSeeker, bucketPath string, lspId string) {
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
		n, err := (*file).Read(buf)
		if err != nil && err != io.EOF {
			log.Errorf("Failed to upload video to course topic: %v", err.Error())
			return
		}
		if n == 0 {
			break
		}
		chunk := make([]byte, n)
		copy(chunk, buf[:n])
		chunkChan <- chunk
	}

	close(chunkChan)
	wg.Wait()
}
