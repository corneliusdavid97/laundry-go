package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/corneliusdavid97/laundry-go/src/transaction"
	"github.com/corneliusdavid97/laundry-go/tools/httputil"
	"github.com/corneliusdavid97/laundry-go/tools/timer"
)

type HTTPHandler struct {
	svc transaction.Service
	cfg Config
}

type Config struct {
	Timeout time.Duration
}

type NewTransactionParam struct {
	transaction.Transaction
	DueDateStr string `json:"due_date_str"`
}

func (h *HTTPHandler) GetTransactionDataByID(w http.ResponseWriter, r *http.Request) {
	t := timer.NewTimer()

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	httpErr := httputil.ValidateRequest(r, http.MethodGet, httputil.ContentTypeForm)
	if !httpErr.Empty() {
		httputil.WriteErrorResponse(w, []httputil.ErrorResponse{httpErr})
		return
	}

	var respErrs []httputil.ErrorResponse

	r.ParseForm()
	id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
		respErrs = append(respErrs, httputil.ErrorResponse{
			HttpStatus: http.StatusBadRequest,
			Title:      http.StatusText(http.StatusBadRequest),
			Detail:     err.Error(),
		})
		httputil.WriteErrorResponse(w, respErrs)
		return
	}
	res, err := h.svc.GetTransactionDataByID(ctx, id)
	if err != nil {
		respErrs = append(respErrs, httputil.ErrorResponse{
			HttpStatus: http.StatusInternalServerError,
			Title:      http.StatusText(http.StatusInternalServerError),
			Detail:     err.Error(),
		})
	}

	if len(respErrs) > 0 {
		httputil.WriteErrorResponse(w, respErrs)
		return
	}

	resp := httputil.Response{
		Data: parseTransactionResponse(res),
		Meta: &httputil.Meta{
			DataCount:   1,
			ProcessTime: t.GetElapsedTime().Seconds(),
		},
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		httputil.WriteErrorResponse(w, []httputil.ErrorResponse{
			{
				HttpStatus: http.StatusInternalServerError,
				Title:      http.StatusText(http.StatusInternalServerError),
				Detail:     "Failed to marshal API response",
			},
		})
		return
	}

	httputil.WriteResponse(w, respJson)
}

func (h *HTTPHandler) HandleNewTransaction(w http.ResponseWriter, r *http.Request) {
	t := timer.NewTimer()

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	httpErr := httputil.ValidateRequest(r, http.MethodPost, httputil.ContentTypeJson)
	if !httpErr.Empty() {
		httputil.WriteErrorResponse(w, []httputil.ErrorResponse{httpErr})
		return
	}

	var respErrs []httputil.ErrorResponse

	var request NewTransactionParam

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respErrs = append(respErrs, httputil.ErrorResponse{
			HttpStatus: http.StatusBadRequest,
			Title:      http.StatusText(http.StatusBadRequest),
			Detail:     err.Error(),
		})
		httputil.WriteErrorResponse(w, respErrs)
		return
	}

	err = json.Unmarshal(data, &request)
	if err != nil {
		respErrs = append(respErrs, httputil.ErrorResponse{
			HttpStatus: http.StatusBadRequest,
			Title:      http.StatusText(http.StatusBadRequest),
			Detail:     err.Error(),
		})
		httputil.WriteErrorResponse(w, respErrs)
		return
	}

	parsedReq, err := parseNewTransactionRequest(request)
	if err != nil {
		respErrs = append(respErrs, httputil.ErrorResponse{
			HttpStatus: http.StatusInternalServerError,
			Title:      http.StatusText(http.StatusInternalServerError),
			Detail:     err.Error(),
		})
	}

	err = h.svc.NewTransaction(ctx, parsedReq)
	if err != nil {
		respErrs = append(respErrs, httputil.ErrorResponse{
			HttpStatus: http.StatusInternalServerError,
			Title:      http.StatusText(http.StatusInternalServerError),
			Detail:     err.Error(),
		})
	}

	res, err := h.svc.GetTransactionDataByID(ctx, request.ID)
	if err != nil {
		respErrs = append(respErrs, httputil.ErrorResponse{
			HttpStatus: http.StatusInternalServerError,
			Title:      http.StatusText(http.StatusInternalServerError),
			Detail:     err.Error(),
		})
	}

	if len(respErrs) > 0 {
		httputil.WriteErrorResponse(w, respErrs)
		return
	}

	resp := httputil.Response{
		Data: parseTransactionResponse(res),
		Meta: &httputil.Meta{
			DataCount:   1,
			ProcessTime: t.GetElapsedTime().Seconds(),
		},
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		httputil.WriteErrorResponse(w, []httputil.ErrorResponse{
			{
				HttpStatus: http.StatusInternalServerError,
				Title:      http.StatusText(http.StatusInternalServerError),
				Detail:     "Failed to marshal API response",
			},
		})
		return
	}

	httputil.WriteResponse(w, respJson)
}

func parseNewTransactionRequest(param NewTransactionParam) (transaction.Transaction, error) {
	dueDate, err := time.Parse("2006-01-02 15:04:05", param.DueDateStr)
	if err != nil {
		return transaction.Transaction{}, err
	}
	return transaction.Transaction{
		ID:            param.ID,
		CustomerID:    param.CustomerID,
		CashierName:   param.CashierName,
		GrandTotal:    param.GrandTotal,
		Paid:          param.Paid,
		PaymentMethod: param.PaymentMethod,
		DueDate:       &dueDate,
		Details:       param.Details,
	}, nil
}

func parseTransactionResponse(trans transaction.Transaction) transaction.Transaction {
	var dueDateStr, dateTakenStr, transactionTimeStr string
	if trans.DueDate != nil {
		dueDateStr = trans.DueDate.Format("2006-01-02")
		trans.DueDateStr = &dueDateStr
	}
	if trans.DateTaken != nil {
		dateTakenStr = trans.DateTaken.Format("2006-01-02 15:04:05")
		trans.DateTakenStr = &dateTakenStr
	}
	if trans.TransactionTime != nil {
		trans.TransactionTimeStr = &transactionTimeStr
		transactionTimeStr = trans.TransactionTime.Format("2006-01-02 15:04:05")
	}
	return trans
}

func NewHandler(svc transaction.Service, cfg Config) *HTTPHandler {
	return &HTTPHandler{
		svc: svc,
		cfg: cfg,
	}
}
