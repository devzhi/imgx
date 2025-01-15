package main

import (
	"flag"
	"fmt"
	"github.com/devzhi/imgx/internal/load"
	"github.com/devzhi/imgx/internal/pull"
	"os"
)

var (
	imageName   = flag.String("name", "", "name of the image")
	tag         = flag.String("tag", "latest", "tag of the image")
	arch        = flag.String("arch", "amd64", "architecture of the image")
	osFlag      = flag.String("os", "linux", "operating system of the image")
	showVersion = flag.Bool("version", false, "show version")
	protocol    = flag.String("protocol", "tcp", "protocol of the remote host")
	host        = flag.String("host", "", "host of the remote host")
	port        = flag.Int("port", 22, "port of the remote host")
	username    = flag.String("username", "", "username of the remote host")
	password    = flag.String("password", "", "password of the remote host")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}
}

func main() {
	// 解析参数
	flag.Parse()

	if *showVersion {
		version := "0.0.8"
		fmt.Println("imgx version", version)
		return
	}

	if *imageName == "" {
		fmt.Println("Image name is required")
		flag.Usage()
		return
	}

	fmt.Println("Pulling image", *imageName, "with tag", *tag, "for architecture", *arch)

	// 获取token
	token, err := pull.GetToken(*imageName)
	if err != nil {
		fmt.Println("Error getting token", err)
		return
	}

	// 获取镜像清单
	manifest, err := pull.GetManifest(token, *imageName, *tag, *arch, *osFlag)
	if err != nil {
		fmt.Println("Error getting manifest", err)
		return
	}

	// 下载镜像
	path, err := pull.DownloadImage(token, manifest, *arch, *osFlag, *imageName, *tag)
	if err != nil {
		fmt.Println("Error downloading image", err)
		return
	}
	fmt.Println("Image downloaded to", *path)

	// 删除临时文件
	defer pull.RemoveImageSaveDir(*imageName, *tag, *arch, *osFlag)

	// 打包镜像
	outputFile, err := pull.Package(*path, *imageName, *tag, *arch, *osFlag, nil)
	if err != nil {
		fmt.Println("Error packaging image", err)
		return
	}
	fmt.Println("Image packaged to", *outputFile)

	// 连接远程主机
	client, err := load.GetSSHClient(*protocol, *host, *port, *username, *password)
	if err != nil {
		fmt.Println("Error connecting to remote host", err)
		return
	}
	defer client.Close()
	// 创建临时目录
	tempDir, err := load.CreateTempDir(client)
	if err != nil {
		fmt.Println("Error creating temp dir", err)
		return
	}
	// 上传文件
	err = load.UploadFile(client, "./"+*outputFile, tempDir+"/"+*outputFile)
	if err != nil {
		fmt.Println("Error uploading file", err)
		return
	}
	// 导入镜像
	image, err := load.LoadImage(client, tempDir+"/"+*outputFile, *password)
	if err != nil {
		fmt.Println("Error loading image", err)
		return
	}
	fmt.Println(image)
}
