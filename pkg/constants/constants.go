// Package constants 定义了 SubConverter 中使用的所有常量
package constants

import "time"

// 代理协议默认端口
const (
	DefaultSSPort       uint16 = 443
	DefaultSSRPort      uint16 = 443
	DefaultVMessPort    uint16 = 443
	DefaultVLESSPort    uint16 = 443
	DefaultTrojanPort   uint16 = 443
	DefaultHysteriaPort uint16 = 443
	DefaultTUICPort     uint16 = 443
	DefaultHTTPPort     uint16 = 80
	DefaultHTTPSPort    uint16 = 443
	DefaultSOCKS5Port   uint16 = 1080
	DefaultSnellPort    uint16 = 6160
	DefaultWGPort       uint16 = 51820
)

// 默认分组名称 - 对应 C++ 版本的宏定义
const (
	SSDefaultGroup        = "SSProvider"
	SSRDefaultGroup       = "SSRProvider"
	V2RayDefaultGroup     = "V2RayProvider"
	SocksDefaultGroup     = "SocksProvider"
	HTTPDefaultGroup      = "HTTPProvider"
	TrojanDefaultGroup    = "TrojanProvider"
	SnellDefaultGroup     = "SnellProvider"
	WGDefaultGroup        = "WireGuardProvider"
	XRayDefaultGroup      = "XRayProvider"
	HysteriaDefaultGroup  = "HysteriaProvider"
	Hysteria2DefaultGroup = "Hysteria2Provider"
	TUICDefaultGroup      = "TuicProvider"
	AnyTLSDefaultGroup    = "AnyTLSProvider"
	MieruDefaultGroup     = "MieruProvider"
)

// 支持的目标客户端类型 - 与 C++ 版本保持一致
var SupportedTargets = []string{
	"clash", "clashr", "surge", "quan", "quanx",
	"loon", "ss", "ssr", "v2ray", "trojan", "singbox",
	"auto", "mixed", "surfboard", "mellow",
}

// 目标客户端显示名称映射
var TargetDisplayNames = map[string]string{
	"clash":     "Clash",
	"clashr":    "Clash.R",
	"surge":     "Surge",
	"quan":      "Quantumult",
	"quanx":     "QuantumultX",
	"loon":      "Loon",
	"ss":        "Shadowsocks",
	"ssr":       "ShadowsocksR",
	"v2ray":     "V2Ray",
	"trojan":    "Trojan",
	"singbox":   "SingBox",
	"surfboard": "Surfboard",
	"mellow":    "Mellow",
	"auto":      "Auto Detect",
	"mixed":     "Mixed",
}

// HTTP 头部常量
const (
	ContentTypeYAML          = "application/x-yaml"
	ContentTypeJSON          = "application/json"
	ContentTypeText          = "text/plain"
	ContentTypeConf          = "application/octet-stream"
	ContentTypeFormData      = "application/x-www-form-urlencoded"
	ContentTypeMultipart     = "multipart/form-data"
	HeaderContentType        = "Content-Type"
	HeaderContentLength      = "Content-Length"
	HeaderContentDisposition = "Content-Disposition"
	HeaderUserAgent          = "User-Agent"
	HeaderAuthorization      = "Authorization"
	HeaderAccessToken        = "Access-Token"
)

// 配置文件格式
const (
	ConfigFormatYAML = "yaml"
	ConfigFormatTOML = "toml"
	ConfigFormatINI  = "ini"
	ConfigFormatJSON = "json"
)

// 模板文件扩展名
const (
	TemplateExtYAML = ".yaml"
	TemplateExtTOML = ".toml"
	TemplateExtINI  = ".ini"
	TemplateExtJSON = ".json"
	TemplateExtConf = ".conf"
	TemplateExtTpl  = ".tpl"
	TemplateExtJ2   = ".j2"
)

// 默认配置值
const (
	DefaultListenAddress    = "0.0.0.0"
	DefaultListenPort       = 25500
	DefaultMaxPendingConns  = 10240
	DefaultMaxConcurThreads = 4
	DefaultRulesetInterval  = 86400 // 24小时
	DefaultTimeout          = 15    // 15秒
	DefaultRequestTimeout   = 30    // 30秒
	DefaultReadTimeout      = 30    // 30秒
	DefaultWriteTimeout     = 30    // 30秒
)

// 缓存相关常量
const (
	DefaultCacheTTL     = 10 * time.Minute
	DefaultMaxCacheSize = 1000
	DefaultCleanupTime  = 1 * time.Minute
	MaxSubscriptionSize = 10 * 1024 * 1024 // 10MB
	MaxConfigSize       = 5 * 1024 * 1024  // 5MB
	MaxRulesetSize      = 1 * 1024 * 1024  // 1MB
)

// URL 验证相关
const (
	MaxURLLength       = 8192
	MinURLLength       = 7 // http://
	MaxSubscriptionURL = 10
	MaxConfigURL       = 5
)

// 代理节点限制
const (
	MaxProxyNodes     = 10000
	MaxProxyGroups    = 100
	MaxRulesets       = 100
	MaxIncludeFilters = 20
	MaxExcludeFilters = 20
)

// 正则表达式常量
const (
	// URL 格式验证
	URLPattern = `^https?://[^\s/$.?#].[^\s]*$`

	// IPv4 地址验证
	IPv4Pattern = `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

	// IPv6 地址验证
	IPv6Pattern = `^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`

	// 域名验证
	DomainPattern = `^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`

	// 端口验证
	PortPattern = `^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`

	// UUID 验证
	UUIDPattern = `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`

	// Base64 验证
	Base64Pattern = `^[A-Za-z0-9+/]*={0,2}$`
)

// 协议前缀常量
const (
	PrefixSS       = "ss://"
	PrefixSSR      = "ssr://"
	PrefixVMess    = "vmess://"
	PrefixVLESS    = "vless://"
	PrefixTrojan   = "trojan://"
	PrefixHysteria = "hysteria://"
	PrefixHy2      = "hysteria2://"
	PrefixTUIC     = "tuic://"
	PrefixHTTP     = "http://"
	PrefixHTTPS    = "https://"
	PrefixSOCKS    = "socks://"
	PrefixSOCKS5   = "socks5://"
)

// 加密算法常量
var SupportedCiphers = map[string][]string{
	"shadowsocks": {
		"aes-128-gcm", "aes-192-gcm", "aes-256-gcm",
		"aes-128-cfb", "aes-192-cfb", "aes-256-cfb",
		"aes-128-ctr", "aes-192-ctr", "aes-256-ctr",
		"chacha20-ietf", "xchacha20-ietf",
		"chacha20-ietf-poly1305", "xchacha20-ietf-poly1305",
		"2022-blake3-aes-128-gcm", "2022-blake3-aes-256-gcm",
		"2022-blake3-chacha20-poly1305",
	},
	"shadowsocksr": {
		"aes-128-cfb", "aes-192-cfb", "aes-256-cfb",
		"aes-128-ctr", "aes-192-ctr", "aes-256-ctr",
		"rc4-md5", "chacha20", "chacha20-ietf",
	},
	"vmess": {
		"auto", "aes-128-gcm", "chacha20-poly1305", "none",
	},
}

// 支持的协议方法
var SupportedProtocols = map[string][]string{
	"shadowsocksr": {
		"origin", "auth_sha1_v4", "auth_aes128_md5", "auth_aes128_sha1",
		"auth_chain_a", "auth_chain_b", "auth_chain_c", "auth_chain_d",
	},
}

// 支持的混淆方法
var SupportedObfs = map[string][]string{
	"shadowsocksr": {
		"plain", "http_simple", "http_post", "random_head", "tls1.2_ticket_auth",
		"tls1.2_ticket_fastauth",
	},
}

// 应用相关常量
const (
	AppName        = "SubConverter Go"
	AppDescription = "High-performance subscription converter written in Go"
	DefaultAgent   = "SubConverter-Go/1.0"
)

// 日志相关常量
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogFormatText = "text"
	LogFormatJSON = "json"
)

// 缓存键前缀
const (
	CacheKeySubscription = "sub:"
	CacheKeyConfig       = "config:"
	CacheKeyRuleset      = "ruleset:"
	CacheKeyTemplate     = "template:"
)

// 监控相关常量
const (
	MetricsPath = "/metrics"
	HealthPath  = "/health"
	ReadyPath   = "/ready"
)