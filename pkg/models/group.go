// Package models 定义代理组相关的数据模型
package models

import (
	"subconverter-go/pkg/types"
	"time"
)

// ProxyGroupConfig 代理组配置 - 对应 C++ 版本的 ProxyGroupConfig
type ProxyGroupConfig struct {
	Name               string                    `json:"name" yaml:"name" validate:"required"`
	Type               types.ProxyGroupType     `json:"type" yaml:"type"`
	Proxies            []string                 `json:"proxies,omitempty" yaml:"proxies,omitempty"`
	UsingProvider      []string                 `json:"use,omitempty" yaml:"use,omitempty"`
	URL                string                   `json:"url,omitempty" yaml:"url,omitempty"`
	Interval           int                      `json:"interval,omitempty" yaml:"interval,omitempty"`
	Timeout            int                      `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Tolerance          int                      `json:"tolerance,omitempty" yaml:"tolerance,omitempty"`
	Strategy           types.BalanceStrategy    `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	Lazy               *bool                    `json:"lazy,omitempty" yaml:"lazy,omitempty"`
	DisableUdp         *bool                    `json:"disable_udp,omitempty" yaml:"disable_udp,omitempty"`
	Persistent         *bool                    `json:"persistent,omitempty" yaml:"persistent,omitempty"`
	EvaluateBeforeUse  *bool                    `json:"evaluate_before_use,omitempty" yaml:"evaluate_before_use,omitempty"`

	// 元数据
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
}

// IsValid 验证代理组配置是否有效
func (pgc *ProxyGroupConfig) IsValid() bool {
	if pgc == nil || pgc.Name == "" {
		return false
	}

	// URL Test 和 Fallback 需要 URL
	if pgc.Type.RequiresURL() && pgc.URL == "" {
		return false
	}

	// 必须有代理或使用提供商
	if len(pgc.Proxies) == 0 && len(pgc.UsingProvider) == 0 {
		return false
	}

	return true
}

// GetTypeString 获取代理组类型字符串
func (pgc *ProxyGroupConfig) GetTypeString() string {
	return pgc.Type.String()
}

// GetStrategyString 获取负载均衡策略字符串
func (pgc *ProxyGroupConfig) GetStrategyString() string {
	return pgc.Strategy.String()
}

// Clone 深拷贝代理组配置
func (pgc *ProxyGroupConfig) Clone() *ProxyGroupConfig {
	if pgc == nil {
		return nil
	}

	clone := *pgc

	// 深拷贝切片
	if pgc.Proxies != nil {
		clone.Proxies = make([]string, len(pgc.Proxies))
		copy(clone.Proxies, pgc.Proxies)
	}
	if pgc.UsingProvider != nil {
		clone.UsingProvider = make([]string, len(pgc.UsingProvider))
		copy(clone.UsingProvider, pgc.UsingProvider)
	}

	// 深拷贝指针字段
	if pgc.Lazy != nil {
		lazy := *pgc.Lazy
		clone.Lazy = &lazy
	}
	if pgc.DisableUdp != nil {
		disableUdp := *pgc.DisableUdp
		clone.DisableUdp = &disableUdp
	}
	if pgc.Persistent != nil {
		persistent := *pgc.Persistent
		clone.Persistent = &persistent
	}
	if pgc.EvaluateBeforeUse != nil {
		evaluateBeforeUse := *pgc.EvaluateBeforeUse
		clone.EvaluateBeforeUse = &evaluateBeforeUse
	}

	return &clone
}

// SetDefaults 设置默认值
func (pgc *ProxyGroupConfig) SetDefaults() {
	if pgc.Interval <= 0 {
		pgc.Interval = 300 // 5分钟
	}
	if pgc.Timeout <= 0 {
		pgc.Timeout = 5000 // 5秒
	}
	if pgc.Tolerance <= 0 {
		pgc.Tolerance = 50 // 50ms
	}
	if pgc.CreatedAt.IsZero() {
		pgc.CreatedAt = time.Now()
	}
	pgc.UpdatedAt = time.Now()
}

// HasProxy 检查是否包含指定代理
func (pgc *ProxyGroupConfig) HasProxy(proxyName string) bool {
	for _, proxy := range pgc.Proxies {
		if proxy == proxyName {
			return true
		}
	}
	return false
}

// AddProxy 添加代理
func (pgc *ProxyGroupConfig) AddProxy(proxyName string) {
	if !pgc.HasProxy(proxyName) {
		pgc.Proxies = append(pgc.Proxies, proxyName)
		pgc.UpdatedAt = time.Now()
	}
}

// RemoveProxy 移除代理
func (pgc *ProxyGroupConfig) RemoveProxy(proxyName string) {
	for i, proxy := range pgc.Proxies {
		if proxy == proxyName {
			pgc.Proxies = append(pgc.Proxies[:i], pgc.Proxies[i+1:]...)
			pgc.UpdatedAt = time.Now()
			break
		}
	}
}

// ProxyGroupList 代理组列表
type ProxyGroupList []*ProxyGroupConfig

// Len 返回代理组列表长度
func (pgl ProxyGroupList) Len() int {
	return len(pgl)
}

// FilterByType 按类型过滤代理组
func (pgl ProxyGroupList) FilterByType(groupType types.ProxyGroupType) ProxyGroupList {
	var filtered ProxyGroupList
	for _, group := range pgl {
		if group.Type == groupType {
			filtered = append(filtered, group)
		}
	}
	return filtered
}

// GetNames 获取所有代理组名称
func (pgl ProxyGroupList) GetNames() []string {
	names := make([]string, len(pgl))
	for i, group := range pgl {
		names[i] = group.Name
	}
	return names
}

// FindByName 根据名称查找代理组
func (pgl ProxyGroupList) FindByName(name string) *ProxyGroupConfig {
	for _, group := range pgl {
		if group.Name == name {
			return group
		}
	}
	return nil
}

// Validate 验证所有代理组
func (pgl ProxyGroupList) Validate() []error {
	var errors []error
	names := make(map[string]bool)

	for _, group := range pgl {
		// 检查代理组是否有效
		if !group.IsValid() {
			errors = append(errors, types.NewConvertError(
				types.ErrorCodeValidationError,
				"invalid proxy group",
				group.Name,
			))
		}

		// 检查名称是否重复
		if names[group.Name] {
			errors = append(errors, types.NewConvertError(
				types.ErrorCodeValidationError,
				"duplicate proxy group name",
				group.Name,
			))
		}
		names[group.Name] = true

		// 验证 URL 格式
		if group.Type.RequiresURL() && group.URL != "" {
			if !isValidURL(group.URL) {
				errors = append(errors, types.NewConvertError(
					types.ErrorCodeValidationError,
					"invalid URL in proxy group",
					group.Name,
				))
			}
		}
	}

	return errors
}

// isValidURL 简单的 URL 验证
func isValidURL(url string) bool {
	return len(url) > 7 && (url[:7] == "http://" || url[:8] == "https://")
}