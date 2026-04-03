package load

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// ExecuteCommand 执行命令
func ExecuteCommand(client *ssh.Client, cmd string, sudoPassword string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return "", fmt.Errorf("request for pseudo terminal failed: %v", err)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	stdin, err := session.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	if err := session.Start(cmd); err != nil {
		return "", fmt.Errorf("failed to start command: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if strings.Contains(stdoutBuf.String(), "[sudo]") ||
		strings.Contains(stderrBuf.String(), "[sudo]") ||
		strings.Contains(stdoutBuf.String(), "Password:") ||
		strings.Contains(stderrBuf.String(), "Password:") {
		_, _ = stdin.Write([]byte(sudoPassword + "\n"))
	}

	err = session.Wait()
	if err != nil {
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
	output, err := ExecuteCommand(client, dockerLoadCommand(dockerPath, filePath), password)
	if err != nil {
		return output, false, err
	}
	if strings.Contains(output, "Loaded image") {
		return output, true, nil
	}
	return output, false, nil
}

func dockerLoadCommand(dockerPath string, filePath string) string {
	return fmt.Sprintf("sudo %s load -i %s", shellQuote(dockerPath), shellQuote(filePath))
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}

// GetDockerPath 获取Docker路径
func GetDockerPath(client *ssh.Client, sudoPassword string) (string, error) {
	dockerPaths := []string{
		"/usr/local/bin/docker",
		"/usr/bin/docker",
		"/bin/docker",
	}

	whichOutput, whichErr := ExecuteCommand(client, "which docker", sudoPassword)
	if whichErr == nil {
		dockerPath := strings.TrimSpace(whichOutput)
		if dockerPath != "" {
			return dockerPath, nil
		}
	}

	for _, path := range dockerPaths {
		checkCmd := fmt.Sprintf("test -f %s && echo %s", shellQuote(path), shellQuote(path))
		output, err := ExecuteCommand(client, checkCmd, sudoPassword)
		if err == nil && strings.TrimSpace(output) == path {
			return path, nil
		}
	}

	return "", fmt.Errorf("docker executable not found on remote host")
}
