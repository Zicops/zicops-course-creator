package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/graph/model"
)

func AddCategory(ctx context.Context, category []*string) (*bool, error) {
	log.Infof("AddCategory called")
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session

	catgories := make([]string, len(category))
	isSuccess := false
	for i, c := range category {
		catgories[i] = *c
		cassandraCategory := coursez.Cat{
			Name: catgories[i],
		}
		insertQuery := CassSession.Query(coursez.CatTable.Insert()).BindStruct(cassandraCategory)
		if err := insertQuery.ExecRelease(); err != nil {
			return &isSuccess, err
		}
	}
	isSuccess = true
	return &isSuccess, nil
}

func AddSubCategory(ctx context.Context, category []*string) (*bool, error) {
	log.Infof("AddSubCategory called")
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session

	catgories := make([]string, len(category))
	isSuccess := false
	for i, c := range category {
		catgories[i] = *c
		cassandraCategory := coursez.SubCategory{
			Name: catgories[i],
		}
		insertQuery := CassSession.Query(coursez.SubCatTable.Insert()).BindStruct(cassandraCategory)
		if err := insertQuery.ExecRelease(); err != nil {
			return &isSuccess, err
		}
	}
	isSuccess = true
	return &isSuccess, nil
}

func AddCategorySubMap(ctx context.Context, category *string, subCategory []*string) (*bool, error) {
	log.Infof("AddCategorySubMap called")
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session

	isSuccess := false
	for _, subCat := range subCategory {
		guid := xid.New()
		currentSubCat := *subCat
		currentCat := *category
		cassandraCategory := coursez.CatSubMap{
			ID:          guid.String(),
			Category:    currentCat,
			SubCategory: currentSubCat,
		}
		insertQuery := CassSession.Query(coursez.CatSubMapTable.Insert()).BindStruct(cassandraCategory)
		if err := insertQuery.ExecRelease(); err != nil {
			return &isSuccess, err
		}
	}
	isSuccess = true
	return &isSuccess, nil
}

func AddCatMain(ctx context.Context, input []*model.CatMainInput) ([]*model.CatMain, error) {
	log.Infof("AddCatMain called")
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session

	catMain := make([]*model.CatMain, len(input))
	for i, cc := range input {
		c := cc
		guid := xid.New()
		cassandraCategory := coursez.CatMain{
			ID:          guid.String(),
			Name:        *c.Name,
			Description: *c.Description,
			Code:        *c.Code,
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
			CreatedBy:   *c.CreatedBy,
			UpdatedBy:   *c.UpdatedBy,
			IsActive:    *c.IsActive,
		}
		insertQuery := CassSession.Query(coursez.CatMainTable.Insert()).BindStruct(cassandraCategory)
		if err := insertQuery.ExecRelease(); err != nil {
			return nil, err
		}
		created := strconv.FormatInt(cassandraCategory.CreatedAt, 10)
		updated := strconv.FormatInt(cassandraCategory.UpdatedAt, 10)
		catMain[i] = &model.CatMain{
			ID:          &cassandraCategory.ID,
			Name:        &cassandraCategory.Name,
			Description: &cassandraCategory.Description,
			Code:        &cassandraCategory.Code,
			CreatedAt:   &created,
			UpdatedAt:   &updated,
			CreatedBy:   &cassandraCategory.CreatedBy,
			UpdatedBy:   &cassandraCategory.UpdatedBy,
			IsActive:    &cassandraCategory.IsActive,
		}
	}
	return catMain, nil
}

func AddSubCatMain(ctx context.Context, input []*model.SubCatMainInput) ([]*model.SubCatMain, error) {
	log.Infof("AddSubCatMain called")
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session

	subCatMain := make([]*model.SubCatMain, len(input))
	for i, cc := range input {
		c := cc
		guid := xid.New()
		cassandraCategory := coursez.SubCatMain{
			ID:          guid.String(),
			Name:        *c.Name,
			Description: *c.Description,
			Code:        *c.Code,
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
			CreatedBy:   *c.CreatedBy,
			UpdatedBy:   *c.UpdatedBy,
			IsActive:    *c.IsActive,
			ParentID:    *c.CatID,
		}
		insertQuery := CassSession.Query(coursez.SubCatMainTable.Insert()).BindStruct(cassandraCategory)
		if err := insertQuery.ExecRelease(); err != nil {
			return nil, err
		}
		created := strconv.FormatInt(cassandraCategory.CreatedAt, 10)
		updated := strconv.FormatInt(cassandraCategory.UpdatedAt, 10)
		subCatMain[i] = &model.SubCatMain{
			ID:          &cassandraCategory.ID,
			Name:        &cassandraCategory.Name,
			Description: &cassandraCategory.Description,
			Code:        &cassandraCategory.Code,
			CreatedAt:   &created,
			UpdatedAt:   &updated,
			CreatedBy:   &cassandraCategory.CreatedBy,
			UpdatedBy:   &cassandraCategory.UpdatedBy,
			IsActive:    &cassandraCategory.IsActive,
			CatID:       &cassandraCategory.ParentID,
		}
	}
	return subCatMain, nil
}
