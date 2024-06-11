package example

import (
	"assistant_reply_message/internal/dingtalk/models"
)

func NewAiCardContent(content string) models.MessageContent {
	result := models.MessageContent{
		ContentType: models.ContentTypeAiCard,
		Content: models.CardContent{
			TemplateID: models.DingTalkAiCardTemplateId,
			CardData: map[string]string{
				"content": content,
			},
		},
	}
	return result
}
