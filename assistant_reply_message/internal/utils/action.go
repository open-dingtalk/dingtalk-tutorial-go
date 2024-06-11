package utils

import (
	"encoding/json"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/plugin"
)

func CreateSuccessResponse(payload any) *plugin.GraphResponse {
	body, _ := json.Marshal(payload)
	resp := &plugin.GraphResponse{
		StatusLine: plugin.GraphStatusLine{200, "OK"},
		Body:       string(body),
	}
	return resp
}
