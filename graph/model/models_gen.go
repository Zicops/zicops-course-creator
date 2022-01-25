// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"github.com/99designs/gqlgen/graphql"
)

type Course struct {
	ID           *string   `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Instructor   *string   `json:"instructor"`
	Image        *string   `json:"image"`
	PreviewVideo *string   `json:"previewVideo"`
	Owner        *string   `json:"owner"`
	Duration     *int      `json:"duration"`
	Level        *string   `json:"level"`
	Language     []*string `json:"language"`
	Takeaways    []*string `json:"takeaways"`
	CreatedAt    *string   `json:"createdAt"`
	UpdatedAt    *string   `json:"updatedAt"`
	Type         *string   `json:"type"`
	Prequisites  []*string `json:"prequisites"`
	GoodFor      []*string `json:"goodFor"`
	MustFor      []*string `json:"mustFor"`
	CreatedBy    *string   `json:"createdBy"`
	UpdatedBy    *string   `json:"updatedBy"`
	Status       *string   `json:"status"`
}

type CourseFile struct {
	File     graphql.Upload `json:"file"`
	CourseID string         `json:"courseId"`
}

type CourseInput struct {
	ID           *string   `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Instructor   *string   `json:"instructor"`
	Image        *string   `json:"image"`
	PreviewVideo *string   `json:"previewVideo"`
	Owner        *string   `json:"owner"`
	Duration     *int      `json:"duration"`
	Level        *string   `json:"level"`
	Language     []*string `json:"language"`
	Takeaways    []*string `json:"takeaways"`
	CreatedAt    *string   `json:"createdAt"`
	UpdatedAt    *string   `json:"updatedAt"`
	Type         *string   `json:"type"`
	Prequisites  []*string `json:"prequisites"`
	GoodFor      []*string `json:"goodFor"`
	MustFor      []*string `json:"mustFor"`
	CreatedBy    *string   `json:"createdBy"`
	UpdatedBy    *string   `json:"updatedBy"`
	Status       *string   `json:"status"`
}
