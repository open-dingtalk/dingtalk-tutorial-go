package main

import (
	"context"
	"flag"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dingtalkcard_1_0 "github.com/alibabacloud-go/dingtalk/card_1_0"
	dingtalkim_1_0 "github.com/alibabacloud-go/dingtalk/im_1_0"
	dingtalkoauth2_1_0 "github.com/alibabacloud-go/dingtalk/oauth2_1_0"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/google/uuid"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/card"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
)

type DingTalkClient struct {
	ClientID     string
	clientSecret string
	accessToken  string
	imClient     *dingtalkim_1_0.Client
	oauthClient  *dingtalkoauth2_1_0.Client
	cardClient   *dingtalkcard_1_0.Client
}

var (
	dingtalkClient *DingTalkClient = nil
)

const (
	CARD_TEMPLATE_ID   string = "7f4fe5d8-9acd-4aa6-936d-8bf922af9e8c.schema"
	CARD_TITLE                = "张三提交的财务报销"
	CARD_TYPE                 = "差旅费"
	CARD_AMOUNT               = "1000.0"
	CARD_REASON               = "出差费用"
	CARD_STATUS_TODO          = "待处理"
	CARD_STATUS_ACCEPT        = "已同意"
	CARD_STATUS_REJECT        = "被拒绝"
)

func NewDingTalkClient(clientId, clientSecret string) *DingTalkClient {
	config := &openapi.Config{}
	config.Protocol = tea.String("https")
	config.RegionId = tea.String("central")
	oauthClient, _ := dingtalkoauth2_1_0.NewClient(config)
	imClient, _ := dingtalkim_1_0.NewClient(config)
	cardClient, _ := dingtalkcard_1_0.NewClient(config)
	return &DingTalkClient{
		ClientID:     clientId,
		clientSecret: clientSecret,
		oauthClient:  oauthClient,
		imClient:     imClient,
		cardClient:   cardClient,
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

func (c *DingTalkClient) SendCard(request *dingtalkcard_1_0.CreateAndDeliverRequest) (*dingtalkcard_1_0.CreateAndDeliverResponse, error) {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}
	headers := &dingtalkcard_1_0.CreateAndDeliverHeaders{}
	headers.XAcsDingtalkAccessToken = tea.String(accessToken)

	resp, tryErr := func() (resp *dingtalkcard_1_0.CreateAndDeliverResponse, _e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		result, _err := c.cardClient.CreateAndDeliverWithOptions(request, headers, &util.RuntimeOptions{})
		if _err != nil {
			return nil, _err
		}

		return result, nil
	}()
	if tryErr != nil {
		var sdkError = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			sdkError = _t
		} else {
			sdkError.Message = tea.String(tryErr.Error())
		}
		if !tea.BoolValue(util.Empty(sdkError.Code)) && !tea.BoolValue(util.Empty(sdkError.Message)) {
			logger.GetLogger().Errorf("CreateAndDeliverWithOptions failed, clientId=%s, err=%+v", c.ClientID, sdkError)
		}
		return nil, tryErr
	}

	return resp, nil
}

func OnCardCallback(ctx context.Context, request *card.CardRequest) (*card.CardResponse, error) {
	logger.GetLogger().Infof("receive card data: %v", request)
	action := request.GetActionString("action")
	status := CARD_STATUS_TODO
	if action == "accept" {
		status = CARD_STATUS_ACCEPT
	} else if action == "reject" {
		status = CARD_STATUS_REJECT
	}
	response := &card.CardResponse{
		CardData: &card.CardDataDto{
			CardParamMap: map[string]string{
				"title":  CARD_TITLE,
				"type":   CARD_TYPE,
				"amount": CARD_AMOUNT,
				"reason": CARD_REASON,
				"status": status,
			},
		},
	}
	return response, nil
}

func OnChatBotMessageReceived(ctx context.Context, data *chatbot.BotCallbackDataModel) ([]byte, error) {
	cardData := &dingtalkcard_1_0.CreateAndDeliverRequestCardData{
		CardParamMap: make(map[string]*string),
	}
	cardData.CardParamMap["title"] = tea.String(CARD_TITLE)
	cardData.CardParamMap["type"] = tea.String(CARD_TYPE)
	cardData.CardParamMap["reason"] = tea.String(CARD_REASON)
	cardData.CardParamMap["amount"] = tea.String(CARD_AMOUNT)
	cardData.CardParamMap["status"] = tea.String(CARD_STATUS_TODO)
	imGroupOpenSpaceModel := &dingtalkcard_1_0.CreateAndDeliverRequestImGroupOpenSpaceModel{
		SupportForward: tea.Bool(true),
	}
	imGroupOpenDeliverModel := &dingtalkcard_1_0.CreateAndDeliverRequestImGroupOpenDeliverModel{
		Extension: make(map[string]*string),
		RobotCode: tea.String(dingtalkClient.ClientID),
	}
	imRobotOpenSpaceModel := &dingtalkcard_1_0.CreateAndDeliverRequestImRobotOpenSpaceModel{
		SupportForward: tea.Bool(true),
	}
	imRobotOpenDeliverModel := &dingtalkcard_1_0.CreateAndDeliverRequestImRobotOpenDeliverModel{
		Extension: make(map[string]*string),
		RobotCode: tea.String(dingtalkClient.ClientID),
		SpaceType: tea.String("IM_ROBOT"),
	}
	u, _ := uuid.NewUUID()
	outTrackId := u.String()
	var openSpaceId string = ""
	if data.ConversationType == "2" { // 群聊
		openSpaceId = fmt.Sprintf("dtv1.card//IM_GROUP.%s", data.ConversationId)
	} else {
		openSpaceId = fmt.Sprintf("dtv1.card//IM_ROBOT.%s", data.SenderStaffId)
	}
	sendCardRequest := &dingtalkcard_1_0.CreateAndDeliverRequest{
		UserIdType:     tea.Int32(1), // 1（默认）：userid模式；2：unionId模式;
		CardTemplateId: tea.String(CARD_TEMPLATE_ID),
		OutTrackId:     tea.String(outTrackId),
		CallbackType:   tea.String("STREAM"), // 采用 Stream 模式接收回调事件
		CardData:       cardData,
		OpenSpaceId:    tea.String(openSpaceId),
	}
	if data.ConversationType == "2" { // 群聊
		sendCardRequest.ImGroupOpenSpaceModel = imGroupOpenSpaceModel
		sendCardRequest.ImGroupOpenDeliverModel = imGroupOpenDeliverModel
	} else {
		sendCardRequest.ImRobotOpenSpaceModel = imRobotOpenSpaceModel
		sendCardRequest.ImRobotOpenDeliverModel = imRobotOpenDeliverModel
	}
	sendCardResponse, err := dingtalkClient.SendCard(sendCardRequest)
	if err != nil {
		return nil, err
	}
	logger.GetLogger().Infof("response=%+v", sendCardResponse)
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
	cli.RegisterCardCallbackRouter(OnCardCallback)

	err := cli.Start(context.Background())
	if err != nil {
		panic(err)
	}

	defer cli.Close()

	select {}
}
