package load

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

// ExecuteCommand 执行命令
func ExecuteCommand(client *ssh.Client, cmd string, sudoPassword string) (string, error) {
	// 获取session
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	// 创建管道用于处理输入输出
	stdin, err := session.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	var outputBuffer bytes.Buffer
	session.Stdout = &outputBuffer
	session.Stderr = &outputBuffer

	// 启动命令
	err = session.Start(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to start command: %v", err)
	}

	// 监听输出中是否包含密码提示
	go func() {
		for {
			if strings.Contains(outputBuffer.String(), "[sudo] password") {
				stdin.Write([]byte(sudoPassword + "\n"))
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 等待命令完成
	err = session.Wait()
	if err != nil {
		if strings.Contains(outputBuffer.String(), "Sorry, try again") {
			return "", fmt.Errorf("incorrect sudo password")
		}
		return outputBuffer.String(), fmt.Errorf("command failed: %v", err)
	}

	return outputBuffer.String(), nil
}

// LoadImage 载入镜像
func LoadImage(client *ssh.Client, filePath string, password string) (string, error) {
	return ExecuteCommand(client, "sudo docker load -i "+filePath, password)
}
