package action

import (
	"assistant_reply_message/internal/dingtalk/client"
	"assistant_reply_message/internal/dingtalk/example"
	"assistant_reply_message/internal/dingtalk/models"
	"assistant_reply_message/internal/dingtalk/webhook"
	"assistant_reply_message/internal/utils"
	"context"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/plugin"
	"net/http"
)

type CallbackHandler struct {
	ClientID       string
	ClientSecret   string
	dingtalkClient *client.Client
}

func NewCallbackHandler(client_id, client_secret string) *CallbackHandler {
	return &CallbackHandler{
		ClientID:       client_id,
		ClientSecret:   client_secret,
		dingtalkClient: client.NewClient(client_id, client_secret),
	}
}

func (h *CallbackHandler) OnPluginCallback(c context.Context, req *plugin.GraphRequest) (*plugin.GraphResponse, error) {
	exampleRequest, err := example.NewExampleRequest(req.Body)
	if err != nil {
		return nil, err
	}

	caseIndex := 0

	replyMessage := func(tag, sessionWebhook string, message any) {
		caseIndex += 1
		resp, err := webhook.ReplyMessage(sessionWebhook, message)
		if err != nil || resp.StatusCode != http.StatusOK {
			logger.GetLogger().Errorf("❌ [%d. %s] webhook.ReplyMessage failed, err=%s, resp=%+v", caseIndex, tag, err, resp)
		} else {
			logger.GetLogger().Infof("✅ [%d. %s] webhook.ReplyMessage success", caseIndex, tag)
		}
	}

	if false {
		msg1 := example.NewAiCardContent(models.DingTalkAiCardNonStream, "白日依山尽，黄河入海流。欲穷千里目，更上一层楼。")
		replyMessage("webhook, ai card", exampleRequest.CurrentConversation.SessionWebhook, msg1)
		msg2 := example.NewAiCardStreamContent(models.DingTalkAiCardStream, "content", "白日依山尽，黄河入海流。欲穷千里目，更上一层楼。", true)
		replyMessage("webhook, stream ai card", exampleRequest.CurrentConversation.SessionWebhook, msg2)
	}
	h.dingtalkClient.SendAiCardStream(exampleRequest.CurrentUser.UnionId, models.DingTalkAiCardStream, "content", "白日依山尽，黄河入海流。欲穷千里目，更上一层楼。白日依山尽，黄河入海流。欲穷千里目，更上一层楼。白日依山尽，黄河入海流。欲穷千里目，更上一层楼。白日依山尽，黄河入海流。欲穷千里目，更上一层楼。")

	weather := utils.GetWeather()
	return utils.CreateSuccessResponse(weather), nil
}
