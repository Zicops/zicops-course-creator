package bucket

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/zicops/zicops-cass-pool/redis"
	"github.com/zicops/zicops-course-creator/constants"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/helpers"
	"google.golang.org/api/iterator"
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

func (sc *Client) GetSignedURLForObject(ctx context.Context, object string) string {
	key := "signed_url_" + base64.StdEncoding.EncodeToString([]byte(object))
	res, err := redis.GetRedisValue(ctx, key)
	if err == nil && res != "" {
		return res
	}
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(24 * time.Hour),
	}
	url, err := sc.bucket.SignedURL(object, opts)
	if err != nil {
		return ""
	}
	allBut10Secsto24Hours := 24*time.Hour - 10*time.Second
	redis.SetRedisValue(ctx, key, url)
	redis.SetTTL(ctx, key, int(allBut10Secsto24Hours.Seconds()))
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

func (sc *Client) DeleteObjectsFromBucket(ctx context.Context, fileName string) string {
	o := sc.bucket.Object(fileName)

	if err := o.Delete(ctx); err != nil {
		return err.Error()
	}

	return "1"
}

func (sc *Client) CheckIfBucketHasItems(ctx context.Context, bucketName string) bool {
	bkt := sc.Client.Bucket(bucketName)
	it := bkt.Objects(ctx, nil)
	for {
		_, err := it.Next()
		if err == iterator.Done {
			return false
		}
		if err != nil {
			return false
		}
		return true
	}
}
