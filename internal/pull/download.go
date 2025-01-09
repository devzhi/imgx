package pull

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// DownloadImage 根据提供的manifests下载镜像层和配置
func DownloadImage(token *TokenResponse, manifests *ManifestsResp, arch string, operateSystem string, imageName string, tag string) (path *string, err error) {
	registryUrl := "https://registry.hub.docker.com/v2"
	if manifests.MediaType != "application/vnd.oci.image.manifest.v1+json" {
		fmt.Println("Unsupported manifest type", manifests.MediaType)
		return nil, err
	}

	// 创建镜像目录
	dir := "./" + imageName + "_" + tag + "_" + arch + "_" + operateSystem
	err = os.Mkdir(dir, 0755)
	if err != nil {
		fmt.Println("Error while creating directory", err)
		return nil, err
	}

	var layers []Layer
	for _, layer := range manifests.Layers {
		// 下载每一层
		err := downloadLayer(token, registryUrl, imageName, layer, dir)
		if err != nil {
			return nil, err
		}
		layers = append(layers, layer)
	}

	// 下载配置文件
	err = downloadConfig(token, registryUrl, manifests.Config, dir, imageName)
	if err != nil {
		return nil, err
	}
	// 创建manifest文件
	err = createManifestFile(manifests.Config, layers, imageName, tag, dir)
	if err != nil {
		return nil, err
	}
	return &dir, nil
}

// downloadLayer 下载单个层并显示进度条
func downloadLayer(token *TokenResponse, registryUrl, imageName string, layer Layer, dir string) error {
	req, err := http.NewRequest("GET", registryUrl+"/library/"+imageName+"/blobs/"+layer.Digest, nil)
	if err != nil {
		fmt.Println("Error while creating layer request", err)
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token.Token)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.list.v2+json,application/vnd.docker.distribution.manifest.v2+json")
	fmt.Println()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error while fetching layer", err)
		return err
	}
	defer res.Body.Close()
	// 创建文件保存层
	file, err := os.Create(dir + "/" + strings.Split(layer.Digest, ":")[1] + ".tar")
	if err != nil {
		fmt.Println("Error while creating layer file", err)
		return err
	}
	defer file.Close()

	// 创建进度条
	bar := progressbar.DefaultBytes(
		res.ContentLength,
		"Downloading layer "+string(layer.Digest[7:19]),
	)

	// 使用进度条复制响应体到文件
	_, err = io.Copy(io.MultiWriter(file, bar), res.Body)
	if err != nil {
		fmt.Println("Error while writing layer file", err)
		return err
	}
	return nil
}

// downloadConfig 下载配置文件
func downloadConfig(token *TokenResponse, registryUrl string, config Config, dir string, imageName string) error {
	req, err := http.NewRequest("GET", registryUrl+"/library/"+imageName+"/blobs/"+config.Digest, nil)
	if err != nil {
		fmt.Println("Error while creating config request", err)
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token.Token)
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.list.v2+json,application/vnd.docker.distribution.manifest.v2+json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error while fetching config", err)
		return err
	}
	defer res.Body.Close()

	// 创建文件保存配置
	file, err := os.Create(dir + "/" + strings.Split(config.Digest, ":")[1] + ".json")
	if err != nil {
		fmt.Println("Error while creating config file", err)
		return err
	}
	defer file.Close()
	// 创建进度条
	bar := progressbar.DefaultBytes(
		res.ContentLength,
		"Downloading config "+string(config.Digest[7:19]),
	)

	// 使用进度条复制响应体到文件
	_, err = io.Copy(io.MultiWriter(file, bar), res.Body)
	if err != nil {
		fmt.Println("Error while writing config file", err)
		return err
	}
	return nil
}

// createManifestFile 创建下载镜像的manifest文件
func createManifestFile(config Config, layers []Layer, imageName, tag, dir string) error {
	manifestFile, err := os.Create(dir + "/manifest.json")
	if err != nil {
		fmt.Println("Error while creating manifest file", err)
		return err
	}
	defer manifestFile.Close()

	_, err = manifestFile.WriteString("{\"Config\":\"" + config.Digest + ".json\",\"RepoTags\":[\"" + imageName + ":" + tag + "\"],\"Layers\":[")
	if err != nil {
		fmt.Println("Error while writing manifest file", err)
		return err
	}

	for i, layer := range layers {
		_, err = manifestFile.WriteString("\"" + layer.Digest + ".tar\"")
		if err != nil {
			fmt.Println("Error while writing manifest file", err)
			return err
		}
		if i < len(layers)-1 {
			_, err = manifestFile.WriteString(",")
			if err != nil {
				fmt.Println("Error while writing manifest file", err)
				return err
			}
		}
	}
	_, err = manifestFile.WriteString("]}")
	if err != nil {
		fmt.Println("Error while writing manifest file", err)
		return err
	}
	return nil
}

func RemoveImageSaveDir(imageName, tag, arch, operateSystem string) {
	dir := "./" + imageName + "_" + tag + "_" + arch + "_" + operateSystem
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Println("Error while removing directory", err)
	}
}
