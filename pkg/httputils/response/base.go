package response

import (
	"encoding/json"
	"net/http"
)

type BaseResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func (response *BaseResponse) Base() *BaseResponse {
	return response
}

type ResponseBody interface {
	Base() *BaseResponse
}

func Send(w http.ResponseWriter, message string, body ResponseBody, statusCode int, err error) {
	if body == nil {
		body = &BaseResponse{}
	}
	base := body.Base()
	base.Message = message
	if err != nil {
		base.Error = err.Error()
	}
	bodyJson, _ := json.Marshal(body)
	w.WriteHeader(statusCode)
	w.Write(bodyJson)
}
