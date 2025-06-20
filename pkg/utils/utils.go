package utils

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/flosch/pongo2/v6"
	"io"
	"log"
	"net"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
	"time"
)

func GetLocalIPv4Address() (string, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addr {

		ipNet, isIpNet := addr.(*net.IPNet)
		if isIpNet && !ipNet.IP.IsLoopback() {
			ipv4 := ipNet.IP.To4()
			if ipv4 != nil {
				return ipv4.String(), nil
			}
		}
	}
	return "", fmt.Errorf("not found ipv4 address")
}

// GetOutBoundIP net.Dial("udp", "8.8.8.8:53")
func GetOutBoundIP(network, addr string) (ip string) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		//log.Errorf("get out bound ip err: %s\n", err)
		panic(any(err))
	}
	var localAddr net.Addr
	if network == "tcp" {
		localAddr = conn.LocalAddr().(*net.TCPAddr) // .(*net.TCPAddr)
	} else {
		localAddr = conn.LocalAddr().(*net.UDPAddr) // .(*net.UDPAddr)
	}
	//fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func EnsureDirExist(name string) error {
	if !FileExists(name) {
		return os.MkdirAll(name, os.ModePerm)
	}
	return nil
}

// DeleteFile 删除文件 name abc/文件名
func DeleteFile(name string) error {
	if FileExists(name) {
		return os.RemoveAll(name)
	}
	return nil
}

func GzipCompressFile(srcPath, dstPath string) error {
	sf, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func(sf *os.File) {
		err := sf.Close()
		if err != nil {

		}
	}(sf)
	df, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func(df *os.File) {
		err := df.Close()
		if err != nil {

		}
	}(df)
	writer := gzip.NewWriter(df)
	writer.Name = dstPath
	writer.ModTime = time.Now().UTC()
	_, err = io.Copy(writer, sf)
	if err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}

func Sum(i []int) int {
	sum := 0
	for _, v := range i {
		sum += v
	}
	return sum
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func CurrentUTCTime() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05 +0000")
}
func CurrentLocalTime() string {
	return time.Now().Local().Format("2006-01-02 15:04:05")
}

// Format 模板字符串替换 例如： I'm is {{ var }}
func Format(text string, args interface{}) (string, error) {
	tpl, err := pongo2.FromString(text)

	if err != nil {
		log.Println(err)
		return "", err
	}
	ctx := pongo2.Context{}
	if err = AnyToAny(args, &ctx); err != nil {
		log.Println(err)
		return "", err
	}
	res, err := tpl.Execute(ctx)

	return res, err
}

// Any2Json 格式化JSON
func Any2Json(s interface{}, ident ...int) string {
	var err error
	var bts []byte
	if len(ident) > 0 {
		bts, err = json.Marshal(s)
	} else {
		bts, err = json.MarshalIndent(s, "", "    ")
	}
	//bts, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		log.Println("obj to json err:", err)
	}
	return string(bts)
}

func Any2Yaml(s interface{}) string {
	bts, err := yaml.JSONToYAML([]byte(Any2Json(s)))
	if err != nil {
		log.Println("obj to yaml err:", err)
	}
	return string(bts)
}

/*
FormattedOsVarString
template: "Hello, $proxy_tracing_id ${method}!"
m: {"proxy_tracing_id": "xx", "method":"get"}
*/
func FormattedOsVarString(template string, m map[string]any) string {
	//template := "Hello, $proxy_tracing_id ${}!"
	mapper := func(placeholderName string) string {
		s, ok := m[placeholderName]
		if ok {
			return fmt.Sprintf("%v", s)
		}
		return ""
	}
	formatted := os.Expand(template, mapper)
	return formatted
}

// IsDomainContainedInWildcard 检查给定的域名是否被泛域名包含
func IsDomainContainedInWildcard(domain, wildcardDomain string) bool {
	// 移除域名和泛域名前面的点（如果存在）
	domain = strings.TrimPrefix(domain, ".")
	wildcardDomain = strings.TrimPrefix(wildcardDomain, ".")

	// 检查泛域名是否以`*.`开头
	if !strings.HasPrefix(wildcardDomain, "*.") {
		// 如果没有，那么它不是一个泛域名，直接比较
		return domain == wildcardDomain
	}

	// 泛域名以`*.`开头，移除这个前缀
	wildcardDomain = wildcardDomain[2:]

	// 检查域名是否以泛域名的后缀结尾
	return strings.HasSuffix(domain, wildcardDomain)
}

func Of[T any](v T) *T {
	return &v
}
