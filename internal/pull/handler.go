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
	if flag == nil || flag.Image == "" {
		return nil, fmt.Errorf("image is required")
	}

	token, err := GetToken(flag.Image)
	if err != nil {
		fmt.Println("Error getting token", err)
		return nil, err
	}

	manifest, err := GetManifest(token, flag.Image, flag.Tag, flag.Arch, flag.OsFlag)
	if err != nil {
		fmt.Println("Error getting manifest", err)
		return nil, err
	}

	path, err := DownloadImage(token, manifest, flag.Arch, flag.OsFlag, flag.Image, flag.Tag)
	if err != nil {
		fmt.Println("Error downloading image", err)
		return nil, err
	}
	fmt.Println("Image downloaded to", *path)
	defer os.RemoveAll(*path)

	outputFile, err := Package(*path, flag.Image, flag.Tag, flag.Arch, flag.OsFlag, flag.Path, nil)
	if err != nil {
		fmt.Println("Error packaging image", err)
		return nil, err
	}
	fmt.Println("\nImage packaged to", *outputFile)
	return outputFile, nil
}
