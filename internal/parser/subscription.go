// internal/parser/subscription.go
package parser

import (
	"encoding/base64"
	"fmt"
	"strings"

	"subconverter-go/pkg/models"
)

// Base64SubscriptionParser Base64订阅解析器
type Base64SubscriptionParser struct{}

// CanParseSubscription 检查是否能解析指定的订阅内容
func (p *Base64SubscriptionParser) CanParseSubscription(content string) bool {
	// 尝试Base64解码
	if _, err := base64.StdEncoding.DecodeString(content); err == nil {
		return true
	}
	if _, err := base64.URLEncoding.DecodeString(content); err == nil {
		return true
	}
	return false
}

// GetFormat 获取订阅格式
func (p *Base64SubscriptionParser) GetFormat() string {
	return "base64"
}

// ParseSubscription 解析Base64订阅内容
func (p *Base64SubscriptionParser) ParseSubscription(content string) ([]*models.Proxy, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("empty subscription content")
	}

	// 尝试Base64解码
	var decoded []byte
	var err error
	
	decoded, err = base64.StdEncoding.DecodeString(content)
	if err != nil {
		decoded, err = base64.URLEncoding.DecodeString(content)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 content: %w", err)
		}
	}

	// 解码后的内容应该是代理链接的列表
	lines := strings.Split(string(decoded), "\n")
	var validLinks []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			validLinks = append(validLinks, line)
		}
	}

	if len(validLinks) == 0 {
		return nil, fmt.Errorf("no valid proxy links found in base64 subscription")
	}

	// 使用通用解析器解析每个链接
	manager := NewManager()
	return manager.ParseMultiple(validLinks)
}

// ClashSubscriptionParser Clash订阅解析器
type ClashSubscriptionParser struct{}

// CanParseSubscription 检查是否能解析指定的订阅内容
func (p *ClashSubscriptionParser) CanParseSubscription(content string) bool {
	// 检查是否包含Clash配置关键词
	content = strings.ToLower(content)
	return strings.Contains(content, "proxies:") || 
		   strings.Contains(content, "proxy-groups:") ||
		   strings.Contains(content, "rules:")
}

// GetFormat 获取订阅格式
func (p *ClashSubscriptionParser) GetFormat() string {
	return "clash"
}

// ParseSubscription 解析Clash订阅内容
func (p *ClashSubscriptionParser) ParseSubscription(content string) ([]*models.Proxy, error) {
	// 这里应该实现完整的Clash YAML解析
	// 为了简化，这里只是返回一个错误
	return nil, fmt.Errorf("clash subscription parsing not implemented yet")
}

// SurgeSubscriptionParser Surge订阅解析器
type SurgeSubscriptionParser struct{}

// CanParseSubscription 检查是否能解析指定的订阅内容
func (p *SurgeSubscriptionParser) CanParseSubscription(content string) bool {
	// 检查是否包含Surge配置关键词
	content = strings.ToLower(content)
	return strings.Contains(content, "[proxy]") || 
		   strings.Contains(content, "[proxy group]") ||
		   strings.Contains(content, "[rule]")
}

// GetFormat 获取订阅格式
func (p *SurgeSubscriptionParser) GetFormat() string {
	return "surge"
}

// ParseSubscription 解析Surge订阅内容
func (p *SurgeSubscriptionParser) ParseSubscription(content string) ([]*models.Proxy, error) {
	// 这里应该实现完整的Surge配置解析
	// 为了简化，这里只是返回一个错误
	return nil, fmt.Errorf("surge subscription parsing not implemented yet")
}

// V2raySubscriptionParser V2ray订阅解析器
type V2raySubscriptionParser struct{}

// CanParseSubscription 检查是否能解析指定的订阅内容
func (p *V2raySubscriptionParser) CanParseSubscription(content string) bool {
	// 检查是否包含V2ray配置关键词
	content = strings.ToLower(content)
	return strings.Contains(content, "\"outbounds\"") || 
		   strings.Contains(content, "\"protocol\"") &&
		   (strings.Contains(content, "vmess") || strings.Contains(content, "vless"))
}

// GetFormat 获取订阅格式
func (p *V2raySubscriptionParser) GetFormat() string {
	return "v2ray"
}

// ParseSubscription 解析V2ray订阅内容
func (p *V2raySubscriptionParser) ParseSubscription(content string) ([]*models.Proxy, error) {
	// 这里应该实现完整的V2ray JSON配置解析
	// 为了简化，这里只是返回一个错误
	return nil, fmt.Errorf("v2ray subscription parsing not implemented yet")
}