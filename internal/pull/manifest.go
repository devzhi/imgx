package pull

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/devzhi/imgx/internal/util"
)

type ManifestsResp struct {
	Manifests     []Manifest  `json:"manifests"`
	MediaType     string      `json:"mediaType"`
	SchemaVersion int         `json:"schemaVersion"`
	Config        Config      `json:"config"`
	Layers        []Layer     `json:"layers"`
	Annotations   Annotations `json:"annotations"`
}

type Manifest struct {
	Annotations Annotations `json:"annotations"`
	Digest      string      `json:"digest"`
	MediaType   string      `json:"mediaType"`
	Platform    Platform    `json:"platform"`
	Size        int         `json:"size"`
}

type Annotations struct {
	ComDockerOfficialImagesBashbrewArch string `json:"com.docker.official-images.bashbrew.arch"`
	OrgOpencontainersImageBaseName      string `json:"org.opencontainers.image.base.name"`
	OrgOpencontainersImageCreated       string `json:"org.opencontainers.image.created"`
	OrgOpencontainersImageRevision      string `json:"org.opencontainers.image.revision"`
	OrgOpencontainersImageSource        string `json:"org.opencontainers.image.source"`
	OrgOpencontainersImageUrl           string `json:"org.opencontainers.image.url"`
	OrgOpencontainersImageVersion       string `json:"org.opencontainers.image.version"`
	VndDockerReferenceDigest            string `json:"vnd.docker.reference.digest"`
	VndDockerReferenceType              string `json:"vnd.docker.reference.type"`
}

type Platform struct {
	Architecture string `json:"architecture"`
	Os           string `json:"os"`
	Variant      string `json:"variant"`
}

type Config struct {
	MediaType string `json:"mediaType"`
	Digest    string `json:"digest"`
	Size      int    `json:"size"`
}

type Layer struct {
	MediaType string `json:"mediaType"`
	Digest    string `json:"digest"`
	Size      int    `json:"size"`
}

func GetManifest(token *TokenResponse, imageName, tag string, arch string, os string) (*ManifestsResp, error) {
	registryURL := "https://registry.hub.docker.com/v2/"
	if util.IsOfficialImage(imageName) {
		registryURL = registryURL + "library/" + imageName
	} else {
		registryURL = registryURL + imageName
	}
	req, err := http.NewRequest("GET", registryURL+"/manifests/"+tag, nil)
	if err != nil {
		fmt.Print("Error while creating manifest request", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.Token)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.list.v2+json,application/vnd.docker.distribution.manifest.v2+json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print("Error while fetching manifest", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch manifest: %s: %s", res.Status, strings.TrimSpace(string(body)))
	}

	var resp ManifestsResp
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Print("Error while parsing manifest response ", err)
		return nil, err
	}
	if isManifestList(resp.MediaType) {
		for _, manifest := range resp.Manifests {
			if manifest.Platform.Architecture == arch && manifest.Platform.Os == os {
				return GetManifestByDigest(token, imageName, manifest.Digest)
			}
		}
		fmt.Printf("No matching manifest found for %s/%s\n", arch, os)
		return nil, fmt.Errorf("no matching manifest found for %s/%s", arch, os)
	}
	return &resp, nil
}

func GetManifestByDigest(token *TokenResponse, imageName, digest string) (*ManifestsResp, error) {
	registryURL := "https://registry.hub.docker.com/v2/"
	if util.IsOfficialImage(imageName) {
		registryURL = registryURL + "library/" + imageName + "/manifests/" + digest
	} else {
		registryURL = registryURL + imageName + "/manifests/" + digest
	}
	req, err := http.NewRequest("GET", registryURL, nil)
	if err != nil {
		fmt.Print("Error while creating manifest request", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.Token)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print("Error while fetching manifest", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch manifest by digest: %s: %s", res.Status, strings.TrimSpace(string(body)))
	}

	var resp ManifestsResp
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Print("Error while parsing manifest response", err)
		return nil, err
	}
	return &resp, nil
}

func isManifestList(mediaType string) bool {
	return mediaType == "application/vnd.oci.image.index.v1+json" || mediaType == "application/vnd.docker.distribution.manifest.list.v2+json"
}
