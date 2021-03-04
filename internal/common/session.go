package common

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
)

const (
	AwsAccessKeyId     = "AWS_ACCESS_KEY_ID"
	AwsAccessSecretKey = "AWS_SECRET_ACCESS_KEY" //nolint:gosec
	AwsDefaultRegion   = "AWS_DEFAULT_REGION"
)

// CreateSession returns a standard session with the AWS CLI environment
//
// This looks into two different config files, `~/.aws/credentials` for a secret
// to the specified profile, and `~/.aws/config` for the region name detection.
//
// These are the standard config files, used by AWS CLI and other SDKs.
//
// Also see https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html
func CreateSession(profile, region string) (sess *session.Session, err error) {
	if region == "" {
		region = os.Getenv(AwsDefaultRegion)
	}

	if region == "" {
		region = DetectRegionFromConfig(os.Getenv("HOME") + "/.aws/config")
	}

	if region == "" {
		err = fmt.Errorf("could not detect AWS region, please specify")
		return
	}

	sess, err = session.NewSessionWithOptions(session.Options{
		Profile: profile,
		Config: aws.Config{
			Region:                        aws.String(region),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	})

	if err != nil {
		err = fmt.Errorf("could not create session: %w", err)
		return
	}

	return
}

func CreateTestSession(region string) (sess *session.Session) {
	sess, err := session.NewSession(&aws.Config{
		Region:      &region,
		Credentials: credentials.NewStaticCredentials("fake", "fake", "fake"),
	})
	if err != nil {
		panic(err)
	}

	return
}

func DetectRegionFromConfig(filename string) (region string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	i, err := ini.Load(data)
	if err != nil {
		return
	}

	section, err := i.GetSection("default")
	if err != nil {
		return
	}

	key, err := section.GetKey("region")
	if err != nil {
		return
	}

	return key.Value()
}
