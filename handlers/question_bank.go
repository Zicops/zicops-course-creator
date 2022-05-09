package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/xid"
	"github.com/scylladb/gocqlx/v2/qb"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/qbankz"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
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

	insertQuery := global.CassSession.Session.Query(qbankz.QuestionBankMainTable.Insert()).BindStruct(cassandraQuestionBank)
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
	getQuery := global.CassSession.Session.Query(qbankz.QuestionBankMainTable.Get()).BindMap(qb.M{"id": cassandraQuestionBank.ID})
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
	if len(updatedCols) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	upStms, uNames := qbankz.QuestionBankMainTable.Update(updatedCols...)
	updateQuery := global.CassSession.Session.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
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
