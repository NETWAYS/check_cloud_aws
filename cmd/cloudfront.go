package cmd

import (
	"fmt"
	"github.com/NETWAYS/check_cloud_aws/internal/cloudfront"
	"github.com/NETWAYS/check_cloud_aws/internal/common"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
	c "github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/spf13/cobra"
)

var (
	ETags []string
)

var cloudfrontCmd = &cobra.Command{
	Use:   "cloudfront",
	Short: "Checks in the Cloudfront context",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			output        string
			summary       string
			totalCrit     int
			totalWarn     int
			rc            int
			states        []int
			distributions []*c.GetDistributionOutput
			perf          perfdata.PerfdataList
		)

		client := RequireCloudfrontClient()

		if ETags == nil {
			distributionsList, err := client.LoadAllDistributions()
			if err != nil {
				check.ExitError(err)
			}

			for _, dist := range distributionsList.DistributionList.Items {
				distribution, err := client.LoadDistributionByETag(*dist.Id)
				if err != nil {
					check.ExitError(err)
				}

				distributions = append(distributions, distribution)
			}
		} else {
			for _, etag := range ETags {
				distribution, err := client.LoadDistributionByETag(etag)
				if err != nil {
					check.ExitError(err)
				}

				distributions = append(distributions, distribution)
			}
		}

		summary += fmt.Sprintf("Found %d Distributions - ", len(distributions))

		for _, distribution := range distributions {
			if *distribution.Distribution.DistributionConfig.Enabled == false {
				rc = 2
				totalCrit++
			} else if *distribution.Distribution.Status == "InProgress" {
				rc = 1
				totalWarn++
			} else {
				rc = 0
			}

			states = append(states, rc)

			if rc != 0 {
				output += client.GetOutput(rc, distribution)
			}

			p := perfdata.Perfdata{
				Label: *distribution.Distribution.Id,
				Value: *distribution.Distribution.Status,
			}

			perf.Add(&p)

			if len(distributions) > 1 {
				output += "\n"
			}
		}

		if result.WorstState(states...) == 0 {
			output = fmt.Sprintf("")
		}

		summary += fmt.Sprintf("critical %d - warning %d\n", totalCrit, totalWarn)

		if len(distributions) > 1 {
			if result.WorstState(states...) != 0 {
				summary += "\n"
			}
		}

		check.ExitRaw(result.WorstState(states...), summary+output, "|", perf.String())
	},
}

func init() {
	cloudfrontFlags := cloudfrontCmd.Flags()
	cloudfrontFlags.StringSliceVarP(&ETags, "etag", "e", nil,
		"Etag name of one or multiple distributions. If '--etag' is empty, all distributions will be evaluated.")

	cloudfrontFlags.SortFlags = false
}

func RequireCloudfrontClient() *cloudfront.CloudfrontClient {
	session, err := common.CreateSession(Profile, Region)
	if err != nil {
		check.ExitError(fmt.Errorf("could not setup AWS API session: %w", err))
	}

	return cloudfront.NewCloudfrontClient(session)
}
