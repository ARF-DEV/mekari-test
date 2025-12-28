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

// CreateExpense godoc
//
//	@Summary	Create an expense
//	@Tags		expenses
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		request	body		model.CreateExpenseRequest	true	"request body"
//	@Success	200		{object}	model.CreateExpenseResponse
//	@Failure	400		{object}	response.BaseResponse
//	@Failure	404		{object}	response.BaseResponse
//	@Failure	500		{object}	response.BaseResponse
//	@Router		/expenses [post]
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

// GetExpenseDetail godoc
//
//	@Summary	Show expense's detail info
//	@Tags		expenses
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path		integer	true	"expense id"
//	@Success	200	{object}	model.GetExpenseResponse
//	@Failure	400	{object}	response.BaseResponse
//	@Failure	404	{object}	response.BaseResponse
//	@Failure	500	{object}	response.BaseResponse
//	@Router		/expenses/{id} [get]
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

// GetExpenseList godoc
//
//	@Summary	Show list of expenses
//	@Tags		expenses
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		page	query		integer	false	"page"
//	@Param		size	query		integer	false	"size"
//	@Success	200		{object}	model.GetExpenseListResponse
//	@Failure	400		{object}	response.BaseResponse
//	@Failure	404		{object}	response.BaseResponse
//	@Failure	500		{object}	response.BaseResponse
//	@Router		/expenses [get]
func (handler *Handler) HandleGetExpenseList(w http.ResponseWriter, r *http.Request) {
	req := model.GetExpenseListRequest{}
	if err := request.ParseQueryParam(r, &req); err != nil {
		log.Log().Err(err).Msg("error when parsing GetExpenseListRequest")
		response.Send(w, "", nil, apierror.ErrBadRequest)
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Page = 20
	}

	resp, err := handler.expenseServ.GetExpenseList(r.Context(), req)
	if err != nil {
		response.Send(w, "", nil, err)
		return
	}
	response.Send(w, "Success", &resp, nil)
}

// UpdateExpense godoc
//
//	@Summary	approve/reject expense
//	@Tags		expenses
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id		path		integer	true	"expense id"
//	@Param		status	path		string	true	"approve / reject"	Enums(approve, reject)
//	@Success	200		{object}	model.BaseResponse
//	@Failure	400		{object}	response.BaseResponse
//	@Failure	404		{object}	response.BaseResponse
//	@Failure	500		{object}	response.BaseResponse
//	@Router		/expenses/{id}/{status} [put]
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
