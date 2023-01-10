package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/helpers"
)

func AddCourseDiscussion(ctx context.Context, inp model.Discussion) (string, error) {
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		log.Printf("Got error while getting claims %v", err)
	}
	id := claims["user_id"].(string)
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return "", err
	}
	CassSession := session
	var likesArray, dislikesArray []int
	for _, l := range inp.Likes {
		likesArray = append(likesArray, *l)
	}
	for _, d := range inp.Dislike {
		dislikesArray = append(dislikesArray, *d)
	}
	discussionId := fmt.Sprint(uuid.New().ID())
	discussionData := coursez.Discussion{
		DiscussionId:   discussionId,
		CourseId:       inp.CourseID,
		UserId:         id,
		Time:           time.Now().Unix(),
		Content:        inp.Content,
		Likes:          likesArray,
		Dislike:        dislikesArray,
		IsAnonymous:    *inp.IsAnonymous,
		IsPinned:       *inp.IsPinned,
		IsAnnouncement: *inp.IsAnnouncement,
		ReplyCount:     *inp.ReplyCount,
		CreatedBy:      inp.CreatedBy,
		CreatedAt:      int64(inp.CreatedAt),
		UpdatedBy:      inp.UpdatedBy,
		UpdatedAt:      int64(inp.UpdatedAt),
		Status:         inp.Status,
	}
	if inp.ReplyID != nil && *inp.ReplyID != "" {
		discussionData.ReplyId = *inp.ReplyID
	}
	if inp.IsAnonymous != nil {
		discussionData.IsAnonymous = *inp.IsAnonymous
	}
	if inp.Topic != nil && *inp.Topic != "" {
		discussionData.Topic = *inp.Topic
	}
	if inp.Chapter != nil && *inp.Chapter != "" {
		discussionData.Chapter = *inp.Chapter
	}
	if inp.Module != nil && *inp.Module != "" {
		discussionData.Module = *inp.Module
	}
	if inp.IsPinned != nil {
		discussionData.IsPinned = *inp.IsPinned
	}
	if inp.IsAnnouncement != nil {
		discussionData.IsAnnouncement = *inp.IsAnnouncement
	}
	if inp.ReplyCount != nil {
		discussionData.ReplyCount = *inp.ReplyCount
	}

	insertQuery := CassSession.Query(coursez.DiscussionTable.Insert()).BindStruct(discussionData)
	if err := insertQuery.ExecRelease(); err != nil {
		return "", err
	}
	return "Data added successfully", nil
}
