// internal/generator/generator.go
package generator

import (
	"fmt"
	"sort"
	"strings"

	"subconverter-go/pkg/models"
)

// Generator 生成器接口
type Generator interface {
	// Generate 生成配置内容
	Generate(proxies []*models.Proxy, config *GenerateConfig) ([]byte, error)
	
	// GetTarget 获取目标客户端类型
	GetTarget() string
	
	// GetFormat 获取输出格式
	GetFormat() string
	
	// Validate 验证代理节点是否兼容该生成器
	Validate(proxy *models.Proxy) error
	
	// SupportsProxyType 检查是否支持指定的代理类型
	SupportsProxyType(proxyType models.ProxyType) bool
}

// GenerateConfig 生成配置
type GenerateConfig struct {
	Target     string `json:"target"`     // 目标客户端
	Format     string `json:"format"`     // 输出格式
	Include    string `json:"include"`    // 包含规则
	Exclude    string `json:"exclude"`    // 排除规则
	Sort       bool   `json:"sort"`       // 是否排序
	UDP        bool   `json:"udp"`        // 是否启用UDP
	TLS13      bool   `json:"tls13"`      // 是否启用TLS1.3
	Scv        bool   `json:"scv"`        // 是否跳过证书验证
	Tfo        bool   `json:"tfo"`        // 是否启用TCP Fast Open
	
	// 节点过滤
	NodeList       bool     `json:"nodelist"`
	IncludeRemarks []string `json:"include_remarks"`
	ExcludeRemarks []string `json:"exclude_remarks"`
	
	// 规则和策略组
	EnableRule   bool   `json:"enable_rule"`
	EnableInsert bool   `json:"enable_insert"`
	RuleProvider string `json:"rule_provider"`
	Config       string `json:"config"`
	
	// 客户端特定选项
	ClashOptions       *ClashOptions       `json:"clash_options,omitempty"`
	SurgeOptions       *SurgeOptions       `json:"surge_options,omitempty"`
	QuantumultXOptions *QuantumultXOptions `json:"quantumultx_options,omitempty"`
	
	// 模板和自定义
	Template      string            `json:"template"`
	CustomOptions map[string]string `json:"custom_options"`
}

// ClashOptions Clash特定选项
type ClashOptions struct {
	NewName          bool `json:"new_name"`
	ClashDns         bool `json:"clash_dns"`
	AppendType       bool `json:"append_type"`
	AppendInfo       bool `json:"append_info"`
	PrependInsert    bool `json:"prepend_insert"`
	ClassicalRuleset bool `json:"classical_ruleset"`
	TLS13Support     bool `json:"tls13_support"`
}

// SurgeOptions Surge特定选项
type SurgeOptions struct {
	SurgeVer        int  `json:"surge_ver"`
	V2rayPlugin     bool `json:"v2ray_plugin"`
	ResolveHostname bool `json:"resolve_hostname"`
}

// QuantumultXOptions QuantumultX特定选项
type QuantumultXOptions struct {
	AddEmoji       bool `json:"add_emoji"`
	RemoveOldEmoji bool `json:"remove_old_emoji"`
	AppendInfo     bool `json:"append_info"`
}

// Manager 生成器管理器
type Manager struct {
	generators map[string]Generator
}

// NewManager 创建生成器管理器
func NewManager() *Manager {
	m := &Manager{
		generators: make(map[string]Generator),
	}
	
	// 注册所有生成器
	m.registerGenerators()
	
	return m
}

// registerGenerators 注册所有生成器
func (m *Manager) registerGenerators() {
	// 注册各种客户端生成器
	m.RegisterGenerator(&ClashGenerator{})
	m.RegisterGenerator(&SurgeGenerator{})
	m.RegisterGenerator(&QuantumultXGenerator{})
	m.RegisterGenerator(&LoonGenerator{})
	m.RegisterGenerator(&SingBoxGenerator{})
	m.RegisterGenerator(&V2rayGenerator{})
	m.RegisterGenerator(&SSGenerator{})
}

// RegisterGenerator 注册生成器
func (m *Manager) RegisterGenerator(generator Generator) {
	m.generators[generator.GetTarget()] = generator
}

// Generate 生成配置
func (m *Manager) Generate(proxies []*models.Proxy, config *GenerateConfig) ([]byte, error) {
	if config == nil {
		return nil, fmt.Errorf("generate config is nil")
	}
	
	if config.Target == "" {
		return nil, fmt.Errorf("target client not specified")
	}
	
	generator, exists := m.generators[config.Target]
	if !exists {
		return nil, fmt.Errorf("unsupported target client: %s", config.Target)
	}
	
	// 过滤代理节点
	filteredProxies, err := m.filterProxies(proxies, config)
	if err != nil {
		return nil, fmt.Errorf("failed to filter proxies: %w", err)
	}
	
	// 验证代理节点兼容性
	var validProxies []*models.Proxy
	for _, proxy := range filteredProxies {
		if err := generator.Validate(proxy); err != nil {
			// 跳过不兼容的代理节点
			continue
		}
		validProxies = append(validProxies, proxy)
	}
	
	if len(validProxies) == 0 {
		return nil, fmt.Errorf("no valid proxies found for target client: %s", config.Target)
	}
	
	// 生成配置
	return generator.Generate(validProxies, config)
}

// filterProxies 过滤代理节点
func (m *Manager) filterProxies(proxies []*models.Proxy, config *GenerateConfig) ([]*models.Proxy, error) {
	var filtered []*models.Proxy
	
	for _, proxy := range proxies {
		// 应用包含/排除规则
		if m.shouldIncludeProxy(proxy, config) {
			filtered = append(filtered, proxy)
		}
	}
	
	// 排序（如果启用）
	if config.Sort {
		filtered = m.sortProxies(filtered)
	}
	
	return filtered, nil
}

// shouldIncludeProxy 检查是否应该包含该代理节点
func (m *Manager) shouldIncludeProxy(proxy *models.Proxy, config *GenerateConfig) bool {
	// 检查包含列表
	if len(config.IncludeRemarks) > 0 {
		found := false
		for _, include := range config.IncludeRemarks {
			if strings.Contains(proxy.Remark, include) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// 检查排除列表
	for _, exclude := range config.ExcludeRemarks {
		if strings.Contains(proxy.Remark, exclude) {
			return false
		}
	}
	
	return true
}

// sortProxies 排序代理节点
func (m *Manager) sortProxies(proxies []*models.Proxy) []*models.Proxy {
	sort.Slice(proxies, func(i, j int) bool {
		return proxies[i].Remark < proxies[j].Remark
	})
	return proxies
}

// GetSupportedTargets 获取支持的目标客户端
func (m *Manager) GetSupportedTargets() []string {
	var targets []string
	for target := range m.generators {
		targets = append(targets, target)
	}
	sort.Strings(targets)
	return targets
}

// GetGenerator 获取指定目标的生成器
func (m *Manager) GetGenerator(target string) (Generator, bool) {
	generator, exists := m.generators[target]
	return generator, exists
}

// ValidateConfig 验证生成配置
func (m *Manager) ValidateConfig(config *GenerateConfig) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}
	
	if config.Target == "" {
		return fmt.Errorf("target is required")
	}
	
	if _, exists := m.generators[config.Target]; !exists {
		return fmt.Errorf("unsupported target: %s", config.Target)
	}
	
	return nil
}