// Package models 定义规则集相关的数据模型
package models

import (
	"time"
)

// RulesetConfig 规则集配置
type RulesetConfig struct {
	Name         string   `json:"name" yaml:"name" validate:"required"`
	Type         string   `json:"type" yaml:"type" validate:"required,oneof=DOMAIN DOMAIN-SUFFIX DOMAIN-KEYWORD IP-CIDR IP-CIDR6 GEOIP SRC-IP-CIDR SRC-PORT DST-PORT PROCESS-NAME RULE-SET"`
	Rule         string   `json:"rule" yaml:"rule" validate:"required"`
	Policy       string   `json:"policy" yaml:"policy" validate:"required"`
	URL          string   `json:"url,omitempty" yaml:"url,omitempty"`
	Path         string   `json:"path,omitempty" yaml:"path,omitempty"`
	Interval     int      `json:"interval,omitempty" yaml:"interval,omitempty"`
	Group        string   `json:"group,omitempty" yaml:"group,omitempty"`
	NoResolve    bool     `json:"no_resolve,omitempty" yaml:"no_resolve,omitempty"`
	SourceIPCIDR []string `json:"source_ip_cidr,omitempty" yaml:"source_ip_cidr,omitempty"`
	IPCIDR       []string `json:"ip_cidr,omitempty" yaml:"ip_cidr,omitempty"`
	Domain       []string `json:"domain,omitempty" yaml:"domain,omitempty"`
	DomainSuffix []string `json:"domain_suffix,omitempty" yaml:"domain_suffix,omitempty"`
	DomainKeyword []string `json:"domain_keyword,omitempty" yaml:"domain_keyword,omitempty"`
	
	// 元数据
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
}

// IsValid 验证规则集配置是否有效
func (rc *RulesetConfig) IsValid() bool {
	if rc == nil {
		return false
	}
	if rc.Name == "" || rc.Type == "" || rc.Policy == "" {
		return false
	}
	
	// 验证规则类型
	validTypes := []string{"DOMAIN", "DOMAIN-SUFFIX", "DOMAIN-KEYWORD", "IP-CIDR", "IP-CIDR6", "GEOIP", "SRC-IP-CIDR", "SRC-PORT", "DST-PORT", "PROCESS-NAME", "RULE-SET"}
	isValidType := false
	for _, validType := range validTypes {
		if rc.Type == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return false
	}
	
	return true
}

// SetDefaults 设置默认值
func (rc *RulesetConfig) SetDefaults() {
	if rc.Interval <= 0 {
		rc.Interval = 86400 // 默认24小时更新一次
	}
	if rc.CreatedAt.IsZero() {
		rc.CreatedAt = time.Now()
	}
	rc.UpdatedAt = time.Now()
}

// Clone 深拷贝规则集配置
func (rc *RulesetConfig) Clone() *RulesetConfig {
	if rc == nil {
		return nil
	}
	
	clone := *rc
	
	// 深拷贝切片
	if rc.SourceIPCIDR != nil {
		clone.SourceIPCIDR = make([]string, len(rc.SourceIPCIDR))
		copy(clone.SourceIPCIDR, rc.SourceIPCIDR)
	}
	if rc.IPCIDR != nil {
		clone.IPCIDR = make([]string, len(rc.IPCIDR))
		copy(clone.IPCIDR, rc.IPCIDR)
	}
	if rc.Domain != nil {
		clone.Domain = make([]string, len(rc.Domain))
		copy(clone.Domain, rc.Domain)
	}
	if rc.DomainSuffix != nil {
		clone.DomainSuffix = make([]string, len(rc.DomainSuffix))
		copy(clone.DomainSuffix, rc.DomainSuffix)
	}
	if rc.DomainKeyword != nil {
		clone.DomainKeyword = make([]string, len(rc.DomainKeyword))
		copy(clone.DomainKeyword, rc.DomainKeyword)
	}
	
	return &clone
}