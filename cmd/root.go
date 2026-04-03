package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var showVersion bool

var rootCmd = &cobra.Command{
	Use:   "imgx",
	Short: "imgx is a Docker image transfer tool",
	Long:  `imgx is a Docker image transport tool that pulls images from docker hub and pushes them to a target server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			versionCmd.Run(cmd, args)
			return nil
		}

		return errors.New("unrecognized command")
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "show version info")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
