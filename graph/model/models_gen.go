// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"

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
	TileImage    *string   `json:"tileImage"`
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
	Status       Status    `json:"status"`
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
	TileImage    *string   `json:"tileImage"`
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
	Status       Status    `json:"status"`
}

type Module struct {
	ID          *string `json:"id"`
	Name        string  `json:"name"`
	IsChapter   bool    `json:"isChapter"`
	Description string  `json:"description"`
	CourseID    string  `json:"courseId"`
	Owner       *string `json:"owner"`
	Duration    *int    `json:"duration"`
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
	Duration    *int    `json:"duration"`
	Level       *string `json:"level"`
	Sequence    *int    `json:"sequence"`
	SetGlobal   *bool   `json:"setGlobal"`
}

type Quiz struct {
	ID          *string `json:"id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Type        string  `json:"type"`
	IsMandatory bool    `json:"isMandatory"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	TopicID     string  `json:"topicId"`
	Sequence    *int    `json:"sequence"`
	StartTime   *int    `json:"startTime"`
}

type QuizDescriptive struct {
	QuizID        string `json:"quizId"`
	Question      string `json:"question"`
	CorrectAnswer string `json:"correctAnswer"`
	Explanation   string `json:"explanation"`
}

type QuizFile struct {
	QuizID string         `json:"quizId"`
	Type   string         `json:"type"`
	Name   string         `json:"name"`
	File   graphql.Upload `json:"file"`
}

type QuizInput struct {
	ID          *string `json:"id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Type        string  `json:"type"`
	IsMandatory bool    `json:"isMandatory"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	TopicID     string  `json:"topicId"`
	Sequence    *int    `json:"sequence"`
	StartTime   *int    `json:"startTime"`
}

type QuizMcq struct {
	QuizID        string    `json:"quizId"`
	Question      string    `json:"question"`
	Options       []*string `json:"options"`
	CorrectOption string    `json:"correctOption"`
	Explanation   string    `json:"explanation"`
}

type StaticContent struct {
	Type     Type           `json:"type"`
	File     graphql.Upload `json:"file"`
	CourseID string         `json:"courseId"`
	TopicID  string         `json:"topicId"`
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

type TopicContent struct {
	Language     string  `json:"language"`
	TopicID      string  `json:"topicId"`
	StartTime    *int    `json:"startTime"`
	Duration     *int    `json:"duration"`
	SkipIntro    *bool   `json:"skipIntro"`
	NextShowTime *int    `json:"nextShowTime"`
	FromEndTime  *int    `json:"fromEndTime"`
	CreatedAt    *string `json:"created_at"`
	UpdatedAt    *string `json:"updated_at"`
}

type TopicContentInput struct {
	Language     string `json:"language"`
	TopicID      string `json:"topicId"`
	StartTime    *int   `json:"startTime"`
	Duration     *int   `json:"duration"`
	SkipIntro    *bool  `json:"skipIntro"`
	NextShowTime *int   `json:"nextShowTime"`
	FromEndTime  *int   `json:"fromEndTime"`
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

type TopicResourceInput struct {
	Type      string         `json:"type"`
	TopicID   string         `json:"topicId"`
	CreatedAt *string        `json:"created_at"`
	UpdatedAt *string        `json:"updated_at"`
	CreatedBy *string        `json:"created_by"`
	UpdatedBy *string        `json:"updated_by"`
	URL       string         `json:"url"`
	File      graphql.Upload `json:"file"`
}

type TopicVideo struct {
	File     graphql.Upload `json:"file"`
	CourseID string         `json:"courseId"`
	TopicID  string         `json:"topicId"`
}

type Status string

const (
	StatusSaved           Status = "SAVED"
	StatusApprovalPending Status = "APPROVAL_PENDING"
	StatusOnHold          Status = "ON_HOLD"
	StatusApproved        Status = "APPROVED"
	StatusPublsihed       Status = "PUBLSIHED"
	StatusRejected        Status = "REJECTED"
)

var AllStatus = []Status{
	StatusSaved,
	StatusApprovalPending,
	StatusOnHold,
	StatusApproved,
	StatusPublsihed,
	StatusRejected,
}

func (e Status) IsValid() bool {
	switch e {
	case StatusSaved, StatusApprovalPending, StatusOnHold, StatusApproved, StatusPublsihed, StatusRejected:
		return true
	}
	return false
}

func (e Status) String() string {
	return string(e)
}

func (e *Status) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Status(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

func (e Status) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Type string

const (
	TypeScrom  Type = "SCROM"
	TypeTincan Type = "TINCAN"
	TypeCmi5   Type = "CMI5"
	TypeHTML5  Type = "HTML5"
)

var AllType = []Type{
	TypeScrom,
	TypeTincan,
	TypeCmi5,
	TypeHTML5,
}

func (e Type) IsValid() bool {
	switch e {
	case TypeScrom, TypeTincan, TypeCmi5, TypeHTML5:
		return true
	}
	return false
}

func (e Type) String() string {
	return string(e)
}

func (e *Type) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Type(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Type", str)
	}
	return nil
}

func (e Type) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
