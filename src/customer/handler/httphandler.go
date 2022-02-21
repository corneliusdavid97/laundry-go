package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/corneliusdavid97/laundry-go/src/customer"
	"github.com/corneliusdavid97/laundry-go/tools/httputil"
	"github.com/corneliusdavid97/laundry-go/tools/timer"
)

type HTTPHandler struct {
	svc customer.Service
	cfg Config
}

type Config struct {
	Timeout time.Duration
}

type Customer struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Active      bool   `json:"active"`
}

func (h *HTTPHandler) HandleGetAllActiveCustomer(w http.ResponseWriter, r *http.Request) {
	t := timer.NewTimer()

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	httpErr := httputil.ValidateRequest(r, http.MethodGet, httputil.ContentTypeJson)
	if !httpErr.Empty() {
		httputil.WriteErrorResponse(w, []httputil.ErrorResponse{httpErr})
		return
	}

	var respErrs []httputil.ErrorResponse

	custs, err := h.svc.GetAllActiveCustomer(ctx)
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

	res := make([]Customer, 0)

	for _, c := range custs {
		res = append(res, parseCustomer(c))
	}

	resp := httputil.Response{
		Data: res,
		Meta: &httputil.Meta{
			DataCount:   len(custs),
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

func (h *HTTPHandler) HandleInsertNewCustomer(w http.ResponseWriter, r *http.Request) {
	t := timer.NewTimer()

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	httpErr := httputil.ValidateRequest(r, http.MethodPost, httputil.ContentTypeForm)
	if !httpErr.Empty() {
		httputil.WriteErrorResponse(w, []httputil.ErrorResponse{httpErr})
		return
	}

	var respErrs []httputil.ErrorResponse

	r.ParseForm()
	err := h.svc.InsertNewCustomer(ctx, customer.Customer{
		Name:        r.FormValue("name"),
		PhoneNumber: r.FormValue("phone_number"),
		Address:     r.FormValue("address"),
	})
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

	respData := struct {
		Success bool   `json:"success"`
		Detail  string `json:"detail"`
	}{
		Success: true,
		Detail:  "Insert new customer successful",
	}

	resp := httputil.Response{
		Data: respData,
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

func parseCustomer(cust customer.Customer) Customer {
	return Customer{
		ID:          cust.ID,
		Name:        cust.Name,
		PhoneNumber: cust.PhoneNumber,
		Address:     cust.Address,
		Active:      cust.Active,
	}
}

func NewHandler(svc customer.Service, cfg Config) *HTTPHandler {
	return &HTTPHandler{
		svc: svc,
		cfg: cfg,
	}
}
