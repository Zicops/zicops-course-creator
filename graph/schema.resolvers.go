package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/zicops-course-creator/graph/generated"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/handlers"
)

func (r *mutationResolver) AddCourse(ctx context.Context, course *model.CourseInput) (*model.Course, error) {
	resp, err := handlers.CourseCreator(ctx, course)
	if err != nil {
		log.Errorf("error creating course: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UploadCourseImage(ctx context.Context, file graphql.Upload) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UploadCoursePreviewVideo(ctx context.Context, file graphql.Upload) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
