package cmd

import (
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
		// 构造pull参数
		command := &pull.Flag{
			Image:  args[0],
			Tag:    tag,
			Arch:   arch,
			OsFlag: osFlag,
		}
		// 执行pull命令
		_, err := pull.Execute(command)
		if err != nil {
			fmt.Println("Error pulling image", err)
			return
		}
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
