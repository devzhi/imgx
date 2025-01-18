package cmd

import (
	"fmt"
	"github.com/devzhi/imgx/internal/load"
	"github.com/devzhi/imgx/internal/pull"
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
		remove, _ := cmd.Flags().GetBool("rm")
		// 构造pull参数
		pullFlags := &pull.Flag{
			Image:  image,
			Tag:    tag,
			Arch:   arch,
			OsFlag: osFlag,
		}
		// 拉取镜像
		output, err := pull.Execute(pullFlags)
		if err != nil {
			fmt.Println("Error pulling image", err)
			return
		}
		// 构造load参数
		loadFlags := &load.Flag{
			InputFile: *output,
			Host:      host,
			Port:      port,
			Username:  username,
			Password:  password,
			Protocol:  protocol,
			Remove:    remove,
		}
		// 执行load命令
		load.Execute(loadFlags)
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
	xCommand.Flags().BoolP("rm", "r", true, "remove the image file after successful loading")
}
