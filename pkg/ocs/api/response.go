package api

import "encoding/json"

type Response struct {
	Result struct {
		Status     string       `json:"status"`
		Code       string       `json:"code"`
		Message    string       `json:"message"`
		Violations []Violations `json:"violations"`
	} `json:"result"`
	Name      string         `json:"name"`
	RequestID string         `json:"requestId"`
	Data      []ResponseData `json:"data"`
}

type ResponseData struct {
	Name       string          `json:"name"`
	SeqNo      uint            `json:"seqno"`
	ObjectType string          `json:"objectType"`
	Object     json.RawMessage `json:"object"`
}

type Violations struct {
	Path string `json:"path"`
	Code string `json:"code"`
}
