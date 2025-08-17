// pkg/models/proxy.go
package models

import (
	"fmt"
	"net/url"
	"strings"
)

// ProxyType 代理类型枚举
type ProxyType int

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

// String 返回代理类型的字符串表示
func (t ProxyType) String() string {
	switch t {
	case ProxyTypeShadowsocks:
		return "SS"
	case ProxyTypeShadowsocksR:
		return "SSR"
	case ProxyTypeVMess:
		return "VMess"
	case ProxyTypeTrojan:
		return "Trojan"
	case ProxyTypeSnell:
		return "Snell"
	case ProxyTypeHTTP:
		return "HTTP"
	case ProxyTypeHTTPS:
		return "HTTPS"
	case ProxyTypeSOCKS5:
		return "SOCKS5"
	case ProxyTypeWireGuard:
		return "WireGuard"
	case ProxyTypeVLESS:
		return "VLESS"
	case ProxyTypeHysteria:
		return "Hysteria"
	case ProxyTypeHysteria2:
		return "Hysteria2"
	case ProxyTypeTUIC:
		return "TUIC"
	case ProxyTypeAnyTLS:
		return "AnyTLS"
	case ProxyTypeMieru:
		return "Mieru"
	default:
		return "Unknown"
	}
}

// Tribool 三态布尔值
type Tribool *bool

func True() Tribool  { t := true; return &t }
func False() Tribool { f := false; return &f }
func Nil() Tribool   { return nil }

// Proxy 代理节点结构
type Proxy struct {
	// 基础信息
	ID       uint32    `json:"id"`
	GroupID  uint32    `json:"group_id"`
	Type     ProxyType `json:"type"`
	Group    string    `json:"group"`
	Remark   string    `json:"remark"`
	Hostname string    `json:"hostname"`
	Port     uint16    `json:"port"`

	// 认证信息
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	AlterID  uint16 `json:"alter_id,omitempty"`

	// 加密信息
	EncryptMethod string `json:"encrypt_method,omitempty"`
	Plugin        string `json:"plugin,omitempty"`
	PluginOption  string `json:"plugin_option,omitempty"`

	// SSR 特有
	Protocol      string `json:"protocol,omitempty"`
	ProtocolParam string `json:"protocol_param,omitempty"`
	OBFS          string `json:"obfs,omitempty"`
	OBFSParam     string `json:"obfs_param,omitempty"`
	OBFSPassword  string `json:"obfs_password,omitempty"`

	// 传输协议
	TransferProtocol string `json:"transfer_protocol,omitempty"`
	FakeType         string `json:"fake_type,omitempty"`
	
	// TLS 相关
	TLS           string   `json:"tls,omitempty"`
	TLSSecure     bool     `json:"tls_secure"`
	ServerName    string   `json:"server_name,omitempty"`
	SNI           string   `json:"sni,omitempty"`
	ALPN          []string `json:"alpn,omitempty"`
	Fingerprint   string   `json:"fingerprint,omitempty"`
	AllowInsecure Tribool  `json:"allow_insecure,omitempty"`
	TLS13         Tribool  `json:"tls13,omitempty"`

	// WebSocket/HTTP2/gRPC 相关
	Host            string `json:"host,omitempty"`
	Path            string `json:"path,omitempty"`
	Edge            string `json:"edge,omitempty"`
	GRPCServiceName string `json:"grpc_service_name,omitempty"`
	GRPCMode        string `json:"grpc_mode,omitempty"`

	// 网络优化
	UDP         Tribool `json:"udp,omitempty"`
	XUDP        Tribool `json:"xudp,omitempty"`
	TCPFastOpen Tribool `json:"tcp_fast_open,omitempty"`

	// 其他协议特有字段（根据需要扩展）
	// Hysteria, TUIC, VLESS, WireGuard 等特有字段
	AuthStr           string `json:"auth_str,omitempty"`
	Flow              string `json:"flow,omitempty"`
	ShortID           string `json:"short_id,omitempty"`
	UpMbps            string `json:"up_mbps,omitempty"`
	DownMbps          string `json:"down_mbps,omitempty"`
	CongestionControl string `json:"congestion_control,omitempty"`
	UDPRelayMode      string `json:"udp_relay_mode,omitempty"`
	PublicKey         string `json:"public_key,omitempty"`
	PrivateKey        string `json:"private_key,omitempty"`
	UnderlyingProxy   string `json:"underlying_proxy,omitempty"`
}

// ProxyList 代理节点列表
type ProxyList []*Proxy

// Clone 克隆代理节点
func (p *Proxy) Clone() *Proxy {
	if p == nil {
		return nil
	}
	
	clone := *p
	
	// 深拷贝切片
	if p.ALPN != nil {
		clone.ALPN = make([]string, len(p.ALPN))
		copy(clone.ALPN, p.ALPN)
	}
	
	return &clone
}

// GetProxyTypeName 根据字符串获取代理类型
func GetProxyTypeName(name string) ProxyType {
	switch strings.ToLower(name) {
	case "ss", "shadowsocks":
		return ProxyTypeShadowsocks
	case "ssr", "shadowsocksr":
		return ProxyTypeShadowsocksR
	case "vmess":
		return ProxyTypeVMess
	case "vless":
		return ProxyTypeVLESS
	case "trojan":
		return ProxyTypeTrojan
	case "snell":
		return ProxyTypeSnell
	case "http":
		return ProxyTypeHTTP
	case "https":
		return ProxyTypeHTTPS
	case "socks5", "socks":
		return ProxyTypeSOCKS5
	case "wireguard", "wg":
		return ProxyTypeWireGuard
	case "hysteria":
		return ProxyTypeHysteria
	case "hysteria2", "hy2":
		return ProxyTypeHysteria2
	case "tuic":
		return ProxyTypeTUIC
	case "anytls":
		return ProxyTypeAnyTLS
	case "mieru":
		return ProxyTypeMieru
	default:
		return ProxyTypeUnknown
	}
}

// Validate 验证代理节点配置
func (p *Proxy) Validate() error {
	if p == nil {
		return fmt.Errorf("proxy is nil")
	}
	
	if p.Hostname == "" {
		return fmt.Errorf("hostname is required")
	}
	
	if p.Port == 0 {
		return fmt.Errorf("port is required")
	}
	
	if p.Type == ProxyTypeUnknown {
		return fmt.Errorf("proxy type is unknown")
	}
	
	// 根据不同类型验证必需字段
	switch p.Type {
	case ProxyTypeShadowsocks:
		if p.Password == "" {
			return fmt.Errorf("shadowsocks password is required")
		}
		if p.EncryptMethod == "" {
			return fmt.Errorf("shadowsocks encrypt method is required")
		}
	case ProxyTypeShadowsocksR:
		if p.Password == "" {
			return fmt.Errorf("shadowsocksr password is required")
		}
		if p.EncryptMethod == "" {
			return fmt.Errorf("shadowsocksr encrypt method is required")
		}
	case ProxyTypeVMess, ProxyTypeVLESS:
		if p.UserID == "" {
			return fmt.Errorf("vmess/vless user id is required")
		}
	case ProxyTypeTrojan:
		if p.Password == "" {
			return fmt.Errorf("trojan password is required")
		}
	case ProxyTypeHTTP, ProxyTypeHTTPS:
		// HTTP/HTTPS 可能有用户名和密码，但不是必需的
	case ProxyTypeSOCKS5:
		// SOCKS5 可能有用户名和密码，但不是必需的
	}
	
	return nil
}

// ToURL 将代理节点转换为URL格式
func (p *Proxy) ToURL() (string, error) {
	if err := p.Validate(); err != nil {
		return "", err
	}
	
	var scheme string
	var userinfo *url.Userinfo
	
	switch p.Type {
	case ProxyTypeShadowsocks:
		scheme = "ss"
		if p.Password != "" && p.EncryptMethod != "" {
			userinfo = url.UserPassword(p.EncryptMethod, p.Password)
		}
	case ProxyTypeShadowsocksR:
		scheme = "ssr"
		// SSR URL格式较复杂，这里简化处理
		if p.Password != "" && p.EncryptMethod != "" {
			userinfo = url.UserPassword(p.EncryptMethod, p.Password)
		}
	case ProxyTypeVMess:
		scheme = "vmess"
		if p.UserID != "" {
			userinfo = url.User(p.UserID)
		}
	case ProxyTypeVLESS:
		scheme = "vless"
		if p.UserID != "" {
			userinfo = url.User(p.UserID)
		}
	case ProxyTypeTrojan:
		scheme = "trojan"
		if p.Password != "" {
			userinfo = url.User(p.Password)
		}
	case ProxyTypeHTTP:
		scheme = "http"
		if p.Username != "" && p.Password != "" {
			userinfo = url.UserPassword(p.Username, p.Password)
		}
	case ProxyTypeHTTPS:
		scheme = "https"
		if p.Username != "" && p.Password != "" {
			userinfo = url.UserPassword(p.Username, p.Password)
		}
	case ProxyTypeSOCKS5:
		scheme = "socks5"
		if p.Username != "" && p.Password != "" {
			userinfo = url.UserPassword(p.Username, p.Password)
		}
	default:
		return "", fmt.Errorf("unsupported proxy type for URL conversion: %s", p.Type.String())
	}
	
	u := &url.URL{
		Scheme: scheme,
		User:   userinfo,
		Host:   fmt.Sprintf("%s:%d", p.Hostname, p.Port),
	}
	
	// 添加查询参数
	query := url.Values{}
	if p.Remark != "" {
		query.Set("name", p.Remark)
	}
	if p.Group != "" {
		query.Set("group", p.Group)
	}
	
	u.RawQuery = query.Encode()
	
	return u.String(), nil
}