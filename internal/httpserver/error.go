package httpserver

import (
	"net/http"
)

const (
	CodeInvalidArgument = 400001
	CodeResponseError   = 400002
)

type errorResponse struct {
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func apiError(w http.ResponseWriter, err string, status, code int) {
	errString, _ := errorResponse{
		Status:  status,
		Code:    code,
		Message: err,
	}.MarshalJSON()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	_, _ = w.Write(errString)
}
