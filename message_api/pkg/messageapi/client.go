package messageapi

import (
	"encoding/json"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"io"
	"message_api/pkg/dingapi"
	"message_api/pkg/messageapi/models"
	"net/http"
	"strings"
)

type MessageClient struct {
	dingtalkClient *dingapi.DingTalkClient
}

func NewMessageClient(client *dingapi.DingTalkClient) *MessageClient {
	return &MessageClient{dingtalkClient: client}
}

func (c *MessageClient) SendAICard(msg models.IncomingMessage, card *models.AICard) error {
	return nil
}

func (c *MessageClient) SendAICard4Webhook(msg models.IncomingMessage, card *models.AICard4Webhook) error {
	requestBody, err := json.Marshal(map[string]any{
		"contentType": "ai_card",
		"content":     card,
	})
	if err != nil {
		return nil
	}
	requestBodyStr := string(requestBody)
	requestBodyStr = `
{
    "contentType": "ai_card",
    "content": {
        "templateId": "e135583f-b9fc-4090-996d-f3119bd47f1a.schema",
        "cardData": {
            "content": "# hello world"
		}
    }
}
`
	logger.GetLogger().Infof("SendAICard4Webhook, url=%s, body=%s", msg.SessionWebhook, requestBodyStr)
	response, err := http.Post(msg.SessionWebhook, "application/json", strings.NewReader(requestBodyStr))
	if err != nil {
		return nil
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil
	}
	logger.GetLogger().Infof("SendAICard4Webhook, code=%d, body=%s", response.StatusCode, string(responseBody))
	return nil
}
