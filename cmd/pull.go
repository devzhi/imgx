package cmd

import (
	"context"
	"fmt"
	"github.com/devzhi/imgx/internal/pull"
	"github.com/spf13/cobra"
)

var pullCommand = &cobra.Command{
	Use:   "pull [image]",
	Short: "Pulling images from Docker hub locally",
	Run: func(cmd *cobra.Command, args []string) {
		// 获取flag参数
		tag, _ := cmd.Flags().GetString("tag")
		arch, _ := cmd.Flags().GetString("arch")
		osFlag, _ := cmd.Flags().GetString("os")
		// 执行pull命令
		// 获取token
		token, err := pull.GetToken(args[0])
		if err != nil {
			fmt.Println("Error getting token", err)
			return
		}

		// 获取镜像清单
		manifest, err := pull.GetManifest(token, args[0], tag, arch, osFlag)
		if err != nil {
			fmt.Println("Error getting manifest", err)
			return
		}

		// 下载镜像
		path, err := pull.DownloadImage(token, manifest, arch, osFlag, args[0], tag)
		if err != nil {
			fmt.Println("Error downloading image", err)
			return
		}
		fmt.Println("Image downloaded to", *path)

		// 删除临时文件
		defer pull.RemoveImageSaveDir(args[0], tag, arch, osFlag)

		// 打包镜像
		outputFile, err := pull.Package(*path, args[0], tag, arch, osFlag, nil)
		if err != nil {
			fmt.Println("Error packaging image", err)
			return
		}
		fmt.Println("\nImage packaged to", *outputFile)
		// 打包后的镜像写入context
		ctx := context.WithValue(cmd.Context(), "outputFile", *outputFile)
		cmd.SetContext(ctx)
	},
}

func init() {
	// 添加pull命令
	rootCmd.AddCommand(pullCommand)
	// 添加pull命令的flag
	pullCommand.Flags().StringP("tag", "t", "latest", "pull image tag")
	pullCommand.Flags().StringP("arch", "a", "amd64", "pull image arch")
	pullCommand.Flags().StringP("os", "o", "linux", "pull image os")
}
