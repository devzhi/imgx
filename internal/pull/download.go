package pull

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/devzhi/imgx/internal/util"
	"github.com/schollz/progressbar/v3"
)

// DownloadImage 根据提供的manifests下载镜像层和配置
func DownloadImage(token *TokenResponse, manifests *ManifestsResp, arch string, operateSystem string, imageName string, tag string) (path *string, err error) {
	registryURL := "https://registry.hub.docker.com/v2"
	if manifests.MediaType != "application/vnd.oci.image.manifest.v1+json" && manifests.MediaType != "application/vnd.docker.distribution.manifest.v2+json" {
		return nil, fmt.Errorf("unsupported manifest type: %s", manifests.MediaType)
	}

	temp, err := os.MkdirTemp("", "tmp-")
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(temp, strings.ReplaceAll(imageName, "/", "_")+"_"+tag+"_"+arch+"_"+operateSystem)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Println("Error while creating directory", err)
		return nil, err
	}

	layers := make([]Layer, len(manifests.Layers))
	errCh := make(chan error, len(manifests.Layers))
	var wg sync.WaitGroup
	wg.Add(len(manifests.Layers))
	for i, layer := range manifests.Layers {
		go func(index int, currentLayer Layer) {
			defer wg.Done()
			if downloadErr := downloadLayer(token, registryURL, imageName, currentLayer, dir); downloadErr != nil {
				errCh <- downloadErr
				return
			}
			layers[index] = currentLayer
		}(i, layer)
	}
	wg.Wait()
	close(errCh)
	for downloadErr := range errCh {
		if downloadErr != nil {
			_ = os.RemoveAll(temp)
			return nil, downloadErr
		}
	}

	err = downloadConfig(token, registryURL, manifests.Config, dir, imageName)
	if err != nil {
		return nil, err
	}
	if err := createManifestFile(manifests.Config, layers, imageName, tag, dir); err != nil {
		return nil, err
	}
	return &dir, nil
}

// downloadLayer 下载单个层并显示进度条
func downloadLayer(token *TokenResponse, registryURL, imageName string, layer Layer, dir string) error {
	requestURL := registryURL
	if util.IsOfficialImage(imageName) {
		requestURL = requestURL + "/library/" + imageName + "/blobs/" + layer.Digest
	} else {
		requestURL = requestURL + "/" + imageName + "/blobs/" + layer.Digest
	}
	req, err := http.NewRequest("GET", requestURL, nil)
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
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("fetch layer %s: %s: %s", layer.Digest, res.Status, strings.TrimSpace(string(body)))
	}

	file, err := os.Create(dir + "/" + strings.Split(layer.Digest, ":")[1] + ".tar")
	if err != nil {
		fmt.Println("Error while creating layer file", err)
		return err
	}
	defer file.Close()

	bar := progressbar.DefaultBytes(
		res.ContentLength,
		"Downloading layer "+layer.Digest[7:19],
	)

	_, err = io.Copy(io.MultiWriter(file, bar), res.Body)
	if err != nil {
		fmt.Println("Error while writing layer file", err)
		return err
	}
	return nil
}

// downloadConfig 下载配置文件
func downloadConfig(token *TokenResponse, registryURL string, config Config, dir string, imageName string) error {
	requestURL := registryURL
	if util.IsOfficialImage(imageName) {
		requestURL = requestURL + "/library/" + imageName + "/blobs/" + config.Digest
	} else {
		requestURL = requestURL + "/" + imageName + "/blobs/" + config.Digest
	}
	req, err := http.NewRequest("GET", requestURL, nil)
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
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("fetch config %s: %s: %s", config.Digest, res.Status, strings.TrimSpace(string(body)))
	}

	file, err := os.Create(dir + "/" + strings.Split(config.Digest, ":")[1] + ".json")
	if err != nil {
		fmt.Println("Error while creating config file", err)
		return err
	}
	defer file.Close()
	bar := progressbar.DefaultBytes(
		res.ContentLength,
		"Downloading config "+config.Digest[7:19],
	)

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

	_, err = manifestFile.WriteString("[{\"Config\":\"" + strings.Split(config.Digest, ":")[1] + ".json\",\"RepoTags\":[\"" + imageName + ":" + tag + "\"],\"Layers\":[")
	if err != nil {
		fmt.Println("Error while writing manifest file", err)
		return err
	}

	for i, layer := range layers {
		_, err = manifestFile.WriteString("\"" + strings.Split(layer.Digest, ":")[1] + ".tar\"")
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
	_, err = manifestFile.WriteString("]}]")
	if err != nil {
		fmt.Println("Error while writing manifest file", err)
		return err
	}
	return nil
}

func RemoveImageSaveDir(imageName, tag, arch, operateSystem string) {
	dir := "./" + strings.ReplaceAll(imageName, "/", "_") + "_" + tag + "_" + arch + "_" + operateSystem
	if err := os.RemoveAll(dir); err != nil {
		fmt.Println("Error while removing directory", err)
	}
}
