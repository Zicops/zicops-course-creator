package bucket

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/zicops/zicops-course-creator/constants"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/helpers"
	"google.golang.org/api/option"
)

// Client ....
type Client struct {
	projectID    string
	Client       *storage.Client
	bucket       *storage.BucketHandle
	bucketPublic *storage.BucketHandle
}

// NewStorageHandler return new database action
func NewStorageHandler() *Client {
	return &Client{projectID: "", Client: nil}
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
	sc.Client = client
	sc.projectID = projectID
	sc.bucket, _ = sc.CreateBucket(ctx, bucket)
	sc.bucketPublic, _ = sc.CreateBucketPublic(ctx, constants.COURSES_PUBLIC_BUCKET)

	//initialize firebase and firestore
	global.Ct = ctx
	opt := option.WithCredentials(currentCreds)
	global.App, err = firebase.NewApp(global.Ct, nil, opt)
	if err != nil {
		log.Printf("error initializing app: %v", err)
	}

	global.Client, err = global.App.Firestore(global.Ct)
	if err != nil {
		log.Printf("Error while initialising firestore %v", err)
	}

	return nil
}

// CreateBucket  ...........
func (sc *Client) CreateBucket(ctx context.Context, bucketName string) (*storage.BucketHandle, error) {
	bkt := sc.Client.Bucket(bucketName)
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
	bkt := sc.Client.Bucket(bucketName)
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
		Expires: time.Now().Add(7 * 24 * time.Hour),
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
	url := "https://storage.googleapis.com/" + "courses-public-zicops-deploy" + "/" + object
	return url
}

func (sc *Client) DeleteObjectsFromBucket(ctx context.Context, fileName string, lang *string) string {
	// l := "en"
	// if lang != nil {
	// 	l = *lang
	// }
	// metadata := map[string]string{
	// 	"language": l,
	// }
	o := sc.bucket.Object(fileName)

	attrs, err := o.Attrs(ctx)
	if err != nil {
		return err.Error()
	}
	o = o.If(storage.Conditions{MetagenerationMatch: attrs.Metageneration})
	if err := o.Delete(ctx); err != nil {
		return ""
	}

	// bucketObject := sc.bucket.Objects(ctx, nil)
	// for {
	// 	attrs, err := bucketObject.Next()
	// 	if err == iterator.Done {
	// 		break
	// 	}
	// 	if err != nil {
	// 		return err.Error()
	// 	}
	// 	res := attrs.Metadata
	// 	if metadata["language"] == res["langauge"] {

	// 		if err := o.Delete(ctx); err != nil {
	// 			return ""
	// 		}
	// 	}
	// }

	return "1"
}
