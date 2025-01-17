package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var xCommand = &cobra.Command{
	Use:   "x [image]",
	Short: "Pulling and loading images to remote host",
	Run: func(cmd *cobra.Command, args []string) {
		image := args[0]
		if image == "" {
			fmt.Println("Error: image is required")
			return
		}
		tag, _ := cmd.Flags().GetString("tag")
		arch, _ := cmd.Flags().GetString("arch")
		osFlag, _ := cmd.Flags().GetString("os")
		// 获取输入文件
		if len(args) == 0 {
			fmt.Println("Error: input file is required")
			return
		}
		// 获取flag参数
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			fmt.Println("Error getting host", err)
			return
		}
		if host == "" {
			fmt.Println("Error: host is required")
			return
		}
		port, _ := cmd.Flags().GetInt("port")
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			fmt.Println("Error getting username", err)
			return
		}
		if username == "" {
			fmt.Println("Error: username is required")
			return
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			fmt.Println("Error getting password", err)
			return
		}
		if password == "" {
			fmt.Println("Error: password is required")
			return
		}
		protocol, _ := cmd.Flags().GetString("protocol")
		// 构造pull参数
		pullArgs := []string{image, "-t", tag, "-a", arch, "-o", osFlag}
		// 拉取镜像
		pullCommand.Run(cmd, pullArgs)
		// 获取拉取的镜像
		// 获取拉取的镜像
		ctx := cmd.Context()
		outputFile, ok := ctx.Value("outputFile").(string)
		if !ok || outputFile == "" {
			fmt.Println("Error: outputFile is empty or not found in context")
			return
		}
		// 构造load参数
		loadArgs := []string{outputFile, "-r", "-H", host, "-P", fmt.Sprintf("%d", port), "-u", username, "-p", password, "--protocol", protocol}
		loadCommand.Run(cmd, loadArgs)
	},
}

func init() {
	// 添加x命令
	rootCmd.AddCommand(xCommand)
	// 添加flag
	xCommand.Flags().StringP("tag", "t", "latest", "pull image tag")
	xCommand.Flags().StringP("arch", "a", "amd64", "pull image arch")
	xCommand.Flags().StringP("os", "o", "linux", "pull image os")
	xCommand.Flags().StringP("host", "H", "", "load image host")
	xCommand.Flags().IntP("port", "P", 22, "load image host's port")
	xCommand.Flags().StringP("username", "u", "", "load image host's username")
	xCommand.Flags().StringP("password", "p", "", "load image host's password")
	xCommand.Flags().String("protocol", "tcp", "load image host's ssh protocol")
	xCommand.Flags().BoolP("rm", "r", false, "remove the image file after successful loading")
}
