package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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

		currentSubCat := *subCat
		currentCat := *category
		cassandraCategory := coursez.CatSubMap{
			ID:          uuid.New().String(),
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
		words := []string{}
		if c.Name != nil {
			name := strings.ToLower(*c.Name)
			wordsLocal := strings.Split(name, " ")
			words = append(words, wordsLocal...)
		}
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
			imageBucket = guid + "/catimages/" + base64.URLEncoding.EncodeToString([]byte(strings.ToLower(c.ImageFile.Filename)))
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
			Words:       words,
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
		words := []string{}
		if c.Name != nil {
			name := strings.ToLower(*c.Name)
			wordsLocal := strings.Split(name, " ")
			words = append(words, wordsLocal...)
		}
		guid := base64.URLEncoding.EncodeToString([]byte(strings.ToLower(*c.Name)))
		parentID := cc.CatID
		if parentID == nil {
			continue
		}
		imageUrl := ""
		imageBucket := ""
		qryStr := fmt.Sprintf(`SELECT * from coursez.sub_cat_main where id='%s' AND parent_id='%s' AND is_active=true`, guid, *parentID)
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
			imageBucket = guid + "/subcatimages/" + base64.URLEncoding.EncodeToString([]byte(c.ImageFile.Filename))
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
			Words:       words,
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

func UpdateCatMain(ctx context.Context, input *model.CatMainInput) (*model.CatMain, error) {
	log.Infof("UpdateCatMain called")
	_, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	if input.ID == nil {
		return nil, errors.New("id is required")
	}
	updateCols := []string{}
	imageUrl := ""
	imageBucket := ""
	qryStr := fmt.Sprintf(`SELECT * from coursez.cat_main where id='%s' AND is_active=true`, *input.ID)
	getCats := func() (banks []coursez.CatMain, err error) {
		q := CassSession.Query(qryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return banks, iter.Select(&banks)
	}
	cats, err := getCats()
	currentLspIds := []string{}
	words := []string{}
	currentSavedCat := coursez.CatMain{}
	if err == nil && len(cats) > 0 {
		currentSavedCat = cats[0]
		if input.LspID != nil {
			if !Contains(currentSavedCat.LspIDs, *input.LspID) {
				currentLspIds = append(currentLspIds, *input.LspID)
				updateCols = append(updateCols, "lsps")
				currentSavedCat.LspIDs = currentLspIds
			}
		}
		if input.Name != nil && *input.Name != currentSavedCat.Name {
			name := strings.ToLower(*input.Name)
			wordsLocal := strings.Split(name, " ")
			words = append(words, wordsLocal...)
			currentSavedCat.Name = *input.Name
			currentSavedCat.Words = words
			updateCols = append(updateCols, "name")
			updateCols = append(updateCols, "words")
		}
		if input.Description != nil && *input.Description != currentSavedCat.Description {
			currentSavedCat.Description = *input.Description
			updateCols = append(updateCols, "description")
		}
		if input.Code != nil && *input.Code != currentSavedCat.Code {
			currentSavedCat.Code = *input.Code
			updateCols = append(updateCols, "code")
		}
		if input.ImageFile != nil {
			storageC := bucket.NewStorageHandler()
			gproject := googleprojectlib.GetGoogleProjectID()
			err = storageC.InitializeStorageClient(ctx, gproject, "coursez-catimages")
			if err != nil {
				log.Errorf("Failed to upload image to course: %v", err.Error())
				return nil, err
			}
			imageBucket = *input.ID + "/catimages/" + base64.URLEncoding.EncodeToString([]byte(input.ImageFile.Filename))
			writer, err := storageC.UploadToGCS(ctx, imageBucket, map[string]string{})
			if err != nil {
				log.Errorf("Failed to upload image to course: %v", err.Error())
				return nil, err
			}
			defer writer.Close()
			fileBuffer := bytes.NewBuffer(nil)
			if _, err := io.Copy(fileBuffer, input.ImageFile.File); err != nil {
				return nil, err
			}
			currentBytes := fileBuffer.Bytes()
			_, err = io.Copy(writer, bytes.NewReader(currentBytes))
			if err != nil {
				return nil, err
			}
			imageUrl = storageC.GetSignedURLForObject(imageBucket)
			currentSavedCat.ImageBucket = imageBucket
			currentSavedCat.ImageURL = imageUrl
			updateCols = append(updateCols, "image_bucket")
			updateCols = append(updateCols, "image_url")
		}
		if len(updateCols) > 0 {
			currentSavedCat.UpdatedAt = time.Now().Unix()
			updateCols = append(updateCols, "updated_at")
			upStms, uNames := coursez.CatMainTable.Update(updateCols...)
			updateQuery := CassSession.Query(upStms, uNames).BindStruct(&currentSavedCat)
			if err := updateQuery.ExecRelease(); err != nil {
				return nil, err
			}
		}
	}
	created := strconv.FormatInt(currentSavedCat.CreatedAt, 10)
	updated := strconv.FormatInt(currentSavedCat.UpdatedAt, 10)
	catMain := &model.CatMain{
		ID:          &currentSavedCat.ID,
		Name:        &currentSavedCat.Name,
		Description: &currentSavedCat.Description,
		Code:        &currentSavedCat.Code,
		CreatedAt:   &created,
		UpdatedAt:   &updated,
		CreatedBy:   &currentSavedCat.CreatedBy,
		UpdatedBy:   &currentSavedCat.UpdatedBy,
		IsActive:    &currentSavedCat.IsActive,
		ImageURL:    &currentSavedCat.ImageURL,
		LspID:       input.LspID,
	}
	return catMain, nil
}

func UpdateSubCatMain(ctx context.Context, input *model.SubCatMainInput) (*model.SubCatMain, error) {
	log.Infof("UpdateSubCatMain called")
	_, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	log.Infof("UpdateSubCatMain: %v", *input)
	if input.ID == nil || input.CatID == nil {
		return nil, errors.New("id and cat_id are required")
	}
	updateCols := []string{}
	imageUrl := ""
	imageBucket := ""
	qryStr := fmt.Sprintf(`SELECT * from coursez.sub_cat_main where id='%s' AND parent_id='%s' AND is_active=true`, *input.ID, *input.CatID)
	getCats := func() (banks []coursez.SubCatMain, err error) {
		q := CassSession.Query(qryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return banks, iter.Select(&banks)
	}
	cats, err := getCats()
	currentLspIds := []string{}
	words := []string{}
	currentSavedCat := coursez.SubCatMain{}
	if err == nil && len(cats) > 0 {
		currentSavedCat = cats[0]
		if input.LspID != nil {
			if !Contains(currentSavedCat.LspIDs, *input.LspID) {
				currentLspIds = append(currentLspIds, *input.LspID)
				updateCols = append(updateCols, "lsps")
				currentSavedCat.LspIDs = currentLspIds
			}
		}
		if input.Name != nil && *input.Name != currentSavedCat.Name {
			name := strings.ToLower(*input.Name)
			wordsLocal := strings.Split(name, " ")
			words = append(words, wordsLocal...)
			currentSavedCat.Name = *input.Name
			currentSavedCat.Words = words
			updateCols = append(updateCols, "name")
			updateCols = append(updateCols, "words")
		}
		if input.Description != nil && *input.Description != currentSavedCat.Description {
			currentSavedCat.Description = *input.Description
			updateCols = append(updateCols, "description")
		}
		if input.Code != nil && *input.Code != currentSavedCat.Code {
			currentSavedCat.Code = *input.Code
			updateCols = append(updateCols, "code")
		}
		if input.ImageFile != nil {
			storageC := bucket.NewStorageHandler()
			gproject := googleprojectlib.GetGoogleProjectID()
			err = storageC.InitializeStorageClient(ctx, gproject, "coursez-catimages")
			if err != nil {
				log.Errorf("Failed to upload image to course: %v", err.Error())
				return nil, err
			}
			imageBucket = *input.ID + "/catimages/" + base64.URLEncoding.EncodeToString([]byte(input.ImageFile.Filename))
			writer, err := storageC.UploadToGCS(ctx, imageBucket, map[string]string{})
			if err != nil {
				log.Errorf("Failed to upload image to course: %v", err.Error())
				return nil, err
			}
			defer writer.Close()
			fileBuffer := bytes.NewBuffer(nil)
			if _, err := io.Copy(fileBuffer, input.ImageFile.File); err != nil {
				return nil, err
			}
			currentBytes := fileBuffer.Bytes()
			_, err = io.Copy(writer, bytes.NewReader(currentBytes))
			if err != nil {
				return nil, err
			}
			imageUrl = storageC.GetSignedURLForObject(imageBucket)
			currentSavedCat.ImageBucket = imageBucket
			currentSavedCat.ImageURL = imageUrl
			updateCols = append(updateCols, "image_bucket")
			updateCols = append(updateCols, "image_url")
		}
		if len(updateCols) > 0 {
			currentSavedCat.UpdatedAt = time.Now().Unix()
			updateCols = append(updateCols, "updated_at")
			upStms, uNames := coursez.SubCatMainTable.Update(updateCols...)
			updateQuery := CassSession.Query(upStms, uNames).BindStruct(&currentSavedCat)
			if err := updateQuery.ExecRelease(); err != nil {
				return nil, err
			}
		}
	}
	created := strconv.FormatInt(currentSavedCat.CreatedAt, 10)
	updated := strconv.FormatInt(currentSavedCat.UpdatedAt, 10)
	catMain := &model.SubCatMain{
		ID:          &currentSavedCat.ID,
		Name:        &currentSavedCat.Name,
		Description: &currentSavedCat.Description,
		Code:        &currentSavedCat.Code,
		CreatedAt:   &created,
		UpdatedAt:   &updated,
		CreatedBy:   &currentSavedCat.CreatedBy,
		UpdatedBy:   &currentSavedCat.UpdatedBy,
		IsActive:    &currentSavedCat.IsActive,
		ImageURL:    &currentSavedCat.ImageURL,
		LspID:       input.LspID,
	}
	return catMain, nil
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
