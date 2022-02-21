package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/corneliusdavid97/laundry-go/src/product"
	"github.com/corneliusdavid97/laundry-go/tools/httputil"
	"github.com/corneliusdavid97/laundry-go/tools/timer"
)

type HTTPHandler struct {
	svc product.Service
	cfg Config
}

type Config struct {
	Timeout time.Duration
}

func (h *HTTPHandler) HandleGetAllActiveProduct(w http.ResponseWriter, r *http.Request) {
	t := timer.NewTimer()

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	httpErr := httputil.ValidateRequest(r, http.MethodGet, httputil.ContentTypeJson)
	if !httpErr.Empty() {
		httputil.WriteErrorResponse(w, []httputil.ErrorResponse{httpErr})
		return
	}

	var respErrs []httputil.ErrorResponse

	// construct queries
	var filter product.Filter
	sActive := r.URL.Query().Get("active")
	if len(sActive) > 0 {
		b, _ := strconv.ParseBool(sActive)
		filter.Active = &b
	}
	sSatuan := r.URL.Query().Get("is_satuan")
	if len(sSatuan) > 0 {
		b, _ := strconv.ParseBool(sSatuan)
		filter.IsSatuan = &b
	}

	products, err := h.svc.GetAllActiveProducts(ctx, filter)
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
		Data: products,
		Meta: &httputil.Meta{
			DataCount:   len(products),
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

func NewHandler(svc product.Service, cfg Config) *HTTPHandler {
	return &HTTPHandler{
		svc: svc,
		cfg: cfg,
	}
}
