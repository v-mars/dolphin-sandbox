package conf

import (
	"fmt"
	"github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"sync"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	Env            string
	Authentication Authentication `yaml:"authentication"` //登录验证配置
	Server         Server         `yaml:"server"`         //服务器配置
	Ext            Ext            `yaml:"ext"`
	JWT            JWT            `yaml:"jwt"`
	Acl            Acl            `yaml:"acl"`
}

type Acl struct {
	Enabled   bool     `yaml:"enabled"`
	AllowCIDR []string `yaml:"allow_cidr"`
	DenyCIDR  []string `yaml:"deny_cidr"`
}

type JWT struct {
	Type       string `json:"type" yaml:"type"` // AUTHORIZATION  Bearer Authorization
	Key        string `json:"key" yaml:"key"`
	RefreshKey string `json:"refresh_key" yaml:"refresh_key"`
	Age        int    `json:"age" yaml:"age"`
}

type Authentication struct {
	MaxAge        int    `yaml:"max_age"`
	AuthSecret    string `yaml:"auth_secret"` // session key
	EnableSession bool   `yaml:"enable_session"`
	MFAType       string `json:"mfa_type" yaml:"mfa_type"` // code, email, sms, otp,none
	EnableCode    bool   `yaml:"enable_code"`              // 启用验证码
	EnableMFA     bool   `yaml:"enable_mfa"`               // 启用双因子认证
	EnableEmail   bool   `yaml:"enable_email"`             // 启用邮箱验证码
	EnableSMS     bool   `yaml:"enable_sms"`               // 启用短信验证码

	LoginFailForbid  int    `yaml:"login_fail_forbid"`
	LoginFailCaptcha int    `yaml:"login_fail_captcha"`
	CaptchaType      string `yaml:"captcha_type"`
	Timeout          int64  `yaml:"timeout"`
	LocalAuth        bool   `yaml:"local_auth"`
	LoginTitle       string `yaml:"login_title"`
	LoginDesc        string `yaml:"login_desc"`
}

type Oauth2Srv struct {
	AccessKey  string `yaml:"access_key" json:"access_key"`
	RefreshKey string `yaml:"refresh_key" json:"refresh_key"`
}

type License struct {
	Data      string `yaml:"data" json:"data"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	IsExpire  bool   `json:"is_expire"`
}

type Ext struct {
	ErrReceivers    string   `json:"err_receivers" yaml:"err_receivers"` // for nodejs web
	WebAesKey       string   `json:"web_aes_key" yaml:"web_aes_key"`     // for nodejs web
	DataAesKey      string   `json:"data_aes_key" yaml:"data_aes_key"`   // for db data
	EnableDataEnc   bool     `yaml:"enable_data_enc"`                    // 启用数据加密
	EnableWebEnc    bool     `yaml:"enable_web_enc"`                     // 启用web数据加密
	WechatOpsId     string   `json:"wechat_ops_id" yaml:"wechat_ops_id"` // for db data
	TelemetryEp     string   `json:"telemetry_ep" yaml:"telemetry_ep"`   // TelemetryProvider ExportEndpoint
	OpsMailReceiver []string `json:"ops_mail_receiver" yaml:"ops_mail_receiver"`
}

type Server struct {
	Service        string `yaml:"service"`
	Domain         string `yaml:"domain"` // 服务域名：http://svc.xx.com
	HttpAddress    string `yaml:"http_address"`
	GrpcAddress    string `yaml:"grpc_address"`
	EnablePprof    bool   `yaml:"enable_pprof"`
	EnableGzip     bool   `yaml:"enable_gzip"`
	EnableAudit    bool   `yaml:"enable_audit"`
	LogLevel       string `yaml:"log_level"`
	LogFileName    string `yaml:"log_file_name"`
	LogFormat      string `yaml:"log_format"` // json or console
	LogMaxSize     int    `yaml:"log_max_size"`
	LogMaxBackups  int    `yaml:"log_max_backups"`
	LogMaxAge      int    `yaml:"log_max_age"`
	EnableRegistry bool   `yaml:"enable_registry"`
	EnableSwagger  bool   `yaml:"enable_swagger"`
	AutoUpdateApi  bool   `yaml:"auto_update_api"`
	RegistryCenter string `yaml:"registry_center"`
	CronType       string `yaml:"cron_type"`
	CertCrt        string `yaml:"cert_crt"`
	CertKey        string `yaml:"cert_key"`
	DBType         string `yaml:"dbtype"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	//prefix := "conf"
	//confFileRelPath := filepath.Join(prefix, filepath.Join("", "config.yaml"))
	confFileRelPath := FlagConf
	content, err := os.ReadFile(confFileRelPath)
	if err != nil {
		panic(err)
	}

	conf = new(Config)
	conf.Server.Service = AppName
	conf.Server.HttpAddress = fmt.Sprintf(":%d", AppPort)
	conf.Server.LogLevel = "./log/std.log"
	conf.Server.LogLevel = "debug"
	conf.Server.EnableSwagger = true
	conf.Ext.WebAesKey = "Webkbmon12g3ntSP"
	conf.Ext.DataAesKey = "AY3b5Z72206GorMa"
	conf.Authentication.LoginFailCaptcha = 50
	conf.Authentication.LoginFailForbid = 10
	conf.Authentication.Timeout = 300
	conf.Authentication.LocalAuth = true
	err = yaml.Unmarshal(content, conf)
	if err != nil {
		hlog.Error("parse yaml error - %v", err)
		panic(err)
	}
	if err := validator.Validate(conf); err != nil {
		hlog.Error("validate config error - %v", err)
		panic(err)
	}

	conf.Env = GetEnv()

	//pretty.Printf("%+v\n", conf)
}

func GetEnv() string {
	e := os.Getenv("APP_ENV")
	if len(e) == 0 {
		return "dev"
	}
	return e
}
func GetDyeing() string {
	e := os.Getenv("dyeing")
	if len(e) == 0 {
		return "default"
	}
	return e
}

func LogLevel() hlog.Level {
	level := GetConf().Server.LogLevel
	switch strings.ToLower(level) {
	case "trace":
		return hlog.LevelTrace
	case "debug":
		return hlog.LevelDebug
	case "info":
		return hlog.LevelInfo
	case "notice":
		return hlog.LevelNotice
	case "warn":
		return hlog.LevelWarn
	case "error":
		return hlog.LevelError
	case "fatal":
		return hlog.LevelFatal
	default:
		return hlog.LevelInfo
	}
}

var FlagConf string
