package response

import (
	"encoding/json"
	"net/http"

	"github.com/arf-dev/mekari-test/pkg/httputils/apierror"
)

type BaseResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (response *BaseResponse) Base() *BaseResponse {
	return response
}

type ResponseBody interface {
	Base() *BaseResponse
}

func Send(w http.ResponseWriter, message string, body ResponseBody, err error) {
	statusCode := http.StatusOK
	code := "success"
	if body == nil {
		body = &BaseResponse{}
	}
	base := body.Base()
	if err != nil {
		apiErr, ok := err.(apierror.APIError)
		if !ok {
			apiErr = apierror.ErrInternalServer
		}
		message = apiErr.Message
		code = apiErr.Code
		statusCode = apiErr.StatusCode
	}
	base.Message = message
	base.Code = code

	bodyJson, _ := json.Marshal(body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(bodyJson)
}
