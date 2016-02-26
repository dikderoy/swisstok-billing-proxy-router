package app

import (
	"net/http"
)

type AppHandler interface {
	Handle(http.ResponseWriter, *http.Request) *AppError
	SetApp(*App)
}

type AppKernel struct {
	AppHandler
	app *App
}

func (self AppKernel) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	return self.Handle(res, req)
}

func (self *AppKernel) SetApp(app *App) {
	self.app = app
}

func (self *AppKernel) GetApp() *App {
	return self.app
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
