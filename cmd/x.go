package cmd

import (
	"errors"
	"fmt"

	"github.com/devzhi/imgx/internal/load"
	"github.com/devzhi/imgx/internal/pull"
	"github.com/devzhi/imgx/internal/util"
	"github.com/spf13/cobra"
)

var xCommand = &cobra.Command{
	Use:   "x [image]",
	Short: "Pulling and loading images to remote host",
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
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			return err
		}
		if host == "" {
			return errors.New("host is required")
		}
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		if username == "" {
			return errors.New("username is required")
		}

		password, err := util.ReadPassword()
		if err != nil {
			return fmt.Errorf("read password: %w", err)
		}
		if password == "" {
			return errors.New("password is required")
		}

		protocol, err := cmd.Flags().GetString("protocol")
		if err != nil {
			return err
		}
		dockerPath, err := cmd.Flags().GetString("docker-path")
		if err != nil {
			return err
		}
		save, err := cmd.Flags().GetBool("save")
		if err != nil {
			return err
		}

		pullFlags := &pull.Flag{
			Image:  args[0],
			Tag:    tag,
			Arch:   arch,
			OsFlag: osFlag,
			Path:   ".",
		}
		output, err := pull.Execute(pullFlags)
		if err != nil {
			return fmt.Errorf("pull image: %w", err)
		}

		loadFlags := &load.Flag{
			InputFile:  *output,
			Host:       host,
			Port:       port,
			Username:   username,
			Password:   password,
			Protocol:   protocol,
			Remove:     !save,
			DockerPath: dockerPath,
		}
		if err := load.Execute(loadFlags); err != nil {
			return fmt.Errorf("load image: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(xCommand)
	xCommand.Flags().StringP("tag", "t", "latest", "pull image tag")
	xCommand.Flags().StringP("arch", "a", "amd64", "pull image arch")
	xCommand.Flags().StringP("os", "o", "linux", "pull image os")
	xCommand.Flags().StringP("host", "H", "", "load image host")
	xCommand.Flags().IntP("port", "P", 22, "load image host's port")
	xCommand.Flags().StringP("username", "u", "", "load image host's username")
	xCommand.Flags().BoolP("password", "p", false, "load image host's password")
	xCommand.Flags().String("protocol", "tcp", "load image host's ssh protocol")
	xCommand.Flags().BoolP("save", "s", false, "save image to disk")
	xCommand.Flags().String("docker-path", "docker", "remote host's docker path")
}
