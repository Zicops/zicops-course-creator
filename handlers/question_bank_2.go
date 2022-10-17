package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/rs/xid"
	"github.com/scylladb/gocqlx/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/qbankz"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/helpers"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
	"github.com/zicops/zicops-course-creator/lib/googleprojectlib"
)

func AddQuestionOptions(ctx context.Context, input *model.QuestionOptionInput) (*model.QuestionOption, error) {
	log.Info("AddQuestionOptions called")
	if input.QmID == nil {
		return nil, fmt.Errorf("question id not found")
	}
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspID := claims["lsp_id"].(string)
	if lspID == "" {
		return nil, fmt.Errorf("lsp id not found")
	}
	email_creator := claims["email"].(string)
	guid := xid.New()
	cassandraQuestionBank := qbankz.OptionsMain{
		ID:             guid.String(),
		Description:    *input.Description,
		IsCorrect:      *input.IsCorrect,
		QmId:           *input.QmID,
		IsActive:       *input.IsActive,
		AttachmentType: *input.AttachmentType,
		CreatedBy:      email_creator,
		UpdatedBy:      email_creator,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
		LspId:          lspID,
	}
	getUrl := ""
	if input.File != nil {
		bucketPath := "question_banks/" + cassandraQuestionBank.QmId + "/" + cassandraQuestionBank.ID + "/" + input.File.Filename
		storageC := bucket.NewStorageHandler()
		gproject := googleprojectlib.GetGoogleProjectID()
		err := storageC.InitializeStorageClient(ctx, gproject, lspID)
		if err != nil {
			log.Errorf("Failed to upload question option: %v", err.Error())
			return nil, err
		}
		writer, err := storageC.UploadToGCS(ctx, bucketPath, map[string]string{})
		if err != nil {
			log.Errorf("Failed to upload question option: %v", err.Error())
			return nil, err
		}
		defer writer.Close()
		fileBuffer := bytes.NewBuffer(nil)
		if _, err := io.Copy(fileBuffer, input.File.File); err != nil {
			return nil, err
		}
		currentBytes := fileBuffer.Bytes()
		_, err = io.Copy(writer, bytes.NewReader(currentBytes))
		if err != nil {
			return nil, err
		}
		getUrl = storageC.GetSignedURLForObject(bucketPath)
		cassandraQuestionBank.Attachment = getUrl
		cassandraQuestionBank.AttachmentBucket = bucketPath
	}
	insertQuery := CassSession.Query(qbankz.OptionsMainTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.QuestionOption{
		ID:             &cassandraQuestionBank.ID,
		Description:    &cassandraQuestionBank.Description,
		IsCorrect:      &cassandraQuestionBank.IsCorrect,
		QmID:           &cassandraQuestionBank.QmId,
		IsActive:       &cassandraQuestionBank.IsActive,
		AttachmentType: &cassandraQuestionBank.AttachmentType,
		CreatedBy:      &cassandraQuestionBank.CreatedBy,
		UpdatedBy:      &cassandraQuestionBank.UpdatedBy,
		CreatedAt:      &created,
		UpdatedAt:      &updated,
		Attachment:     &cassandraQuestionBank.Attachment,
	}
	return &responseModel, nil
}

func UpdateQuestionOptions(ctx context.Context, input *model.QuestionOptionInput) (*model.QuestionOption, error) {
	log.Info("UpdateQuestionOptions called")
	if input.ID == nil {
		return nil, fmt.Errorf("options id not found")
	}
	if input.QmID == nil {
		return nil, fmt.Errorf("question id not found")
	}
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	email_creator := claims["email"].(string)
	lspID := claims["lsp_id"].(string)
	if lspID == "" {
		return nil, fmt.Errorf("lsp id not found")
	}
	cassandraQuestionBank := *GetOptionsMain(ctx, *input.ID, lspID, CassSession)
	updatedCols := []string{}
	if input.Description != nil && *input.Description != cassandraQuestionBank.Description {
		cassandraQuestionBank.Description = *input.Description
		updatedCols = append(updatedCols, "description")
	}
	if input.IsCorrect != nil && *input.IsCorrect != cassandraQuestionBank.IsCorrect {
		cassandraQuestionBank.IsCorrect = *input.IsCorrect
		updatedCols = append(updatedCols, "is_correct")
	}
	if input.AttachmentType != nil && *input.AttachmentType != cassandraQuestionBank.AttachmentType {
		cassandraQuestionBank.AttachmentType = *input.AttachmentType
		updatedCols = append(updatedCols, "attachment_type")
	}
	if email_creator != "" && email_creator != cassandraQuestionBank.UpdatedBy {
		cassandraQuestionBank.UpdatedBy = email_creator
		updatedCols = append(updatedCols, "updated_by")
	}
	if input.QmID != nil && *input.QmID != cassandraQuestionBank.QmId {
		cassandraQuestionBank.QmId = *input.QmID
		updatedCols = append(updatedCols, "qm_id")
	}
	updatedAt := time.Now().Unix()
	if input.File != nil {
		bucketPath := "question_banks/" + cassandraQuestionBank.QmId + "/" + cassandraQuestionBank.ID + "/" + input.File.Filename
		storageC := bucket.NewStorageHandler()
		gproject := googleprojectlib.GetGoogleProjectID()
		err := storageC.InitializeStorageClient(ctx, gproject, lspID)
		if err != nil {
			log.Errorf("Failed to update options: %v", err.Error())
			return nil, err
		}
		writer, err := storageC.UploadToGCS(ctx, bucketPath, map[string]string{})
		if err != nil {
			log.Errorf("Failed to update options: %v", err.Error())
			return nil, err
		}
		defer writer.Close()
		fileBuffer := bytes.NewBuffer(nil)
		if _, err := io.Copy(fileBuffer, input.File.File); err != nil {
			return nil, err
		}
		currentBytes := fileBuffer.Bytes()
		_, err = io.Copy(writer, bytes.NewReader(currentBytes))
		if err != nil {
			return nil, err
		}
		getUrl := storageC.GetSignedURLForObject(bucketPath)
		cassandraQuestionBank.Attachment = getUrl
		updatedCols = append(updatedCols, "attachment_url")
		cassandraQuestionBank.AttachmentBucket = bucketPath
		updatedCols = append(updatedCols, "attachment_bucket")
	}
	if len(updatedCols) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	cassandraQuestionBank.UpdatedAt = updatedAt
	updatedCols = append(updatedCols, "updated_at")
	upStms, uNames := qbankz.OptionsMainTable.Update(updatedCols...)
	updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.QuestionOption{
		ID:             &cassandraQuestionBank.ID,
		Description:    &cassandraQuestionBank.Description,
		IsCorrect:      &cassandraQuestionBank.IsCorrect,
		QmID:           &cassandraQuestionBank.QmId,
		IsActive:       &cassandraQuestionBank.IsActive,
		AttachmentType: &cassandraQuestionBank.AttachmentType,
		CreatedBy:      &cassandraQuestionBank.CreatedBy,
		UpdatedBy:      &cassandraQuestionBank.UpdatedBy,
		CreatedAt:      &created,
		UpdatedAt:      &updated,
		Attachment:     &cassandraQuestionBank.Attachment,
	}
	return &responseModel, nil
}

func GetOptionsMain(ctx context.Context, courseID string, lspID string, session *gocqlx.Session) *qbankz.OptionsMain {
	chapters := []qbankz.OptionsMain{}
	getQueryStr := fmt.Sprintf("SELECT * FROM qbankz.options_main WHERE id='%s' and lsp_id='%s' and is_active=true", courseID, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}
