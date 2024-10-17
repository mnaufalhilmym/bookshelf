package model

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code int
	Err  error
}

func (e Error) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Err.Error())
}

func BadRequest(err error) error {
	return &Error{
		Code: http.StatusBadRequest,
		Err:  err,
	}
}

func NotFound(err error) error {
	return &Error{
		Code: http.StatusNotFound,
		Err:  err,
	}
}

func InternalServerError(err error) error {
	return &Error{
		Code: http.StatusInternalServerError,
		Err:  err,
	}
}
