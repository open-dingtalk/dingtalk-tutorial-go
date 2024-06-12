package client

import (
	"assistant_reply_message/internal/dingtalk/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const defaultTimeout = time.Second * 60

type Client struct {
	clientId     string
	clientSecret string
	mutex        sync.Mutex
	expireAt     int64
	AccessToken  string
}

func NewClient(clientId, clientSecret string) *Client {
	return &Client{
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func (c *Client) GetUserByAuthCode(code string) (*models.User, error) {
	getUserAccessResponse, err := c.GetUserAccessToken(code)
	if err != nil {
		return nil, err
	}
	contactUserResponse, err := c.GetUserByAccessToken(getUserAccessResponse.AccessToken)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		UnionID: contactUserResponse.UnionID,
		Avatar:  contactUserResponse.Avatar,
		Nick:    contactUserResponse.Nick,
		Name:    contactUserResponse.Nick,
	}
	return user, nil
}

func (c *Client) GetUserAccessToken(code string) (*models.GetUserAccessResponse, error) {
	request := &models.GetUserAccessRequest{
		ClientID:     c.clientId,
		ClientSecret: c.clientSecret,
		Code:         code,
		GrantType:    "authorization_code",
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	httpClient := http.Client{
		Timeout: defaultTimeout,
	}
	const url = "https://api.dingtalk.com/v1.0/oauth2/userAccessToken"
	resp, err := httpClient.Post(url, "application/json", bytes.NewReader(requestBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	response := &models.GetUserAccessResponse{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		logger.GetLogger().Errorf("dingtalk.Client, GetUserAccessToken failed, statusCode=%d, response=%+v", resp.StatusCode, response)
		return nil, errors.New("http error status")
	}
	return response, nil
}

func (c *Client) GetUserByAccessToken(userAccessToken string) (*models.GetContactUserResponse, error) {
	const queryUrl = "https://api.dingtalk.com/v1.0/contact/users/me"
	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-acs-dingtalk-access-token", userAccessToken)
	req.Header.Set("Content-Type", "application/json")
	httpClient := http.Client{Timeout: defaultTimeout}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := &models.GetContactUserResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		logger.GetLogger().Errorf("dingtalk.Client, GetUserByAccessToken failed, statusCode=%d, response=%+v",
			resp.StatusCode, response)
		return nil, errors.New("http error status")
	}
	return response, nil
}

func (c *Client) GetAccessToken() (string, error) {
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

func (c *Client) getAccessTokenFromAPI() (*models.GetTokenResponse, error) {
	// OpenAPI doc: https://open.dingtalk.com/document/orgapp/obtain-orgapp-token
	const apiUrl = "https://oapi.dingtalk.com/gettoken"
	query := url.Values{}
	query.Add("appkey", c.clientId)
	query.Add("appsecret", c.clientSecret)
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
		logger.GetLogger().Errorf("dingtalk.Client, getAccessTokenFromAPI failed, statusCode=%d, response=%+v",
			res.StatusCode, response)
		return nil, errors.New(response.ErrorMessage)
	}
	return response, nil
}

func (c *Client) SendAiCardStream(unionId, templateId, key, value string) error {
	accessToken, err := c.GetAccessToken()

	cardContent := models.CardContent{
		TemplateID: templateId,
	}
	cardContentBytes, err := json.Marshal(cardContent)
	if err != nil {
		return err
	}
	prepareRequest := map[string]any{
		"unionId":     unionId,
		"contentType": "ai_card",
		"content":     string(cardContentBytes),
	}
	requestBodyBytes, err := json.Marshal(prepareRequest)
	httpClient := http.Client{Timeout: defaultTimeout}
	req, err := http.NewRequest("POST", "https://api.dingtalk.com/v1.0/aiInteraction/prepare", bytes.NewReader(requestBodyBytes))
	req.Header.Set("x-acs-dingtalk-access-token", accessToken)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("%s", string(responseBody))
	return nil
}
