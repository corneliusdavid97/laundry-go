package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/corneliusdavid97/laundry-go/src/user"
	"github.com/corneliusdavid97/laundry-go/tools/httputil"
	"github.com/corneliusdavid97/laundry-go/tools/timer"
)

type UserResponse struct {
	UserID   int64        `json:"user_id"`
	Username string       `json:"username"`
	Name     string       `json:"name"`
	Role     RoleResponse `json:"role"`
}

type RoleResponse struct {
	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name"`
}

type HTTPHandler struct {
	svc user.Service
	cfg Config
}

type Config struct {
	Timeout time.Duration
}

func (h *HTTPHandler) HandleAuthUser(w http.ResponseWriter, r *http.Request) {
	t := timer.NewTimer()

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	httpErr := httputil.ValidateRequest(r, http.MethodPost, httputil.ContentTypeJson)
	if !httpErr.Empty() {
		httputil.WriteErrorResponse(w, []httputil.ErrorResponse{httpErr})
		return
	}

	var respErrs []httputil.ErrorResponse

	request := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

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

	user, err := h.svc.AuthUser(ctx, request.Username, request.Password)
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
		Data: parseResponse(user),
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

func parseResponse(u user.User) UserResponse {
	return UserResponse{
		UserID:   u.UserID,
		Name:     u.Name,
		Username: u.Username,
		Role: RoleResponse{
			RoleID:   int(u.Role.RoleID),
			RoleName: u.Role.RoleName,
		},
	}
}

func NewHandler(svc user.Service, cfg Config) *HTTPHandler {
	return &HTTPHandler{
		svc: svc,
		cfg: cfg,
	}
}
