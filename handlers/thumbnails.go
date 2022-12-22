package handlers

import (
	"context"
	"log"

	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
	"google.golang.org/api/iterator"
)

func AddContentThumbail(ctx context.Context, contentID string, thumbnail string) (string, error) {
	global.Ct = ctx
	data := &model.ThumbnailsData{
		ContentID: contentID,
		Thumbnail: thumbnail,
	}
	_, _, err := global.Client.Collection("thumbnails").Add(global.Ct, data)
	if err != nil {
		log.Printf("Got error while adding data to firestore %v\n", err)
		return "", nil
	}
	return "Data added successfully", nil
}

func GetThumbnails(ctx context.Context, contentID []*string) ([]*model.ThumbnailsData, error) {
	var resp []*model.ThumbnailsData
	for _, v := range contentID {
		//here just get the value from firestore
		iter := global.Client.Collection("thumbnails").Where("ContentID", "==", *v).Documents(ctx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate %v", err)
				return nil, err
			}
			data := doc.Data()

			tmp := &model.ThumbnailsData{
				ContentID: data["ContentID"].(string),
				Thumbnail: data["Thumbnail"].(string),
			}
			log.Println(data["ContentID"].(string), "   ", data["Thumbnail"].(string))
			resp = append(resp, tmp)
		}
	}
	return resp, nil
}
