package ec2_test

import (
	"github.com/NETWAYS/go-check-network/http/mock"

	"github.com/jarcoal/httpmock"
)

var (
	TestRegion   = "eu-central-1"
	TestQueryMap = checkhttpmock.QueryMap{
		"Action=DescribeInstances&Filter.1.Name=instance-id&Filter.1.Value.1=i-good":     "instance.xml",
		"Action=DescribeInstances&Filter.1.Name=tag%3AName&Filter.1.Value.1=instance001": "instance.xml",
		"Action=DescribeInstances&Filter.1.Name=tag%3AName&Filter.1.Value.1=instance%2A": "instance.xml",
		"Action=DescribeInstanceStatus&InstanceId.1=i-good":                              "instanceStatus.xml",
		"Action=DescribeInstances&Filter.1.Name=instance-id&Filter.1.Value.1=i-stopped":  "instance-stopped.xml",
		"Action=DescribeInstanceStatus&InstanceId.1=i-stopped":                           "instanceStatus-stopped.xml",
	}
)

func enableMocking() func() {
	httpmock.Activate()

	checkhttpmock.RegisterQueryMapResponder("POST", "https://ec2."+TestRegion+".amazonaws.com/", TestQueryMap)

	return httpmock.DeactivateAndReset
}
