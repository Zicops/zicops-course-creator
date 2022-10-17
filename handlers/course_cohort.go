package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/xid"
	"github.com/scylladb/gocqlx/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/helpers"
)

func AddCourseCohort(ctx context.Context, input *model.CourseCohortInput) (*model.CourseCohort, error) {
	log.Info("AddCourseCohort called")
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	email_creator := claims["email"].(string)
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	expectedComp := 0
	if input.ExpectedCompletion != nil {
		expectedComp = *input.ExpectedCompletion
	}
	guid := xid.New()
	cassandraQuestionBank := coursez.CourseCohortMapping{
		ID:                     guid.String(),
		CourseID:               *input.CourseID,
		CourseType:             *input.CourseType,
		CohortID:               *input.CohortID,
		CourseStatus:           *input.CourseStatus,
		LspId:                  *input.LspID,
		IsMandatory:            *input.IsMandatory,
		AddedBy:                *input.AddedBy,
		IsActive:               *input.IsActive,
		CreatedBy:              email_creator,
		UpdatedBy:              email_creator,
		CreatedAt:              time.Now().Unix(),
		UpdatedAt:              time.Now().Unix(),
		CohortCode:             *input.CohortCode,
		ExpectedCompletionDays: expectedComp,
	}
	insertQuery := CassSession.Query(coursez.CourseCohortTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.CourseCohort{
		ID:                 &cassandraQuestionBank.ID,
		CourseID:           input.CourseID,
		CourseType:         input.CourseType,
		CohortID:           input.CohortID,
		CourseStatus:       input.CourseStatus,
		LspID:              input.LspID,
		IsMandatory:        input.IsMandatory,
		AddedBy:            input.AddedBy,
		IsActive:           input.IsActive,
		CreatedBy:          input.CreatedBy,
		UpdatedBy:          input.UpdatedBy,
		CreatedAt:          &created,
		UpdatedAt:          &updated,
		CohortCode:         input.CohortCode,
		ExpectedCompletion: input.ExpectedCompletion,
	}

	return &responseModel, nil
}

func UpdateCourseCohort(ctx context.Context, input *model.CourseCohortInput) (*model.CourseCohort, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("course cohort id is required")
	}
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	log.Info("AddCourseCohort called")
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	email_creator := claims["email"].(string)
	lspId := claims["lsp_id"].(string)
	CassSession := session
	cassandraQuestionBank := *GetCourseCohort(ctx, *input.ID, lspId, CassSession)
	updatedCols := []string{}
	if input.CohortID != nil && cassandraQuestionBank.CohortID != *input.CohortID {
		cassandraQuestionBank.CohortID = *input.CohortID
		updatedCols = append(updatedCols, "cohortid")
	}
	if input.CourseID != nil && cassandraQuestionBank.CourseID != *input.CourseID {
		cassandraQuestionBank.CourseID = *input.CourseID
		updatedCols = append(updatedCols, "courseid")
	}
	if email_creator != "" && cassandraQuestionBank.UpdatedBy != email_creator {
		cassandraQuestionBank.UpdatedBy = email_creator
		updatedCols = append(updatedCols, "updated_by")
	}
	if input.CourseStatus != nil && cassandraQuestionBank.CourseStatus != *input.CourseStatus {
		cassandraQuestionBank.CourseStatus = *input.CourseStatus
		updatedCols = append(updatedCols, "course_status")
	}
	if input.CourseType != nil && cassandraQuestionBank.CourseType != *input.CourseType {
		cassandraQuestionBank.CourseType = *input.CourseType
		updatedCols = append(updatedCols, "course_type")
	}
	if input.IsMandatory != nil && cassandraQuestionBank.IsMandatory != *input.IsMandatory {
		cassandraQuestionBank.IsMandatory = *input.IsMandatory
		updatedCols = append(updatedCols, "is_mandatory")
	}
	if input.AddedBy != nil && cassandraQuestionBank.AddedBy != *input.AddedBy {
		cassandraQuestionBank.AddedBy = *input.AddedBy
		updatedCols = append(updatedCols, "added_by")
	}
	if input.CohortCode != nil && cassandraQuestionBank.CohortCode != *input.CohortCode {
		cassandraQuestionBank.CohortCode = *input.CohortCode
		updatedCols = append(updatedCols, "cohort_code")
	}
	if input.ExpectedCompletion != nil && cassandraQuestionBank.ExpectedCompletionDays != *input.ExpectedCompletion {
		cassandraQuestionBank.ExpectedCompletionDays = *input.ExpectedCompletion
		updatedCols = append(updatedCols, "expected_completion_days")
	}
	if len(updatedCols) > 0 {
		updatedAt := time.Now().Unix()
		cassandraQuestionBank.UpdatedAt = updatedAt
		updatedCols = append(updatedCols, "updated_at")
		upStms, uNames := coursez.CourseCohortTable.Update(updatedCols...)
		updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
		if err := updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.CourseCohort{
		ID:                 &cassandraQuestionBank.ID,
		CourseID:           input.CourseID,
		CourseType:         input.CourseType,
		CohortID:           input.CohortID,
		CourseStatus:       input.CourseStatus,
		LspID:              input.LspID,
		IsMandatory:        input.IsMandatory,
		AddedBy:            input.AddedBy,
		IsActive:           input.IsActive,
		CreatedBy:          input.CreatedBy,
		UpdatedBy:          input.UpdatedBy,
		CreatedAt:          &created,
		UpdatedAt:          &updated,
		CohortCode:         input.CohortCode,
		ExpectedCompletion: input.ExpectedCompletion,
	}

	return &responseModel, nil
}

func GetCourseCohort(ctx context.Context, courseID string, lspID string, session *gocqlx.Session) *coursez.CourseCohortMapping {
	chapters := []coursez.CourseCohortMapping{}
	getQueryStr := fmt.Sprintf("SELECT * FROM coursez.course_cohort_mapping WHERE id='%s' and lsp_id='%s' and is_active=true", courseID, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}
