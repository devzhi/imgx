package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show imgx version info.",
	Run: func(cmd *cobra.Command, args []string) {
		version := "1.1.2"
		fmt.Println("imgx version", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
