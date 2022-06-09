package cmd

import (
	"fmt"
	"github.com/NETWAYS/check_cloud_aws/internal/common"
	"github.com/NETWAYS/check_cloud_aws/internal/s3"
	"github.com/NETWAYS/go-check"

	"github.com/spf13/cobra"
)

var BucketNames []string

var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Checks in the S3 context",
	Run:   Help,
}

func RequireS3Client() *s3.S3Client {
	session, err := common.CreateSession(Profile, Region)
	if err != nil {
		check.ExitError(fmt.Errorf("could not setup AWS API session: %w", err))
	}

	return s3.NewS3Client(session)
}
