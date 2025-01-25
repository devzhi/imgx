package load

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Flag struct {
	InputFile  string
	Host       string
	Port       int
	Username   string
	Password   string
	Protocol   string
	Remove     bool
	DockerPath string
}

// Execute 执行load命令
func Execute(flag *Flag) error {
	// 连接远程主机
	client, err := GetSSHClient(flag.Protocol, flag.Host, flag.Port, flag.Username, flag.Password)
	if err != nil {
		fmt.Println("Error connecting to remote host", err)
		return err
	}
	defer client.Close()
	// 创建临时目录
	tempDir, err := CreateTempDir(client)
	if err != nil {
		fmt.Println("Error creating temp dir", err)
		return err
	}
	// 上传文件
	remotePath := strings.ReplaceAll(filepath.Join(tempDir, flag.InputFile), "\\", "/")
	err = UploadFile(client, flag.InputFile, remotePath)
	if err != nil {
		fmt.Println("Error uploading file", err)
		return err
	}
	// 导入镜像
	output, success, err := LoadImage(client, remotePath, flag.Password, flag.DockerPath)
	if err != nil {
		fmt.Println("Error loading image", err)
		fmt.Println("Output:\n", output)
		return err
	} else if success {
		fmt.Println("Image loaded", output)
		if flag.Remove {
			err := os.Remove(flag.InputFile)
			if err != nil {
				fmt.Println("Error removing image file", err)
			} else {
				fmt.Println("Image file removed")
			}
		}
	} else {
		fmt.Println("Error loading image", output)
	}
	return nil
}
