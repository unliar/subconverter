// Package models 定义配置相关的数据模型
package models

import (
	"fmt"
	"subconverter-go/pkg/constants"
	"time"
)

// ServerConfig 服务器配置 - 对应 C++ 版本的服务器设置
type ServerConfig struct {
	ListenAddress    string `json:"listen_address" yaml:"listen_address" validate:"required"`
	ListenPort       int    `json:"listen_port" yaml:"listen_port" validate:"required,min=1,max=65535"`
	MaxPendingConns  int    `json:"max_pending_conns" yaml:"max_pending_conns"`
	MaxConcurThreads int    `json:"max_concur_threads" yaml:"max_concur_threads"`
	APIMode          bool   `json:"api_mode" yaml:"api_mode"`
	AccessToken      string `json:"access_token,omitempty" yaml:"access_token,omitempty"`

	// 超时设置
	RequestTimeout time.Duration `json:"request_timeout" yaml:"request_timeout"`
	ReadTimeout    time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout" yaml:"write_timeout"`

	// 元数据
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
}

// IsValid 验证服务器配置是否有效
func (sc *ServerConfig) IsValid() bool {
	if sc == nil {
		return false
	}
	if sc.ListenAddress == "" {
		return false
	}
	if sc.ListenPort <= 0 || sc.ListenPort > 65535 {
		return false
	}
	return true
}

// SetDefaults 设置默认值
func (sc *ServerConfig) SetDefaults() {
	if sc.ListenAddress == "" {
		sc.ListenAddress = constants.DefaultListenAddress
	}
	if sc.ListenPort <= 0 {
		sc.ListenPort = constants.DefaultListenPort
	}
	if sc.MaxPendingConns <= 0 {
		sc.MaxPendingConns = constants.DefaultMaxPendingConns
	}
	if sc.MaxConcurThreads <= 0 {
		sc.MaxConcurThreads = constants.DefaultMaxConcurThreads
	}
	if sc.RequestTimeout <= 0 {
		sc.RequestTimeout = time.Duration(constants.DefaultRequestTimeout) * time.Second
	}
	if sc.ReadTimeout <= 0 {
		sc.ReadTimeout = time.Duration(constants.DefaultReadTimeout) * time.Second
	}
	if sc.WriteTimeout <= 0 {
		sc.WriteTimeout = time.Duration(constants.DefaultWriteTimeout) * time.Second
	}
	if sc.CreatedAt.IsZero() {
		sc.CreatedAt = time.Now()
	}
	sc.UpdatedAt = time.Now()
}

// GetAddress 获取完整的监听地址
func (sc *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", sc.ListenAddress, sc.ListenPort)
}

// Clone 深拷贝服务器配置
func (sc *ServerConfig) Clone() *ServerConfig {
	if sc == nil {
		return nil
	}
	
	clone := *sc
	return &clone
}

// ConverterConfig 转换器配置
type ConverterConfig struct {
	DefaultConfig            string                `json:"default_config" yaml:"default_config"`
	EnableRuleGenerator      bool                  `json:"enable_rule_generator" yaml:"enable_rule_generator"`
	OverwriteOriginalRules   bool                  `json:"overwrite_original_rules" yaml:"overwrite_original_rules"`
	AddEmoji                 bool                  `json:"add_emoji" yaml:"add_emoji"`
	RemoveEmoji              bool                  `json:"remove_emoji" yaml:"remove_emoji"`
	AppendProxyType          bool                  `json:"append_proxy_type" yaml:"append_proxy_type"`
	FilterDeprecated         bool                  `json:"filter_deprecated" yaml:"filter_deprecated"`
	SortFlag                 bool                  `json:"sort_flag" yaml:"sort_flag"`
	ClashNewFieldName        bool                  `json:"clash_new_field_name" yaml:"clash_new_field_name"`
	ClashScript              bool                  `json:"clash_script" yaml:"clash_script"`
	ClashClassicalRuleset    bool                  `json:"clash_classical_ruleset" yaml:"clash_classical_ruleset"`

	// 路径设置
	SurgeSSRPath             string                `json:"surge_ssr_path" yaml:"surge_ssr_path"`
	ManagedConfigPrefix      string                `json:"managed_config_prefix" yaml:"managed_config_prefix"`
	QuanXDevID               string                `json:"quanx_dev_id" yaml:"quanx_dev_id"`

	// 样式设置
	ClashProxiesStyle        string                `json:"clash_proxies_style" yaml:"clash_proxies_style"`
	ClashProxyGroupsStyle    string                `json:"clash_proxy_groups_style" yaml:"clash_proxy_groups_style"`

	// 脚本设置
	SortScript               string                `json:"sort_script" yaml:"sort_script"`

	// 特性开关
	UDP                      *bool                 `json:"udp,omitempty" yaml:"udp,omitempty"`
	TFO                      *bool                 `json:"tfo,omitempty" yaml:"tfo,omitempty"`
	XUDP                     *bool                 `json:"xudp,omitempty" yaml:"xudp,omitempty"`
	SkipCertVerify           *bool                 `json:"skip_cert_verify,omitempty" yaml:"skip_cert_verify,omitempty"`
	TLS13                    *bool                 `json:"tls13,omitempty" yaml:"tls13,omitempty"`

	// 配置列表
	CustomRulesets           []RulesetConfig       `json:"custom_rulesets" yaml:"custom_rulesets"`
	ProxyGroups              []ProxyGroupConfig    `json:"proxy_groups" yaml:"proxy_groups"`

	// 元数据
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
}

// IsValid 验证转换器配置是否有效
func (cc *ConverterConfig) IsValid() bool {
	// 验证规则集配置
	for _, ruleset := range cc.CustomRulesets {
		if !ruleset.IsValid() {
			return false
		}
	}
	
	// 验证代理组配置
	for _, group := range cc.ProxyGroups {
		if !group.IsValid() {
			return false
		}
	}
	
	return true
}

// SetDefaults 设置默认值
func (cc *ConverterConfig) SetDefaults() {
	if cc.ClashProxiesStyle == "" {
		cc.ClashProxiesStyle = "block"
	}
	if cc.ClashProxyGroupsStyle == "" {
		cc.ClashProxyGroupsStyle = "block"
	}
	if cc.CreatedAt.IsZero() {
		cc.CreatedAt = time.Now()
	}
	cc.UpdatedAt = time.Now()
}

// AddRuleset 添加规则集
func (cc *ConverterConfig) AddRuleset(ruleset RulesetConfig) {
	cc.CustomRulesets = append(cc.CustomRulesets, ruleset)
	cc.UpdatedAt = time.Now()
}

// AddProxyGroup 添加代理组
func (cc *ConverterConfig) AddProxyGroup(group ProxyGroupConfig) {
	cc.ProxyGroups = append(cc.ProxyGroups, group)
	cc.UpdatedAt = time.Now()
}

// Clone 深拷贝转换器配置
func (cc *ConverterConfig) Clone() *ConverterConfig {
	if cc == nil {
		return nil
	}

	clone := *cc

	// 深拷贝切片
	if cc.CustomRulesets != nil {
		clone.CustomRulesets = make([]RulesetConfig, len(cc.CustomRulesets))
		copy(clone.CustomRulesets, cc.CustomRulesets)
	}
	if cc.ProxyGroups != nil {
		clone.ProxyGroups = make([]ProxyGroupConfig, len(cc.ProxyGroups))
		copy(clone.ProxyGroups, cc.ProxyGroups)
	}

	// 深拷贝指针字段
	if cc.UDP != nil {
		udp := *cc.UDP
		clone.UDP = &udp
	}
	if cc.TFO != nil {
		tfo := *cc.TFO
		clone.TFO = &tfo
	}
	if cc.XUDP != nil {
		xudp := *cc.XUDP
		clone.XUDP = &xudp
	}
	if cc.SkipCertVerify != nil {
		scv := *cc.SkipCertVerify
		clone.SkipCertVerify = &scv
	}
	if cc.TLS13 != nil {
		tls13 := *cc.TLS13
		clone.TLS13 = &tls13
	}

	return &clone
}

// TemplateConfig 模板配置
type TemplateConfig struct {
	Name        string            `json:"name" yaml:"name" validate:"required"`
	Path        string            `json:"path" yaml:"path" validate:"required"`
	Type        string            `json:"type" yaml:"type"`
	Variables   map[string]string `json:"variables,omitempty" yaml:"variables,omitempty"`
	Functions   []string          `json:"functions,omitempty" yaml:"functions,omitempty"`

	// 元数据
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
}

// IsValid 验证模板配置是否有效
func (tc *TemplateConfig) IsValid() bool {
	return tc != nil && tc.Name != "" && tc.Path != ""
}

// SetDefaults 设置默认值
func (tc *TemplateConfig) SetDefaults() {
	if tc.Type == "" {
		tc.Type = "jinja2"
	}
	if tc.CreatedAt.IsZero() {
		tc.CreatedAt = time.Now()
	}
	tc.UpdatedAt = time.Now()
}

// GetVariable 获取变量值
func (tc *TemplateConfig) GetVariable(key string) (string, bool) {
	if tc.Variables == nil {
		return "", false
	}
	value, exists := tc.Variables[key]
	return value, exists
}

// SetVariable 设置变量值
func (tc *TemplateConfig) SetVariable(key, value string) {
	if tc.Variables == nil {
		tc.Variables = make(map[string]string)
	}
	tc.Variables[key] = value
	tc.UpdatedAt = time.Now()
}

// Clone 深拷贝模板配置
func (tc *TemplateConfig) Clone() *TemplateConfig {
	if tc == nil {
		return nil
	}

	clone := *tc

	// 深拷贝 map
	if tc.Variables != nil {
		clone.Variables = make(map[string]string)
		for k, v := range tc.Variables {
			clone.Variables[k] = v
		}
	}

	// 深拷贝切片
	if tc.Functions != nil {
		clone.Functions = make([]string, len(tc.Functions))
		copy(clone.Functions, tc.Functions)
	}

	return &clone
}

// ApplicationConfig 应用配置 - 包含所有配置
type ApplicationConfig struct {
	Server    *ServerConfig    `json:"server" yaml:"server"`
	Converter *ConverterConfig `json:"converter" yaml:"converter"`
	Templates []*TemplateConfig `json:"templates" yaml:"templates"`

	// 日志配置
	LogLevel  string `json:"log_level" yaml:"log_level"`
	LogFormat string `json:"log_format" yaml:"log_format"`

	// 缓存配置
	CacheEnabled bool          `json:"cache_enabled" yaml:"cache_enabled"`
	CacheTTL     time.Duration `json:"cache_ttl" yaml:"cache_ttl"`
	MaxCacheSize int           `json:"max_cache_size" yaml:"max_cache_size"`

	// 元数据
	Version   string    `json:"version,omitempty" yaml:"version,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
}

// IsValid 验证应用配置是否有效
func (ac *ApplicationConfig) IsValid() bool {
	if ac.Server != nil && !ac.Server.IsValid() {
		return false
	}
	if ac.Converter != nil && !ac.Converter.IsValid() {
		return false
	}
	for _, template := range ac.Templates {
		if !template.IsValid() {
			return false
		}
	}
	return true
}

// SetDefaults 设置默认值
func (ac *ApplicationConfig) SetDefaults() {
	if ac.Server == nil {
		ac.Server = &ServerConfig{}
	}
	ac.Server.SetDefaults()

	if ac.Converter == nil {
		ac.Converter = &ConverterConfig{}
	}
	ac.Converter.SetDefaults()

	if ac.LogLevel == "" {
		ac.LogLevel = constants.LogLevelInfo
	}
	if ac.LogFormat == "" {
		ac.LogFormat = constants.LogFormatText
	}
	if ac.CacheTTL <= 0 {
		ac.CacheTTL = constants.DefaultCacheTTL
	}
	if ac.MaxCacheSize <= 0 {
		ac.MaxCacheSize = constants.DefaultMaxCacheSize
	}
	if ac.CreatedAt.IsZero() {
		ac.CreatedAt = time.Now()
	}
	ac.UpdatedAt = time.Now()
}

// Clone 深拷贝应用配置
func (ac *ApplicationConfig) Clone() *ApplicationConfig {
	if ac == nil {
		return nil
	}

	clone := *ac

	// 深拷贝嵌套结构
	if ac.Server != nil {
		clone.Server = ac.Server.Clone()
	}
	if ac.Converter != nil {
		clone.Converter = ac.Converter.Clone()
	}
	if ac.Templates != nil {
		clone.Templates = make([]*TemplateConfig, len(ac.Templates))
		for i, template := range ac.Templates {
			clone.Templates[i] = template.Clone()
		}
	}

	return &clone
}