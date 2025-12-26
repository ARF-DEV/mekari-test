package authhandlr

import (
	"net/http"

	"github.com/arf-dev/mekari-test/model"
	"github.com/arf-dev/mekari-test/pkg/httputils/response"
	"github.com/arf-dev/mekari-test/service/authsv"
)

type Handler struct {
	authServ *authsv.Service
}

func New() *Handler {
	return &Handler{}
}
func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	req := model.LoginRequest{}
	if err := response.ParseRequestBody(r, &req); err != nil {
		response.Send(w, "error", nil, http.StatusBadRequest, err)
		return
	}

	// TODO: login logic
	handler.authServ.AuthenticateUser(r.Context(), req)

	response.Send(w, "success", &model.LoginResponse{}, http.StatusOK, nil)
}

func Login(req model.LoginRequest) error {
	// if user exists, then create jwt
	// otherwise, return error
	return nil
}
