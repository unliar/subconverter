// internal/parser/parser.go
package parser

import (
	"fmt"
	"strings"

	"subconverter-go/pkg/models"
)

// Parser 解析器接口
type Parser interface {
	// Parse 解析单个代理链接
	Parse(link string) (*models.Proxy, error)
	
	// CanParse 检查是否能解析指定的链接
	CanParse(link string) bool
	
	// GetType 获取解析器类型
	GetType() models.ProxyType
}

// SubscriptionParser 订阅解析器接口
type SubscriptionParser interface {
	// ParseSubscription 解析订阅链接
	ParseSubscription(content string) ([]*models.Proxy, error)
	
	// CanParseSubscription 检查是否能解析指定的订阅内容
	CanParseSubscription(content string) bool
	
	// GetFormat 获取订阅格式
	GetFormat() string
}

// Manager 解析器管理器
type Manager struct {
	parsers              map[models.ProxyType]Parser
	subscriptionParsers  map[string]SubscriptionParser
}

// NewManager 创建解析器管理器
func NewManager() *Manager {
	m := &Manager{
		parsers:             make(map[models.ProxyType]Parser),
		subscriptionParsers: make(map[string]SubscriptionParser),
	}
	
	// 注册所有解析器
	m.registerParsers()
	
	return m
}

// registerParsers 注册所有解析器
func (m *Manager) registerParsers() {
	// 注册代理解析器
	m.RegisterParser(&ShadowsocksParser{})
	m.RegisterParser(&ShadowsocksRParser{})
	m.RegisterParser(&VMeSSParser{})
	m.RegisterParser(&VLESSParser{})
	m.RegisterParser(&TrojanParser{})
	m.RegisterParser(&HTTPParser{})
	m.RegisterParser(&SOCKS5Parser{})
	
	// 注册订阅解析器
	m.RegisterSubscriptionParser("base64", &Base64SubscriptionParser{})
	m.RegisterSubscriptionParser("clash", &ClashSubscriptionParser{})
	m.RegisterSubscriptionParser("surge", &SurgeSubscriptionParser{})
	m.RegisterSubscriptionParser("v2ray", &V2raySubscriptionParser{})
}

// RegisterParser 注册代理解析器
func (m *Manager) RegisterParser(parser Parser) {
	m.parsers[parser.GetType()] = parser
}

// RegisterSubscriptionParser 注册订阅解析器
func (m *Manager) RegisterSubscriptionParser(format string, parser SubscriptionParser) {
	m.subscriptionParsers[format] = parser
}

// Parse 解析单个代理链接
func (m *Manager) Parse(link string) (*models.Proxy, error) {
	link = strings.TrimSpace(link)
	if link == "" {
		return nil, fmt.Errorf("empty link")
	}
	
	// 尝试所有解析器
	for _, parser := range m.parsers {
		if parser.CanParse(link) {
			return parser.Parse(link)
		}
	}
	
	return nil, fmt.Errorf("unsupported proxy link format: %s", link)
}

// ParseMultiple 解析多个代理链接
func (m *Manager) ParseMultiple(links []string) ([]*models.Proxy, error) {
	var proxies []*models.Proxy
	var errors []error
	
	for _, link := range links {
		proxy, err := m.Parse(link)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to parse link %s: %w", link, err))
			continue
		}
		if proxy != nil {
			proxies = append(proxies, proxy)
		}
	}
	
	if len(proxies) == 0 && len(errors) > 0 {
		return nil, fmt.Errorf("failed to parse any links: %v", errors)
	}
	
	return proxies, nil
}

// ParseSubscription 解析订阅内容
func (m *Manager) ParseSubscription(content string) ([]*models.Proxy, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("empty subscription content")
	}
	
	// 尝试所有订阅解析器
	for _, parser := range m.subscriptionParsers {
		if parser.CanParseSubscription(content) {
			return parser.ParseSubscription(content)
		}
	}
	
	// 如果没有专门的订阅解析器能处理，尝试按行解析
	return m.parseAsLines(content)
}

// parseAsLines 按行解析订阅内容
func (m *Manager) parseAsLines(content string) ([]*models.Proxy, error) {
	lines := strings.Split(content, "\n")
	var validLines []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			validLines = append(validLines, line)
		}
	}
	
	if len(validLines) == 0 {
		return nil, fmt.Errorf("no valid proxy links found in subscription")
	}
	
	return m.ParseMultiple(validLines)
}

// GetSupportedTypes 获取支持的代理类型
func (m *Manager) GetSupportedTypes() []models.ProxyType {
	var types []models.ProxyType
	for proxyType := range m.parsers {
		types = append(types, proxyType)
	}
	return types
}

// GetSupportedFormats 获取支持的订阅格式
func (m *Manager) GetSupportedFormats() []string {
	var formats []string
	for format := range m.subscriptionParsers {
		formats = append(formats, format)
	}
	return formats
}

// ValidateProxy 验证代理配置
func (m *Manager) ValidateProxy(proxy *models.Proxy) error {
	if proxy == nil {
		return fmt.Errorf("proxy is nil")
	}
	
	return proxy.Validate()
}

// ParseFromURL 从URL解析代理链接
func (m *Manager) ParseFromURL(rawURL string) (*models.Proxy, error) {
	// 这里可以添加URL预处理逻辑
	// 比如处理URL编码、验证URL格式等
	return m.Parse(rawURL)
}