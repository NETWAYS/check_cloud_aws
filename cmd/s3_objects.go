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
	CriticalObjectSize string
	WarningObjectSize  string
	ShowPerfdata       bool
	ObjectPrefix       string
	ObjectBucketNames  []string
)

var s3ObjectCmd = &cobra.Command{
	Use:   "object",
	Short: "Checks the size of objects, stored in a single bucket or multiple buckets",
	Example: `
	check_cloud_aws s3 object
	OK - 2 Objects: 0 Critical - 0 Warning - 2 OK
	 \_[my-bucket]:
		 \_[OK] foo.fs: 100MiB
		 \_[OK] bar.fs: 100MiB

	check_cloud_aws s3 object --prefix file
	OK - 3 Objects: 0 Critical - 0 Warning - 3 OK
	 \_[my-bucket]:
		 \_[OK] file_1.fs: 10MiB
		 \_[OK] file_2.fs: 20MiB
		 \_[OK] file_3.fs: 30MiB

	check_cloud_aws s3 object --crit-object-size 10KB --buckets foo --buckets bar
	CRITICAL - 2 Objects: 2 Critical - 0 Warning - 0 OK
	 \_[foo]:
		 \_[CRITICAL] file.fs: 100MiB
	 \_[bar]:
		 \_[CRITICAL] file.fs: 100MiB`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err          error
			summary      string
			output       string
			states       []int
			totalCrit    int
			totalWarn    int
			totalOk      int
			totalObjects int64
			rc           int
			perf         perfdata.PerfdataList
		)

		buckets := s3.ListBucketsOutput{}
		objectsOutput := s3i.V2Output{}

		client := RequireS3Client()

		// Load all buckets of if no buckets are speficied
		if len(ObjectBucketNames) == 0 {
			bk, err := client.LoadAllBuckets()
			if err != nil {
				check.ExitError(err)
			}

			buckets = *bk
		} else {
			// Load specific buckets
			for _, bucketName := range ObjectBucketNames {
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

			output += fmt.Sprintf(" \\_[%s]:\n", *bucket.Name)

			for _, content := range objectsOutput.V2Output.Contents {
				if crit.DoesViolate(float64(*content.Size)) {
					rc = 2
					totalCrit++
				} else if warn.DoesViolate(float64(*content.Size)) {
					rc = 1
					totalWarn++
				} else {
					rc = 0
					totalOk++
				}

				states = append(states, rc)

				output += objectsOutput.GetObjectOutput(*content.Size, rc, *content.Key)

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

		summary += fmt.Sprintf("%d Objects: %d Critical - %d Warning - %d OK\n", totalObjects, totalCrit, totalWarn, totalOk)

		if ShowPerfdata {
			check.ExitRaw(result.WorstState(states...), summary+output, "|", perf.String())
		}

		check.ExitRaw(result.WorstState(states...), summary+output)
	},
}

func init() {
	s3Cmd.AddCommand(s3ObjectCmd)

	s3ObjectFlags := s3ObjectCmd.Flags()
	s3ObjectFlags.StringSliceVarP(&ObjectBucketNames, "buckets", "b", nil,
		"Name of one or multiple S3 buckets. If '--buckets' is empty, all buckets will be evaluated.")
	s3ObjectFlags.StringVar(&ObjectPrefix, "prefix", "",
		"Limits the response to keys that begin with the specified prefix, e.G. '--prefix test' filters all objects which starts with 'test'.\n"+
			"NOTE: Keep in mind, that objects beneath a directory will be ignored.")
	s3ObjectFlags.StringVarP(&CriticalObjectSize, "crit-object-size", "c", "1gb",
		"Critical threshold for the size of the object. Alerts if the size is greater than the critical threshold.\n"+
			"Possible units are MB, GB or TB (defaults to MB if none is specified).")
	s3ObjectFlags.StringVarP(&WarningObjectSize, "warn-object-size", "w", "800mb",
		"Critical threshold for the size of the object. Alerts if the size is greater than the warning threshold.\n"+
			"Possible units are MB, GB or TB (defaults to MB if none is specified).")
	s3ObjectFlags.BoolVarP(&ShowPerfdata, "perfdata", "p", false,
		"Displays perfdata and lists all objects in the specified bucket.")

	s3ObjectFlags.SortFlags = false
}
