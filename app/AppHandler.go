package app

import (
	"net/http"
)

type AppHandler func(http.ResponseWriter, *http.Request) *AppError

func (fn AppHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			http.Error(res, "EFatal Recover", 500)
		}
	}()
	if err := fn(res, req); err != nil {
		http.Error(res, http.StatusText(err.code), err.code)
		Logger.Printf("error[%d %s]:%s", err.code, err.message, err.error)
	}
}

type AppError struct {
	code    int
	message string
	error   error
}

func NewAppError(code int, message string, previous ...error) *AppError {
	var err error
	if len(previous) > 0 {
		err = previous[:1][0]
	}
	return &AppError{code: code, message: message, error: err}
}
