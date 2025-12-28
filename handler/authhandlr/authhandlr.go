package authhandlr

import (
	"net/http"

	"github.com/arf-dev/mekari-test/model"
	"github.com/arf-dev/mekari-test/pkg/httputils/request"
	"github.com/arf-dev/mekari-test/pkg/httputils/response"
	"github.com/arf-dev/mekari-test/service/authsv"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	authServ *authsv.Service
}

func New(authServ *authsv.Service) *Handler {
	return &Handler{authServ: authServ}
}

// Login godoc
//
//	@Summary	User login
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		request	body		model.LoginRequest	true	"request body"
//	@Success	200		{object}	model.LoginResponse
//	@Failure	400		{object}	response.BaseResponse
//	@Failure	404		{object}	response.BaseResponse
//	@Failure	500		{object}	response.BaseResponse
//	@Router		/auth/login [post]
func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	req := model.LoginRequest{}
	if err := request.ParseRequestBody(r, &req); err != nil {
		log.Log().Err(err).Msg("error when parsing body")
		response.Send(w, "error", nil, err)
		return
	}

	token, err := handler.authServ.AuthenticateUser(r.Context(), req)
	if err != nil {
		response.Send(w, "error", nil, err)
		return
	}

	response.Send(
		w,
		"success",
		&model.LoginResponse{
			Data: token,
		},
		nil,
	)
}
