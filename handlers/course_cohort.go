package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/xid"
	"github.com/scylladb/gocqlx/v2/qb"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
)

func AddCourseCohort(ctx context.Context, input *model.CourseCohortInput) (*model.CourseCohort, error) {
	log.Info("AddCourseCohort called")
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	global.CassSessioQBank = session
	guid := xid.New()
	cassandraQuestionBank := coursez.CourseCohortMapping{
		ID:           guid.String(),
		CourseID:     *input.CourseID,
		CourseType:   *input.CourseType,
		CohortID:     *input.CohortID,
		CourseStatus: *input.CourseStatus,
		LspID:        *input.LspID,
		IsMandatory:  *input.IsMandatory,
		AddedBy:      *input.AddedBy,
		IsActive:     *input.IsActive,
		CreatedBy:    *input.CreatedBy,
		UpdatedBy:    *input.UpdatedBy,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
		CohortCode:   *input.CohortCode,
	}
	insertQuery := global.CassSessioQBank.Query(coursez.CourseCohortTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.CourseCohort{
		ID:           &cassandraQuestionBank.ID,
		CourseID:     input.CourseID,
		CourseType:   input.CourseType,
		CohortID:     input.CohortID,
		CourseStatus: input.CourseStatus,
		LspID:        input.LspID,
		IsMandatory:  input.IsMandatory,
		AddedBy:      input.AddedBy,
		IsActive:     input.IsActive,
		CreatedBy:    input.CreatedBy,
		UpdatedBy:    input.UpdatedBy,
		CreatedAt:    &created,
		UpdatedAt:    &updated,
		CohortCode:   input.CohortCode,
	}

	return &responseModel, nil
}

func UpdateCourseCohort(ctx context.Context, input *model.CourseCohortInput) (*model.CourseCohort, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("course cohort id is required")
	}
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	global.CassSessioQBank = session
	cassandraQuestionBank := coursez.CourseCohortMapping{
		ID: *input.ID,
	}
	banks := []coursez.CourseCohortMapping{}
	getQuery := global.CassSessioQBank.Query(coursez.CourseCohortTable.Get()).BindMap(qb.M{"id": cassandraQuestionBank.ID})
	if err := getQuery.SelectRelease(&banks); err != nil {
		return nil, err
	}
	if len(banks) == 0 {
		return nil, fmt.Errorf("course cohorts not found")
	}
	cassandraQuestionBank = banks[0]
	updatedCols := []string{}
	if input.CohortID != nil {
		cassandraQuestionBank.CohortID = *input.CohortID
		updatedCols = append(updatedCols, "cohortid")
	}
	if input.CourseID != nil {
		cassandraQuestionBank.CourseID = *input.CourseID
		updatedCols = append(updatedCols, "courseid")
	}
	if input.IsActive != nil {
		cassandraQuestionBank.IsActive = *input.IsActive
		updatedCols = append(updatedCols, "is_active")
	}
	if input.CreatedBy != nil {
		cassandraQuestionBank.CreatedBy = *input.CreatedBy
		updatedCols = append(updatedCols, "created_by")
	}
	if input.UpdatedBy != nil {
		cassandraQuestionBank.UpdatedBy = *input.UpdatedBy
		updatedCols = append(updatedCols, "updated_by")
	}
	if input.CourseStatus != nil {
		cassandraQuestionBank.CourseStatus = *input.CourseStatus
		updatedCols = append(updatedCols, "course_status")
	}
	if input.CourseType != nil {
		cassandraQuestionBank.CourseType = *input.CourseType
		updatedCols = append(updatedCols, "course_type")
	}
	if input.IsMandatory != nil {
		cassandraQuestionBank.IsMandatory = *input.IsMandatory
		updatedCols = append(updatedCols, "is_mandatory")
	}
	if input.LspID != nil {
		cassandraQuestionBank.LspID = *input.LspID
		updatedCols = append(updatedCols, "lsp_id")
	}
	if input.AddedBy != nil {
		cassandraQuestionBank.AddedBy = *input.AddedBy
		updatedCols = append(updatedCols, "added_by")
	}
	if input.CohortCode != nil {
		cassandraQuestionBank.CohortCode = *input.CohortCode
		updatedCols = append(updatedCols, "cohort_code")
	}
	updatedAt := time.Now().Unix()
	cassandraQuestionBank.UpdatedAt = updatedAt
	updatedCols = append(updatedCols, "updated_at")
	if len(updatedCols) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	upStms, uNames := coursez.CourseCohortTable.Update(updatedCols...)
	updateQuery := global.CassSessioQBank.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.CourseCohort{
		ID:           &cassandraQuestionBank.ID,
		CourseID:     input.CourseID,
		CourseType:   input.CourseType,
		CohortID:     input.CohortID,
		CourseStatus: input.CourseStatus,
		LspID:        input.LspID,
		IsMandatory:  input.IsMandatory,
		AddedBy:      input.AddedBy,
		IsActive:     input.IsActive,
		CreatedBy:    input.CreatedBy,
		UpdatedBy:    input.UpdatedBy,
		CreatedAt:    &created,
		UpdatedAt:    &updated,
		CohortCode:   input.CohortCode,
	}

	return &responseModel, nil
}
