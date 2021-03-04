package cmd

import (
	"fmt"
	"github.com/NETWAYS/check_cloud_aws/internal/ec2"
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
)

var AutoscalingName string

var ec2InstancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Checks multiple EC2 instances by name",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			instances *ec2.Instances
			err       error
		)

		client := RequireEC2Client()

		if InstanceName != "" && AutoscalingName == "" {
			instances, err = client.LoadAllInstancesByFilter(
				ec2.DescribeInstancesInput(ec2.Filter("tag:Name", InstanceName)))
		} else if AutoscalingName != "" && InstanceName == "" {
			instances, err = client.LoadAllInstancesByFilter(
				ec2.DescribeInstancesInput(ec2.Filter("tag:aws:autoscaling:groupName", AutoscalingName)))
		} else {
			instances, err = client.LoadAllInstancesByFilter(nil)
		}

		if err != nil {
			check.ExitError(fmt.Errorf("could not load instances: %w", err))
		}

		if len(instances.Instances) <= 0 {
			check.ExitError(fmt.Errorf("no instances found that matches the filter"))
		}

		states := map[string]int{}
		for _, instance := range instances.Instances {
			states[*instance.Instance.State.Name]++
		}

		summary := fmt.Sprintf("%d Instances found", len(instances.Instances))

		for state, count := range states {
			summary += fmt.Sprintf(" - %d %s", count, state)
		}

		check.Exit(instances.GetStatus(), summary+"\n\n"+instances.GetOutput())
	},
}

func init() {
	ec2InstancesCmd.Flags().StringVarP(&InstanceName, "name", "n", "", "Search for ec2 instances by name (e.g. instance*)")
	ec2InstancesCmd.Flags().StringVarP(&AutoscalingName, "autoscale", "a", "", "Search for ec2 instances by autoscaling group")

	ec2Cmd.AddCommand(ec2InstancesCmd)
}
