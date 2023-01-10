package cmd

import (
	"fmt"
	"github.com/NETWAYS/check_cloud_aws/internal/ec2"
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
)

var (
	InstanceId   string
	InstanceName string
)

var ec2InstanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Checks a single EC2 instance",
	Run: func(cmd *cobra.Command, args []string) {
		if InstanceId != "" && InstanceName != "" {
			check.ExitError(fmt.Errorf("please specify only instance id or name"))
		}

		var (
			instance *ec2.Instance
			err      error
		)

		client := RequireEC2Client()

		if InstanceId != "" {
			instance, err = client.LoadInstance(InstanceId)
		} else if InstanceName != "" {
			instance, err = client.LoadInstanceByName(InstanceName)
		} else {
			check.ExitError(fmt.Errorf("please specify instance id or name"))
		}

		if err != nil {
			check.ExitError(fmt.Errorf("could not load instance: %w", err))
		}

		output := instance.GetOutput()
		output += "\n" + instance.GetLongOutput()

		check.Exit(instance.GetStatus(), output)
	},
}

func init() {
	ec2InstanceCmd.Flags().StringVarP(&InstanceId, "id", "i", "", "Search for ec2 instance by id")
	ec2InstanceCmd.Flags().StringVarP(&InstanceName, "name", "n", "", "Search for ec2 instance by name")

	ec2Cmd.AddCommand(ec2InstanceCmd)
}
