package handlers

import (
	"bytes"
	"context"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
	"github.com/zicops/zicops-course-creator/lib/googleprojectlib"
)

func AddTopicResources(ctx context.Context, courseID string, resource *model.TopicResourceInput) (*model.UploadResult, error) {
	log.Infof("AddTopicResources Called")
	isSuccess := model.UploadResult{}
	storageC := bucket.NewStorageHandler()
	gproject := googleprojectlib.GetGoogleProjectID()
	err := storageC.InitializeStorageClient(ctx, gproject)
	if err != nil {
		log.Errorf("Failed to upload video to course topic: %v", err.Error())
		return &isSuccess, nil
	}
	if courseID == "" || resource.TopicID == nil {
		return &isSuccess, nil
	}
	bucketPath := courseID + "/" + *resource.TopicID + "/" + resource.File.Filename
	writer, err := storageC.UploadToGCS(ctx, bucketPath)
	if err != nil {
		log.Errorf("Failed to upload video to course topic: %v", err.Error())
		return &isSuccess, nil
	}
	defer writer.Close()
	fileBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBuffer, resource.File.File); err != nil {
		return &isSuccess, nil
	}
	currentBytes := fileBuffer.Bytes()
	_, err = io.Copy(writer, bytes.NewReader(currentBytes))
	if err != nil {
		return &isSuccess, err
	}
	getUrl := storageC.GetSignedURLForObject(bucketPath)
	sourceName := ""
	if resource.Name != nil {
		sourceName = *resource.Name
	}
	cassandraResource := coursez.Resource{
		Name:       sourceName,
		TopicId:    *resource.TopicID,
		BucketPath: bucketPath,
		Url:        getUrl,
		IsActive:   false,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}
	if resource.Type != nil {
		cassandraResource.Type = *resource.Type
	}
	// update course image in cassandra
	resourceAdd := global.CassSession.Session.Query(coursez.ResourceTable.Insert()).BindStruct(cassandraResource)
	if err := resourceAdd.ExecRelease(); err != nil {
		return &isSuccess, err
	}
	isSuccessRes := true
	isSuccess.Success = &isSuccessRes
	isSuccess.URL = &getUrl
	return &isSuccess, nil
}
