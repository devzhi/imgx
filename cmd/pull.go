package cmd

import (
	"fmt"

	"github.com/devzhi/imgx/internal/pull"
	"github.com/spf13/cobra"
)

var pullCommand = &cobra.Command{
	Use:   "pull [image]",
	Short: "Pulling images from Docker hub locally",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tag, err := cmd.Flags().GetString("tag")
		if err != nil {
			return err
		}
		arch, err := cmd.Flags().GetString("arch")
		if err != nil {
			return err
		}
		osFlag, err := cmd.Flags().GetString("os")
		if err != nil {
			return err
		}
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}

		command := &pull.Flag{
			Image:  args[0],
			Tag:    tag,
			Arch:   arch,
			OsFlag: osFlag,
			Path:   path,
		}

		if _, err := pull.Execute(command); err != nil {
			return fmt.Errorf("pull image: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pullCommand)
	pullCommand.Flags().StringP("tag", "t", "latest", "pull image tag")
	pullCommand.Flags().StringP("arch", "a", "amd64", "pull image arch")
	pullCommand.Flags().StringP("os", "o", "linux", "pull image os")
	pullCommand.Flags().StringP("path", "p", "./", "pull image path")
}
