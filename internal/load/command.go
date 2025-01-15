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
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// 设置伪终端以处理 sudo 提示
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // 禁用回显
		ssh.TTY_OP_ISPEED: 14400, // 输入速度
		ssh.TTY_OP_OSPEED: 14400, // 输出速度
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return "", fmt.Errorf("request for pseudo terminal failed: %v", err)
	}

	// 创建输出缓冲区
	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	// 获取输入管道
	stdin, err := session.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	// 启动命令
	if err := session.Start(cmd); err != nil {
		return "", fmt.Errorf("failed to start command: %v", err)
	}

	// 等待片刻以让可能的密码提示出现
	time.Sleep(100 * time.Millisecond)

	// 如果看到密码提示，输入密码
	if strings.Contains(stdoutBuf.String(), "[sudo]") || strings.Contains(stderrBuf.String(), "[sudo]") {
		stdin.Write([]byte(sudoPassword + "\n"))
	}

	// 等待命令完成
	err = session.Wait()
	if err != nil {
		// 检查是否是密码错误
		if strings.Contains(stdoutBuf.String(), "Sorry, try again") ||
			strings.Contains(stderrBuf.String(), "Sorry, try again") {
			return "", fmt.Errorf("incorrect sudo password")
		}
		return stdoutBuf.String() + stderrBuf.String(), fmt.Errorf("command failed: %v", err)
	}

	return stdoutBuf.String() + stderrBuf.String(), nil
}

// LoadImage 载入镜像
func LoadImage(client *ssh.Client, filePath string, password string) (string, error) {
	return ExecuteCommand(client, "sudo docker load -i "+filePath, password)
}
