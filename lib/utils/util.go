package utils

import (
	"bytes"
	"context"
	"io"

	log "github.com/sirupsen/logrus"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
)

func UploadFileToGCP(storageC bucket.Client, file io.Reader, bucketPath string, lspId string) {
	ctx := context.Background()
	writer, err := storageC.UploadToGCS(ctx, bucketPath, map[string]string{})
	if err != nil {
		log.Errorf("Failed to upload video to course topic: %v", err.Error())
	}
	defer writer.Close()
	fileBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBuffer, file); err != nil {
		log.Errorf("Failed to upload video to course topic: %v", err.Error())
	}
	currentBytes := fileBuffer.Bytes()
	_, err = io.Copy(writer, bytes.NewReader(currentBytes))
	if err != nil {
		log.Errorf("Failed to upload video to course topic: %v", err.Error())
	}
}
