package models

type InteractResult[T any] struct {
	Result T `json:"result"`
}

type PrepareResponse struct {
	ConversationToken string `json:"conversationToken"`
}

type UpdateResponse struct {
	Success bool `json:"success"`
}

type FinishResponse struct {
	Success bool `json:"success"`
}
