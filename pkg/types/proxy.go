// Package types 定义了 SubConverter 中使用的所有类型和枚举
package types

import (
	"fmt"
	"strings"
)

// ProxyType 代理类型枚举 - 完全对应 C++ 版本的代理类型
type ProxyType int

// 代理类型常量定义 - 与 C++ 版本保持一致
const (
	ProxyTypeUnknown ProxyType = iota
	ProxyTypeShadowsocks
	ProxyTypeShadowsocksR
	ProxyTypeVMess
	ProxyTypeTrojan
	ProxyTypeSnell
	ProxyTypeHTTP
	ProxyTypeHTTPS
	ProxyTypeSOCKS5
	ProxyTypeWireGuard
	ProxyTypeVLESS
	ProxyTypeHysteria
	ProxyTypeHysteria2
	ProxyTypeTUIC
	ProxyTypeAnyTLS
	ProxyTypeMieru
)

// 代理类型字符串映射
var proxyTypeStrings = map[ProxyType]string{
	ProxyTypeUnknown:      "Unknown",
	ProxyTypeShadowsocks:  "SS",
	ProxyTypeShadowsocksR: "SSR",
	ProxyTypeVMess:        "VMess",
	ProxyTypeTrojan:       "Trojan",
	ProxyTypeSnell:        "Snell",
	ProxyTypeHTTP:         "HTTP",
	ProxyTypeHTTPS:        "HTTPS",
	ProxyTypeSOCKS5:       "SOCKS5",
	ProxyTypeWireGuard:    "WireGuard",
	ProxyTypeVLESS:        "VLESS",
	ProxyTypeHysteria:     "Hysteria",
	ProxyTypeHysteria2:    "Hysteria2",
	ProxyTypeTUIC:         "TUIC",
	ProxyTypeAnyTLS:       "AnyTLS",
	ProxyTypeMieru:        "Mieru",
}

// 字符串到代理类型的反向映射
var stringToProxyType = map[string]ProxyType{
	"unknown":    ProxyTypeUnknown,
	"ss":         ProxyTypeShadowsocks,
	"shadowsocks": ProxyTypeShadowsocks,
	"ssr":        ProxyTypeShadowsocksR,
	"shadowsocksr": ProxyTypeShadowsocksR,
	"vmess":      ProxyTypeVMess,
	"trojan":     ProxyTypeTrojan,
	"snell":      ProxyTypeSnell,
	"http":       ProxyTypeHTTP,
	"https":      ProxyTypeHTTPS,
	"socks5":     ProxyTypeSOCKS5,
	"socks":      ProxyTypeSOCKS5,
	"wireguard":  ProxyTypeWireGuard,
	"wg":         ProxyTypeWireGuard,
	"vless":      ProxyTypeVLESS,
	"hysteria":   ProxyTypeHysteria,
	"hysteria2":  ProxyTypeHysteria2,
	"hy2":        ProxyTypeHysteria2,
	"tuic":       ProxyTypeTUIC,
	"anytls":     ProxyTypeAnyTLS,
	"mieru":      ProxyTypeMieru,
}

// String 返回代理类型的字符串表示
func (pt ProxyType) String() string {
	if str, ok := proxyTypeStrings[pt]; ok {
		return str
	}
	return "Unknown"
}

// IsValid 检查代理类型是否有效
func (pt ProxyType) IsValid() bool {
	return pt > ProxyTypeUnknown && pt <= ProxyTypeMieru
}

// SupportsTLS 检查代理类型是否支持 TLS
func (pt ProxyType) SupportsTLS() bool {
	switch pt {
	case ProxyTypeVMess, ProxyTypeVLESS, ProxyTypeTrojan, 
		 ProxyTypeHysteria, ProxyTypeHysteria2, ProxyTypeTUIC:
		return true
	default:
		return false
	}
}

// SupportsUDP 检查代理类型是否支持 UDP
func (pt ProxyType) SupportsUDP() bool {
	switch pt {
	case ProxyTypeShadowsocks, ProxyTypeShadowsocksR, ProxyTypeVMess,
		 ProxyTypeVLESS, ProxyTypeTrojan, ProxyTypeWireGuard,
		 ProxyTypeHysteria, ProxyTypeHysteria2, ProxyTypeTUIC:
		return true
	default:
		return false
	}
}

// RequiresPassword 检查代理类型是否需要密码
func (pt ProxyType) RequiresPassword() bool {
	switch pt {
	case ProxyTypeShadowsocks, ProxyTypeShadowsocksR, ProxyTypeTrojan:
		return true
	default:
		return false
	}
}

// RequiresUUID 检查代理类型是否需要 UUID
func (pt ProxyType) RequiresUUID() bool {
	switch pt {
	case ProxyTypeVMess, ProxyTypeVLESS:
		return true
	default:
		return false
	}
}

// GetDefaultPort 获取代理类型的默认端口
func (pt ProxyType) GetDefaultPort() uint16 {
	switch pt {
	case ProxyTypeShadowsocks, ProxyTypeShadowsocksR, ProxyTypeVMess,
		 ProxyTypeVLESS, ProxyTypeTrojan, ProxyTypeHysteria, 
		 ProxyTypeHysteria2, ProxyTypeTUIC:
		return 443
	case ProxyTypeHTTP:
		return 80
	case ProxyTypeHTTPS:
		return 443
	case ProxyTypeSOCKS5:
		return 1080
	case ProxyTypeSnell:
		return 6160
	case ProxyTypeWireGuard:
		return 51820
	default:
		return 443
	}
}

// ParseProxyType 从字符串解析代理类型
func ParseProxyType(s string) ProxyType {
	if pt, ok := stringToProxyType[strings.ToLower(s)]; ok {
		return pt
	}
	return ProxyTypeUnknown
}

// MarshalText 实现 encoding.TextMarshaler 接口
func (pt ProxyType) MarshalText() ([]byte, error) {
	return []byte(pt.String()), nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler 接口
func (pt *ProxyType) UnmarshalText(text []byte) error {
	*pt = ParseProxyType(string(text))
	if *pt == ProxyTypeUnknown && string(text) != "Unknown" {
		return fmt.Errorf("unknown proxy type: %s", string(text))
	}
	return nil
}

// GetSupportedProxyTypes 获取所有支持的代理类型
func GetSupportedProxyTypes() []ProxyType {
	return []ProxyType{
		ProxyTypeShadowsocks,
		ProxyTypeShadowsocksR,
		ProxyTypeVMess,
		ProxyTypeTrojan,
		ProxyTypeSnell,
		ProxyTypeHTTP,
		ProxyTypeHTTPS,
		ProxyTypeSOCKS5,
		ProxyTypeWireGuard,
		ProxyTypeVLESS,
		ProxyTypeHysteria,
		ProxyTypeHysteria2,
		ProxyTypeTUIC,
		ProxyTypeAnyTLS,
		ProxyTypeMieru,
	}
}

// GetProxyTypeNames 获取所有代理类型的名称
func GetProxyTypeNames() []string {
	types := GetSupportedProxyTypes()
	names := make([]string, len(types))
	for i, pt := range types {
		names[i] = pt.String()
	}
	return names
}