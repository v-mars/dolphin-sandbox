package srv_http

import (
	"context"
	"dolphin-sandbox/biz/dal"
	"dolphin-sandbox/biz/mw"
	"dolphin-sandbox/conf"
	"fmt"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/requestid"
	"github.com/spf13/cobra"
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	hc "github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
)

var (
	//cfg string

	RootCmd = &cobra.Command{
		Use:     "",
		Short:   "InitWebStart DemoRun aigc Server",
		Long:    `This is dolphin aigc server`,
		Example: `## 启动命令 ./app -c ./conf/config.yaml`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(conf.FlagConf) == 0 {
				_ = cmd.Help()
				os.Exit(0)
			}
			// 加载配置
			conf.GetConf()

			// 初始化基础依赖
			dal.Init()

			//start grpc server service
			//go func() {
			//	hlog.Fatal(srv_rpc.StartSrvRpc())
			//}()

			// start http server service
			Start()
		},
	}
)

func init() {
	//var c = &conf.AgentConfig{}
	//DefaultIP := "0.0.0.0"
	//DefaultPort := int32(20052)
	RootCmd.Flags().StringVarP(&conf.FlagConf, "conf", "c", "", "config file, example: ./conf/config.yaml")
	//rootCmd.Flags().StringVarP(&c.IP, "ip", "i", DefaultIP, "服务IP")
	//rootCmd.Flags().Int32VarP(&c.Port, "port", "p", DefaultPort, "服务启动的端口: 20052 e.g")
	if conf.FlagConf == "" {
		conf.FlagConf = "./conf/config.yaml"
		//fmt.Println("请使用\"-c\"指定配置文件")
		//os.Exit(-1)
	}
	return
}

var sessionName = fmt.Sprintf("%s_session", conf.Dom)

func Start() {
	hlog.Infof("start http server service")
	//flag.Parse()
	startTime := time.Now()
	var options = []hc.Option{
		server.WithHostPorts(conf.GetConf().Server.HttpAddress),
		server.WithMaxRequestBodySize(20 << 20),
		server.WithTransport(standard.NewTransporter),
		server.WithMaxRequestBodySize(1000 << 20), // 10MB

		//server.WithTracer(prometheus.NewServerTracer(":9091", "/metric")),
	}
	//discover.RegistryWeb(&options)
	//discover.RegistryRPC(&options)
	h := server.Default(options...)
	//dal.Init()
	h.Use(requestid.New())
	if conf.GetConf().Acl.Enabled {
		h.Use(mw.AccessAcl())
	}

	dal.H = h
	h.Use(mw.AccessLog())
	hlog.Debugf("init time：%v", time.Since(startTime).String())

	if len(conf.GetConf().Ext.TelemetryEp) > 0 {
		p := provider.NewOpenTelemetryProvider(
			provider.WithServiceName(conf.AppName),
			// Support setting ExportEndpoint via environment variables: OTEL_EXPORTER_OTLP_ENDPOINT
			provider.WithExportEndpoint(conf.GetConf().Ext.TelemetryEp),
			provider.WithInsecure(),
		)
		defer p.Shutdown(context.Background())
	}

	//gin.DefaultWriter = io.MultiWriter(logg.Gin.Writer())
	// 跨域请求
	h.Use(mw.Cors())
	// 解密加密请求
	h.Use(mw.ReqAesDec())
	h.Use(gzip.Gzip(
		gzip.DefaultCompression,
		// This WithExcludedPaths takes as its parameter the file path
		//gzip.WithExcludedPathRegexes([]string{"/"}),
	))

	// 捕获异常，并返回500
	h.Use(mw.Recovery())

	// 认证过滤
	h.Use(mw.UserAuthMw(mw.AllowPathPrefixSkipper(dal.NoLogin...)))

	register(h)
	print(banner)
	hlog.Infof("startup time：%v", time.Since(startTime).String())
	h.Spin()
}

var banner = ``
