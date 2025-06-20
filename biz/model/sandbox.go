package model

import (
	"context"
	"dolphin-sandbox/biz/response"
	"github.com/cloudwego/hertz/pkg/app"
)

func BindRequest[T any](ctx context.Context, c *app.RequestContext, success func(T)) {
	var request T
	var err error

	contextType := string(c.GetHeader("Content-Type"))
	if contextType == "application/json" {
		err = c.BindJSON(&request)
	} else {
		err = c.Bind(&request)
	}

	if err != nil {
		resp := response.ErrorResponse(-400, err.Error())
		c.JSON(200, resp)
		return
	}
	success(request)
}

type RunRequest struct {
	Lang          string `json:"lang" form:"lang:required" binding:"required"`
	Code          string `json:"code" form:"code:required" binding:"required"`
	Preload       string `json:"preload" form:"preload"`
	EnableNetwork bool   `json:"enable_network" form:"enable_network"`
}
