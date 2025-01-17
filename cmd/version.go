package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show imgx version info.",
	Run: func(cmd *cobra.Command, args []string) {
		version := "v0.1.0"
		fmt.Println("imgx version", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
