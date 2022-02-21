package httputil

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Meta   *Meta           `json:"meta,omitempty"`
	Data   interface{}     `json:"data,omitempty"`
	Errors []ErrorResponse `json:"errors,omitempty"`
}

type Meta struct {
	DataCount   int     `json:"data_count"`
	ProcessTime float64 `json:"process_time"`
}

func WriteErrorResponse(w http.ResponseWriter, errors []ErrorResponse) {
	resp := Response{
		Errors: errors,
	}
	respJson, _ := json.Marshal(resp)
	WriteResponse(w, respJson)
	return
}

func WriteResponse(w http.ResponseWriter, data json.RawMessage) {
	DecorateHeader(w)
	w.Write(data)
	return
}
