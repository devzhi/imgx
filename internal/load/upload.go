package load

import (
	"context"
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/pkg/sftp"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/crypto/ssh"
	"os"
)

// UploadFile 上传文件
func UploadFile(client *ssh.Client, localPath string, remotePath string) error {
	// 打开本地文件
	localFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	// 获取本地文件大小
	localFileInfo, err := localFile.Stat()
	if err != nil {
		return err
	}
	localFileSize := localFileInfo.Size()

	// 建立Sftp客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	// 创建远程文件
	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Try using scp to upload file")
		return UploadBySCP(client, localPath, remotePath)
	}
	defer remoteFile.Close()

	// 创建进度条
	bar := progressbar.DefaultBytes(
		localFileSize,
		"Uploading",
	)

	// 创建带进度条的Reader
	progressReader := progressbar.NewReader(localFile, bar)

	// 上传文件
	_, err = remoteFile.ReadFrom(&progressReader)
	return err
}

func UploadBySCP(client *ssh.Client, localPath string, remotePath string) error {
	// 创建scp客户端
	scpClient, err := scp.NewClientBySSH(client)
	if err != nil {
		return err
	}
	defer scpClient.Close()
	// 打开本地文件
	localFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()
	// 上传文件
	fmt.Printf("Uploading %s to %s\n", localPath, remotePath)
	err = scpClient.CopyFromFile(context.Background(), *localFile, remotePath, "0655")
	if err != nil {
		return err
	}
	fmt.Printf("\nSuccessfully uploaded %s to %s\n", localPath, remotePath)
	return nil
}
