package cmd

import (
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
	"os"
)

var (
	Timeout         = 30
	Profile         string
	Region          string
	CredentialsFile string
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

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.DisableAutoGenTag = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	homedir, err := os.UserHomeDir()
	if err != nil {
		check.ExitError(err)
	}

	p := rootCmd.PersistentFlags()
	p.StringVarP(&CredentialsFile, "credentials-file", "C", homedir+"/.aws/credentials", "Path to the credentials file")
	p.StringVarP(&Region, "region", "R", "eu-central-1", "The AWS region to send requests to")
	p.StringVarP(&Profile, "profile", "P", "default", "The AWS profile name, which represents a separate credential profile in the credential file")
	p.IntVarP(&Timeout, "timeout", "t", Timeout, "Timeout for the check")

	rootCmd.Flags().SortFlags = false
	p.SortFlags = false
}
