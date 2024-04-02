package models

type PluginRequest struct {
	SenderUnionId      string `json:"senderUnionId"`
	SenderCorpId       string `json:"senderCorpId"`
	OpenConversationId string `json:"openConversationId"`
	SessionWebhook     string `json:"sessionWebhook"`
}
