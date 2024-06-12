package example

import (
	"assistant_reply_message/internal/dingtalk/models"
)

func NewAiCardContent(templateId, content string) models.MessageContent {
	result := models.MessageContent{
		ContentType: models.ContentTypeAiCard,
		Content: models.CardContent{
			TemplateID: templateId,
			CardData: map[string]any{
				"content": content,
			},
		},
	}
	return result
}

func NewAiCardStreamContent(templateId, key, value string, isFinalize bool) models.MessageContent {
	result := models.MessageContent{
		ContentType: models.ContentTypeAiCard,
		Content: models.CardContent{
			TemplateID: templateId,
			CardData: map[string]any{
				"key":        key,
				"value":      value,
				"isFinalize": isFinalize,
			},
			Options: map[string]any{
				"componentTag": "streamingComponent",
			},
		},
	}
	return result
}
