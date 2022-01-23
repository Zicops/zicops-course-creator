package helpers

import (
	"context"
	"io/ioutil"

	"golang.org/x/oauth2/google"
)

// ReadCredentialsFile .... simply to initialize GCP credentials
func ReadCredentialsFile(ctx context.Context, filename string, scopes []string) (*google.Credentials, []byte, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	creds, err := google.CredentialsFromJSON(ctx, b, scopes...)
	if err != nil {
		return nil, nil, err
	}
	return creds, b, nil
}
