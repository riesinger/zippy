package api

import (
	"encoding/json"
)

type APIResponse struct {
	Success   bool        `json:"success"`
	ErrorText string      `json:"error"`
	ErrorCode int32       `json:"errorCode"`
	Payload   interface{} `json:"payload"`
}

func (a *APIResponse) Marshal() []byte {
	data, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return data
}

func SuccessResponse() *APIResponse {
	return &APIResponse{true, "", 0, nil}
}

func PayloadResponse(payload interface{}) *APIResponse {
	return &APIResponse{true, "", 0, payload}
}

func ErrorResponse(text string, code int32) *APIResponse {
	return &APIResponse{false, text, code, nil}
}
