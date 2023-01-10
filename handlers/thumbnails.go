package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"log"

	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
	"google.golang.org/api/iterator"
)

func AddContentThumbail(ctx context.Context, data *model.ThumbnailsDataInput) (string, error) {
	global.Ct = ctx

	for _, v := range data.Thumbnail {

		val := enc(*v)
		inp := &model.ThumbnailsData{
			ContentID: data.ContentID,
			Thumbnail: val,
		}
		_, _, err := global.Client.Collection("thumbnails").Add(global.Ct, inp)
		if err != nil {
			log.Printf("Got error while adding data to firestore %v\n", err)
			return "", nil
		}
	}
	return "Data added successfully", nil
}

func enc(inp string) string {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(inp)); err != nil {
		log.Println(err)
	}
	gz.Close()
	return b.String()
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
