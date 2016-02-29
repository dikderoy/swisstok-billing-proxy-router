package app

import "net/http"

type Handler interface {
	SetApp(*App)
	GetApp() *App
	ServeHTTP(res http.ResponseWriter, req *http.Request) *AppError
}

type AppHandler struct {
	handler Handler
}

func (self *AppHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	if err := self.handler.ServeHTTP(res, req); err != nil {
		http.Error(res, http.StatusText(err.code), err.code)
		Logger.Printf("error[%d %s]:%s", err.code, err.message, err.error)
	}
}

func NewAppHandler(h Handler) http.Handler {
	return &AppHandler{handler:h}
}

type Kernel struct {
	app *App
}

func (self *Kernel) SetApp(app *App) {
	self.app = app
}

func (self *Kernel) GetApp() *App {
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
