package dal

import (
	"crypto/tls"
	"dolphin-sandbox/biz/dal/alog"
	"dolphin-sandbox/biz/dal/bind"
	"dolphin-sandbox/biz/response"
	"dolphin-sandbox/cmd"
	"dolphin-sandbox/conf"
	"dolphin-sandbox/docs"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
)

var (
	H *server.Hertz

	// NoLogin 登录验证
	NoLogin = []string{
		"/login", "/404", "/static", "/favicon.ico",
		"/iconfont", "/index.html", "/vite.svg", "/assets", fmt.Sprintf("/%s", conf.Dom),
		"/favicon.ico", "/ping", "/swagger/", "/debug/pprof",
		fmt.Sprintf("/api/v1/%s/swagger", conf.Dom),
	}

	// NoAuthorized 权限验证
	NoAuthorized = []string{
		"/api/v1/dolphin-sandbox/enum",
	}

	TlsCfg   *tls.Config
	CertHash string
)

func init() {
	NoAuthorized = append(NoAuthorized, NoLogin...)
}

func Init() {
	//startTime := time.Now()
	alog.InitZeroLog(conf.GetConf().Server.LogFileName, conf.GetConf().Server.LogFormat, conf.LogLevel())

	response.DefaultAesKey = conf.GetConf().Ext.WebAesKey

	// 参数绑定报错自定义初始化
	bind.Init()

}
func setSwag() {
	docs.SwaggerInfo.Title = ""
	docs.SwaggerInfo.Description = ""
	docs.SwaggerInfo.Version = cmd.Ver
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
