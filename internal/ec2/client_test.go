package ec2_test

import (
	"github.com/NETWAYS/check_cloud_aws/internal/common"
	"github.com/NETWAYS/check_cloud_aws/internal/ec2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createTestClient() (*ec2.Client, func()) {
	return ec2.NewClient(common.CreateTestSession(TestRegion)), enableMocking()
}

func TestClient_LoadInstance(t *testing.T) {
	client, cleanup := createTestClient()
	defer cleanup()

	i, err := client.LoadInstance("i-good")
	assert.NoError(t, err)
	assert.NotNil(t, i.Instance)
	assert.Equal(t, 0, i.GetStatus())

	i, err = client.LoadInstance("i-stopped")
	assert.NoError(t, err)
	assert.NotNil(t, i.Instance)
	assert.Equal(t, 2, i.GetStatus())
	assert.Contains(t, i.GetOutput(), "stopped")
}

func TestClient_LoadInstanceByName(t *testing.T) {
	client, cleanup := createTestClient()
	defer cleanup()

	i, err := client.LoadInstanceByName("instance001")
	assert.NoError(t, err)
	assert.NotNil(t, i.Instance)
	assert.Equal(t, 0, i.GetStatus())
}

func TestClient_LoadAllInstancesByFilter(t *testing.T) {
	client, cleanup := createTestClient()
	defer cleanup()

	i, err := client.LoadAllInstancesByFilter(
		ec2.DescribeInstancesInput(ec2.Filter("tag:Name", "instance*")))
	assert.NoError(t, err)
	assert.Len(t, i.Instances, 1)
	assert.Equal(t, 0, i.GetStatus())

	output := i.GetOutput()
	assert.Contains(t, output, "running")
	assert.Contains(t, output, "[OK]")
}
