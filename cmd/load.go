package cmd

import (
	"errors"
	"fmt"

	"github.com/devzhi/imgx/internal/load"
	"github.com/devzhi/imgx/internal/util"
	"github.com/spf13/cobra"
)

var loadCommand = &cobra.Command{
	Use:   "load",
	Short: "Load the image to the remote host",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := cmd.Flags().GetString("input")
		if err != nil {
			return err
		}
		if input == "" {
			return errors.New("input is required")
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

		flag := &load.Flag{
			InputFile:  input,
			Host:       host,
			Port:       port,
			Username:   username,
			Password:   password,
			Protocol:   protocol,
			Remove:     false,
			DockerPath: dockerPath,
		}

		if err := load.Execute(flag); err != nil {
			return fmt.Errorf("load image: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loadCommand)
	loadCommand.Flags().StringP("input", "i", "", "load image input file")
	loadCommand.Flags().StringP("host", "H", "", "load image host")
	loadCommand.Flags().IntP("port", "P", 22, "load image host's port")
	loadCommand.Flags().StringP("username", "u", "", "load image host's username")
	loadCommand.Flags().BoolP("password", "p", false, "load image host's password")
	loadCommand.Flags().String("protocol", "tcp", "load image host's ssh protocol")
	loadCommand.Flags().String("docker-path", "docker", "remote host's docker path")
}
