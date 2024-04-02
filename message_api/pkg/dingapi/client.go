package dingapi

import (
	"encoding/json"
	"errors"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"io"
	"message_api/pkg/dingapi/models"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const defaultTimeout = time.Second * 60

type DingTalkClient struct {
	CorpId       string
	ClientId     string
	ClientSecret string

	mutex       sync.Mutex
	expireAt    int64
	AccessToken string
}

func NewDingTalkClient(corpId, clientId, clientSecret string) *DingTalkClient {
	return &DingTalkClient{
		CorpId:       corpId,
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}

func (c *DingTalkClient) GetAccessToken() (string, error) {
	accessToken := ""
	{
		// 先查询缓存
		c.mutex.Lock()
		now := time.Now().Unix()
		if c.expireAt > 0 && c.AccessToken != "" && (now+60) < c.expireAt {
			// 预留一分钟有效期避免在Token过期的临界点调用接口出现401错误
			accessToken = c.AccessToken
		}
		c.mutex.Unlock()
	}
	if accessToken != "" {
		return accessToken, nil
	}

	tokenResult, err := c.getAccessTokenFromAPI()
	if err != nil {
		return "", err
	}

	{
		// 更新缓存
		c.mutex.Lock()
		c.AccessToken = tokenResult.AccessToken
		c.expireAt = time.Now().Unix() + int64(tokenResult.ExpiresIn)
		c.mutex.Unlock()
	}
	return tokenResult.AccessToken, nil
}

func (c *DingTalkClient) getAccessTokenFromAPI() (*models.GetTokenResponse, error) {
	// OpenAPI doc: https://open.dingtalk.com/document/orgapp/obtain-orgapp-token
	const apiUrl = "https://oapi.dingtalk.com/gettoken"
	query := url.Values{}
	query.Add("appkey", c.ClientId)
	query.Add("appsecret", c.ClientSecret)
	fullUrl := apiUrl + "?" + query.Encode()

	// Send the HTTP request and parse the response body as JSON
	httpClient := http.Client{Timeout: defaultTimeout}
	res, err := httpClient.Get(fullUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	response := &models.GetTokenResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK || response.ErrorCode != 0 {
		logger.GetLogger().Errorf("dingtalk.Client, getAccessTokenFromAPI failed, statucCode=%d", res.StatusCode)
		return nil, errors.New(response.ErrorMessage)
	}
	return response, nil
}
