package pull

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/devzhi/imgx/internal/util"
)

type TokenResponse struct {
	Token       string    `json:"token"`
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	IssuedAt    time.Time `json:"issued_at"`
}

func GetToken(imageName string) (*TokenResponse, error) {
	authURL := "https://auth.docker.io/token"
	params := map[string]string{
		"service": "registry.docker.io",
	}
	if util.IsOfficialImage(imageName) {
		params["scope"] = "repository:library/" + imageName + ":pull"
	} else {
		params["scope"] = "repository:" + imageName + ":pull"
	}
	resp, err := http.Get(authURL + "?service=" + params["service"] + "&scope=" + params["scope"])
	if err != nil {
		fmt.Print("Error while fetching token", err)
		return nil, err
	}
	defer func(body io.ReadCloser) {
		if closeErr := body.Close(); closeErr != nil {
			fmt.Print("Error while closing response body", closeErr)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch token: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		fmt.Print("Error while parsing token response", err)
		return nil, err
	}
	return &tokenResponse, nil
}
