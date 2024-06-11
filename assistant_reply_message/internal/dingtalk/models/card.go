package models

type CardData map[string]string

type CardContent struct {
	TemplateID string   `json:"templateId"`
	CardData   CardData `json:"cardData"`
}

type MessageContent struct {
	ContentType string      `json:"contentType"`
	Content     CardContent `json:"content"`
}
