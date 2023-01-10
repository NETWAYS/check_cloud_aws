package cmd

import (
	"fmt"
	"github.com/NETWAYS/check_cloud_aws/internal/common"
	"github.com/NETWAYS/check_cloud_aws/internal/ec2"
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
)

var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Checks in the EC2 context",
	Run:   Help,
}

func RequireEC2Client() *ec2.EC2Client {
	session, err := common.CreateSession(CredentialsFile, Profile, Region)
	if err != nil {
		check.ExitError(fmt.Errorf("could not setup AWS API session: %w", err))
	}

	return ec2.NewEC2Client(session)
}
