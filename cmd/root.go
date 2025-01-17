package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "imgx",
	Short: "imgx is a Docker image transfer tool",
	Long:  `imgx is a Docker image transport tool that pulls images from docker hub and pushes them to a target server.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.PrintErrln(errors.New("unrecognized command"))
	},
}

func Execute() {
	rootCmd.Execute()
}
