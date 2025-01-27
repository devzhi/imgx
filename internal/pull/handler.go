package pull

import (
	"fmt"
	"os"
)

// Flag pull命令参数
type Flag struct {
	Image  string
	Tag    string
	Arch   string
	OsFlag string
	Path   string
}

// Execute 执行pull命令
func Execute(flag *Flag) (*string, error) {
	// 执行pull命令
	// 获取token
	token, err := GetToken(flag.Image)
	if err != nil {
		fmt.Println("Error getting token", err)
		return nil, err
	}

	// 获取镜像清单
	manifest, err := GetManifest(token, flag.Image, flag.Tag, flag.Arch, flag.OsFlag)
	if err != nil {
		fmt.Println("Error getting manifest", err)
		return nil, err
	}

	// 下载镜像
	path, err := DownloadImage(token, manifest, flag.Arch, flag.OsFlag, flag.Image, flag.Tag)
	if err != nil {
		fmt.Println("Error downloading image", err)
		return nil, err
	}
	fmt.Println("Image downloaded to", *path)

	// 删除临时文件
	defer os.RemoveAll(*path)

	// 打包镜像
	outputFile, err := Package(*path, flag.Image, flag.Tag, flag.Arch, flag.OsFlag, flag.Path, nil)
	if err != nil {
		fmt.Println("Error packaging image", err)
		return nil, err
	}
	fmt.Println("\nImage packaged to", *outputFile)
	// 返回打包的镜像
	return outputFile, nil
}
