package api

import (
	"encoding/json"
	"github.com/brcodingdev/chat-app/service/errors"
	"github.com/brcodingdev/chat-app/service/port/http/response"
	"io"
	"net/http"
)

type JWTProps string

// ParseBody ...
func parseBody(r *http.Request, x interface{}) error {
	if body, err := io.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal(body, x); err != nil {
			return err
		}
	}
	return nil
}

func ok(res []byte, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// ErrResponse ...
func errResponse(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	errCode := codeFrom(err)
	w.WriteHeader(errCode)
	res := response.ErrorResponse{
		Message: err.Error(),
		Status:  false,
		Code:    errCode,
	}
	data, _ := json.Marshal(res)
	w.Write(data)
}

// codeFrom returns the http status code from service errors
func codeFrom(err error) int {
	switch err {
	case errors.ErrInCredentials:
		return http.StatusBadRequest
	case errors.ErrDupEmail:
		return http.StatusBadRequest
	case errors.ErrInRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
