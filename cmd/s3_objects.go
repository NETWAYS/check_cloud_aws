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
	CriticalObjectSize string
	WarningObjectSize  string
	ShowPerfdata       bool
	ObjectPrefix       string
)

var s3ObjectCmd = &cobra.Command{
	Use:   "object",
	Short: "Checks the size of objects, stored in a single bucket or multiple buckets",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err          error
			summary      string
			output       string
			states       []int
			totalCrit    int
			totalWarn    int
			totalObjects int64
			rc           int
			perf         perfdata.PerfdataList
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

		critical, err := internal.ParseThreshold(CriticalObjectSize)
		if err != nil {
			check.ExitError(err)
		}

		warning, err := internal.ParseThreshold(WarningObjectSize)
		if err != nil {
			check.ExitError(err)
		}

		for _, bucket := range buckets.Buckets {
			objectsOutput.V2Output, err = client.LoadAllObjectsFromBucket(*bucket.Name, ObjectPrefix)
			if err != nil {
				check.ExitError(err)
			}

			crit, err := check.ParseThreshold(strconv.FormatUint(critical, 10))
			if err != nil {
				check.ExitError(err)
			}

			warn, err := check.ParseThreshold(strconv.FormatUint(warning, 10))
			if err != nil {
				check.ExitError(err)
			}

			output += fmt.Sprintf("[%s]:\n", *bucket.Name)

			for _, content := range objectsOutput.V2Output.Contents {
				if crit.DoesViolate(float64(*content.Size)) {
					rc = 2
					totalCrit++
				} else if warn.DoesViolate(float64(*content.Size)) {
					rc = 1
					totalWarn++
				} else {
					rc = 0
				}

				states = append(states, rc)

				if rc != 0 {
					output += objectsOutput.GetObjectOutput(*content.Size, rc, *content.Key)
				}

				p := perfdata.Perfdata{
					Label: *content.Key,
					Value: uint64(*content.Size) / internal.MebiByte,
					Uom:   "MB",
					Warn:  &check.Threshold{Upper: float64(warning / internal.MebiByte)},
					Crit:  &check.Threshold{Upper: float64(critical / internal.MebiByte)},
				}

				perf.Add(&p)

				totalObjects++
			}
		}

		if result.WorstState(states...) == 0 {
			output = fmt.Sprintf("")
		}

		summary += fmt.Sprintf("Found %d objects - critical %d - warning %d", totalObjects, totalCrit, totalWarn)

		if len(buckets.Buckets) > 1 {
			if result.WorstState(states...) != 0 {
				summary += "\n"
			}
		}

		if ShowPerfdata {
			check.ExitRaw(result.WorstState(states...), summary+output, "|", perf.String())
		} else {
			check.ExitRaw(result.WorstState(states...), summary+output)
		}
	},
}

func init() {
	s3Cmd.AddCommand(s3ObjectCmd)

	s3ObjectFlags := s3ObjectCmd.Flags()
	s3ObjectFlags.StringSliceVarP(&BucketNames, "buckets", "b", nil,
		"Name of one or multiple S3 buckets. If '--buckets' is empty, all buckets will be evaluated.")
	s3ObjectFlags.StringVar(&ObjectPrefix, "prefix", "",
		"Limits the response to keys that begin with the specified prefix, e.G. '--prefix test' filters all objects which starts with 'test'.\n"+
			"NOTE: Keep in mind, that objects beneath a directory will be ignored!")
	s3ObjectFlags.StringVarP(&CriticalObjectSize, "crit-object-size", "c", "1gb",
		"Critical threshold for the size of the object. Alerts if size is greater than critical threshold.\n"+
			"Possible  values are MB, GB or TB. Without any identifier specified MB is used.")
	s3ObjectFlags.StringVarP(&WarningObjectSize, "warn-object-size", "w", "800mb",
		"Critical threshold for the size of the object. Alerts if size is greater than warning threshold.\n"+
			"Possible  values are MB, GB or TB. Without any identifier specified MB is used.")
	s3ObjectFlags.BoolVarP(&ShowPerfdata, "perfdata", "p", false,
		"Displays perfdata and lists ALL objects in the specified bucket.")

	s3ObjectFlags.SortFlags = false
}
