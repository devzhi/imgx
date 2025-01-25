package load

import (
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
)

// GetSSHClient 获取SSH客户端
func GetSSHClient(protocol string, host string, port int, username string, password string) (*ssh.Client, error) {
	return ssh.Dial(protocol, host+":"+strconv.Itoa(port), &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
}

// CreateTempDir 创建临时目录
func CreateTempDir(client *ssh.Client) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	// 创建临时目录
	output, err := session.CombinedOutput("mktemp -d")
	if err != nil {
		return "", err
	}
	splitOutput := strings.Split(strings.Trim(string(output), "\n"), "\n")

	return splitOutput[len(splitOutput)-1], nil
}
