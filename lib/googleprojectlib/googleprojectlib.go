package googleprojectlib

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/compute/metadata"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

// GetGoogleProjectID - Get name from metadata
func GetGoogleProjectID() string {
	ctx := context.Background()
	credentials, err := google.FindDefaultCredentials(ctx, compute.ComputeScope)
	if err != nil {
		return "zicops-one"
	}
	return credentials.ProjectID
}

// GetGoogleProjectDefaultID - Get name from metadata
func GetGoogleProjectDefaultID() string {
	srcfilename := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if srcfilename != "" {
		return GetGoogleProjectID()

	}
	id, err := metadata.ProjectID()
	if err != nil {
		return ""
	}
	return id
}

//GetInstanceID current meta instance
func GetInstanceID() string {

	instanceID, err := metadata.InstanceID()
	if err != nil {
		return "unknownInstanceId"
	}
	return instanceID

}

//GetZone Project Zone
func GetZone() string {

	zone, err := metadata.Zone()
	if err != nil {
		return "unknownZone"
	}
	return zone

}

// GetKubernetesNodeName - Get info from metadata
func GetKubernetesNodeName() string {
	in, _ := metadata.InstanceName()
	return fmt.Sprintf("KubNode:%s", in)
}
