package deploy_static

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amplify"
	"github.com/aws/aws-sdk-go-v2/service/amplify/types"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
)

var (
	bucketClient  *bucket.Client
	amplifyClient *amplify.Client
)

func Initialize(storageClient *bucket.Client) error {
	var err error
	bucketClient = storageClient
	cfg := aws.NewConfig()
	cfg.Region = "ap-south-1"
	creds := aws.Credentials{}
	creds.AccessKeyID = "AKIA3NKOMRI4HEGRUSUS"
	creds.SecretAccessKey = "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"
	funcCfg := aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return creds, nil
	})
	cfg.Credentials = funcCfg
	ampOpts := amplify.Options{}
	ampOpts.Credentials = funcCfg
	amplifyClient = amplify.NewFromConfig(*cfg)
	if err != nil {
		log.Fatalf("amplify.NewFromConfig: %v", err)
	}
	return err
}

// StorageObjectData contains metadata of the Cloud Storage object.
type StorageObjectData struct {
	ContentID string `json:"contentId"`
	Bucket    string `json:"bucket,omitempty"`
	Name      string `json:"name,omitempty"`
}

// DeployStatic consumes a CloudEvent message and logs details about the changed object.
func DeployStatic(ctx context.Context, data StorageObjectData) (string, error) {
	log.Printf("Bucket: %s", data.Bucket)
	rgcs, err := bucketClient.Client.Bucket(data.Bucket).Object(data.Name).NewReader(ctx)
	if err != nil {
		return "", err
	}
	{
		b, err := ioutil.ReadAll(rgcs)
		if err != nil {
			return "", err
		}
		br := bytes.NewBuffer(b)

		var r io.Reader
		r, err = gzip.NewReader(br)
		if err != nil {
			return "", err
		}
		var resB bytes.Buffer
		_, err = resB.ReadFrom(r)
		if err != nil {
			return "", err
		}
		zipBytes := resB.Bytes()
		// write zip file to disk
		{
			err = ioutil.WriteFile("static.zip", zipBytes, 0644)
			if err != nil {
				return "", err
			}
		}
	}
	ampDepInput := amplify.CreateDeploymentInput{}
	appID := "d" + data.ContentID
	// create app
	ampCreateAppInput := amplify.CreateAppInput{}
	ampCreateAppInput.Name = &appID
	ampAppOutput, err := amplifyClient.CreateApp(ctx, &ampCreateAppInput)
	if err != nil {
		return "", err
	}
	ampCreateBranchInput := amplify.CreateBranchInput{}
	ampCreateBranchInput.AppId = ampAppOutput.App.AppId
	ampCreateBranchInput.BranchName = aws.String("master")
	ampCreateBranchInput.EnableAutoBuild = aws.Bool(true)
	ampCreateBranchInput.EnableBasicAuth = aws.Bool(false)
	ampCreateBranchInput.EnableNotification = aws.Bool(false)
	ampCreateBranchInput.EnablePerformanceMode = aws.Bool(false)
	ampCreateBranchInput.EnablePullRequestPreview = aws.Bool(false)
	ampCreateBranchInput.EnvironmentVariables = map[string]string{}
	ampCreateBranchInput.PullRequestEnvironmentName = aws.String("PR-{prNumber}")
	ampCreateBranchInput.Stage = types.StageProduction
	ampCreateBranchInput.Tags = map[string]string{}
	ampCreateBranchOutput, err := amplifyClient.CreateBranch(ctx, &ampCreateBranchInput)
	if err != nil {
		return "", err
	}

	ampDepInput.AppId = ampAppOutput.App.AppId
	ampDepInput.BranchName = ampCreateBranchOutput.Branch.BranchName
	// create deployment
	ampDepOutput, err := amplifyClient.CreateDeployment(ctx, &ampDepInput)
	if err != nil {
		return "", err
	}
	ampStartDepInput := amplify.StartDeploymentInput{}
	ampStartDepInput.AppId = ampAppOutput.App.AppId
	ampStartDepInput.BranchName = ampCreateBranchOutput.Branch.BranchName
	zipUploadUrl := *ampDepOutput.ZipUploadUrl
	// upload zip file to zipUploadUrl
	// do a put request to zipUploadUrl and upload the zip file
	httpClient := &http.Client{}
	// read static.zip file
	zipFile, err := ioutil.ReadFile("static.zip")
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("PUT", zipUploadUrl, bytes.NewBuffer(zipFile))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/zip")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ampStartDepInput.JobId = ampDepOutput.JobId
	_, err = amplifyClient.StartDeployment(ctx, &ampStartDepInput)
	if err != nil {
		return "", err
	}
	//delete zip file
	ioutil.WriteFile("static.zip", []byte{}, 0644)
	amplifyUrl := fmt.Sprintf("https://%s.%s.amplifyapp.com/", *ampCreateBranchOutput.Branch.BranchName, appID)
	return amplifyUrl, nil
}
