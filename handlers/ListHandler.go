package handlers

import (
	"net/http"
	"strconv"
	"io"
	"fmt"
	"../app"
	"../models"
)

type ListHandler struct {
	app.Kernel
}

func (self *ListHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) *app.AppError {
	defer req.Body.Close()
	if id := req.URL.Query().Get("id"); id != "" {
		id, err := strconv.ParseInt(id, 10, 2)
		if err != nil {
			return app.NewAppError(400, "Invalid Param Type")
		}
		if r, _ := models.GlobalBucket().Find(int(id)); r != nil {
			io.WriteString(res, fmt.Sprintf("%+v", *r))
			return nil
		}
		return app.NewAppError(404, fmt.Sprintf("Record [%d] not found", id))
	}

	for _, v := range *models.GlobalBucket().GetInternalMap() {
		io.WriteString(res, fmt.Sprintf("\n%+v", *v))
	}
	return nil
}