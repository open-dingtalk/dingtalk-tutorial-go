package models

type AICard4Webhook struct {
	TemplateId string         `json:"templateId"`
	CardData   map[string]any `json:"cardData"`
}
