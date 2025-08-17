// internal/parser/vmess.go
package parser

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"subconverter-go/pkg/models"
)

// VMeSSParser VMess解析器
type VMeSSParser struct{}

// CanParse 检查是否能解析指定的链接
func (p *VMeSSParser) CanParse(link string) bool {
	return strings.HasPrefix(link, "vmess://")
}

// GetType 获取解析器类型
func (p *VMeSSParser) GetType() models.ProxyType {
	return models.ProxyTypeVMess
}

// VMeSSConfig VMess配置结构
type VMeSSConfig struct {
	Version   string      `json:"v"`
	Remarks   string      `json:"ps"`
	Address   string      `json:"add"`
	Port      interface{} `json:"port"` // 可能是字符串或数字
	ID        string      `json:"id"`
	AlterID   interface{} `json:"aid"` // 可能是字符串或数字
	Security  string      `json:"scy"`
	Network   string      `json:"net"`
	Type      string      `json:"type"`
	Host      string      `json:"host"`
	Path      string      `json:"path"`
	TLS       string      `json:"tls"`
	SNI       string      `json:"sni"`
	ALPN      string      `json:"alpn"`
}

// Parse 解析VMess链接
func (p *VMeSSParser) Parse(link string) (*models.Proxy, error) {
	if !p.CanParse(link) {
		return nil, fmt.Errorf("not a vmess link")
	}

	// 移除前缀
	link = strings.TrimPrefix(link, "vmess://")

	// Base64解码
	decoded, err := base64.StdEncoding.DecodeString(link)
	if err != nil {
		// 尝试URL安全的Base64解码
		decoded, err = base64.URLEncoding.DecodeString(link)
		if err != nil {
			return nil, fmt.Errorf("failed to decode vmess link: %w", err)
		}
	}

	// 解析JSON
	var config VMeSSConfig
	if err := json.Unmarshal(decoded, &config); err != nil {
		return nil, fmt.Errorf("failed to parse vmess config: %w", err)
	}

	proxy := &models.Proxy{
		Type:     models.ProxyTypeVMess,
		Group:    "VMeSSProvider",
		Remark:   config.Remarks,
		Hostname: config.Address,
		UserID:   config.ID,
	}

	// 解析端口
	if port, err := parsePort(config.Port); err == nil {
		proxy.Port = port
	} else {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	// 解析AlterID
	if alterID, err := parseAlterID(config.AlterID); err == nil {
		proxy.AlterID = alterID
	}

	// 解析传输协议
	proxy.TransferProtocol = config.Network
	if proxy.TransferProtocol == "" {
		proxy.TransferProtocol = "tcp"
	}

	// 解析伪装类型
	proxy.FakeType = config.Type
	if proxy.FakeType == "" {
		proxy.FakeType = "none"
	}

	// 解析TLS
	if config.TLS == "tls" {
		proxy.TLS = "tls"
		proxy.TLSSecure = true
	}

	// 解析Host
	proxy.Host = config.Host

	// 解析Path
	proxy.Path = config.Path

	// 解析SNI
	proxy.SNI = config.SNI

	// 解析ALPN
	if config.ALPN != "" {
		proxy.ALPN = strings.Split(config.ALPN, ",")
	}

	// 解析加密方法
	proxy.EncryptMethod = config.Security
	if proxy.EncryptMethod == "" {
		proxy.EncryptMethod = "auto"
	}

	return proxy, nil
}

// VLESSParser VLESS解析器
type VLESSParser struct{}

// CanParse 检查是否能解析指定的链接
func (p *VLESSParser) CanParse(link string) bool {
	return strings.HasPrefix(link, "vless://")
}

// GetType 获取解析器类型
func (p *VLESSParser) GetType() models.ProxyType {
	return models.ProxyTypeVLESS
}

// Parse 解析VLESS链接
func (p *VLESSParser) Parse(link string) (*models.Proxy, error) {
	if !p.CanParse(link) {
		return nil, fmt.Errorf("not a vless link")
	}

	// 解析URL
	u, err := url.Parse(link)
	if err != nil {
		return nil, fmt.Errorf("invalid vless URL: %w", err)
	}

	proxy := &models.Proxy{
		Type:     models.ProxyTypeVLESS,
		Group:    "VLESSProvider",
		Hostname: u.Hostname(),
		UserID:   u.User.Username(),
	}

	// 解析端口
	portStr := u.Port()
	if portStr == "" {
		return nil, fmt.Errorf("missing port")
	}
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}
	proxy.Port = uint16(port)

	// 解析查询参数
	query := u.Query()

	// 解析加密方法
	proxy.EncryptMethod = query.Get("encryption")
	if proxy.EncryptMethod == "" {
		proxy.EncryptMethod = "none"
	}

	// 解析传输协议
	proxy.TransferProtocol = query.Get("type")
	if proxy.TransferProtocol == "" {
		proxy.TransferProtocol = "tcp"
	}

	// 解析安全类型
	security := query.Get("security")
	if security == "tls" {
		proxy.TLS = "tls"
		proxy.TLSSecure = true
	} else if security == "reality" {
		proxy.TLS = "reality"
		proxy.TLSSecure = true
	}

	// 解析SNI
	proxy.SNI = query.Get("sni")

	// 解析Host
	proxy.Host = query.Get("host")

	// 解析Path
	proxy.Path = query.Get("path")

	// 解析流控
	proxy.Flow = query.Get("flow")

	// 解析指纹
	proxy.Fingerprint = query.Get("fp")

	// 解析ALPN
	if alpn := query.Get("alpn"); alpn != "" {
		proxy.ALPN = strings.Split(alpn, ",")
	}

	// 解析Short ID
	proxy.ShortID = query.Get("sid")

	// 解析Public Key
	proxy.PublicKey = query.Get("pbk")

	// 解析名称
	if fragment := u.Fragment; fragment != "" {
		decoded, err := url.QueryUnescape(fragment)
		if err == nil {
			proxy.Remark = decoded
		} else {
			proxy.Remark = fragment
		}
	}

	return proxy, nil
}

// parsePort 解析端口（支持字符串和数字）
func parsePort(port interface{}) (uint16, error) {
	switch v := port.(type) {
	case string:
		if v == "" {
			return 0, fmt.Errorf("empty port")
		}
		p, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return 0, err
		}
		return uint16(p), nil
	case float64:
		return uint16(v), nil
	case int:
		return uint16(v), nil
	default:
		return 0, fmt.Errorf("unsupported port type: %T", port)
	}
}

// parseAlterID 解析AlterID（支持字符串和数字）
func parseAlterID(alterID interface{}) (uint16, error) {
	switch v := alterID.(type) {
	case string:
		if v == "" {
			return 0, nil
		}
		aid, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return 0, err
		}
		return uint16(aid), nil
	case float64:
		return uint16(v), nil
	case int:
		return uint16(v), nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported alterid type: %T", alterID)
	}
}