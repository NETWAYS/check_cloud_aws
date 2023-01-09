package cmd

import (
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var (
	Timeout         = 30
	Profile         string
	Region          string
	CredentialsFile string
)

var rootCmd = &cobra.Command{
	Use:   "check_cloud_aws",
	Short: "An Icinga check plugin to check Amazon Web Services",
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

	p := rootCmd.PersistentFlags()
	// Default is empty and set later with the user's Home dir
	// If we set the default before the help text will show the full user's home dir which we decided against
	p.StringVarP(&CredentialsFile, "credentials-file", "C", "", "Path to the credentials file (default \"$HOME/.aws/credentials\")")

	if CredentialsFile == "" {
		CredentialsFile = filepath.Join(userHomeDir(), ".aws", "credentials")
	}

	p.StringVarP(&Region, "region", "R", "eu-central-1", "The AWS region to send requests to")
	p.StringVarP(&Profile, "profile", "P", "default", "The AWS profile name, which represents a separate credential profile in the credential file")
	p.IntVarP(&Timeout, "timeout", "t", Timeout, "Timeout for the check")

	rootCmd.Flags().SortFlags = false
	p.SortFlags = false
}

// Wrapper for os.UserHomeDir when getting config files for example
func userHomeDir() string {
	var home string

	if runtime.GOOS == "windows" { // Windows
		home = os.Getenv("USERPROFILE")
	} else { // Linux/macOS
		home = os.Getenv("HOME")
	}

	if len(home) > 0 {
		return home
	}

	currUser, _ := user.Current()
	if currUser != nil {
		home = currUser.HomeDir
	}

	return home
}
