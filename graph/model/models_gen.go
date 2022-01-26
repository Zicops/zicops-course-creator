// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"github.com/99designs/gqlgen/graphql"
)

type Chapter struct {
	ID          *string `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ModuleID    *string `json:"moduleId"`
	CourseID    string  `json:"courseId"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	Sequence    *int    `json:"sequence"`
}

type ChapterInput struct {
	ID          *string `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ModuleID    *string `json:"moduleId"`
	CourseID    string  `json:"courseId"`
	Sequence    *int    `json:"sequence"`
}

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
	CreatedAt    *string   `json:"created_at"`
	UpdatedAt    *string   `json:"updated_at"`
	Type         *string   `json:"type"`
	Prequisites  []*string `json:"prequisites"`
	GoodFor      []*string `json:"goodFor"`
	MustFor      []*string `json:"mustFor"`
	CreatedBy    *string   `json:"created_by"`
	UpdatedBy    *string   `json:"updated_by"`
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
	Type         *string   `json:"type"`
	Prequisites  []*string `json:"prequisites"`
	GoodFor      []*string `json:"goodFor"`
	MustFor      []*string `json:"mustFor"`
	CreatedBy    *string   `json:"created_by"`
	UpdatedBy    *string   `json:"updated_by"`
	Status       *string   `json:"status"`
}

type Module struct {
	ID          *string `json:"id"`
	Name        string  `json:"name"`
	IsChapter   bool    `json:"isChapter"`
	Description string  `json:"description"`
	CourseID    string  `json:"courseId"`
	Owner       *string `json:"owner"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	Level       *string `json:"level"`
	Sequence    *int    `json:"sequence"`
	SetGlobal   *bool   `json:"setGlobal"`
}

type ModuleInput struct {
	ID          *string `json:"id"`
	Name        string  `json:"name"`
	IsChapter   bool    `json:"isChapter"`
	Description string  `json:"description"`
	CourseID    string  `json:"courseId"`
	Owner       *string `json:"owner"`
	Level       *string `json:"level"`
	Sequence    *int    `json:"sequence"`
	SetGlobal   *bool   `json:"setGlobal"`
}

type Topic struct {
	ID          *string `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	ModuleID    *string `json:"moduleId"`
	ChapterID   *string `json:"chapterId"`
	CourseID    string  `json:"courseId"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	Sequence    *int    `json:"sequence"`
	CreatedBy   *string `json:"created_by"`
	UpdatedBy   *string `json:"updated_by"`
}

type TopicInput struct {
	ID          *string `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	ModuleID    *string `json:"moduleId"`
	ChapterID   *string `json:"chapterId"`
	CourseID    string  `json:"courseId"`
	Sequence    *int    `json:"sequence"`
	CreatedBy   *string `json:"created_by"`
	UpdatedBy   *string `json:"updated_by"`
}
