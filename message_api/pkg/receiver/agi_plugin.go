package receiver

import (
	"context"
	"encoding/json"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/plugin"
	"message_api/pkg/dingapi"
	"message_api/pkg/messageapi"
	models2 "message_api/pkg/messageapi/models"
	"message_api/pkg/receiver/models"
	"net/http"
)

type AgiPlugin struct {
	CorpId         string
	ClientId       string
	ClientSecret   string
	dingTalkClient *dingapi.DingTalkClient
	messageClient  *messageapi.MessageClient
}

func NewAgiPlugin(corpId, clientId, clientSecret string) *AgiPlugin {
	dingtalkClient := dingapi.NewDingTalkClient(corpId, clientId, clientSecret)
	messageClient := messageapi.NewMessageClient(dingtalkClient)
	return &AgiPlugin{
		CorpId:         corpId,
		ClientId:       clientId,
		ClientSecret:   clientSecret,
		dingTalkClient: dingtalkClient,
		messageClient:  messageClient,
	}
}

func (p *AgiPlugin) OnIncomingRequest(c context.Context, request *plugin.GraphRequest) (*plugin.GraphResponse, error) {
	var pluginRequest models.PluginRequest
	if err := json.Unmarshal([]byte(request.Body), &pluginRequest); err != nil {
		return nil, err
	}

	incomingMessage := models2.IncomingMessage{
		SenderUnionId:  pluginRequest.SenderUnionId,
		SessionWebhook: pluginRequest.SessionWebhook,
	}

	card_data := make(map[string]any, 0)
	card_data["content"] = "# hello world"
	card_data_bytes, err := json.Marshal(card_data)
	if err != nil {
		return nil, err
	}
	aicard := &models2.AICard{
		TemplateId: "e135583f-b9fc-4090-996d-f3119bd47f1a.schema",
		CardData:   string(card_data_bytes),
	}

	aicard4webhook := &models2.AICard4Webhook{
		TemplateId: "e135583f-b9fc-4090-996d-f3119bd47f1a.schema",
		CardData: map[string]any{
			"content": "# hello world",
		},
	}
	err = p.messageClient.SendAICard(incomingMessage, aicard)
	if err != nil {
		return nil, err
	}
	err = p.messageClient.SendAICard4Webhook(incomingMessage, aicard4webhook)
	if err != nil {
		return nil, err
	}

	return &plugin.GraphResponse{
		StatusLine: plugin.GraphStatusLine{
			Code: http.StatusOK,
		},
	}, nil
}
