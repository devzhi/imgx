package pull

import (
	"encoding/json"
	"fmt"
	"github.com/devzhi/imgx/internal/util"
	"io"
	"net/http"
	"time"
)

type TokenResponse struct {
	Token       string    `json:"token"`
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	IssuedAt    time.Time `json:"issued_at"`
}

func GetToken(imageName string) (*TokenResponse, error) {
	// 构建请求
	authUrl := "https://auth.docker.io/token"
	params := map[string]string{
		"service": "registry.docker.io",
	}
	if util.IsOfficialImage(imageName) {
		params["scope"] = "repository:library/" + imageName + ":pull"
	} else {
		params["scope"] = "repository:" + imageName + ":pull"
	}
	// 发送请求
	resp, err := http.Get(authUrl + "?service=" + params["service"] + "&scope=" + params["scope"])
	if err != nil {
		fmt.Print("Error while fetching token", err)
		return nil, err
	}
	// 读取响应
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Print("Error while closing response body", err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		fmt.Print("Error while parsing token response", err)
		return nil, err
	}
	// 返回结果
	return &tokenResponse, nil
}
