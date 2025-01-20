package cmd

import (
	"fmt"
	"github.com/devzhi/imgx/internal/load"
	"github.com/devzhi/imgx/internal/util"
	"github.com/spf13/cobra"
)

var loadCommand = &cobra.Command{
	Use:   "load",
	Short: "Load the image to the remote host",
	Run: func(cmd *cobra.Command, args []string) {
		// 获取输入文件
		input, err := cmd.Flags().GetString("input")
		if err != nil {
			fmt.Println("Error getting input file", err)
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
		password, err := util.ReadPassword()
		if err != nil {
			fmt.Println("Error reading password", err)
		}
		if password == "" {
			fmt.Println("Error: password is required")
			return
		}
		protocol, _ := cmd.Flags().GetString("protocol")
		// 构造load参数
		flag := &load.Flag{
			InputFile: input,
			Host:      host,
			Port:      port,
			Username:  username,
			Password:  password,
			Protocol:  protocol,
			Remove:    false,
		}
		// 执行load命令
		err = load.Execute(flag)
	},
}

func init() {
	// 添加load命令
	rootCmd.AddCommand(loadCommand)
	// 添加load命令的flag
	loadCommand.Flags().StringP("input", "i", "", "load image input file")
	loadCommand.Flags().StringP("host", "H", "", "load image host")
	loadCommand.Flags().IntP("port", "P", 22, "load image host's port")
	loadCommand.Flags().StringP("username", "u", "", "load image host's username")
	loadCommand.Flags().BoolP("password", "p", false, "load image host's password")
	loadCommand.Flags().String("protocol", "tcp", "load image host's ssh protocol")
}
