package expensehandlr

import (
	"net/http"

	"github.com/arf-dev/mekari-test/model"
	"github.com/arf-dev/mekari-test/pkg/httputils/apierror"
	"github.com/arf-dev/mekari-test/pkg/httputils/request"
	"github.com/arf-dev/mekari-test/pkg/httputils/response"
	"github.com/arf-dev/mekari-test/pkg/validate"
	"github.com/arf-dev/mekari-test/service/expensesv"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	expenseServ *expensesv.Service
}

func New(expenseServ *expensesv.Service) *Handler {
	return &Handler{
		expenseServ: expenseServ,
	}
}

func (handler *Handler) HandleCreateExpense(w http.ResponseWriter, r *http.Request) {
	req := model.CreateExpenseRequest{}
	if err := request.ParseRequestBody(r, &req); err != nil {
		log.Log().Err(err).Msg("error when parsing CreateExpenseRequest")
		response.Send(w, "", nil, apierror.ErrBadRequest)
		return
	}

	if err := validate.ValidateStruct(req); err != nil {
		log.Log().Err(err).Msg("validation error on HandleCreateExpense")
		response.Send(w, "", nil, apierror.ErrBadRequest)
		return
	}

	expenseId, err := handler.expenseServ.CreateExpense(r.Context(), req)
	if err != nil {
		log.Log().Err(err).Msg("error on CreateExpense")
		response.Send(w, "", nil, err)
		return
	}

	response.Send(
		w,
		"Expense created",
		&model.CreateExpenseResponse{
			Data: expenseId,
		},
		nil,
	)
}

func (handler *Handler) HandleGetExpense(w http.ResponseWriter, r *http.Request) {
	req := model.GetExpenseRequest{}
	if err := request.ParsePathParam(r, &req); err != nil {
		log.Log().Err(err).Msg("error when parsing GetExpenseRequest")
		response.Send(w, "", nil, apierror.ErrBadRequest)
		return
	}

	resp, err := handler.expenseServ.GetExpense(r.Context(), req)
	if err != nil {
		response.Send(w, "", nil, err)
		return
	}

	response.Send(w, "Success", &resp, nil)
}

func (handler *Handler) HandleGetExpenseList(w http.ResponseWriter, r *http.Request) {
	req := model.GetExpenseListRequest{}
	if err := request.ParseQueryParam(r, &req); err != nil {
		log.Log().Err(err).Msg("error when parsing GetExpenseListRequest")
		response.Send(w, "", nil, apierror.ErrBadRequest)
		return
	}

	resp, err := handler.expenseServ.GetExpenseList(r.Context(), req)
	if err != nil {
		response.Send(w, "", nil, err)
		return
	}
	response.Send(w, "Success", &resp, nil)
}

func (handler *Handler) HandleUpdateExpense(w http.ResponseWriter, r *http.Request) {
	req := model.UpdateExpenseRequest{}
	if err := request.ParsePathParam(r, &req); err != nil {
		log.Log().Err(err).Msg("error when parsing UpdateExpenseRequest")
		response.Send(w, "", nil, apierror.ErrBadRequest)
		return
	}
	if err := request.ParseRequestBody(r, &req); err != nil {
		log.Log().Err(err).Msg("error when parsing UpdateExpenseRequest")
		response.Send(w, "", nil, apierror.ErrBadRequest)
		return
	}

	if err := validate.ValidateStruct(req); err != nil {
		log.Log().Err(err).Msg("validation error on HandleUpdateExpense")
		response.Send(w, "", nil, apierror.ErrBadRequest)
		return
	}

	if err := handler.expenseServ.UpdateExpense(r.Context(), req); err != nil {
		response.Send(w, "", nil, err)
		return
	}
	response.Send(w, "Expense updated", nil, nil)
}
