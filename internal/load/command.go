package load

import (
	"golang.org/x/crypto/ssh"
	"io"
	"os"
)

// ExecCommand 执行命令
func ExecCommand(session *ssh.Session, command string) error {
	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return err
	}

	if err := session.Start(command); err != nil {
		return err
	}

	go func() {
		_, err := io.Copy(os.Stdout, stdout)
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		_, err := io.Copy(os.Stderr, stderr)
		if err != nil {
			panic(err)
		}
	}()

	return session.Wait()
}
