package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/rs/xid"
	"github.com/scylladb/gocqlx/v2/qb"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/qbankz"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
	"github.com/zicops/zicops-course-creator/lib/googleprojectlib"
)

func QuestionBankCreate(ctx context.Context, input *model.QuestionBankInput) (*model.QuestionBank, error) {
	log.Info("QuestionBankCreate called")
	guid := xid.New()
	cassandraQuestionBank := qbankz.QuestionBankMain{
		ID:          guid.String(),
		Name:        *input.Name,
		Category:    *input.Category,
		SubCategory: *input.SubCategory,
		IsActive:    *input.IsActive,
		IsDefault:   *input.IsDefault,
		Owner:       *input.Owner,
		CreatedBy:   *input.CreatedBy,
		UpdatedBy:   *input.UpdatedBy,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	insertQuery := global.CassSessioQBank.Session.Query(qbankz.QuestionBankMainTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	responseModel := model.QuestionBank{
		ID:          &cassandraQuestionBank.ID,
		Name:        input.Name,
		Owner:       input.Owner,
		Category:    input.Category,
		SubCategory: input.SubCategory,
		IsActive:    input.IsActive,
		IsDefault:   input.IsDefault,
		CreatedBy:   input.CreatedBy,
		UpdatedBy:   input.UpdatedBy,
		CreatedAt:   &created,
		UpdatedAt:   &created,
	}
	return &responseModel, nil
}

func QuestionBankUpdate(ctx context.Context, input *model.QuestionBankInput) (*model.QuestionBank, error) {
	log.Info("QuestionBankUpdate called")
	if input.ID == nil {
		return nil, fmt.Errorf("question bank not found")
	}
	cassandraQuestionBank := qbankz.QuestionBankMain{
		ID: *input.ID,
	}
	banks := []qbankz.QuestionBankMain{}
	getQuery := global.CassSessioQBank.Session.Query(qbankz.QuestionBankMainTable.Get()).BindMap(qb.M{"id": cassandraQuestionBank.ID})
	if err := getQuery.SelectRelease(&banks); err != nil {
		return nil, err
	}
	if len(banks) == 0 {
		return nil, fmt.Errorf("question bank not found")
	}
	cassandraQuestionBank = banks[0]
	updatedCols := []string{}
	if input.Name != nil {
		cassandraQuestionBank.Name = *input.Name
		updatedCols = append(updatedCols, "name")
	}
	if input.Category != nil {
		cassandraQuestionBank.Category = *input.Category
		updatedCols = append(updatedCols, "category")
	}
	if input.SubCategory != nil {
		cassandraQuestionBank.SubCategory = *input.SubCategory
		updatedCols = append(updatedCols, "sub_category")
	}
	if input.IsActive != nil {
		cassandraQuestionBank.IsActive = *input.IsActive
		updatedCols = append(updatedCols, "is_active")
	}
	if input.IsDefault != nil {
		cassandraQuestionBank.IsDefault = *input.IsDefault
		updatedCols = append(updatedCols, "is_default")
	}
	if input.Owner != nil {
		cassandraQuestionBank.Owner = *input.Owner
		updatedCols = append(updatedCols, "owner")
	}
	if input.UpdatedBy != nil {
		cassandraQuestionBank.UpdatedBy = *input.UpdatedBy
		updatedCols = append(updatedCols, "updated_by")
	}
	updatedAt := time.Now().Unix()
	cassandraQuestionBank.UpdatedAt = updatedAt
	updatedCols = append(updatedCols, "updated_at")
	if len(updatedCols) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	upStms, uNames := qbankz.QuestionBankMainTable.Update(updatedCols...)
	updateQuery := global.CassSessioQBank.Session.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	responseModel := model.QuestionBank{
		ID:          &cassandraQuestionBank.ID,
		Name:        input.Name,
		Owner:       input.Owner,
		Category:    input.Category,
		SubCategory: input.SubCategory,
		IsActive:    input.IsActive,
		IsDefault:   input.IsDefault,
		CreatedBy:   input.CreatedBy,
		UpdatedBy:   input.UpdatedBy,
		CreatedAt:   &created,
		UpdatedAt:   &created,
	}
	return &responseModel, nil
}

func AddQuestionBankQuestion(ctx context.Context, input *model.QuestionBankQuestionInput) (*model.QuestionBankQuestion, error) {
	log.Info("AddQuestionBankQuestion called")
	if input.QbmID == nil {
		return nil, fmt.Errorf("question bank main id not found")
	}
	guid := xid.New()
	cassandraQuestionBank := qbankz.QuestionMain{
		ID:             guid.String(),
		Description:    *input.Description,
		QbmId:          *input.QbmID,
		Type:           *input.Type,
		Difficulty:     *input.Difficulty,
		Hint:           *input.Hint,
		Status:         *input.Status,
		AttachmentType: *input.AttachmentType,
		CreatedBy:      *input.CreatedBy,
		UpdatedBy:      *input.UpdatedBy,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
	}
	bucketPath := "question_banks/" + cassandraQuestionBank.QbmId + "/" + cassandraQuestionBank.ID + "/" + input.File.Filename
	storageC := bucket.NewStorageHandler()
	gproject := googleprojectlib.GetGoogleProjectID()
	err := storageC.InitializeStorageClient(ctx, gproject)
	if err != nil {
		log.Errorf("Failed to upload video to course topic: %v", err.Error())
		return nil, err
	}
	writer, err := storageC.UploadToGCS(ctx, bucketPath, map[string]string{})
	if err != nil {
		log.Errorf("Failed to upload video to course topic: %v", err.Error())
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
	cassandraQuestionBank.AttachmentBucket = bucketPath
	insertQuery := global.CassSessioQBank.Session.Query(qbankz.QuestionMainTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	responseModel := model.QuestionBankQuestion{
		ID:             &cassandraQuestionBank.ID,
		QbmID:          input.QbmID,
		Type:           input.Type,
		Difficulty:     input.Difficulty,
		Hint:           input.Hint,
		Status:         input.Status,
		Attachment:     &getUrl,
		AttachmentType: input.AttachmentType,
		CreatedBy:      input.CreatedBy,
		UpdatedBy:      input.UpdatedBy,
		CreatedAt:      &created,
		UpdatedAt:      &created,
		Description:   input.Description,
	}
	return &responseModel, nil
}

func UpdateQuestionBankQuestion(ctx context.Context, input *model.QuestionBankQuestionInput) (*model.QuestionBankQuestion, error) {
	log.Info("UpdateQuestionBankQuestion called")
	if input.ID == nil {
		return nil, fmt.Errorf("question id not found")
	}
	if input.QbmID == nil {
		return nil, fmt.Errorf("question bank main id not found")
	}

	cassandraQuestionBank := qbankz.QuestionMain{
		ID: *input.ID,
	}
	banks := []qbankz.QuestionMain{}
	getQuery := global.CassSessioQBank.Session.Query(qbankz.QuestionMainTable.Get()).BindMap(qb.M{"id": cassandraQuestionBank.ID})
	if err := getQuery.SelectRelease(&banks); err != nil {
		return nil, err
	}
	if len(banks) == 0 {
		return nil, fmt.Errorf("question bank not found")
	}
	cassandraQuestionBank = banks[0]
	updatedCols := []string{}
	if input.Description != nil {
		cassandraQuestionBank.Description = *input.Description
		updatedCols = append(updatedCols, "description")
	}
	if input.Type != nil {
		cassandraQuestionBank.Type = *input.Type
		updatedCols = append(updatedCols, "type")
	}
	if input.Difficulty != nil {
		cassandraQuestionBank.Difficulty = *input.Difficulty
		updatedCols = append(updatedCols, "difficulty_score")
	}
	if input.Hint != nil {
		cassandraQuestionBank.Hint = *input.Hint
		updatedCols = append(updatedCols, "hint")
	}
	if input.Status != nil {
		cassandraQuestionBank.Status = *input.Status
		updatedCols = append(updatedCols, "status")
	}
	if input.AttachmentType != nil {
		cassandraQuestionBank.AttachmentType = *input.AttachmentType
		updatedCols = append(updatedCols, "attachment_type")
	}
	if input.UpdatedBy != nil {
		cassandraQuestionBank.UpdatedBy = *input.UpdatedBy
		updatedCols = append(updatedCols, "updated_by")
	}
	if input.QbmID != nil {
		cassandraQuestionBank.QbmId = *input.QbmID
		updatedCols = append(updatedCols, "qbm_id")
	}
	updatedAt := time.Now().Unix()
	cassandraQuestionBank.UpdatedAt = updatedAt
	updatedCols = append(updatedCols, "updated_at")
	if input.File != nil {
		bucketPath := "question_banks/" + cassandraQuestionBank.QbmId + "/" + cassandraQuestionBank.ID + "/" + input.File.Filename
		storageC := bucket.NewStorageHandler()
		gproject := googleprojectlib.GetGoogleProjectID()
		err := storageC.InitializeStorageClient(ctx, gproject)
		if err != nil {
			log.Errorf("Failed to upload video to course topic: %v", err.Error())
			return nil, err
		}
		writer, err := storageC.UploadToGCS(ctx, bucketPath, map[string]string{})
		if err != nil {
			log.Errorf("Failed to upload video to course topic: %v", err.Error())
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
	upStms, uNames := qbankz.QuestionMainTable.Update(updatedCols...)
	updateQuery := global.CassSessioQBank.Session.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)

	responseModel := model.QuestionBankQuestion{
		ID:             &cassandraQuestionBank.ID,
		QbmID:          input.QbmID,
		Type:           input.Type,
		Difficulty:     input.Difficulty,
		Hint:           input.Hint,
		Status:         input.Status,
		Attachment:     &cassandraQuestionBank.Attachment,
		AttachmentType: input.AttachmentType,
		CreatedBy:      input.CreatedBy,
		UpdatedBy:      input.UpdatedBy,
		CreatedAt:      &created,
		UpdatedAt:      &updated,
	}
	return &responseModel, nil
}
