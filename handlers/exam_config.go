package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/scylladb/gocqlx/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/qbankz"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/helpers"
)

func AddExamConfiguration(ctx context.Context, input *model.ExamConfigurationInput) (*model.ExamConfiguration, error) {
	log.Info("ExamInstructionsCreate called")

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
	email_creator := claims["email"].(string)
	cassandraQuestionBank := qbankz.ExamConfig{
		ID:           uuid.New().String(),
		ExamID:       *input.ExamID,
		IsActive:     *input.IsActive,
		CreatedBy:    email_creator,
		UpdatedBy:    email_creator,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
		Shuffle:      *input.Shuffle,
		DisplayHints: *input.DisplayHints,
		ShowAnswer:   *input.ShowAnswer,
		ShowResult:   *input.ShowResult,
		LspId:        lspID,
	}
	insertQuery := CassSession.Query(qbankz.ExamConfigTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.ExamConfiguration{
		ID:           &cassandraQuestionBank.ID,
		ExamID:       input.ExamID,
		IsActive:     input.IsActive,
		CreatedBy:    input.CreatedBy,
		UpdatedBy:    input.UpdatedBy,
		CreatedAt:    &created,
		UpdatedAt:    &updated,
		Shuffle:      input.Shuffle,
		DisplayHints: input.DisplayHints,
		ShowAnswer:   input.ShowAnswer,
		ShowResult:   input.ShowResult,
	}
	return &responseModel, nil
}

func UpdateExamConfiguration(ctx context.Context, input *model.ExamConfigurationInput) (*model.ExamConfiguration, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("exam schedule id not found")
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
	cassandraQuestionBank := *GetExamConfig(ctx, *input.ID, lspID, CassSession)
	updatedCols := []string{}
	if input.ExamID != nil && cassandraQuestionBank.ExamID != *input.ExamID {
		cassandraQuestionBank.ExamID = *input.ExamID
		updatedCols = append(updatedCols, "exam_id")
	}
	if input.Shuffle != nil && cassandraQuestionBank.Shuffle != *input.Shuffle {
		cassandraQuestionBank.Shuffle = *input.Shuffle
		updatedCols = append(updatedCols, "shuffle_questions")
	}
	if input.DisplayHints != nil && cassandraQuestionBank.DisplayHints != *input.DisplayHints {
		cassandraQuestionBank.DisplayHints = *input.DisplayHints
		updatedCols = append(updatedCols, "display_hints")
	}
	if input.ShowAnswer != nil && cassandraQuestionBank.ShowAnswer != *input.ShowAnswer {
		cassandraQuestionBank.ShowAnswer = *input.ShowAnswer
		updatedCols = append(updatedCols, "show_answer")
	}
	if input.ShowResult != nil && cassandraQuestionBank.ShowResult != *input.ShowResult {
		cassandraQuestionBank.ShowResult = *input.ShowResult
		updatedCols = append(updatedCols, "show_result")
	}
	if email_creator != "" && cassandraQuestionBank.UpdatedBy != email_creator {
		cassandraQuestionBank.UpdatedBy = email_creator
		updatedCols = append(updatedCols, "updated_by")
	}

	if len(updatedCols) > 0 {
		updatedAt := time.Now().Unix()
		cassandraQuestionBank.UpdatedAt = updatedAt
		updatedCols = append(updatedCols, "updated_at")
		upStms, uNames := qbankz.ExamConfigTable.Update(updatedCols...)
		updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
		if err := updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.ExamConfiguration{
		ID:           &cassandraQuestionBank.ID,
		ExamID:       input.ExamID,
		IsActive:     input.IsActive,
		CreatedBy:    input.CreatedBy,
		UpdatedBy:    input.UpdatedBy,
		CreatedAt:    &created,
		UpdatedAt:    &updated,
		Shuffle:      input.Shuffle,
		DisplayHints: input.DisplayHints,
		ShowAnswer:   input.ShowAnswer,
		ShowResult:   input.ShowResult,
	}
	return &responseModel, nil
}

func GetExamConfig(ctx context.Context, id string, lspID string, session *gocqlx.Session) *qbankz.ExamConfig {
	chapters := []qbankz.ExamConfig{}
	getQueryStr := fmt.Sprintf("SELECT * FROM qbankz.exam_config WHERE id='%s' and lsp_id='%s' and is_active=true", id, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}
