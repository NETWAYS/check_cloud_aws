package cmd

import (
	"errors"
	"fmt"
	"github.com/NETWAYS/check_cloud_aws/internal"
	s3i "github.com/NETWAYS/check_cloud_aws/internal/s3"
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
	BucketNames        []string
)

var s3BucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "Checks the size of a single bucket or multiple buckets",
	Example: `
	check_cloud_aws s3 bucket
	OK - 1 Buckets: 0 Critical - 0 Warning - 1 OK
	 \_[OK] my-bucket - value: 100MiB | my-bucket=100MB;10240;20480

	check_cloud_aws s3 bucket --crit-bucket-size 10
	CRITICAL - 1 Buckets: 1 Critical - 0 Warning - 0 OK
	 \_[CRITICAL] my-bucket - value: 100MiB | my-bucket=100MB;10240;10

	check_cloud_aws s3 bucket --crit-bucket-size 5GB
	CRITICAL - 1 Buckets: 1 Critical - 0 Warning - 0 OK
	 \_[CRITICAL] my-bucket - value: 100MiB | my-bucket=100MB;10240;10`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err       error
			summary   string
			output    string
			states    []int
			totalCrit int
			totalWarn int
			totalOk   int
			rc        int
			perf      perfdata.PerfdataList
		)

		buckets := s3.ListBucketsOutput{}
		objectsOutput := s3i.V2Output{}

		client := RequireS3Client()

		// Load all buckets of if no buckets are speficied
		if len(BucketNames) == 0 {
			bk, err := client.LoadAllBuckets()
			if err != nil {
				check.ExitError(err)
			}

			buckets = *bk
		} else {
			// Load specific buckets
			for _, bucketName := range BucketNames {
				b, err := client.LoadBucketByName(bucketName)
				// Requested bucket does not exist, skip
				if errors.Is(err, s3i.ErrBucketNotFound) {
					continue
				}
				buckets.Buckets = append(buckets.Buckets, b)
			}
		}

		if len(buckets.Buckets) == 0 {
			check.ExitError(fmt.Errorf("No buckets available"))
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
				totalOk++
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

		summary += fmt.Sprintf("%d Buckets: %d Critical - %d Warning - %d OK\n", len(buckets.Buckets), totalCrit, totalWarn, totalOk)

		check.ExitRaw(result.WorstState(states...), summary+output, "|", perf.String())
	},
}

func init() {
	s3Cmd.AddCommand(s3BucketCmd)

	s3BucketFlags := s3BucketCmd.Flags()
	s3BucketFlags.StringSliceVarP(&BucketNames, "buckets", "b", nil,
		"Name of the S3 bucket. If '--buckets' is empty, all buckets will be evaluated.")
	s3BucketFlags.StringVarP(&CriticalBucketSize, "crit-bucket-size", "c", "20Gb",
		"Critical threshold for the size of the specified bucket. Alerts if the size is greater than the critical threshold.\n"+
			"Possible units are MB, GB or TB (defaults to MB if none is specified).")
	s3BucketFlags.StringVarP(&WarningBucketSize, "warn-bucket-size", "w", "10Gb",
		"Warning threshold for the size of the specified bucket. Alerts if the size is greater than the warning threshold.\n"+
			"Possible units are MB, GB or TB (defaults to MB if none is specified).")

	s3BucketFlags.SortFlags = false
}
