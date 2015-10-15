package models

import (
	"net/http"
	"strconv"
	"io"
	"fmt"
	"../app"
)

func ListHandler(res http.ResponseWriter, req *http.Request) *app.AppError {
	if id := req.URL.Query().Get("id"); id != "" {
		id, err := strconv.ParseInt(id, 10, 2)
		if err != nil {
			return app.NewAppError(400, "Invalid Param Type")
		}
		if r := Requests[int(id)]; r != nil {
			io.WriteString(res, fmt.Sprintf("%+v", *r))
			return nil
		}
		return app.NewAppError(404, fmt.Sprintf("Record [%d] not found", id))
	}
	return app.NewAppError(400, "NF")
}