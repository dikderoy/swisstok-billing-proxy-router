package models

import (
	"fmt"
	"net/http"
	"../app"
)

func RequestHandler(res http.ResponseWriter, req *http.Request) *app.AppError {
	defer req.Body.Close()

	addr := req.Header.Get("X-Callback-Path")
	if addr == "" {
		return app.NewAppError(400, "Callback path not given")
	}

	r := NewRequest(*req)
	response, err := r.requestTarget()
	if err != nil {
		return app.NewAppError(500, "Failed to proxy request to target")
	}
	GlobalBucket().add(r)
	app.Logger.Printf("serving:", *r)
	res.Header().Add("X-Request-Id", fmt.Sprint(r.id))
	res.Write(response)
	return nil
}
