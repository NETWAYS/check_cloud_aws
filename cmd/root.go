package cmd

import (
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
	"os"
)

var (
	Timeout = 30
	Profile string
	Region  string
)

var rootCmd = &cobra.Command{
	Use:   "check_cloud_aws",
	Short: "Icinga check plugin to check Amazon EC2 instances",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		go check.HandleTimeout(Timeout)
	},
	Run: Help,
}

func Execute(version string) {
	defer check.CatchPanic()

	rootCmd.Version = version

	if err := rootCmd.Execute(); err != nil {
		check.ExitError(err)
	}
}

func Help(cmd *cobra.Command, strings []string) {
	_ = cmd.Usage()

	os.Exit(3)
}

func init() {
	rootCmd.AddCommand(ec2Cmd)
	rootCmd.AddCommand(s3Cmd)
	rootCmd.AddCommand(cloudfrontCmd)
	rootCmd.SetHelpFunc(Help)

	p := rootCmd.PersistentFlags()
	p.IntVarP(&Timeout, "timeout", "t", Timeout, "Timeout for the check")
	p.StringVarP(&Profile, "profile", "P", Profile, "AWS credential profile (~/.aws/credentials)")
	p.StringVarP(&Region, "region", "R", Region, "AWS region name (e.g. eu-central-1)")
}
