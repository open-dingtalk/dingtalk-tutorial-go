package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func DumpRequest(r *http.Request, requestBody []byte) string {
	lines := make([]string, 0)
	lines = append(lines, fmt.Sprintf("%s %s", r.Method, r.URL.String()))
	for name, values := range r.Header {
		for _, value := range values {
			lines = append(lines, fmt.Sprintf("%s: %s", name, value))
		}
	}
	lines = append(lines, "\n")

	body := map[string]any{}
	json.Unmarshal(requestBody, &body)
	beautifyBody, _ := json.MarshalIndent(body, " ", "    ")
	//beautifyBody = requestBody
	lines = append(lines, string(beautifyBody))

	return strings.Join(lines, "\n")
}
