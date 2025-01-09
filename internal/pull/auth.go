package pull

import (
	"encoding/json"
	"fmt"
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
		"scope":   "repository:library/" + imageName + ":pull",
	}
	// 发送请求
	resp, err := http.Get(authUrl + "?service=" + params["service"] + "&scope=" + params["scope"])
	if err != nil {
		fmt.Print("Error while fetching token", err)
		return nil, err
	}
	// 读取响应
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		fmt.Print("JSON unmarshal error", err)
		return nil, err
	}
	// 返回结果
	return &tokenResponse, nil
}
