package graphproc

import (
	"fmt"
	"testing"
)

func TestAwsS3Download(t *testing.T) {

	payload := NewPayload()
	state := NewState()
	state.SetConfig("region", "us-east-2")
	state.SetConfig("bucket", "lambdapipeprocessor")
	state.SetConfig("file", "bulkdata.txt")
	AWSDownloadS3Bucket(state, payload)
	fmt.Println(string(payload.Raw))
}
