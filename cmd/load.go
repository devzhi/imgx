package cmd

import (
	"fmt"
	"github.com/devzhi/imgx/internal/load"
	"github.com/spf13/cobra"
	"os"
)

var loadCommand = &cobra.Command{
	Use:   "load [input]",
	Short: "Load the image to the remote host",
	Run: func(cmd *cobra.Command, args []string) {
		// 获取输入文件
		if len(args) == 0 {
			fmt.Println("Error: input file is required")
			return
		}
		inputFile := &args[0]
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
		// 连接远程主机
		client, err := load.GetSSHClient(protocol, host, port, username, password)
		if err != nil {
			fmt.Println("Error connecting to remote host", err)
			return
		}
		defer client.Close()
		// 创建临时目录
		tempDir, err := load.CreateTempDir(client)
		if err != nil {
			fmt.Println("Error creating temp dir", err)
			return
		}
		// 上传文件
		err = load.UploadFile(client, "./"+*inputFile, tempDir+"/"+*inputFile)
		if err != nil {
			fmt.Println("Error uploading file", err)
			return
		}
		// 导入镜像
		image, success, err := load.LoadImage(client, tempDir+"/"+*inputFile, password)
		if err != nil {
			fmt.Println("Error loading image", err)
			return
		} else if success {
			fmt.Println("Image loaded", image)
			if remove {
				err := os.Remove(*inputFile)
				if err != nil {
					fmt.Println("Error removing image file", err)
				} else {
					fmt.Println("Image file removed")
				}
			}
		} else {
			fmt.Println("Error loading image", image)
		}
	},
}

func init() {
	// 添加load命令
	rootCmd.AddCommand(loadCommand)
	// 添加load命令的flag
	loadCommand.Flags().StringP("host", "H", "", "load image host")
	loadCommand.Flags().IntP("port", "P", 22, "load image host's port")
	loadCommand.Flags().StringP("username", "u", "", "load image host's username")
	loadCommand.Flags().StringP("password", "p", "", "load image host's password")
	loadCommand.Flags().String("protocol", "tcp", "load image host's ssh protocol")
	loadCommand.Flags().BoolP("rm", "r", false, "remove the image file after successful loading")
}
