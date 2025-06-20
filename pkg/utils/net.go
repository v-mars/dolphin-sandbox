package utils

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// IPInSubnet ; 支持IPv4和IPv6
// ipAddr: 192.168.0.1;
// cidr: 192.168.0.0/24;
func IPInSubnet(ipAddr, cidr string) (bool, error) {
	ip := net.ParseIP(ipAddr)
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, err
	}
	if ipNet.Contains(ip) {
		return true, nil
	}
	return false, nil
}

func IPsInSubnets(ips, cidrs []string) (bool, error) {
	for _, ip := range ips {
		ip := ip
		for _, cidr := range cidrs {
			cidr := cidr
			inSubnet, err := IPInSubnet(ip, cidr)
			if err != nil {
				return false, err
			}
			if inSubnet {
				return true, nil
			}
		}
	}
	return false, nil
}

// Long2ip int转ip
func Long2ip(i uint32) net.IP {
	ip := make([]byte, 4)
	binary.BigEndian.PutUint32(ip, i)
	return ip
}

// Ip2long ip转int
func Ip2long(ip net.IP) uint32 {
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

// ConvertCIDRToNetmask 将 CIDR 表示法转换为点分十进制格式的子网掩码（IPv4）或十六进制格式的子网掩码（IPv6）
func ConvertCIDRToNetmask(cidr string) (ipmask string, err error) {
	// 解析 CIDR 表示法
	ipo, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}

	// 获取子网掩码
	mask := ipNet.Mask

	// 根据 IP 地址的长度判断是 IPv4 还是 IPv6
	if len(mask) == net.IPv4len {
		// IPv4: 将子网掩码转换为点分十进制格式
		maskStr := fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])
		return fmt.Sprintf("%s/%s", ipo, maskStr), nil
	} else if len(mask) == net.IPv6len {
		// IPv6: 将子网掩码转换为十六进制格式
		maskStr := fmt.Sprintf("%x:%x:%x:%x:%x:%x:%x:%x",
			mask[0]<<8|mask[1], mask[2]<<8|mask[3], mask[4]<<8|mask[5], mask[6]<<8|mask[7],
			mask[8]<<8|mask[9], mask[10]<<8|mask[11], mask[12]<<8|mask[13], mask[14]<<8|mask[15])
		return fmt.Sprintf("%s/%s", ipo, maskStr), nil
	}

	return "", fmt.Errorf("unknown IP length")
}

// IsCIDR 判断输入是否为有效的 CIDR 表示法
func IsCIDR(input string) (ipOk, isIpv4, isCIDR bool) {
	// 检查是否包含斜杠
	if !strings.Contains(input, "/") {
		return ipOk, false, false
	}

	// 分割 IP 地址和前缀长度
	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return ipOk, false, false
	}

	// 解析 IP 地址
	ip := net.ParseIP(parts[0])
	if ip == nil {
		return ipOk, false, false
	}

	ipOk = true

	// 解析前缀长度
	prefixLength, err := strconv.Atoi(parts[1])
	if err != nil {
		return ipOk, false, false
	}

	// 检查前缀长度是否在有效范围内
	if ip.To4() != nil {
		isIpv4 = true
		// IPv4
		if prefixLength < 0 || prefixLength > 32 {
			return ipOk, isIpv4, false
		}
	} else {
		// IPv6
		if prefixLength < 0 || prefixLength > 128 {
			return ipOk, isIpv4, false
		}
	}

	return ipOk, isIpv4, true
}

// ConvertIPv4ToIPv6 将 IPv4 地址转换为 IPv4 映射的 IPv6 地址
func ConvertIPv4ToIPv6(ipv4Addr string) (net.IP, error) {
	// 解析 IPv4 地址
	ip := net.ParseIP(ipv4Addr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IPv4 address")
	}

	// 确保是 IPv4 地址
	ip = ip.To4()
	if ip == nil {
		return nil, fmt.Errorf("not an IPv4 address")
	}

	// 创建 IPv4 映射的 IPv6 地址
	ipv6Addr := net.IPv6zero
	copy(ipv6Addr[12:], ip)

	return ipv6Addr, nil
}

// ConvertIPv4MaskToIPv6PrefixLength 将 IPv4 子网掩码转换为 IPv6 前缀长度
func ConvertIPv4MaskToIPv6PrefixLength(ipv4Mask string) (int, error) {
	// 解析 IPv4 掩码
	mask := net.ParseIP(ipv4Mask).To4()
	if mask == nil {
		return 0, fmt.Errorf("invalid IPv4 mask")
	}

	// 计算前缀长度
	prefixLength := 0
	for _, octet := range mask {
		for i := 7; i >= 0; i-- {
			if octet&(1<<i) != 0 {
				prefixLength++
			} else {
				break
			}
		}
	}

	return prefixLength, nil
}

// 0: invalid ip
// 4: IPv4
// 6: IPv6
// 检查IP 用于 net.ParseIP(ipv4)
func ParseIP(s string) (net.IP, int) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, 0
	}
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return ip, 4
		case ':':
			return ip, 6
		}
	}
	return nil, 0
}

// IPv4ByLong
// 将 uint32 长整型转换成IPV4 地址
// converts a uint32 represented by a string into an ipv4 address string
// 168427779 => "10.10.1.2"
func IPv4ByLong(ipv4long string) (string, error) {

	ipv4Int, err := strconv.ParseInt(ipv4long, 10, 64)
	if err != nil {
		return "", errors.New(fmt.Sprintf("fail to convert string to Int64 :%s", err.Error()))
	}

	ipv4 := uint32(ipv4Int)

	return fmt.Sprintf("%d.%d.%d.%d", ipv4>>24, ipv4<<8>>24, ipv4<<16>>24, ipv4<<24>>24), nil
}

// IPv6ByLong
// 将 big.Int 长整型转换成IPV6 地址
// converts a big integer represented by a string into an IPv6 address string
// 53174336768213711679990085974688268287=> "2801:0137:0000:0000:0000:ffff:ffff:ffff"
func IPv6ByLong(ipv6long string) (string, error) {
	bi, ok := new(big.Int).SetString(ipv6long, 10)
	if !ok {
		return "", errors.New("fail to convert string to big.Int")
	}

	b255 := new(big.Int).SetBytes([]byte{255})
	var buf = make([]byte, 2)
	p := make([]string, 8)
	j := 0
	var i uint
	tmpint := new(big.Int)
	for i = 0; i < 16; i += 2 {
		tmpint.Rsh(bi, 120-i*8).And(tmpint, b255)
		bytes := tmpint.Bytes()
		if len(bytes) > 0 {
			buf[0] = bytes[0]
		} else {
			buf[0] = 0
		}
		tmpint.Rsh(bi, 120-(i+1)*8).And(tmpint, b255)
		bytes = tmpint.Bytes()
		if len(bytes) > 0 {
			buf[1] = bytes[0]
		} else {
			buf[1] = 0
		}
		p[j] = hex.EncodeToString(buf)
		j++
	}

	return strings.Join(p, ":"), nil
}

// IPv4ToInt 将IPV4 转换成 uint32 长整型
func IPv4ToInt(ipv4 string) (ip uint32) {
	r := `^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})`
	reg, err := regexp.Compile(r)
	if err != nil {
		return 0
	}
	ips := reg.FindStringSubmatch(ipv4)
	if ips == nil {
		return 0
	}

	//上面正则做了判断，这里就不报错了
	ip1, _ := strconv.Atoi(ips[1])
	ip2, _ := strconv.Atoi(ips[2])
	ip3, _ := strconv.Atoi(ips[3])
	ip4, _ := strconv.Atoi(ips[4])

	if ip1 > 255 || ip2 > 255 || ip3 > 255 || ip4 > 255 {
		return 0
	}

	ip += uint32(ip1 * 0x1000000) // 左移24位
	ip += uint32(ip2 * 0x10000)   // 左移16位
	ip += uint32(ip3 * 0x100)     // 左移8位
	ip += uint32(ip4)             // 左移0位

	return ip
}

// IPv6ToInt 将IPV6 转换成 big.Int 长整型
func IPv6ToInt(ipv6 string) (*big.Int, error) {
	ip := net.ParseIP(ipv6)
	return NetIpv6ToInt(ip)
}

// NetIpv6ToInt 将net.IP 类型 转换成 big.Int 长整型
func NetIpv6ToInt(ip net.IP) (*big.Int, error) {
	if ip == nil {
		return nil, errors.New("invalid ipv6")
	}
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(ip)
	return IPv6Int, nil
}
