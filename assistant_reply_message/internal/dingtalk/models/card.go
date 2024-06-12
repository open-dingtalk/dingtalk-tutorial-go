package models

type CardData map[string]any
type CardOptions map[string]any

type CardContent struct {
	TemplateID string      `json:"templateId"`
	CardData   CardData    `json:"cardData,omitempty"`
	Options    CardOptions `json:"options,omitempty"`
}

type MessageContent struct {
	ContentType string      `json:"contentType"`
	Content     CardContent `json:"content"`
}
