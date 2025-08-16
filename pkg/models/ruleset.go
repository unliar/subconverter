// Package models 定义规则集相关的数据模型
package models

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"subconverter-go/pkg/types"
	"time"
)

// RulesetConfig 规则集配置 - 对应 C++ 版本的 RulesetConfig
type RulesetConfig struct {
	Group    string           `json:"group" yaml:"group" validate:"required"`
	URL      string           `json:"url" yaml:"url" validate:"required,url"`
	Interval int              `json:"interval,omitempty" yaml:"interval,omitempty"`
	Type     types.RulesetType `json:"type,omitempty" yaml:"type,omitempty"`

	// 元数据
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
}

// IsValid 验证规则集配置是否有效
func (rc *RulesetConfig) IsValid() bool {
	if rc == nil || rc.Group == "" || rc.URL == "" {
		return false
	}

	// 验证 URL 格式
	if _, err := url.Parse(rc.URL); err != nil {
		return false
	}

	// Interval 应该大于 0
	if rc.Interval <= 0 {
		rc.Interval = 86400 // 默认 24 小时
	}

	return true
}

// SetDefaults 设置默认值
func (rc *RulesetConfig) SetDefaults() {
	if rc.Interval <= 0 {
		rc.Interval = 86400 // 24小时
	}
	if rc.CreatedAt.IsZero() {
		rc.CreatedAt = time.Now()
	}
	rc.UpdatedAt = time.Now()
}

// GetKey 获取规则集的唯一键
func (rc *RulesetConfig) GetKey() string {
	return fmt.Sprintf("%s:%s", rc.Group, rc.URL)
}

// Clone 深拷贝规则集配置
func (rc *RulesetConfig) Clone() *RulesetConfig {
	if rc == nil {
		return nil
	}
	
	clone := *rc
	return &clone
}

// RulesetContent 规则集内容
type RulesetContent struct {
	Group     string           `json:"group" yaml:"group"`
	URL       string           `json:"url" yaml:"url"`
	Content   string           `json:"content" yaml:"content"`
	Type      types.RulesetType `json:"type" yaml:"type"`
	UpdatedAt time.Time        `json:"updated_at" yaml:"updated_at"`
	Hash      string           `json:"hash,omitempty" yaml:"hash,omitempty"`
}

// GetHash 计算内容哈希
func (rc *RulesetContent) GetHash() string {
	if rc.Hash == "" {
		hasher := md5.New()
		hasher.Write([]byte(rc.Content))
		rc.Hash = fmt.Sprintf("%x", hasher.Sum(nil))
	}
	return rc.Hash
}

// IsExpired 检查内容是否过期
func (rc *RulesetContent) IsExpired(interval int) bool {
	if interval <= 0 {
		interval = 86400 // 默认24小时
	}
	return time.Since(rc.UpdatedAt) > time.Duration(interval)*time.Second
}

// GetSize 获取内容大小
func (rc *RulesetContent) GetSize() int {
	return len(rc.Content)
}

// Validate 验证规则集内容
func (rc *RulesetContent) Validate() error {
	if rc.Group == "" {
		return types.NewConvertError(types.ErrorCodeValidationError, "ruleset group is required")
	}
	if rc.URL == "" {
		return types.NewConvertError(types.ErrorCodeValidationError, "ruleset URL is required")
	}
	if rc.Content == "" {
		return types.NewConvertError(types.ErrorCodeValidationError, "ruleset content is empty")
	}
	return nil
}

// RulesetList 规则集列表
type RulesetList []*RulesetConfig

// Len 返回规则集列表长度
func (rl RulesetList) Len() int {
	return len(rl)
}

// FilterByGroup 按分组过滤规则集
func (rl RulesetList) FilterByGroup(group string) RulesetList {
	var filtered RulesetList
	for _, ruleset := range rl {
		if ruleset.Group == group {
			filtered = append(filtered, ruleset)
		}
	}
	return filtered
}

// FilterByType 按类型过滤规则集
func (rl RulesetList) FilterByType(rulesetType types.RulesetType) RulesetList {
	var filtered RulesetList
	for _, ruleset := range rl {
		if ruleset.Type == rulesetType {
			filtered = append(filtered, ruleset)
		}
	}
	return filtered
}

// GetGroups 获取所有分组名称（去重）
func (rl RulesetList) GetGroups() []string {
	groups := make(map[string]bool)
	for _, ruleset := range rl {
		groups[ruleset.Group] = true
	}
	
	result := make([]string, 0, len(groups))
	for group := range groups {
		result = append(result, group)
	}
	return result
}

// GetURLs 获取所有 URL
func (rl RulesetList) GetURLs() []string {
	urls := make([]string, len(rl))
	for i, ruleset := range rl {
		urls[i] = ruleset.URL
	}
	return urls
}

// FindByURL 根据 URL 查找规则集
func (rl RulesetList) FindByURL(url string) *RulesetConfig {
	for _, ruleset := range rl {
		if ruleset.URL == url {
			return ruleset
		}
	}
	return nil
}

// Validate 验证所有规则集
func (rl RulesetList) Validate() []error {
	var errors []error
	keys := make(map[string]bool)

	for _, ruleset := range rl {
		// 检查规则集是否有效
		if !ruleset.IsValid() {
			errors = append(errors, types.NewConvertError(
				types.ErrorCodeValidationError,
				"invalid ruleset configuration",
				ruleset.Group,
			))
		}

		// 检查键是否重复
		key := ruleset.GetKey()
		if keys[key] {
			errors = append(errors, types.NewConvertError(
				types.ErrorCodeValidationError,
				"duplicate ruleset",
				key,
			))
		}
		keys[key] = true
	}

	return errors
}

// GroupByType 按类型分组规则集
func (rl RulesetList) GroupByType() map[types.RulesetType]RulesetList {
	groups := make(map[types.RulesetType]RulesetList)
	for _, ruleset := range rl {
		groups[ruleset.Type] = append(groups[ruleset.Type], ruleset)
	}
	return groups
}

// Clone 深拷贝规则集列表
func (rl RulesetList) Clone() RulesetList {
	clone := make(RulesetList, len(rl))
	for i, ruleset := range rl {
		clone[i] = ruleset.Clone()
	}
	return clone
}