package action

import (
	"assistant_reply_message/internal/dingtalk/example"
	"assistant_reply_message/internal/dingtalk/models"
	"assistant_reply_message/internal/dingtalk/webhook"
	"assistant_reply_message/internal/utils"
	"context"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/plugin"
)

type CallbackHandler struct {
	ClientID     string
	ClientSecret string
}

func NewCallbackHandler(client_id, client_secret string) *CallbackHandler {
	return &CallbackHandler{
		ClientID:     client_id,
		ClientSecret: client_secret,
	}
}

func (h *CallbackHandler) OnPluginCallback(c context.Context, req *plugin.GraphRequest) (*plugin.GraphResponse, error) {
	logger.GetLogger().Infof("request body=%s", req.Body)
	exampleRequest, err := example.NewExampleRequest(req.Body)
	if err != nil {
		return nil, err
	}

	msg1 := models.MessageContent{
		ContentType: models.ContentTypeAiCard,
		Content: models.CardContent{
			TemplateID: models.DingTalkAiCardNonStream,
			CardData: models.CardData{
				"content": "白日依山尽，黄河入海流。欲穷千里目，更上一层楼。",
			},
		},
	}
	if resp, err := webhook.ReplyMessage(exampleRequest.CurrentConversation.SessionWebhook, msg1); err != nil {
		logger.GetLogger().Errorf("webhook.ReplyMessage failed, err=%s, resp=%+v", err, resp)
		return nil, err
	}

	logger.GetLogger().Infof("exampleRequest=%+v", exampleRequest)
	weather := utils.GetWeather()
	return utils.CreateSuccessResponse(weather), nil
}
