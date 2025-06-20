package mw

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

var (
	HeaderAuthorization = "Authorization"
	HeaderTag           = "userinfo"
)

// UserAuthMw 用户授权中间件
func UserAuthMw(skipper ...SkipperFunc) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {

		c.Next(ctx)
	}
}
