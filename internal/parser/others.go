// internal/parser/others.go
package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"subconverter-go/pkg/models"
)

// TrojanParser Trojan解析器
type TrojanParser struct{}

// CanParse 检查是否能解析指定的链接
func (p *TrojanParser) CanParse(link string) bool {
	return strings.HasPrefix(link, "trojan://")
}

// GetType 获取解析器类型
func (p *TrojanParser) GetType() models.ProxyType {
	return models.ProxyTypeTrojan
}

// Parse 解析Trojan链接
func (p *TrojanParser) Parse(link string) (*models.Proxy, error) {
	if !p.CanParse(link) {
		return nil, fmt.Errorf("not a trojan link")
	}

	// 解析URL
	u, err := url.Parse(link)
	if err != nil {
		return nil, fmt.Errorf("invalid trojan URL: %w", err)
	}

	proxy := &models.Proxy{
		Type:     models.ProxyTypeTrojan,
		Group:    "TrojanProvider",
		Hostname: u.Hostname(),
		Password: u.User.Username(),
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

	// 解析传输协议
	proxy.TransferProtocol = query.Get("type")
	if proxy.TransferProtocol == "" {
		proxy.TransferProtocol = "tcp"
	}

	// 解析安全类型
	security := query.Get("security")
	if security == "tls" || security == "" {
		proxy.TLS = "tls"
		proxy.TLSSecure = true
	}

	// 解析SNI
	proxy.SNI = query.Get("sni")

	// 解析Host
	proxy.Host = query.Get("host")

	// 解析Path
	proxy.Path = query.Get("path")

	// 解析指纹
	proxy.Fingerprint = query.Get("fp")

	// 解析ALPN
	if alpn := query.Get("alpn"); alpn != "" {
		proxy.ALPN = strings.Split(alpn, ",")
	}

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

// HTTPParser HTTP代理解析器
type HTTPParser struct{}

// CanParse 检查是否能解析指定的链接
func (p *HTTPParser) CanParse(link string) bool {
	return strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://")
}

// GetType 获取解析器类型
func (p *HTTPParser) GetType() models.ProxyType {
	return models.ProxyTypeHTTP
}

// Parse 解析HTTP代理链接
func (p *HTTPParser) Parse(link string) (*models.Proxy, error) {
	if !p.CanParse(link) {
		return nil, fmt.Errorf("not an http/https link")
	}

	// 解析URL
	u, err := url.Parse(link)
	if err != nil {
		return nil, fmt.Errorf("invalid http URL: %w", err)
	}

	var proxyType models.ProxyType
	var group string
	if u.Scheme == "https" {
		proxyType = models.ProxyTypeHTTPS
		group = "HTTPSProvider"
	} else {
		proxyType = models.ProxyTypeHTTP
		group = "HTTPProvider"
	}

	proxy := &models.Proxy{
		Type:     proxyType,
		Group:    group,
		Hostname: u.Hostname(),
	}

	// 解析端口
	portStr := u.Port()
	if portStr != "" {
		port, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("invalid port: %w", err)
		}
		proxy.Port = uint16(port)
	} else {
		// 使用默认端口
		if u.Scheme == "https" {
			proxy.Port = 443
		} else {
			proxy.Port = 80
		}
	}

	// 解析认证信息
	if u.User != nil {
		proxy.Username = u.User.Username()
		if password, ok := u.User.Password(); ok {
			proxy.Password = password
		}
	}

	// 解析查询参数中的名称
	query := u.Query()
	if name := query.Get("name"); name != "" {
		proxy.Remark = name
	}

	// 解析fragment作为名称
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

// SOCKS5Parser SOCKS5代理解析器
type SOCKS5Parser struct{}

// CanParse 检查是否能解析指定的链接
func (p *SOCKS5Parser) CanParse(link string) bool {
	return strings.HasPrefix(link, "socks5://") || strings.HasPrefix(link, "socks://")
}

// GetType 获取解析器类型
func (p *SOCKS5Parser) GetType() models.ProxyType {
	return models.ProxyTypeSOCKS5
}

// Parse 解析SOCKS5代理链接
func (p *SOCKS5Parser) Parse(link string) (*models.Proxy, error) {
	if !p.CanParse(link) {
		return nil, fmt.Errorf("not a socks5 link")
	}

	// 解析URL
	u, err := url.Parse(link)
	if err != nil {
		return nil, fmt.Errorf("invalid socks5 URL: %w", err)
	}

	proxy := &models.Proxy{
		Type:     models.ProxyTypeSOCKS5,
		Group:    "SOCKS5Provider",
		Hostname: u.Hostname(),
	}

	// 解析端口
	portStr := u.Port()
	if portStr == "" {
		proxy.Port = 1080 // SOCKS5默认端口
	} else {
		port, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("invalid port: %w", err)
		}
		proxy.Port = uint16(port)
	}

	// 解析认证信息
	if u.User != nil {
		proxy.Username = u.User.Username()
		if password, ok := u.User.Password(); ok {
			proxy.Password = password
		}
	}

	// 解析查询参数中的名称
	query := u.Query()
	if name := query.Get("name"); name != "" {
		proxy.Remark = name
	}

	// 解析fragment作为名称
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