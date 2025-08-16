// Package models 定义了 SubConverter 的核心数据模型
package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"subconverter-go/pkg/constants"
	"subconverter-go/pkg/types"
)

// Proxy 代理节点模型 - 完全对应 C++ 版本的 Proxy struct
type Proxy struct {
	// 基础信息
	Type    types.ProxyType `json:"type" yaml:"type" validate:"required"`
	ID      uint32          `json:"id" yaml:"id"`
	GroupID uint32          `json:"group_id" yaml:"group_id"`
	Group   string          `json:"group" yaml:"group"`
	Remark  string          `json:"remark" yaml:"remark" validate:"required"`

	// 服务器信息
	Hostname string `json:"hostname" yaml:"hostname" validate:"required"`
	Port     uint16 `json:"port" yaml:"port" validate:"required,min=1,max=65535"`

	// 认证信息
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`

	// 加密和协议设置
	EncryptMethod string `json:"encrypt_method,omitempty" yaml:"encrypt_method,omitempty"`
	Plugin        string `json:"plugin,omitempty" yaml:"plugin,omitempty"`
	PluginOption  string `json:"plugin_option,omitempty" yaml:"plugin_option,omitempty"`
	Protocol      string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	ProtocolParam string `json:"protocol_param,omitempty" yaml:"protocol_param,omitempty"`
	OBFS          string `json:"obfs,omitempty" yaml:"obfs,omitempty"`
	OBFSParam     string `json:"obfs_param,omitempty" yaml:"obfs_param,omitempty"`

	// V2Ray/VMess 专用字段
	UserID           string `json:"user_id,omitempty" yaml:"user_id,omitempty"`
	AlterID          uint16 `json:"alter_id,omitempty" yaml:"alter_id,omitempty"`
	TransferProtocol string `json:"transfer_protocol,omitempty" yaml:"transfer_protocol,omitempty"`
	FakeType         string `json:"fake_type,omitempty" yaml:"fake_type,omitempty"`
	AuthStr          string `json:"auth_str,omitempty" yaml:"auth_str,omitempty"`

	// TLS 相关设置
	TLSStr    string `json:"tls_str,omitempty" yaml:"tls_str,omitempty"`
	TLSSecure bool   `json:"tls_secure" yaml:"tls_secure"`

	// 网络设置
	Host string `json:"host,omitempty" yaml:"host,omitempty"`
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
	Edge string `json:"edge,omitempty" yaml:"edge,omitempty"`

	// gRPC 相关
	GRPCServiceName string `json:"grpc_service_name,omitempty" yaml:"grpc_service_name,omitempty"`
	GRPCMode        string `json:"grpc_mode,omitempty" yaml:"grpc_mode,omitempty"`

	// VLESS 相关
	Flow           string `json:"flow,omitempty" yaml:"flow,omitempty"`
	FlowShow       bool   `json:"flow_show,omitempty" yaml:"flow_show,omitempty"`
	ShortID        string `json:"short_id,omitempty" yaml:"short_id,omitempty"`
	PacketEncoding string `json:"packet_encoding,omitempty" yaml:"packet_encoding,omitempty"`

	// Snell 专用字段
	SnellVersion uint16 `json:"snell_version,omitempty" yaml:"snell_version,omitempty"`
	ServerName   string `json:"server_name,omitempty" yaml:"server_name,omitempty"`

	// 其他扩展字段
	SNI                  string `json:"sni,omitempty" yaml:"sni,omitempty"`
	OBFSPassword         string `json:"obfs_password,omitempty" yaml:"obfs_password,omitempty"`

	// 特性开关 - 使用指针实现三态逻辑 (true/false/nil)
	UDP           *bool `json:"udp,omitempty" yaml:"udp,omitempty"`
	TCPFastOpen   *bool `json:"tfo,omitempty" yaml:"tfo,omitempty"`
	AllowInsecure *bool `json:"allow_insecure,omitempty" yaml:"allow_insecure,omitempty"`

	// 元数据
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
}

// IsValid 验证代理节点是否有效
func (p *Proxy) IsValid() bool {
	if p == nil {
		return false
	}
	if !p.Type.IsValid() {
		return false
	}
	if p.Hostname == "" || p.Port == 0 {
		return false
	}
	if p.Remark == "" {
		return false
	}

	// 根据代理类型验证必要字段
	switch p.Type {
	case types.ProxyTypeShadowsocks, types.ProxyTypeShadowsocksR:
		if p.Password == "" || p.EncryptMethod == "" {
			return false
		}
	case types.ProxyTypeVMess, types.ProxyTypeVLESS:
		if p.UserID == "" {
			return false
		}
	case types.ProxyTypeTrojan:
		if p.Password == "" {
			return false
		}
	}

	return true
}

// GetDefaultGroup 获取代理的默认分组
func (p *Proxy) GetDefaultGroup() string {
	if p.Group != "" {
		return p.Group
	}

	switch p.Type {
	case types.ProxyTypeShadowsocks:
		return constants.SSDefaultGroup
	case types.ProxyTypeShadowsocksR:
		return constants.SSRDefaultGroup
	case types.ProxyTypeVMess, types.ProxyTypeVLESS:
		return constants.V2RayDefaultGroup
	case types.ProxyTypeTrojan:
		return constants.TrojanDefaultGroup
	case types.ProxyTypeSnell:
		return constants.SnellDefaultGroup
	case types.ProxyTypeHTTP, types.ProxyTypeHTTPS:
		return constants.HTTPDefaultGroup
	case types.ProxyTypeSOCKS5:
		return constants.SocksDefaultGroup
	default:
		return "Unknown"
	}
}

// Clone 深拷贝代理对象
func (p *Proxy) Clone() *Proxy {
	if p == nil {
		return nil
	}

	clone := *p

	// 深拷贝指针字段
	if p.UDP != nil {
		udp := *p.UDP
		clone.UDP = &udp
	}
	if p.TCPFastOpen != nil {
		tfo := *p.TCPFastOpen
		clone.TCPFastOpen = &tfo
	}
	if p.AllowInsecure != nil {
		insecure := *p.AllowInsecure
		clone.AllowInsecure = &insecure
	}

	return &clone
}

// String 返回代理的字符串表示
func (p *Proxy) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s://%s:%d [%s]", p.Type.String(), p.Hostname, p.Port, p.Remark)
}

// ToJSON 将代理转换为 JSON
func (p *Proxy) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// FromJSON 从 JSON 创建代理
func FromJSON(data []byte) (*Proxy, error) {
	var proxy Proxy
	err := json.Unmarshal(data, &proxy)
	return &proxy, err
}

// GetKey 获取代理的唯一标识符
func (p *Proxy) GetKey() string {
	return fmt.Sprintf("%s:%d:%s", p.Hostname, p.Port, p.Type.String())
}

// SetDefaults 设置默认值
func (p *Proxy) SetDefaults() {
	if p.Port == 0 {
		p.Port = p.Type.GetDefaultPort()
	}
	if p.Group == "" {
		p.Group = p.GetDefaultGroup()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	p.UpdatedAt = time.Now()
}

// HasFeature 检查代理是否支持某个特性
func (p *Proxy) HasFeature(feature string) bool {
	switch strings.ToLower(feature) {
	case "udp":
		return p.Type.SupportsUDP()
	case "tls":
		return p.Type.SupportsTLS()
	default:
		return false
	}
}

// GetDisplayName 获取显示名称
func (p *Proxy) GetDisplayName() string {
	if p.Remark != "" {
		return p.Remark
	}
	return fmt.Sprintf("%s:%d", p.Hostname, p.Port)
}

// IsSecure 检查代理是否使用安全连接
func (p *Proxy) IsSecure() bool {
	return p.TLSSecure || p.TLSStr == "tls" || p.Type.SupportsTLS()
}

// ProxyList 代理列表类型
type ProxyList []*Proxy

// Len 返回代理列表长度
func (pl ProxyList) Len() int {
	return len(pl)
}

// FilterByType 按类型过滤代理
func (pl ProxyList) FilterByType(proxyType types.ProxyType) ProxyList {
	var filtered ProxyList
	for _, proxy := range pl {
		if proxy.Type == proxyType {
			filtered = append(filtered, proxy)
		}
	}
	return filtered
}

// FilterByGroup 按分组过滤代理
func (pl ProxyList) FilterByGroup(group string) ProxyList {
	var filtered ProxyList
	for _, proxy := range pl {
		if proxy.GetDefaultGroup() == group {
			filtered = append(filtered, proxy)
		}
	}
	return filtered
}

// GroupByType 按类型分组代理
func (pl ProxyList) GroupByType() map[types.ProxyType]ProxyList {
	groups := make(map[types.ProxyType]ProxyList)
	for _, proxy := range pl {
		groups[proxy.Type] = append(groups[proxy.Type], proxy)
	}
	return groups
}

// GetRemarks 获取所有备注名称
func (pl ProxyList) GetRemarks() []string {
	remarks := make([]string, len(pl))
	for i, proxy := range pl {
		remarks[i] = proxy.Remark
	}
	return remarks
}

// Validate 验证所有代理
func (pl ProxyList) Validate() []error {
	var errors []error
	for i, proxy := range pl {
		if !proxy.IsValid() {
			errors = append(errors, fmt.Errorf("proxy[%d] is invalid: %s", i, proxy.String()))
		}
	}
	return errors
}