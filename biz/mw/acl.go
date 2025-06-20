package mw

import (
	"context"
	"dolphin-sandbox/biz/response"
	"dolphin-sandbox/conf"
	"dolphin-sandbox/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"strings"
)

// AccessAcl 访问控制
func AccessAcl() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ipstr := c.ClientIP()
		ips := strings.Split(ipstr, ",")
		hitAllow, _ := utils.IPsInSubnets(ips, conf.GetConf().Acl.AllowCIDR)
		hitDeny, _ := utils.IPsInSubnets(ips, conf.GetConf().Acl.DenyCIDR)
		if (len(conf.GetConf().Acl.AllowCIDR) > 0 && !hitAllow) || (len(conf.GetConf().Acl.DenyCIDR) > 0 && hitDeny) {
			response.SendBaseResp(ctx, c, response.AclErr)
			c.Abort()
			return
		}
		// 处理请求
		c.Next(ctx)
	}
}
