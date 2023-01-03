package common_test

import (
	"github.com/NETWAYS/check_cloud_aws/internal/common"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

const testConfig = `[default]
region = eu-central-1
`

// nolint: gosec
const testCredentials = `[default]
aws_access_key_id=FAKE
aws_secret_access_key=FAKE
`

func TestCreateSession_WithoutConfig(t *testing.T) {
	_ = os.Setenv("HOME", "/nonexistent")

	_, err := common.CreateSession("", "", "")
	assert.Error(t, err)
}

func TestCreateSession_WithEnvironment(t *testing.T) {
	_ = os.Setenv(common.AwsAccessKeyId, "FAKE")
	_ = os.Setenv(common.AwsAccessSecretKey, "FAKE")
	_ = os.Setenv(common.AwsDefaultRegion, "eu-central-1")

	session, err := common.CreateSession("", "", "")
	assert.NoError(t, err)
	assert.Equal(t, "eu-central-1", *session.Config.Region)

	_ = os.Unsetenv(common.AwsAccessKeyId)
	_ = os.Unsetenv(common.AwsAccessSecretKey)
	_ = os.Unsetenv(common.AwsDefaultRegion)
}

func TestCreateSession_WithDefaultConfigFiles(t *testing.T) {
	dir, err := os.MkdirTemp(os.TempDir(), "awstest")
	assert.NoError(t, err)

	defer func() { _ = os.RemoveAll(dir) }()

	_ = os.Mkdir(dir+"/.aws", 0700)
	_ = os.Setenv("HOME", dir)

	_ = os.WriteFile(dir+"/.aws/config", []byte(testConfig), 0600)
	_ = os.WriteFile(dir+"/.aws/credentials", []byte(testCredentials), 0600)

	session, err := common.CreateSession("", "", "")
	assert.NoError(t, err)

	assert.Equal(t, "eu-central-1", *session.Config.Region)
	assert.NotNil(t, session.Config.Credentials)
}

func TestDetectRegionFromConfig(t *testing.T) {
	file, err := os.CreateTemp(os.TempDir(), "awstest")
	assert.NoError(t, err)

	defer func() { _ = os.Remove(file.Name()) }()

	_, err = io.WriteString(file, testConfig)
	assert.NoError(t, err)

	assert.Equal(t, "", common.DetectRegionFromConfig("/nonexistent/aws/config"))
	assert.Equal(t, "eu-central-1", common.DetectRegionFromConfig(file.Name()))
}
