package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRulesetConfigCreation(t *testing.T) {
	t.Run("Create Basic RulesetConfig", func(t *testing.T) {
		ruleset := &RulesetConfig{
			Name:     "test-ruleset",
			Type:     "DOMAIN",
			Rule:     "example.com",
			Policy:   "REJECT",
			URL:      "https://example.com/rules.list",
			Interval: 86400,
			Group:    "ads",
		}
		
		assert.NotNil(t, ruleset)
		assert.Equal(t, "test-ruleset", ruleset.Name)
		assert.Equal(t, "DOMAIN", ruleset.Type)
		assert.Equal(t, "example.com", ruleset.Rule)
		assert.Equal(t, "REJECT", ruleset.Policy)
		assert.Equal(t, "https://example.com/rules.list", ruleset.URL)
		assert.Equal(t, 86400, ruleset.Interval)
		assert.Equal(t, "ads", ruleset.Group)
	})
}

func TestRulesetConfigValidation(t *testing.T) {
	t.Run("Test Validation Logic", func(t *testing.T) {
		tests := []struct {
			name     string
			ruleset  *RulesetConfig
			expected bool
		}{
			{
				name: "valid ruleset",
				ruleset: &RulesetConfig{
					Name:   "test",
					Type:   "DOMAIN",
					Policy: "REJECT",
				},
				expected: true,
			},
			{
				name: "missing name",
				ruleset: &RulesetConfig{
					Type:   "DOMAIN",
					Policy: "REJECT",
				},
				expected: false,
			},
			{
				name: "missing type",
				ruleset: &RulesetConfig{
					Name:   "test",
					Policy: "REJECT",
				},
				expected: false,
			},
			{
				name: "missing policy",
				ruleset: &RulesetConfig{
					Name: "test",
					Type: "DOMAIN",
				},
				expected: false,
			},
			{
				name: "invalid type",
				ruleset: &RulesetConfig{
					Name:   "test",
					Type:   "INVALID-TYPE",
					Policy: "REJECT",
				},
				expected: false,
			},
			{
				name:     "nil ruleset",
				ruleset:  nil,
				expected: false,
			},
		}
		
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := tt.ruleset.IsValid()
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

func TestRulesetConfigTypes(t *testing.T) {
	t.Run("Test Valid Rule Types", func(t *testing.T) {
		validTypes := []string{
			"DOMAIN",
			"DOMAIN-SUFFIX", 
			"DOMAIN-KEYWORD",
			"IP-CIDR",
			"IP-CIDR6",
			"GEOIP",
			"SRC-IP-CIDR",
			"SRC-PORT",
			"DST-PORT",
			"PROCESS-NAME",
			"RULE-SET",
		}
		
		for _, ruleType := range validTypes {
			t.Run(ruleType, func(t *testing.T) {
				ruleset := &RulesetConfig{
					Name:   "test",
					Type:   ruleType,
					Policy: "REJECT",
				}
				
				assert.True(t, ruleset.IsValid(), "Type %s should be valid", ruleType)
			})
		}
	})
}

func TestRulesetConfigFields(t *testing.T) {
	t.Run("Test All Fields", func(t *testing.T) {
		ruleset := &RulesetConfig{}
		
		// 设置必需字段
		ruleset.Name = "test-ruleset"
		assert.Equal(t, "test-ruleset", ruleset.Name)
		
		ruleset.Type = "DOMAIN"
		assert.Equal(t, "DOMAIN", ruleset.Type)
		
		ruleset.Rule = "example.com"
		assert.Equal(t, "example.com", ruleset.Rule)
		
		ruleset.Policy = "REJECT"
		assert.Equal(t, "REJECT", ruleset.Policy)
		
		// 设置可选字段
		ruleset.URL = "https://test.example.com/rules.list"
		assert.Equal(t, "https://test.example.com/rules.list", ruleset.URL)
		
		ruleset.Path = "/local/path/rules.list"
		assert.Equal(t, "/local/path/rules.list", ruleset.Path)
		
		ruleset.Interval = 3600
		assert.Equal(t, 3600, ruleset.Interval)
		
		ruleset.Group = "test-group"
		assert.Equal(t, "test-group", ruleset.Group)
		
		ruleset.NoResolve = true
		assert.True(t, ruleset.NoResolve)
	})
}

func TestRulesetConfigDefaults(t *testing.T) {
	t.Run("Test SetDefaults Method", func(t *testing.T) {
		ruleset := &RulesetConfig{
			Name:   "test",
			Type:   "DOMAIN",
			Policy: "REJECT",
		}
		
		// 默认间隔应该是0（未设置）
		assert.Equal(t, 0, ruleset.Interval)
		assert.True(t, ruleset.CreatedAt.IsZero())
		assert.True(t, ruleset.UpdatedAt.IsZero())
		
		// 调用SetDefaults
		ruleset.SetDefaults()
		
		// 检查默认值
		assert.Equal(t, 86400, ruleset.Interval) // 24小时
		assert.False(t, ruleset.CreatedAt.IsZero())
		assert.False(t, ruleset.UpdatedAt.IsZero())
	})
}

func TestRulesetConfigClone(t *testing.T) {
	t.Run("Test Clone Method", func(t *testing.T) {
		original := &RulesetConfig{
			Name:          "test",
			Type:          "DOMAIN",
			Rule:          "example.com",
			Policy:        "REJECT",
			URL:           "https://example.com/rules.list",
			Group:         "ads",
			SourceIPCIDR:  []string{"192.168.1.0/24"},
			Domain:        []string{"example.com", "test.com"},
			DomainSuffix:  []string{".com"},
			DomainKeyword: []string{"ads"},
		}
		
		cloned := original.Clone()
		
		// 验证值相等
		assert.Equal(t, original.Name, cloned.Name)
		assert.Equal(t, original.Type, cloned.Type)
		assert.Equal(t, original.Rule, cloned.Rule)
		assert.Equal(t, original.Policy, cloned.Policy)
		assert.Equal(t, original.URL, cloned.URL)
		assert.Equal(t, original.Group, cloned.Group)
		
		// 验证切片内容相等
		assert.Equal(t, original.SourceIPCIDR, cloned.SourceIPCIDR)
		assert.Equal(t, original.Domain, cloned.Domain)
		assert.Equal(t, original.DomainSuffix, cloned.DomainSuffix)
		assert.Equal(t, original.DomainKeyword, cloned.DomainKeyword)
		
		// 验证是独立的对象
		assert.NotSame(t, original, cloned)
		
		// 修改克隆不应影响原对象
		cloned.Name = "modified"
		assert.Equal(t, "test", original.Name)
		assert.Equal(t, "modified", cloned.Name)
		
		// 修改切片不应影响原对象
		if len(cloned.Domain) > 0 {
			cloned.Domain[0] = "modified.com"
			assert.Equal(t, "example.com", original.Domain[0])
			assert.Equal(t, "modified.com", cloned.Domain[0])
		}
	})
}

func TestRulesetSliceOperations(t *testing.T) {
	t.Run("Test Ruleset Slice Operations", func(t *testing.T) {
		ruleset1 := &RulesetConfig{
			Name:   "ads-rules",
			Type:   "DOMAIN",
			Policy: "REJECT",
			Group:  "ads",
		}
		
		ruleset2 := &RulesetConfig{
			Name:   "privacy-rules",
			Type:   "DOMAIN-SUFFIX",
			Policy: "REJECT",
			Group:  "privacy",
		}
		
		ruleset3 := &RulesetConfig{
			Name:   "ads-rules-2",
			Type:   "DOMAIN-KEYWORD",
			Policy: "REJECT",
			Group:  "ads",
		}
		
		rulesets := []*RulesetConfig{ruleset1, ruleset2, ruleset3}
		
		// 测试长度
		assert.Equal(t, 3, len(rulesets))
		
		// 测试访问元素
		assert.Equal(t, ruleset1, rulesets[0])
		assert.Equal(t, ruleset2, rulesets[1])
		assert.Equal(t, ruleset3, rulesets[2])
		
		// 测试按分组计数
		adsCount := 0
		privacyCount := 0
		for _, ruleset := range rulesets {
			switch ruleset.Group {
			case "ads":
				adsCount++
			case "privacy":
				privacyCount++
			}
		}
		assert.Equal(t, 2, adsCount)
		assert.Equal(t, 1, privacyCount)
	})
}

func TestRulesetNilHandling(t *testing.T) {
	t.Run("Test Nil Handling", func(t *testing.T) {
		// 测试nil规则集不会导致panic
		var nilRuleset *RulesetConfig
		assert.Nil(t, nilRuleset)
		assert.False(t, nilRuleset.IsValid())
		
		// 测试包含nil元素的列表
		rulesets := []*RulesetConfig{
			nil,
			{Name: "test", Type: "DOMAIN", Policy: "REJECT"},
		}
		assert.Equal(t, 2, len(rulesets))
		assert.Nil(t, rulesets[0])
		assert.NotNil(t, rulesets[1])
	})
}

func BenchmarkRulesetCreation(b *testing.B) {
	b.Run("Create RulesetConfig", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = &RulesetConfig{
				Name:     "benchmark",
				Type:     "DOMAIN",
				Policy:   "REJECT",
				Interval: 86400,
			}
		}
	})
}

func BenchmarkRulesetValidation(b *testing.B) {
	ruleset := &RulesetConfig{
		Name:   "benchmark",
		Type:   "DOMAIN",
		Policy: "REJECT",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ruleset.IsValid()
	}
}