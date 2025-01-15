package load

import (
	"golang.org/x/crypto/ssh"
	"strconv"
)

// GetSSHClient 获取SSH客户端
func GetSSHClient(protocol string, host string, port int, username string, password string) (*ssh.Client, error) {
	return ssh.Dial(protocol, host+":"+strconv.Itoa(port), &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
}
