package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func ReplyMessage(sessionWebhook string, message any) (*http.Response, error) {
	requestBody, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(sessionWebhook, "application/json", bytes.NewReader(requestBody))
	return resp, err
}
