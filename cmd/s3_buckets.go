package cmd

import (
	"fmt"
	"github.com/NETWAYS/check_cloud_aws/internal"
	b "github.com/NETWAYS/check_cloud_aws/internal/s3"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
	"strconv"
)

var (
	CriticalBucketSize string
	WarningBucketSize  string
)

var s3BucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "Checks the size of a single bucket or multiple buckets",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err       error
			summary   string
			output    string
			states    []int
			totalCrit int
			totalWarn int
			rc        int
			perf      perfdata.PerfdataList
		)

		buckets := s3.ListBucketsOutput{}
		objectsOutput := b.V2Output{}

		client := RequireS3Client()

		if BucketNames == nil {
			bk, err := client.LoadAllBuckets()
			if err != nil {
				check.ExitError(err)
			}

			buckets = *bk
		} else {
			for _, bucketName := range BucketNames {
				buckets.Buckets = append(buckets.Buckets, client.LoadBucketByName(bucketName))
			}
		}

		critical, err := internal.ParseThreshold(CriticalBucketSize)
		if err != nil {
			check.ExitError(err)
		}

		warning, err := internal.ParseThreshold(WarningBucketSize)
		if err != nil {
			check.ExitError(err)
		}

		for idx, bucket := range buckets.Buckets {
			objectsOutput.V2Output, err = client.LoadAllObjectsFromBucket(*bucket.Name, "")
			if err != nil {
				check.ExitError(err)
			}

			bucketSize := objectsOutput.CalculateBucketSize()

			crit, err := check.ParseThreshold(strconv.FormatUint(critical, 10))
			if err != nil {
				check.ExitError(err)
			}

			warn, err := check.ParseThreshold(strconv.FormatUint(warning, 10))
			if err != nil {
				check.ExitError(err)
			}

			if crit.DoesViolate(float64(bucketSize)) {
				rc = 2
				totalCrit++
			} else if warn.DoesViolate(float64(bucketSize)) {
				rc = 1
				totalWarn++
			} else {
				rc = 0
			}

			states = append(states, rc)

			output += objectsOutput.GetBucketOutput(bucketSize, rc)

			if len(buckets.Buckets) > 1 && !(len(buckets.Buckets) == idx+1) {
				output += "\n"
			}

			p := perfdata.Perfdata{
				Label: *bucket.Name,
				Value: uint64(bucketSize) / internal.MebiByte,
				Uom:   "MB",
				Warn:  &check.Threshold{Upper: float64(warning / internal.MebiByte)},
				Crit:  &check.Threshold{Upper: float64(critical / internal.MebiByte)},
			}

			perf.Add(&p)
		}

		summary += fmt.Sprintf("Found %d buckets - critical %d - warning %d\n", len(buckets.Buckets), totalCrit, totalWarn)

		check.ExitRaw(result.WorstState(states...), summary+output, "|", perf.String())
	},
}

func init() {
	s3Cmd.AddCommand(s3BucketCmd)

	s3BucketFlags := s3BucketCmd.Flags()
	s3BucketFlags.StringSliceVarP(&BucketNames, "buckets", "b", nil,
		"Name of the S3 bucket. If '--buckets' is empty, all buckets will be evaluated.")
	s3BucketFlags.StringVarP(&CriticalBucketSize, "crit-bucket-size", "c", "20Gb",
		"Critical threshold for the size of the specified bucket. Alerts if size is greater than critical threshold. Values are MB.\n"+
			"Possible  values are MB, GB or TB. Without any identifier specified MB is used.")
	s3BucketFlags.StringVarP(&WarningBucketSize, "warn-bucket-size", "w", "10Gb",
		"Warning threshold for the size of the specified bucket. Alerts if size is greater than warning threshold. Values are MB.\n"+
			"Possible  values are MB, GB or TB. Without any identifier specified MB is used.")

	s3BucketFlags.SortFlags = false
}
