package handlers

import (
	"context"

	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-course-creator/global"
)

func AddCategory(ctx context.Context, category []*string) (*bool, error) {
	log.Infof("AddCategory called")
	catgories := make([]string, len(category))
	isSuccess := false
	for i, c := range category {
		catgories[i] = *c
		cassandraCategory := coursez.Cat{
			Name: catgories[i],
		}
		insertQuery := global.CassSession.Session.Query(coursez.CatTable.Insert()).BindStruct(cassandraCategory)
		if err := insertQuery.ExecRelease(); err != nil {
			return &isSuccess, err
		}
	}
	isSuccess = true
	return &isSuccess, nil
}

func AddSubCategory(ctx context.Context, category []*string) (*bool, error) {
	log.Infof("AddCategory called")
	catgories := make([]string, len(category))
	isSuccess := false
	for i, c := range category {
		catgories[i] = *c
		cassandraCategory := coursez.SubCategory{
			Name: catgories[i],
		}
		insertQuery := global.CassSession.Session.Query(coursez.SubCatTable.Insert()).BindStruct(cassandraCategory)
		if err := insertQuery.ExecRelease(); err != nil {
			return &isSuccess, err
		}
	}
	isSuccess = true
	return &isSuccess, nil
}

func AddCategorySubMap(ctx context.Context, category *string, subCategory []*string) (*bool, error) {
	log.Infof("AddCategorySubMap called")
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
		insertQuery := global.CassSession.Session.Query(coursez.CatSubMapTable.Insert()).BindStruct(cassandraCategory)
		if err := insertQuery.ExecRelease(); err != nil {
			return &isSuccess, err
		}
	}
	isSuccess = true
	return &isSuccess, nil
}
