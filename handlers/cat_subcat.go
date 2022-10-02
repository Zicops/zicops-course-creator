package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/helpers"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
	"github.com/zicops/zicops-course-creator/lib/googleprojectlib"
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
	catMain := make([]*model.CatMain, len(input))
	for i, cc := range input {
		c := cc
		guid := base64.URLEncoding.EncodeToString([]byte(strings.ToLower(*c.Name)))
		imageUrl := ""
		imageBucket := ""
		qryStr := fmt.Sprintf(`SELECT * from coursez.cat_main where id='%s'`, guid)
		getCats := func() (banks []coursez.CatMain, err error) {
			q := CassSession.Query(qryStr, nil)
			defer q.Release()
			iter := q.Iter()
			return banks, iter.Select(&banks)
		}
		cats, err := getCats()
		currentLspIds := []string{}
		if err == nil && len(cats) > 0 {
			currentSavedCat := cats[0]
			if c.LspID != nil {
				currentLspIds = append(currentLspIds, *c.LspID)
			} else {
				currentLspIds = currentSavedCat.LspIDs
			}
		} else {
			if c.LspID != nil {
				currentLspIds = append(currentLspIds, *c.LspID)
			}
		}
		if c.ImageFile != nil {
			imageBucket = guid + "/catimages/" + c.ImageFile.Filename
			storageC := bucket.NewStorageHandler()
			gproject := googleprojectlib.GetGoogleProjectID()
			err = storageC.InitializeStorageClient(ctx, gproject, "coursez-catimages")
			if err != nil {
				log.Errorf("Failed to upload image to course: %v", err.Error())
				continue
			}
			writer, err := storageC.UploadToGCS(ctx, imageBucket, map[string]string{})
			if err != nil {
				log.Errorf("Failed to upload image to course: %v", err.Error())
				continue
			}
			defer writer.Close()
			fileBuffer := bytes.NewBuffer(nil)
			if _, err := io.Copy(fileBuffer, c.ImageFile.File); err != nil {
				continue
			}
			currentBytes := fileBuffer.Bytes()
			_, err = io.Copy(writer, bytes.NewReader(currentBytes))
			if err != nil {
				continue
			}
			imageUrl = storageC.GetSignedURLForObject(imageBucket)
		} else if c.ImageURL != nil {
			imageUrl = *c.ImageURL
		}
		cassandraCategory := coursez.CatMain{
			ID:          guid,
			Name:        *c.Name,
			Description: *c.Description,
			Code:        *c.Code,
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
			CreatedBy:   email_creator,
			UpdatedBy:   email_creator,
			IsActive:    *c.IsActive,
			ImageBucket: imageBucket,
			ImageURL:    imageUrl,
			LspIDs:      currentLspIds,
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
			ImageURL:    &cassandraCategory.ImageURL,
			LspID:       c.LspID,
		}
	}
	return catMain, nil
}

func AddSubCatMain(ctx context.Context, input []*model.SubCatMainInput) ([]*model.SubCatMain, error) {
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	email_creator := claims["email"].(string)
	log.Infof("AddSubCatMain called")
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session

	subCatMain := make([]*model.SubCatMain, len(input))
	for i, cc := range input {
		c := cc
		guid := base64.URLEncoding.EncodeToString([]byte(strings.ToLower(*c.Name)))
		imageUrl := ""
		imageBucket := ""
		qryStr := fmt.Sprintf(`SELECT * from coursez.sub_cat_main where id='%s'`, guid)
		getCats := func() (banks []coursez.SubCatMain, err error) {
			q := CassSession.Query(qryStr, nil)
			defer q.Release()
			iter := q.Iter()
			return banks, iter.Select(&banks)
		}
		cats, err := getCats()
		currentLspIds := []string{}
		if err == nil && len(cats) > 0 {
			currentSavedCat := cats[0]
			if c.LspID != nil {
				currentLspIds = append(currentLspIds, *c.LspID)
			} else {
				currentLspIds = currentSavedCat.LspIDs
			}
		} else {
			if c.LspID != nil {
				currentLspIds = append(currentLspIds, *c.LspID)
			}
		}
		if c.ImageFile != nil {
			storageC := bucket.NewStorageHandler()
			gproject := googleprojectlib.GetGoogleProjectID()
			err = storageC.InitializeStorageClient(ctx, gproject, "coursez-catimages")
			if err != nil {
				log.Errorf("Failed to upload image to course: %v", err.Error())
				continue
			}
			imageBucket = guid + "/subcatimages/" + c.ImageFile.Filename
			writer, err := storageC.UploadToGCS(ctx, imageBucket, map[string]string{})
			if err != nil {
				log.Errorf("Failed to upload image to course: %v", err.Error())
				continue
			}
			defer writer.Close()
			fileBuffer := bytes.NewBuffer(nil)
			if _, err := io.Copy(fileBuffer, c.ImageFile.File); err != nil {
				continue
			}
			currentBytes := fileBuffer.Bytes()
			_, err = io.Copy(writer, bytes.NewReader(currentBytes))
			if err != nil {
				continue
			}
			imageUrl = storageC.GetSignedURLForObject(imageBucket)
		} else if c.ImageURL != nil {
			imageUrl = *c.ImageURL
		}
		cassandraCategory := coursez.SubCatMain{
			ID:          guid,
			Name:        *c.Name,
			Description: *c.Description,
			Code:        *c.Code,
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
			CreatedBy:   email_creator,
			UpdatedBy:   email_creator,
			IsActive:    *c.IsActive,
			ParentID:    *c.CatID,
			ImageBucket: imageBucket,
			ImageURL:    imageUrl,
			LspIDs:      currentLspIds,
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
			ImageURL:    &cassandraCategory.ImageURL,
			LspID:       c.LspID,
		}
	}
	return subCatMain, nil
}
