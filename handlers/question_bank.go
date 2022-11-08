package handlers

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
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

func QuestionBankCreate(ctx context.Context, input *model.QuestionBankInput) (*model.QuestionBank, error) {
	log.Info("QuestionBankCreate called")
	guid := xid.New()
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	words := []string{}
	if input.Name != nil {
		name := strings.ToLower(*input.Name)
		wordsLocal := strings.Split(name, " ")
		words = append(words, wordsLocal...)
	}
	lspID := claims["lsp_id"].(string)
	email_creator := claims["email"].(string)
	cassandraQuestionBank := qbankz.QuestionBankMain{
		ID:          guid.String(),
		Name:        *input.Name,
		Words:       words,
		Description: *input.Description,
		Category:    *input.Category,
		SubCategory: *input.SubCategory,
		IsActive:    *input.IsActive,
		IsDefault:   *input.IsDefault,
		Owner:       *input.Owner,
		CreatedBy:   email_creator,
		UpdatedBy:   email_creator,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		LspId:       lspID,
	}

	insertQuery := CassSession.Query(qbankz.QuestionBankMainTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	responseModel := model.QuestionBank{
		ID:          &cassandraQuestionBank.ID,
		Name:        input.Name,
		Description: input.Description,
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
		return nil, fmt.Errorf("provide question bank id")
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
	cassandraQuestionBank := *GetQBank(ctx, *input.ID, lspID, CassSession)
	updatedCols := []string{}
	if input.Description != nil && cassandraQuestionBank.Description != *input.Description {
		cassandraQuestionBank.Description = *input.Description
		updatedCols = append(updatedCols, "description")
	}
	if input.Name != nil && cassandraQuestionBank.Name != *input.Name {
		words := []string{}
		if input.Name != nil {
			name := strings.ToLower(*input.Name)
			wordsLocal := strings.Split(name, " ")
			words = append(words, wordsLocal...)
		}
		cassandraQuestionBank.Name = *input.Name
		cassandraQuestionBank.Words = words
		updatedCols = append(updatedCols, "words")
		updatedCols = append(updatedCols, "name")
	}
	if input.Category != nil && cassandraQuestionBank.Category != *input.Category {
		cassandraQuestionBank.Category = *input.Category
		updatedCols = append(updatedCols, "category")
	}
	if input.SubCategory != nil && cassandraQuestionBank.SubCategory != *input.SubCategory {
		cassandraQuestionBank.SubCategory = *input.SubCategory
		updatedCols = append(updatedCols, "sub_category")
	}
	if input.IsDefault != nil && cassandraQuestionBank.IsDefault != *input.IsDefault {
		cassandraQuestionBank.IsDefault = *input.IsDefault
		updatedCols = append(updatedCols, "is_default")
	}
	if input.Owner != nil && cassandraQuestionBank.Owner != *input.Owner {
		cassandraQuestionBank.Owner = *input.Owner
		updatedCols = append(updatedCols, "owner")
	}
	if email_creator != "" && cassandraQuestionBank.UpdatedBy != email_creator {
		cassandraQuestionBank.UpdatedBy = email_creator
		updatedCols = append(updatedCols, "updated_by")
	}
	updatedAt := time.Now().Unix()
	if len(updatedCols) > 0 {
		cassandraQuestionBank.UpdatedAt = updatedAt
		updatedCols = append(updatedCols, "updated_at")
		upStms, uNames := qbankz.QuestionBankMainTable.Update(updatedCols...)
		updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
		if err := updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.QuestionBank{
		ID:          &cassandraQuestionBank.ID,
		Name:        input.Name,
		Description: input.Description,
		Owner:       input.Owner,
		Category:    input.Category,
		SubCategory: input.SubCategory,
		IsActive:    input.IsActive,
		IsDefault:   input.IsDefault,
		CreatedBy:   input.CreatedBy,
		UpdatedBy:   input.UpdatedBy,
		CreatedAt:   &created,
		UpdatedAt:   &updated,
	}
	return &responseModel, nil
}

func AddQuestionBankQuestion(ctx context.Context, input *model.QuestionBankQuestionInput) (*model.QuestionBankQuestion, error) {
	log.Info("AddQuestionBankQuestion called")
	if input.QbmID == nil {
		return nil, fmt.Errorf("question bank main id not found")
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
	getUrl := ""
	cassandraQuestionBank := qbankz.QuestionMain{
		ID:             guid.String(),
		Name:           *input.Name,
		Description:    *input.Description,
		QbmId:          *input.QbmID,
		Type:           *input.Type,
		Difficulty:     *input.Difficulty,
		Hint:           *input.Hint,
		Status:         *input.Status,
		AttachmentType: *input.AttachmentType,
		CreatedBy:      email_creator,
		UpdatedBy:      email_creator,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
		LspId:          lspID,
		IsActive:       true,
	}
	if input.File != nil {
		bucketPath := "question_banks/" + cassandraQuestionBank.QbmId + "/" + cassandraQuestionBank.ID + "/" + input.File.Filename
		storageC := bucket.NewStorageHandler()
		gproject := googleprojectlib.GetGoogleProjectID()
		err := storageC.InitializeStorageClient(ctx, gproject, lspID)
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
		getUrl = storageC.GetSignedURLForObject(bucketPath)
		cassandraQuestionBank.Attachment = getUrl
		cassandraQuestionBank.AttachmentBucket = bucketPath
	}
	insertQuery := CassSession.Query(qbankz.QuestionMainTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.QuestionBankQuestion{
		ID:             &cassandraQuestionBank.ID,
		Name:           input.Name,
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
		UpdatedAt:      &updated,
		Description:    input.Description,
	}
	return &responseModel, nil
}

func BulkAddQuestionBankQuestions(ctx context.Context, qbID string, qfile graphql.Upload) (bool, error) {
	log.Info("AddQuestionBankQuestion called")
	if qbID == "" {
		return false, fmt.Errorf("question bank main id not found")
	}
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return false, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return false, err
	}
	lspID := claims["lsp_id"].(string)
	if lspID == "" {
		return false, fmt.Errorf("lsp id not found")
	}
	email_creator := claims["email"].(string)
	// read csv file and parse it as QuestionBankQuestionInput
	fileBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBuffer, qfile.File); err != nil {
		return false, err
	}
	currentBytes := fileBuffer.Bytes()
	csvReader := csv.NewReader(bytes.NewReader(currentBytes))
	csvReader.Comma = ','
	csvReader.FieldsPerRecord = -1
	csvData, err := csvReader.ReadAll()
	if err != nil {
		return false, err
	}
	for i, row := range csvData {
		if i < 3 {
			continue
		}
		go populateQuestionBankQuestions(ctx, row, i, qbID, email_creator, lspID, CassSession)
	}
	return true, nil
}

func populateQuestionBankQuestions(ctx context.Context, row []string, i int, qbID string, email_creator string, lspID string, CassSession *gocqlx.Session) {
	col9 := 9
	correctOptions := strings.ToLower(row[col9])
	guid := xid.New()
	// row 3 to int
	difficulty, err := strconv.Atoi(row[2])
	if err != nil {
		log.Errorf("Failed to convert difficulty to int: %v", err.Error())
		return
	}

	cassandraQuestionBank := qbankz.QuestionMain{
		ID:               guid.String(),
		Name:             row[3],
		Description:      row[3],
		QbmId:            qbID,
		Type:             row[1],
		Difficulty:       difficulty,
		Hint:             row[4],
		Status:           "active",
		AttachmentType:   "",
		Attachment:       "",
		AttachmentBucket: "",
		CreatedBy:        email_creator,
		UpdatedBy:        email_creator,
		CreatedAt:        time.Now().Unix(),
		UpdatedAt:        time.Now().Unix(),
		LspId:            lspID,
		IsActive:         true,
	}
	insertQuery := CassSession.Query(qbankz.QuestionMainTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		log.Errorf("Failed to insert question bank question: %v", err.Error())
		return
	}
	// create optionsmain
	isCorrect1 := false
	isCorrect2 := false
	isCorrect3 := false
	isCorrect4 := false

	if strings.Contains(correctOptions, "a") && i == 5 {
		isCorrect1 = true
	}
	if strings.Contains(correctOptions, "b") && i == 6 {
		isCorrect2 = true
	}
	if strings.Contains(correctOptions, "c") && i == 7 {
		isCorrect3 = true
	}
	if strings.Contains(correctOptions, "d") && i == 8 {
		isCorrect4 = true
	}
	questionOption1 := qbankz.OptionsMain{
		ID:               xid.New().String(),
		QmId:             guid.String(),
		Description:      row[5],
		LspId:            lspID,
		IsActive:         true,
		CreatedBy:        email_creator,
		UpdatedBy:        email_creator,
		CreatedAt:        time.Now().Unix(),
		UpdatedAt:        time.Now().Unix(),
		AttachmentType:   "",
		Attachment:       "",
		AttachmentBucket: "",
		IsCorrect:        isCorrect1,
	}
	questionOption2 := qbankz.OptionsMain{
		ID:               xid.New().String(),
		QmId:             guid.String(),
		Description:      row[6],
		LspId:            lspID,
		IsActive:         true,
		CreatedBy:        email_creator,
		UpdatedBy:        email_creator,
		CreatedAt:        time.Now().Unix(),
		UpdatedAt:        time.Now().Unix(),
		AttachmentType:   "",
		Attachment:       "",
		AttachmentBucket: "",
		IsCorrect:        isCorrect2,
	}
	questionOption3 := qbankz.OptionsMain{
		ID:               xid.New().String(),
		QmId:             guid.String(),
		Description:      row[7],
		LspId:            lspID,
		IsActive:         true,
		CreatedBy:        email_creator,
		UpdatedBy:        email_creator,
		CreatedAt:        time.Now().Unix(),
		UpdatedAt:        time.Now().Unix(),
		AttachmentType:   "",
		Attachment:       "",
		AttachmentBucket: "",
		IsCorrect:        isCorrect3,
	}
	questionOption4 := qbankz.OptionsMain{
		ID:               xid.New().String(),
		QmId:             guid.String(),
		Description:      row[8],
		LspId:            lspID,
		IsActive:         true,
		CreatedBy:        email_creator,
		UpdatedBy:        email_creator,
		CreatedAt:        time.Now().Unix(),
		UpdatedAt:        time.Now().Unix(),
		AttachmentType:   "",
		Attachment:       "",
		AttachmentBucket: "",
		IsCorrect:        isCorrect4,
	}
	insertQuery = CassSession.Query(qbankz.OptionsMainTable.Insert()).BindStruct(questionOption1)
	if err := insertQuery.ExecRelease(); err != nil {
		log.Errorf("Failed to insert question bank question option: %v", err.Error())
		return
	}
	insertQuery = CassSession.Query(qbankz.OptionsMainTable.Insert()).BindStruct(questionOption2)
	if err := insertQuery.ExecRelease(); err != nil {
		log.Errorf("Failed to insert question bank question option: %v", err.Error())
		return
	}
	insertQuery = CassSession.Query(qbankz.OptionsMainTable.Insert()).BindStruct(questionOption3)
	if err := insertQuery.ExecRelease(); err != nil {
		log.Errorf("Failed to insert question bank question option: %v", err.Error())
		return
	}
	insertQuery = CassSession.Query(qbankz.OptionsMainTable.Insert()).BindStruct(questionOption4)
	if err := insertQuery.ExecRelease(); err != nil {
		log.Errorf("Failed to insert question bank question option: %v", err.Error())
		return
	}
}

func UpdateQuestionBankQuestion(ctx context.Context, input *model.QuestionBankQuestionInput) (*model.QuestionBankQuestion, error) {
	log.Info("UpdateQuestionBankQuestion called")
	if input.ID == nil {
		return nil, fmt.Errorf("question id not found")
	}
	if input.QbmID == nil {
		return nil, fmt.Errorf("question bank main id not found")
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
	cassandraQuestionBank := *GetQBankQ(ctx, *input.ID, lspID, CassSession)
	updatedCols := []string{}
	if input.Name != nil && *input.Name != cassandraQuestionBank.Name {
		cassandraQuestionBank.Name = *input.Name
		updatedCols = append(updatedCols, "name")
	}
	if input.Description != nil && *input.Description != cassandraQuestionBank.Description {
		cassandraQuestionBank.Description = *input.Description
		updatedCols = append(updatedCols, "description")
	}
	if input.Type != nil && *input.Type != cassandraQuestionBank.Type {
		cassandraQuestionBank.Type = *input.Type
		updatedCols = append(updatedCols, "type")
	}
	if input.Difficulty != nil && *input.Difficulty != cassandraQuestionBank.Difficulty {
		cassandraQuestionBank.Difficulty = *input.Difficulty
		updatedCols = append(updatedCols, "difficulty_score")
	}
	if input.Hint != nil && *input.Hint != cassandraQuestionBank.Hint {
		cassandraQuestionBank.Hint = *input.Hint
		updatedCols = append(updatedCols, "hint")
	}
	if input.Status != nil && *input.Status != cassandraQuestionBank.Status {
		cassandraQuestionBank.Status = *input.Status
		updatedCols = append(updatedCols, "status")
	}
	if input.AttachmentType != nil && *input.AttachmentType != cassandraQuestionBank.AttachmentType {
		cassandraQuestionBank.AttachmentType = *input.AttachmentType
		updatedCols = append(updatedCols, "attachment_type")
	}
	if email_creator != "" && email_creator != cassandraQuestionBank.CreatedBy {
		cassandraQuestionBank.UpdatedBy = email_creator
		updatedCols = append(updatedCols, "updated_by")
	}
	if input.QbmID != nil && *input.QbmID != cassandraQuestionBank.QbmId {
		cassandraQuestionBank.QbmId = *input.QbmID
		updatedCols = append(updatedCols, "qbm_id")
	}
	updatedAt := time.Now().Unix()
	if input.File != nil {
		bucketPath := "question_banks/" + cassandraQuestionBank.QbmId + "/" + cassandraQuestionBank.ID + "/" + input.File.Filename
		storageC := bucket.NewStorageHandler()
		gproject := googleprojectlib.GetGoogleProjectID()
		err := storageC.InitializeStorageClient(ctx, gproject, lspID)
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
	if len(updatedCols) > 0 {
		cassandraQuestionBank.UpdatedAt = updatedAt
		updatedCols = append(updatedCols, "updated_at")
		upStms, uNames := qbankz.QuestionMainTable.Update(updatedCols...)
		updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
		if err := updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)

	responseModel := model.QuestionBankQuestion{
		ID:             &cassandraQuestionBank.ID,
		Name:           &cassandraQuestionBank.Name,
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
		Description:    input.Description,
	}
	return &responseModel, nil
}

func GetQBank(ctx context.Context, courseID string, lspID string, session *gocqlx.Session) *qbankz.QuestionBankMain {
	chapters := []qbankz.QuestionBankMain{}
	getQueryStr := fmt.Sprintf("SELECT * FROM qbankz.question_bank_main WHERE id='%s' and lsp_id='%s' and is_active=true", courseID, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}

func GetQBankQ(ctx context.Context, courseID string, lspID string, session *gocqlx.Session) *qbankz.QuestionMain {
	chapters := []qbankz.QuestionMain{}
	getQueryStr := fmt.Sprintf("SELECT * FROM qbankz.question_main WHERE id='%s' and lsp_id='%s' and is_active=true", courseID, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}
