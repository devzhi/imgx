package pull

import (
	"fmt"
)

// Flag pull命令参数
type Flag struct {
	Image  string
	Tag    string
	Arch   string
	OsFlag string
}

// Execute 执行pull命令
func Execute(command *Flag) (*string, error) {
	// 执行pull命令
	// 获取token
	token, err := GetToken(command.Image)
	if err != nil {
		fmt.Println("Error getting token", err)
		return nil, err
	}

	// 获取镜像清单
	manifest, err := GetManifest(token, command.Image, command.Tag, command.Arch, command.OsFlag)
	if err != nil {
		fmt.Println("Error getting manifest", err)
		return nil, err
	}

	// 下载镜像
	path, err := DownloadImage(token, manifest, command.Arch, command.OsFlag, command.Image, command.Tag)
	if err != nil {
		fmt.Println("Error downloading image", err)
		return nil, err
	}
	fmt.Println("Image downloaded to", *path)

	// 删除临时文件
	defer RemoveImageSaveDir(command.Image, command.Tag, command.Arch, command.OsFlag)

	// 打包镜像
	outputFile, err := Package(*path, command.Image, command.Tag, command.Arch, command.OsFlag, nil)
	if err != nil {
		fmt.Println("Error packaging image", err)
		return nil, err
	}
	fmt.Println("\nImage packaged to", *outputFile)
	// 返回打包的镜像
	return outputFile, nil
}
