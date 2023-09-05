package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dingtalkim_1_0 "github.com/alibabacloud-go/dingtalk/im_1_0"
	dingtalkoauth2_1_0 "github.com/alibabacloud-go/dingtalk/oauth2_1_0"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/google/uuid"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"strings"
	"time"
)

// doc: https://open-dingtalk.github.io/developerpedia/docs/learn/card/intro
const messageCardTemplate = `
{
  "config": {
    "autoLayout": true,
    "enableForward": true
  },
  "header": {
    "title": {
      "type": "text",
      "text": "打字机模式"
    },
    "logo": "@lALPDfJ6V_FPDmvNAfTNAfQ"
  },
  "contents": [
    {
      "type": "text",
      "text": "%s",
      "id": "text_1693929551595"
    },
    {
      "type": "divider",
      "id": "divider_1693929551595"
    },
    {
      "type": "markdown",
      "text": "%s",
      "id": "markdown_1693929674245"
    }
  ]
}
`

type DingTalkClient struct {
	ClientID     string
	clientSecret string
	accessToken  string
	imClient     *dingtalkim_1_0.Client
	oauthClient  *dingtalkoauth2_1_0.Client
}

var (
	dingtalkClient *DingTalkClient = nil
)

func NewDingTalkClient(clientId, clientSecret string) *DingTalkClient {
	config := &openapi.Config{}
	config.Protocol = tea.String("https")
	config.RegionId = tea.String("central")
	imClient, _ := dingtalkim_1_0.NewClient(config)
	oauthClient, _ := dingtalkoauth2_1_0.NewClient(config)
	return &DingTalkClient{
		ClientID:     clientId,
		clientSecret: clientSecret,
		imClient:     imClient,
		oauthClient:  oauthClient,
	}
}

func (c *DingTalkClient) GetAccessToken() (string, error) {
	request := &dingtalkoauth2_1_0.GetAccessTokenRequest{
		AppKey:    tea.String(c.ClientID),
		AppSecret: tea.String(c.clientSecret),
	}
	response, tryErr := func() (_resp *dingtalkoauth2_1_0.GetAccessTokenResponse, _e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_resp, _err := c.oauthClient.GetAccessToken(request)
		if _err != nil {
			return nil, _err
		}

		return _resp, nil
	}()
	if tryErr != nil {
		return "", tryErr
	}
	return *response.Body.AccessToken, nil
}

func (c *DingTalkClient) SendInteractiveCard(request *dingtalkim_1_0.SendRobotInteractiveCardRequest) (*dingtalkim_1_0.SendRobotInteractiveCardResponse, error) {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}

	headers := &dingtalkim_1_0.SendRobotInteractiveCardHeaders{
		XAcsDingtalkAccessToken: tea.String(accessToken),
	}
	response, tryErr := func() (_resp *dingtalkim_1_0.SendRobotInteractiveCardResponse, _e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_resp, _e = c.imClient.SendRobotInteractiveCardWithOptions(request, headers, &util.RuntimeOptions{})
		if _e != nil {
			return
		}
		return
	}()
	if tryErr != nil {
		return nil, tryErr
	}
	return response, nil
}

func (c *DingTalkClient) UpdateInteractiveCard(request *dingtalkim_1_0.UpdateRobotInteractiveCardRequest) (*dingtalkim_1_0.UpdateRobotInteractiveCardResponse, error) {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}

	headers := &dingtalkim_1_0.UpdateRobotInteractiveCardHeaders{
		XAcsDingtalkAccessToken: tea.String(accessToken),
	}
	response, tryErr := func() (_resp *dingtalkim_1_0.UpdateRobotInteractiveCardResponse, _e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_resp, _e = c.imClient.UpdateRobotInteractiveCardWithOptions(request, headers, &util.RuntimeOptions{})
		if _e != nil {
			return
		}
		return
	}()
	if tryErr != nil {
		return nil, tryErr
	}
	return response, nil
}

func OnChatBotMessageReceived(ctx context.Context, data *chatbot.BotCallbackDataModel) ([]byte, error) {
	// create an uniq card id to identify a card instance while updating
	// see: https://open.dingtalk.com/document/orgapp/robots-send-interactive-cards (cardBizId)
	u, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	cardInstanceId := u.String()

	// send interactive card; 发送交互式卡片
	cardData := fmt.Sprintf(messageCardTemplate, "", "")
	sendOptions := &dingtalkim_1_0.SendRobotInteractiveCardRequestSendOptions{}
	request := &dingtalkim_1_0.SendRobotInteractiveCardRequest{
		CardTemplateId: tea.String("StandardCard"),
		CardBizId:      tea.String(cardInstanceId),
		CardData:       tea.String(cardData),
		RobotCode:      tea.String(dingtalkClient.ClientID),
		SendOptions:    sendOptions,
		PullStrategy:   tea.Bool(false),
	}
	if data.ConversationType == "2" {
		// group chat; 群聊
		request.SetOpenConversationId(data.ConversationId)
	} else {
		// ConversationType == "1": private chat; 单聊
		receiverBytes, err := json.Marshal(map[string]string{"userId": data.SenderStaffId})
		if err != nil {
			return nil, err
		}
		request.SetSingleChatReceiver(string(receiverBytes))
	}
	_, err = dingtalkClient.SendInteractiveCard(request)
	if err != nil {
		return nil, err
	}

	// 持续更新交互式卡片
	fullTitle := []string{"登", "鹳", "雀", "楼"}
	fullContent := []string{"* 白", "日", "依", "山", "尽，", "\n* 黄", "河", "入", "海", "流", "。", "\n* 欲", "穷", "千", "里", "目，", "\n* 更", "上", "一", "层", "楼。"}
	fmt.Println(len(fullTitle))
	for i := 1; i <= len(fullContent); i++ {
		if i > 1 {
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
		title := strings.Join(fullTitle, "")
		if i <= len(fullTitle) {
			title = strings.Join(fullTitle[0:i], "")
		}
		content := strings.Join(fullContent[0:i], "")

		updateRequest := &dingtalkim_1_0.UpdateRobotInteractiveCardRequest{
			CardBizId: tea.String(cardInstanceId),
			CardData:  tea.String(fmt.Sprintf(messageCardTemplate, title, content)),
		}
		_, err = dingtalkClient.UpdateInteractiveCard(updateRequest)
		if err != nil {
			return nil, err
		}
	}

	return []byte(""), nil
}

func main() {
	var clientId, clientSecret string
	flag.StringVar(&clientId, "client_id", "", "your-client-id")
	flag.StringVar(&clientSecret, "client_secret", "", "your-client-secret")
	flag.Parse()
	if len(clientId) == 0 || len(clientSecret) == 0 {
		panic("command line options --client_id and --client_secret required")
	}

	logger.SetLogger(logger.NewStdTestLogger())

	dingtalkClient = NewDingTalkClient(clientId, clientSecret)

	cli := client.NewStreamClient(client.WithAppCredential(client.NewAppCredentialConfig(clientId, clientSecret)))
	cli.RegisterChatBotCallbackRouter(OnChatBotMessageReceived)

	err := cli.Start(context.Background())
	if err != nil {
		panic(err)
	}

	defer cli.Close()

	select {}
}
