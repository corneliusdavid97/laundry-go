package httputil

import (
	"fmt"
	"net/http"
)

const ContentTypeJson = "application/json"
const ContentTypeForm = "application/x-www-form-urlencoded"

func ValidateRequest(r *http.Request, method, contentType string) ErrorResponse {
	if r.Method != method {
		return ErrorResponse{
			HttpStatus: http.StatusMethodNotAllowed,
			Title:      http.StatusText(http.StatusMethodNotAllowed),
			Detail:     fmt.Sprintf("Method %s not supported, only %s allowed", r.Method, method),
		}
	}
	if r.Header.Get("Content-Type") != contentType {
		return ErrorResponse{
			HttpStatus: http.StatusUnsupportedMediaType,
			Title:      http.StatusText(http.StatusUnsupportedMediaType),
			Detail:     "Your request containing unsupported media type. Our service currently only support 'application/json' as media type",
		}
	}
	return ErrorResponse{}
}
