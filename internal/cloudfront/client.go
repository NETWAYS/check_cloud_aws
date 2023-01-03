package cloudfront

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

type CloudfrontClient struct {
	CloudfrontClient *cloudfront.CloudFront
}

func NewCloudfrontClient(session *session.Session) *CloudfrontClient {
	return &CloudfrontClient{cloudfront.New(session)}
}

func (c *CloudfrontClient) LoadAllDistributions() (distributions *cloudfront.ListDistributionsOutput, err error) {
	distributions, err = c.CloudfrontClient.ListDistributions(&cloudfront.ListDistributionsInput{})
	if err != nil {
		err = fmt.Errorf("could not load all distributions: %w", err)
	}

	return
}

func (c *CloudfrontClient) LoadDistributionByETag(etag string) (distribution *cloudfront.GetDistributionOutput, err error) {
	distribution, err = c.CloudfrontClient.GetDistribution(&cloudfront.GetDistributionInput{Id: aws.String(etag)})
	if err != nil {
		err = fmt.Errorf("could not load all distribution: %w", err)
	}

	return
}
