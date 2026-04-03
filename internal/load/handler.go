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
	if flag == nil {
		return fmt.Errorf("load flag is required")
	}
	if flag.InputFile == "" {
		return fmt.Errorf("input file is required")
	}
	if _, err := os.Stat(flag.InputFile); err != nil {
		return fmt.Errorf("input file: %w", err)
	}
	if flag.Host == "" {
		return fmt.Errorf("host is required")
	}
	if flag.Username == "" {
		return fmt.Errorf("username is required")
	}
	if flag.Password == "" {
		return fmt.Errorf("password is required")
	}

	client, err := GetSSHClient(flag.Protocol, flag.Host, flag.Port, flag.Username, flag.Password)
	if err != nil {
		fmt.Println("Error connecting to remote host", err)
		return err
	}
	defer client.Close()

	tempDir, err := CreateTempDir(client)
	if err != nil {
		fmt.Println("Error creating temp dir", err)
		return err
	}

	remotePath := buildRemotePath(tempDir, flag.InputFile)
	err = UploadFile(client, flag.InputFile, remotePath)
	if err != nil {
		fmt.Println("Error uploading file", err)
		fmt.Println("Try using scp to upload image")
		err = UploadBySCP(client, flag.InputFile, remotePath)
		if err != nil {
			fmt.Println("Error uploading file", err)
			return fmt.Errorf("upload file: %w", err)
		}
	}

	output, success, err := LoadImage(client, remotePath, flag.Password, flag.DockerPath)
	if err != nil {
		fmt.Println("Error loading image", err)
		fmt.Println("Output:\n", output)
		return err
	}
	if !success {
		fmt.Println("Error loading image", output)
		return fmt.Errorf("docker load did not report success")
	}

	fmt.Println("Image loaded", output)
	if flag.Remove {
		if err := os.Remove(flag.InputFile); err != nil {
			fmt.Println("Error removing image file", err)
		} else {
			fmt.Println("Image file removed")
		}
	}
	return nil
}

func buildRemotePath(tempDir, inputFile string) string {
	return strings.ReplaceAll(filepath.Join(tempDir, filepath.Base(inputFile)), "\\", "/")
}
