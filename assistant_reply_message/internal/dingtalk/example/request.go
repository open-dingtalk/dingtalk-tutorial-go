package example

import (
	"assistant_reply_message/internal/dingtalk/models"
	"encoding/json"
)

type RawExampleRequest struct {
	Date               string `json:"date"`
	InputAttribute     string `json:"inputAttribute"`
	JobNumber          string `json:"jobNum"`
	UnionId            string `json:"unionId"`
	CorpId             string `json:"corpId"`
	SessionWebhook     string `json:"sessionWebhook"`
	RawInput           string `json:"rawInput"`
	Location           string `json:"location"`
	UserId             string `json:"userId"`
	OpenConversationId string `json:"openConversationId"`
	ConversationToken  string `json:"conversationToken"`
}

type WeatherRequest struct {
	Date     string `json:"date"`
	Location string `json:"location"`
}

type CurrentInput struct {
	RawInput       string                    `json:"rawInput"`
	InputAttribute models.UserInputAttribute `json:"input"`
}
type CurrentUser struct {
	JobNumber string `json:"jobNum"`
	UserId    string `json:"userId"`
	UnionId   string `json:"unionId"`
}

type CurrentOrg struct {
	CorpId string `json:"corpId"`
}

type CurrentConversation struct {
	SessionWebhook     string `json:"sessionWebhook"`
	OpenConversationId string `json:"openConversationId"`
	ConversationToken  string `json:"conversationToken"`
}

type ExampleRequest struct {
	WeatherRequest      WeatherRequest
	CurrentUser         CurrentUser
	CurrentInput        CurrentInput
	CurrentOrg          CurrentOrg
	CurrentConversation CurrentConversation
}

func NewExampleRequest(requestBody string) (*ExampleRequest, error) {
	rawRequest := RawExampleRequest{}
	if err := json.Unmarshal([]byte(requestBody), &rawRequest); err != nil {
		return nil, err
	}
	req := &ExampleRequest{
		WeatherRequest: WeatherRequest{
			Location: rawRequest.Location,
			Date:     rawRequest.Date,
		},
		CurrentInput: CurrentInput{
			RawInput: rawRequest.RawInput,
		},
		CurrentUser: CurrentUser{
			UnionId:   rawRequest.UnionId,
			UserId:    rawRequest.UserId,
			JobNumber: rawRequest.JobNumber,
		},
		CurrentOrg: CurrentOrg{
			CorpId: rawRequest.CorpId,
		},
		CurrentConversation: CurrentConversation{
			SessionWebhook:     rawRequest.SessionWebhook,
			OpenConversationId: rawRequest.OpenConversationId,
			ConversationToken:  rawRequest.ConversationToken,
		},
	}
	if err := json.Unmarshal([]byte(rawRequest.InputAttribute), &req.CurrentInput.InputAttribute); err != nil {
		return nil, err
	}

	return req, nil
}
