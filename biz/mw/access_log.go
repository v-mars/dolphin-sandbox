package mw

import (
	"context"
	"dolphin-sandbox/conf"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

const (
	ApiDataTag = "api_data" //

)

func AccessLog() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		//reqId, _ := ctx.Value(conf.RequestIDHeaderValue).(string)
		//c.Set(conf.RequestIDHeaderValue, reqId)

		defer func() {
			un := c.GetString("username")
			if un == "" {
				un = "Anonymous"
			}
			end := time.Now()
			latency := end.Sub(start).String()
			if conf.GetConf().Server.EnableAudit {
				hlog.CtxDebugf(ctx, "username=%s status=%d cost=%s method=%s full_path=%s client_ip=%s host=%s user_agent=%s",
					un, c.Response.StatusCode(), latency,
					c.Request.Header.Method(), c.Request.URI().RequestURI(), GetReqClientIp(ctx, c), c.Request.Host(),
					c.UserAgent())
			}
		}()
		c.Next(ctx)
		if c.Response.Header.Get("Server") == "hertz" {
			c.Response.Header.Del("Server")
		}
	}
}
