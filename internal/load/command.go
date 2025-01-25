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
	if strings.Contains(stdoutBuf.String(), "[sudo]") ||
		strings.Contains(stderrBuf.String(), "[sudo]") ||
		strings.Contains(stdoutBuf.String(), "Password:") ||
		strings.Contains(stderrBuf.String(), "Password:") {
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
func LoadImage(client *ssh.Client, filePath string, password string, dockerPath string) (string, bool, error) {
	output, err := ExecuteCommand(client, "sudo "+dockerPath+" load -i "+filePath, password)
	if err != nil {
		return output, false, err
	} else if strings.Contains(output, "Loaded image") {
		return output, true, nil
	} else {
		return output, false, nil
	}
}

// GetDockerPath 获取Docker路径
func GetDockerPath(client *ssh.Client, sudoPassword string) (string, error) {
	// Potential Docker paths to check
	dockerPaths := []string{
		"/usr/local/bin/docker",
		"/usr/bin/docker",
		"/bin/docker",
	}

	// Check which command first
	whichCmd := "which docker"
	whichOutput, whichErr := ExecuteCommand(client, whichCmd, sudoPassword)

	if whichErr == nil {
		dockerPath := strings.TrimSpace(whichOutput)
		if dockerPath != "" {
			return dockerPath, nil
		}
	}

	// Manual path checking
	for _, path := range dockerPaths {
		checkCmd := fmt.Sprintf("test -f %s && echo %s", path, path)
		output, err := ExecuteCommand(client, checkCmd, sudoPassword)
		if err == nil && strings.TrimSpace(output) == path {
			return path, nil
		}
	}

	return "", fmt.Errorf("docker executable not found on remote host")
}
