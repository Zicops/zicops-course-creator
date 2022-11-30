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

func AddExamCohort(ctx context.Context, input *model.ExamCohortInput) (*model.ExamCohort, error) {
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
	cassandraQuestionBank := qbankz.ExamCohort{
		ID:        uuid.New().String(),
		ExamID:    *input.ExamID,
		IsActive:  *input.IsActive,
		CreatedBy: email_creator,
		UpdatedBy: email_creator,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		CohortID:  *input.CohortID,
		LspId:     lspID,
	}
	insertQuery := CassSession.Query(qbankz.ExamCohortTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.ExamCohort{
		ID:        &cassandraQuestionBank.ID,
		ExamID:    input.ExamID,
		IsActive:  input.IsActive,
		CreatedBy: input.CreatedBy,
		UpdatedBy: input.UpdatedBy,
		CreatedAt: &created,
		UpdatedAt: &updated,
		CohortID:  input.CohortID,
	}
	return &responseModel, nil
}

func UpdateExamCohort(ctx context.Context, input *model.ExamCohortInput) (*model.ExamCohort, error) {
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
	cassandraQuestionBank := *GetExamCohort(ctx, *input.ID, lspID, CassSession)
	updatedCols := []string{}
	if input.ExamID != nil && cassandraQuestionBank.ExamID != *input.ExamID {
		cassandraQuestionBank.ExamID = *input.ExamID
		updatedCols = append(updatedCols, "exam_id")
	}
	if input.CohortID != nil && cassandraQuestionBank.CohortID != *input.CohortID {
		cassandraQuestionBank.CohortID = *input.CohortID
		updatedCols = append(updatedCols, "cohort_id")
	}
	if email_creator != "" && cassandraQuestionBank.UpdatedBy != email_creator {
		cassandraQuestionBank.UpdatedBy = email_creator
		updatedCols = append(updatedCols, "updated_by")
	}

	if len(updatedCols) > 0 {
		updatedAt := time.Now().Unix()
		cassandraQuestionBank.UpdatedAt = updatedAt
		updatedCols = append(updatedCols, "updated_at")
		upStms, uNames := qbankz.ExamCohortTable.Update(updatedCols...)
		updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
		if err := updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.ExamCohort{
		ID:        &cassandraQuestionBank.ID,
		ExamID:    input.ExamID,
		IsActive:  input.IsActive,
		CreatedBy: input.CreatedBy,
		UpdatedBy: input.UpdatedBy,
		CreatedAt: &created,
		UpdatedAt: &updated,
		CohortID:  input.CohortID,
	}
	return &responseModel, nil
}

func GetExamCohort(ctx context.Context, courseID string, lspID string, session *gocqlx.Session) *qbankz.ExamCohort {
	chapters := []qbankz.ExamCohort{}
	getQueryStr := fmt.Sprintf("SELECT * FROM qbankz.exam_cohort WHERE id='%s' and lsp_id='%s' and is_active=true", courseID, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}
