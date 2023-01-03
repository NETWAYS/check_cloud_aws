package cloudfront

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

type GetDistributionOutput struct {
	GetDistributionOutput *cloudfront.GetDistributionOutput
}

func (c *CloudfrontClient) GetOutput(rc int, distribution *cloudfront.GetDistributionOutput) (output string) {
	output = fmt.Sprintf("[%s] %s status=%s enabled=%t",
		check.StatusText(rc),
		*distribution.Distribution.Id,
		*distribution.Distribution.Status,
		*distribution.Distribution.DistributionConfig.Enabled)

	return
}
