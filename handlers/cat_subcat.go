package handlers

import (
	"context"

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
