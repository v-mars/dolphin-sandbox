package mw

import (
	"context"
	"dolphin-sandbox/biz/response"
	"dolphin-sandbox/conf"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
)

// ReqAesDec header:  X-Enc-Data = yes
func ReqAesDec() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		encryptTag := c.Request.Header.Get(conf.HeaderEncTag)
		if encryptTag == "yes" {
			rawData := c.GetRawData() // body 只能读一次，读出来之后需要重置下 Body
			if len(rawData) > 0 {
				//fmt.Println("aes rawData:", string(rawData))
				//data, err := utils.DeTxtByAesWithErr(string(rawData), response.DefaultAesKey)
				//if err != nil {
				//	response.SendBaseResp(ctx, c, fmt.Errorf("http请求body解密失败：%s", err))
				//	c.Abort()
				//	return
				//}
				c.Request.ResetBody()
				c.Request.SetBodyRaw([]byte("data")) // 重置body
				c.Request.Header.Set("Content-Type", "application/json; charset=utf-8")
			}
		}
		// 处理请求
		c.Next(ctx)
	}
}

func ReqAesCheck() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		encryptTag := c.Request.Header.Get(conf.HeaderEncTag)
		if encryptTag == "yes" {
			// 处理请求
			c.Next(ctx)
		} else {
			response.SendBaseResp(ctx, c, fmt.Errorf("缺少请求头：X-Enc-Data"))
			c.Abort()
			return
		}
	}
}
