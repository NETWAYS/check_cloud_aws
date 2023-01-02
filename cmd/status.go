package cmd

import (
	"encoding/xml"
	"fmt"
	"github.com/NETWAYS/check_cloud_aws/internal/status"
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
	"strings"
)

// To store the CLI parameters
type StatusConfig struct {
	Url     string
	Service string
}

var cliStatusConfig StatusConfig

func contains(s string, list []string) bool {
	// Tiny helper to see if a string is in a list of strings
	for _, elem := range list {
		if s == elem {
			return true
		}
	}

	return false
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Checks the status of AWS services",
	Example: `
	check_cloud_aws status --service cloudfront
	OK - Service cloudfront is operating normally

	check_cloud_aws --region us-west-1 status --service cloudwatch
	WARNING - Information available for cloudwatch in us-west-1

	check_cloud_aws --region eu-west-1 status
	CRITICAL - Service disruption for ec2 in eu-west-1

	check_cloud_aws --region "" status iam
	CRITICAL - WARNING - Information available for iam (Global)`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			feed    status.Rss
			rc      int
			output  string
			feedUrl string
		)

		// These services don't require a region
		// Hint: This list might not be extensive
		var globalServices = []string{"route53",
			"route53domainregistration",
			"route53apprecoverycontroller",
			"chime",
			"health",
			"import-export",
			"iam",
			"awsiotdevicemanagement",
			"marketplace",
			"apipricing",
			"awswaf",
			"trustedadvisor",
			"supportcenter",
			"resourcegroups",
			"organizations",
			"management-console",
			"awsiotdevicemanagement",
			"account",
			"interregionvpcpeering",
			"cloudfront",
			"billingconsole",
			"chatbot"}

		// Creating an client and connecting to the RSS Feed
		// Access the AWS Health Dashboard at health.aws.amazon.com directly,
		// since the AWS Go SDK just supports Personal Dashboards
		c := &http.Client{}

		if Region == "" && !contains(cliStatusConfig.Service, globalServices) {
			check.ExitError(fmt.Errorf("Region required for regional services"))
		}

		// Using + concatenation since the JoinPath will add / inbetween
		if Region == "" && contains(cliStatusConfig.Service, globalServices) {
			feedUrl, _ = url.JoinPath(cliStatusConfig.Url, "/rss/", cliStatusConfig.Service+".rss")
			// Just for the output later
			Region = "Global"
		} else {
			feedUrl, _ = url.JoinPath(cliStatusConfig.Url, "/rss/", cliStatusConfig.Service+"-"+Region+".rss")
		}

		resp, err := c.Get(feedUrl)

		if err != nil {
			check.ExitError(err)
		}

		if resp.StatusCode != http.StatusOK {
			check.ExitError(fmt.Errorf("Could not get %s - Error: %d", feedUrl, resp.StatusCode))
		}

		defer resp.Body.Close()
		err = xml.NewDecoder(resp.Body).Decode(&feed)

		if err != nil {
			check.ExitError(err)
		}

		// Exit if there are no events
		if len(feed.Channel.Items) == 0 {
			rc = check.OK
			output = fmt.Sprintf("No events for %s (%s)", cliStatusConfig.Service, Region)
			check.ExitRaw(rc, output)
		}

		rc = check.Unknown
		output = "Status unknown"

		// Get the latest event
		item := strings.Split(feed.Channel.Items[0].Title, ":")

		// If we couldn't split the title
		if len(item) < 2 {
			output = "Could not determine status."
			check.ExitRaw(rc, output, feed.Channel.Items[0].Title)
		}

		event := item[0]
		details := item[1]

		if strings.Contains(event, "Service disruption") {
			// Service disruptions are Critical
			rc = check.Critical
			output = fmt.Sprintf("Service disruption for %s (%s)", cliStatusConfig.Service, Region)
		}

		if strings.Contains(event, "Performance issue") {
			// Performance issues are Warnings
			rc = check.Warning
			output = fmt.Sprintf("Performance issues for %s (%s)", cliStatusConfig.Service, Region)
		}

		if strings.Contains(event, "Informational message") {
			// There will be no new item if an information is resolved,
			// we need to check if the info is resolved.
			if strings.Contains(details, "[RESOLVED]") {
				rc = check.OK
				output = fmt.Sprintf("Event resolved for %s (%s)", cliStatusConfig.Service, Region)
			} else {
				// An information that should be checked by someone
				rc = check.Warning
				output = fmt.Sprintf("Information available for %s (%s)", cliStatusConfig.Service, Region)
			}
		}

		if strings.Contains(event, "Service is operating normally") {
			rc = check.OK
			output = fmt.Sprintf("Service %s is operating normally (%s)", cliStatusConfig.Service, Region)
		}

		check.ExitRaw(rc, output)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	fs := statusCmd.Flags()

	fs.StringVarP(&cliStatusConfig.Url, "url", "u", "https://status.aws.amazon.com",
		"The AWS Status Page URL")
	fs.StringVarP(&cliStatusConfig.Service, "service", "s", "ec2",
		"The AWS Service to check")

	fs.SortFlags = false
}
