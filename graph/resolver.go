package graph

import "github.com/zicops/zicops-course-creator/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Courses      []*model.Course
	Topics       []*model.Topic
	Modules      []*model.Module
	Chapters     []*model.Chapter
	TopicContent []*model.TopicContent
}
