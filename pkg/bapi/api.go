package bapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// APP_KEY APP_SECRET
	USERNAME = "linda"
	PASSWORD = "123456"
)

type AccessToken struct {
	Token string `json:"token"`
}

type API struct {
	URL string
}

func NewAPI(url string) *API {
	return &API{URL: url}
}

// httpGet 封装基本的 get 网络请求，ctx备用
func (a *API) httpGet(ctx context.Context, path string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", a.URL, path))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body, nil
}

// getAccessToken 基于httpGet封装业务通用接口 token
func (a *API) getAccessToken(ctx context.Context) (string, error) {
	url := fmt.Sprintf(
		"%s?username=%s&password=%s",
		"auth",
		USERNAME,
		PASSWORD)

	body, err := a.httpGet(ctx, url)
	if err != nil {
		return "", err
	}

	var accessToken AccessToken
	_ = json.Unmarshal(body, &accessToken)
	return accessToken.Token, nil
}

// GetTagList 业务接口获取Tag列表
func (a *API) GetTagList(ctx context.Context) ([]byte, error) {
	token, err := a.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	body, err := a.httpGet(ctx, fmt.Sprintf("%s?token=%s", "api/v1/tags", token))
	if err != nil {
		return nil, err
	}

	return body, nil
}
