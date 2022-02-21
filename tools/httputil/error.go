package httputil

type HttpStatus int

type ErrorResponse struct {
	HttpStatus HttpStatus `json:"status"`
	Title      string     `json:"title"`
	Detail     string     `json:"detail"`
}

func (e ErrorResponse) Empty() bool {
	return e == ErrorResponse{}
}
