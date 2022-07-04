package middleware

import (
	"net/http"
)

type Error struct {
	Code    string `json:"code" valid:"required"`
	Message string `json:"message" valid:"required"`
}

func WriteError(w http.ResponseWriter, _ *http.Request, code int, msg, msgCode string) {
	w.WriteHeader(code)
	d, _ := Error{Code: msgCode, Message: msg}.MarshalJSON()
	_, _ = w.Write(d)
}
