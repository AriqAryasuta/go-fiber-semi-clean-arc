package errors

import (
	stdErrors "errors"
	"net/http"
)

var (
	ErrNotFound   = stdErrors.New("resource not found")
	ErrBadRequest = stdErrors.New("bad request")
	ErrInternal   = stdErrors.New("internal server error")
)

func HTTPStatus(err error) int {
	switch {
	case stdErrors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case stdErrors.Is(err, ErrBadRequest):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
