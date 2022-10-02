package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/helpers"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
	"github.com/zicops/zicops-course-creator/lib/googleprojectlib"
)

func AddTopicResources(ctx context.Context, courseID string, resource *model.TopicResourceInput) (*model.UploadResult, error) {
	log.Infof("AddTopicResources Called")
	guid := xid.New().String()
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	if lspId == "" {
		return nil, fmt.Errorf("lsp_id is empty")
	}
	email_creator := claims["email"].(string)
	isSuccess := model.UploadResult{}
	getUrl := ""
	bucketPath := ""
	if resource.URL == nil {
		storageC := bucket.NewStorageHandler()
		gproject := googleprojectlib.GetGoogleProjectID()
		err := storageC.InitializeStorageClient(ctx, gproject, lspId)
		if err != nil {
			log.Errorf("Failed to upload video to course topic: %v", err.Error())
			return &isSuccess, nil
		}
		if courseID == "" || resource.TopicID == nil {
			return &isSuccess, nil
		}
		bucketPath = courseID + "/" + *resource.TopicID + "/" + resource.File.Filename
		writer, err := storageC.UploadToGCS(ctx, bucketPath, map[string]string{})
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
		getUrl = storageC.GetSignedURLForObject(bucketPath)
	} else {
		getUrl = *resource.URL
	}
	sourceName := ""
	if resource.Name != nil {
		sourceName = *resource.Name
	}
	createdBy := email_creator
	updatedBy := email_creator
	cassandraResource := coursez.Resource{
		ID:         guid,
		Name:       sourceName,
		TopicId:    *resource.TopicID,
		CourseId:   courseID,
		BucketPath: bucketPath,
		Url:        getUrl,
		IsActive:   false,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
		CreatedBy:  createdBy,
		UpdatedBy:  updatedBy,
	}
	if resource.Type != nil {
		cassandraResource.Type = *resource.Type
	}
	// update course image in cassandra
	resourceAdd := CassSession.Query(coursez.ResourceTable.Insert()).BindStruct(cassandraResource)
	if err := resourceAdd.ExecRelease(); err != nil {
		return &isSuccess, err
	}
	isSuccessRes := true
	isSuccess.Success = &isSuccessRes
	isSuccess.URL = &getUrl
	return &isSuccess, nil
}
