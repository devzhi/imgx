package load

import (
	"fmt"
	"os"
)

type Flag struct {
	InputFile string
	Host      string
	Port      int
	Username  string
	Password  string
	Protocol  string
	Remove    bool
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
	err = UploadFile(client, "./"+flag.InputFile, tempDir+"/"+flag.InputFile)
	if err != nil {
		fmt.Println("Error uploading file", err)
		return err
	}
	// 导入镜像
	image, success, err := LoadImage(client, tempDir+"/"+flag.InputFile, flag.Password)
	if err != nil {
		fmt.Println("Error loading image", err)
		return err
	} else if success {
		fmt.Println("Image loaded", image)
		if flag.Remove {
			err := os.Remove(flag.InputFile)
			if err != nil {
				fmt.Println("Error removing image file", err)
			} else {
				fmt.Println("Image file removed")
			}
		}
	} else {
		fmt.Println("Error loading image", image)
	}
	return nil
}
