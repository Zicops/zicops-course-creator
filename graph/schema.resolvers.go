package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/zicops/zicops-course-creator/graph/generated"
	"github.com/zicops/zicops-course-creator/graph/model"
)

func (r *mutationResolver) AddCourse(ctx context.Context, course *model.CourseInput) (*model.Course, error) {
	panic(fmt.Errorf("not implemented"))
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
