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
	if inp.UserID != nil {
		id = *inp.UserID
	}
	CreatedBy := claims["user_id"].(string)

	t := 0
	if inp.Time != nil {
		t = *inp.Time
	}
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return "", err
	}
	CassSession := session

	//coursez.Discussion's likes and dislikes were []string, and
	//we had []*string as input so we had to convert them
	var likesArray, dislikesArray []string
	for _, l := range inp.Likes {
		likesArray = append(likesArray, *l)
	}
	for _, d := range inp.Dislike {
		dislikesArray = append(dislikesArray, *d)
	}
	discussionId := fmt.Sprint(uuid.New().ID())
	discussionData := coursez.Discussion{
		DiscussionId: discussionId,
		CourseId:     inp.CourseID,
		UserId:       id,
		Time:         int64(t),
		Content:      inp.Content,
		Likes:        likesArray,
		Dislike:      dislikesArray,
		CreatedBy:    CreatedBy,
		CreatedAt:    time.Now().Unix(),
		UpdatedBy:    CreatedBy,
		UpdatedAt:    time.Now().Unix(),
		Status:       inp.Status,
		ReplyCount:   0,
	}
	if inp.ReplyID != nil && *inp.ReplyID != "" {
		discussionData.ReplyId = *inp.ReplyID
		err := updateReplyCount(inp)
		if err != nil {
			return "", err
		}
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

	insertQuery := CassSession.Query(coursez.DiscussionTable.Insert()).BindStruct(discussionData)
	if err := insertQuery.ExecRelease(); err != nil {
		return "", err
	}
	return discussionId, nil
}

func updateReplyCount(inp model.Discussion) error {
	parentId := *inp.ReplyID
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return err
	}
	CassSession := session

	querystr := fmt.Sprintf(`SELECT * from coursez.discussion where discussion_id='%s' ALLOW FILTERING`, parentId)
	getDiscussion := func() (parent []coursez.Discussion, err error) {
		q := CassSession.Query(querystr, nil)
		defer q.Release()
		iter := q.Iter()
		return parent, iter.Select(&parent)
	}

	parentDiscussions, err := getDiscussion()
	if err != nil {
		return err
	}
	parentDiscussion := parentDiscussions[0]
	// if parentDiscussion.ReplyCount
	parentDiscussion.ReplyCount = parentDiscussion.ReplyCount + 1
	stmt, names := coursez.DiscussionTable.Update("reply_count")
	updatedQuery := CassSession.Query(stmt, names).BindStruct(&parentDiscussion)
	if err = updatedQuery.ExecRelease(); err != nil {
		return err
	}

	return nil
}

func UpdateCourseDiscussion(ctx context.Context, discussionID string, courseID string, content *string, likes []*string, dislikes []*string, isAnonymous *bool, isPinned *bool, isAnnouncement *bool, status *string) (*model.DiscussionData, error) {
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		log.Printf("Got error while getting claims %v", err)
		return nil, err
	}
	uId := claims["user_id"].(string)

	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session

	//if we have any changes given in query, then we will update them
	queryStr := fmt.Sprintf(`SELECT * from coursez.discussion where course_id = '%s' and discussion_id = '%s' ALLOW FILTERING`, courseID, discussionID)
	updatedCols := []string{}

	getDiscussions := func() (discussions []coursez.Discussion, err error) {
		q := CassSession.Query(queryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return discussions, iter.Select(&discussions)
	}

	data, err := getDiscussions()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}

	discussion := data[0]
	if content != nil && uId == discussion.UserId && *content != "" {
		tmp := *content
		discussion.Content = tmp
		updatedCols = append(updatedCols, "user_id")
	}
	//we have discusion.Likes as []string and in input we receive []*string, as ofcourse, likes array is not compulsory
	var ArrayLikes, ArrayDislikes []string
	if likes != nil {
		for _, v := range likes {
			ArrayLikes = append(ArrayLikes, *v)
		}
		discussion.Likes = ArrayLikes
		updatedCols = append(updatedCols, "likes")
	}
	if dislikes != nil {
		for _, v := range dislikes {
			ArrayDislikes = append(ArrayDislikes, *v)
		}
		discussion.Dislike = ArrayDislikes
		updatedCols = append(updatedCols, "dislikes")
	}
	if isAnonymous != nil {
		tmp := *isAnonymous
		discussion.IsAnonymous = tmp
		updatedCols = append(updatedCols, "is_anonymous")
	}
	if isPinned != nil {
		tmp := *isPinned
		discussion.IsPinned = tmp
		updatedCols = append(updatedCols, "is_pinned")
	}
	if isAnnouncement != nil {
		tmp := *isAnnouncement
		discussion.IsAnnouncement = tmp
		updatedCols = append(updatedCols, "is_announcement")
	}
	if status != nil && *status != "" {
		tmp := *status
		discussion.Status = tmp
	}
	updatedBy := claims["user_id"].(string)

	//we have updated all the values, lets put those updates in database table
	if len(updatedCols) > 0 {
		discussion.UpdatedAt = time.Now().Unix()
		updatedCols = append(updatedCols, "updated_at")
		discussion.UpdatedBy = updatedBy
		updatedCols = append(updatedCols, "updated_by")
		stmt, names := coursez.DiscussionTable.Update(updatedCols...)
		updatedQuery := CassSession.Query(stmt, names).BindStruct(&discussion)
		if err = updatedQuery.ExecRelease(); err != nil {
			return nil, err
		}
	}

	var likesArray, dislikesArray []*string
	for _, l := range discussion.Likes {
		likesArray = append(likesArray, &l)
	}
	for _, d := range discussion.Dislike {
		dislikesArray = append(dislikesArray, &d)
	}
	ca := int(discussion.CreatedAt)
	ua := int(discussion.UpdatedAt)
	t := int(discussion.Time)
	result := model.DiscussionData{
		DiscussionID:   &discussion.DiscussionId,
		CourseID:       &discussion.CourseId,
		ReplyID:        &discussion.ReplyId,
		UserID:         &discussion.UserId,
		Time:           &t,
		Content:        &discussion.Content,
		Module:         &discussion.Module,
		Chapter:        &discussion.Chapter,
		Topic:          &discussion.Topic,
		Likes:          likesArray,
		Dislike:        dislikesArray,
		IsAnonymous:    &discussion.IsAnonymous,
		IsPinned:       &discussion.IsPinned,
		IsAnnouncement: &discussion.IsAnnouncement,
		ReplyCount:     &discussion.ReplyCount,
		CreatedBy:      &discussion.CreatedBy,
		CreatedAt:      &ca,
		UpdatedBy:      &discussion.UpdatedBy,
		UpdatedAt:      &ua,
		Status:         &discussion.Status,
	}

	return &result, nil
}
