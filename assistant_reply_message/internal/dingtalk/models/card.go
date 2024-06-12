package models

type CardData map[string]any

type CardContent struct {
	TemplateID string         `json:"templateId"`
	CardData   CardData       `json:"cardData,omitempty"`
	Options    map[string]any `json:"options,omitempty"`
}

type MessageContent struct {
	ContentType string      `json:"contentType"`
	Content     CardContent `json:"content"`
}
