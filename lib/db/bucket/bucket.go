package bucket

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/zicops/zicops-course-creator/constants"
	"github.com/zicops/zicops-course-creator/helpers"
	"google.golang.org/api/option"
)

// Client ....
type Client struct {
	projectID    string
	client       *storage.Client
	bucket       *storage.BucketHandle
	bucketPublic *storage.BucketHandle
}

// NewStorageHandler return new database action
func NewStorageHandler() *Client {
	return &Client{projectID: "", client: nil}
}

// InitializeStorageClient ...........
func (sc *Client) InitializeStorageClient(ctx context.Context, projectID string, bucket string) error {
	serviceAccountZicops := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if serviceAccountZicops == "" {
		return fmt.Errorf("failed to get right credentials for course creator")
	}
	targetScopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/userinfo.email",
	}
	currentCreds, _, err := helpers.ReadCredentialsFile(ctx, serviceAccountZicops, targetScopes)
	if err != nil {
		return err
	}
	client, err := storage.NewClient(ctx, option.WithCredentials(currentCreds))
	if err != nil {
		return err
	}
	sc.client = client
	sc.projectID = projectID
	sc.bucket, _ = sc.CreateBucket(ctx, bucket)
	sc.bucketPublic, _ = sc.CreateBucketPublic(ctx, constants.COURSES_PUBLIC_BUCKET)
	return nil
}

// CreateBucket  ...........
func (sc *Client) CreateBucket(ctx context.Context, bucketName string) (*storage.BucketHandle, error) {
	bkt := sc.client.Bucket(bucketName)
	exists, err := bkt.Attrs(ctx)
	if err != nil && exists == nil {
		if err := bkt.Create(ctx, sc.projectID, nil); err != nil {
			return nil, err
		}
	}
	return bkt, nil
}

// CreateBucket  ...........
func (sc *Client) CreateBucketPublic(ctx context.Context, bucketName string) (*storage.BucketHandle, error) {
	bkt := sc.client.Bucket(bucketName)
	exists, err := bkt.Attrs(ctx)
	if err != nil && exists == nil {
		if err := bkt.Create(ctx, sc.projectID, nil); err != nil {
			return nil, err
		}
	}
	return bkt, nil
}

// UploadToGCS ....
func (sc *Client) UploadToGCS(ctx context.Context, fileName string, tags map[string]string) (*storage.Writer, error) {
	bucketObject := sc.bucket.Object(fileName)
	return bucketObject.NewWriter(ctx), nil
}

// UploadToGCSPub ....
func (sc *Client) UploadToGCSPub(ctx context.Context, fileName string, tags map[string]string) (*storage.Writer, error) {
	bucketObject := sc.bucketPublic.Object(fileName)
	return bucketObject.NewWriter(ctx), nil
}

// BucketReader ....
func (sc *Client) BucketReader(ctx context.Context, objectName string) (*storage.Reader, error) {
	readerObject, _ := sc.bucketPublic.Object(objectName).NewReader(ctx)
	return readerObject, nil
}

// SetTags ....
func (sc *Client) SetTags(ctx context.Context, path string, tags map[string]string) error {
	bucketObject := sc.bucket.Object(path)
	if len(tags) > 0 {
		attrs, err := bucketObject.Attrs(ctx)
		if err != nil {
			return err
		}
		bucketObject = bucketObject.If(storage.Conditions{MetagenerationMatch: attrs.Metageneration})

		// writer metadata
		objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
			Metadata: tags,
		}
		if _, err := bucketObject.Update(ctx, objectAttrsToUpdate); err != nil {
			return err
		}
		return nil
	}
	return nil

}

func (sc *Client) GetSignedURLForObject(object string) string {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(24 * time.Hour),
	}
	url, err := sc.bucket.SignedURL(object, opts)
	if err != nil {
		return ""
	}

	return url
}

func (sc *Client) GetSignedURLForObjectPub(object string) string {
	// opts := &storage.SignedURLOptions{
	// 	Scheme:  storage.SigningSchemeV4,
	// 	Method:  "GET",
	// 	Expires: time.Now().Add(24 * time.Hour),
	// }
	// url, err := sc.bucketPublic.SignedURL(object, opts)
	// if err != nil {
	// 	return ""
	// }
	url := "https://storage.googleapis.com/" + constants.COURSES_PUBLIC_BUCKET + "/" + object
	return url
}
