package mw

import (
	"context"
	"dolphin-sandbox/conf"
	"github.com/cloudwego/hertz/pkg/app"
)

var (
	HeaderAuthorization = "Authorization"
	HeaderTag           = "userinfo"
)

// UserAuthMw 用户授权中间件
func UserAuthMw(skipper ...SkipperFunc) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if len(skipper) > 0 && skipper[0](ctx, c) {
			c.Next(ctx)
			return
		}
		if conf.GetConf().Server.ApiKey != string(c.GetHeader(conf.SchemeKeyValue)) {
			c.AbortWithStatus(401)
			return
		}
		c.Next(ctx)
	}
}
