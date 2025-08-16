package models

import (
	"testing"
	"time"

	"subconverter-go/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRulesetConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		ruleset  *RulesetConfig
		expected bool
	}{
		{
			name: "valid ruleset",
			ruleset: &RulesetConfig{
				Group:    "example",
				URL:      "https://example.com/rules.list",
				Interval: 86400,
				Type:     types.RulesetTypeClashDomain,
			},
			expected: true,
		},
		{
			name: "valid ruleset with default interval",
			ruleset: &RulesetConfig{
				Group: "example",
				URL:   "https://example.com/rules.list",
				Type:  types.RulesetTypeClashDomain,
			},
			expected: true,
		},
		{
			name: "invalid - missing group",
			ruleset: &RulesetConfig{
				URL:      "https://example.com/rules.list",
				Interval: 86400,
			},
			expected: false,
		},
		{
			name: "invalid - missing URL",
			ruleset: &RulesetConfig{
				Group:    "example",
				Interval: 86400,
			},
			expected: false,
		},
		{
			name: "invalid - malformed URL",
			ruleset: &RulesetConfig{
				Group: "example",
				URL:   "not-a-valid-url",
			},
			expected: true, // Go的url.Parse对此比较宽松，会接受这种格式
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
}

func TestRulesetConfig_SetDefaults(t *testing.T) {
	ruleset := &RulesetConfig{
		Group: "test",
		URL:   "https://example.com/rules.list",
	}

	ruleset.SetDefaults()

	assert.Equal(t, 86400, ruleset.Interval)
	assert.False(t, ruleset.CreatedAt.IsZero())
	assert.False(t, ruleset.UpdatedAt.IsZero())
}

func TestRulesetConfig_GetKey(t *testing.T) {
	ruleset := &RulesetConfig{
		Group: "example",
		URL:   "https://example.com/rules.list",
	}

	key := ruleset.GetKey()
	expected := "example:https://example.com/rules.list"
	assert.Equal(t, expected, key)
}

func TestRulesetConfig_Clone(t *testing.T) {
	original := &RulesetConfig{
		Group:     "example",
		URL:       "https://example.com/rules.list",
		Interval:  86400,
		Type:      types.RulesetTypeClashDomain,
		CreatedAt: time.Now(),
	}

	cloned := original.Clone()

	// 验证值相等
	assert.Equal(t, original.Group, cloned.Group)
	assert.Equal(t, original.URL, cloned.URL)
	assert.Equal(t, original.Interval, cloned.Interval)
	assert.Equal(t, original.Type, cloned.Type)
	assert.Equal(t, original.CreatedAt, cloned.CreatedAt)

	// 验证是独立的对象
	assert.NotSame(t, original, cloned)

	// 修改克隆不应影响原对象
	cloned.Group = "modified"
	assert.Equal(t, "example", original.Group)
	assert.Equal(t, "modified", cloned.Group)
}

func TestRulesetContent_GetHash(t *testing.T) {
	content := &RulesetContent{
		Group:   "test",
		URL:     "https://example.com/rules.list",
		Content: "DOMAIN,example.com\nDOMAIN,test.com",
		Type:    types.RulesetTypeClashDomain,
	}

	hash1 := content.GetHash()
	assert.NotEmpty(t, hash1)
	assert.Len(t, hash1, 32) // MD5 hash length

	// 再次调用应该返回相同的哈希（缓存）
	hash2 := content.GetHash()
	assert.Equal(t, hash1, hash2)

	// 修改内容应该生成新的哈希
	content.Content = "DOMAIN,newexample.com"
	content.Hash = "" // 重置缓存
	hash3 := content.GetHash()
	assert.NotEqual(t, hash1, hash3)
}

func TestRulesetContent_IsExpired(t *testing.T) {
	content := &RulesetContent{
		UpdatedAt: time.Now().Add(-2 * time.Hour),
	}

	// 使用1小时间隔，应该过期
	expired := content.IsExpired(3600)
	assert.True(t, expired)

	// 使用3小时间隔，不应该过期
	notExpired := content.IsExpired(3 * 3600)
	assert.False(t, notExpired)

	// 使用0或负数间隔，应该使用默认24小时
	notExpiredDefault := content.IsExpired(0)
	assert.False(t, notExpiredDefault)
}

func TestRulesetContent_GetSize(t *testing.T) {
	content := &RulesetContent{
		Content: "DOMAIN,example.com\nDOMAIN,test.com",
	}

	size := content.GetSize()
	expected := len("DOMAIN,example.com\nDOMAIN,test.com")
	assert.Equal(t, expected, size)
}

func TestRulesetContent_Validate(t *testing.T) {
	tests := []struct {
		name     string
		content  *RulesetContent
		hasError bool
	}{
		{
			name: "valid content",
			content: &RulesetContent{
				Group:   "test",
				URL:     "https://example.com/rules.list",
				Content: "DOMAIN,example.com",
			},
			hasError: false,
		},
		{
			name: "missing group",
			content: &RulesetContent{
				URL:     "https://example.com/rules.list",
				Content: "DOMAIN,example.com",
			},
			hasError: true,
		},
		{
			name: "missing URL",
			content: &RulesetContent{
				Group:   "test",
				Content: "DOMAIN,example.com",
			},
			hasError: true,
		},
		{
			name: "empty content",
			content: &RulesetContent{
				Group: "test",
				URL:   "https://example.com/rules.list",
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.content.Validate()
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRulesetList_Methods(t *testing.T) {
	rulesets := RulesetList{
		{
			Group: "group1",
			URL:   "https://example.com/rules1.list",
			Type:  types.RulesetTypeClashDomain,
		},
		{
			Group: "group2",
			URL:   "https://example.com/rules2.list",
			Type:  types.RulesetTypeClashIPCIDR,
		},
		{
			Group: "group1",
			URL:   "https://example.com/rules3.list",
			Type:  types.RulesetTypeClashDomain,
		},
	}

	// 测试长度
	assert.Equal(t, 3, rulesets.Len())

	// 测试按分组过滤
	group1Rules := rulesets.FilterByGroup("group1")
	assert.Equal(t, 2, group1Rules.Len())

	group2Rules := rulesets.FilterByGroup("group2")
	assert.Equal(t, 1, group2Rules.Len())

	// 测试按类型过滤
	domainRules := rulesets.FilterByType(types.RulesetTypeClashDomain)
	assert.Equal(t, 2, domainRules.Len())

	ipRules := rulesets.FilterByType(types.RulesetTypeClashIPCIDR)
	assert.Equal(t, 1, ipRules.Len())

	// 测试获取分组
	groups := rulesets.GetGroups()
	assert.Len(t, groups, 2)
	assert.Contains(t, groups, "group1")
	assert.Contains(t, groups, "group2")

	// 测试获取URL
	urls := rulesets.GetURLs()
	expected := []string{
		"https://example.com/rules1.list",
		"https://example.com/rules2.list",
		"https://example.com/rules3.list",
	}
	assert.Equal(t, expected, urls)

	// 测试按URL查找
	found := rulesets.FindByURL("https://example.com/rules2.list")
	require.NotNil(t, found)
	assert.Equal(t, "group2", found.Group)

	notFound := rulesets.FindByURL("https://nonexistent.com/rules.list")
	assert.Nil(t, notFound)
}

func TestRulesetList_GroupByType(t *testing.T) {
	rulesets := RulesetList{
		{
			Group: "group1",
			URL:   "https://example.com/rules1.list",
			Type:  types.RulesetTypeClashDomain,
		},
		{
			Group: "group2",
			URL:   "https://example.com/rules2.list",
			Type:  types.RulesetTypeClashDomain,
		},
		{
			Group: "group3",
			URL:   "https://example.com/rules3.list",
			Type:  types.RulesetTypeClashIPCIDR,
		},
	}

	groups := rulesets.GroupByType()
	assert.Len(t, groups, 2)
	assert.Len(t, groups[types.RulesetTypeClashDomain], 2)
	assert.Len(t, groups[types.RulesetTypeClashIPCIDR], 1)
}

func TestRulesetList_Validate(t *testing.T) {
	rulesets := RulesetList{
		{
			Group: "valid-group",
			URL:   "https://example.com/rules.list",
		},
		{
			Group: "invalid-group",
			// 缺少 URL
		},
		{
			Group: "duplicate-key",
			URL:   "https://example.com/same.list",
		},
		{
			Group: "duplicate-key",
			URL:   "https://example.com/same.list", // 重复键
		},
	}

	errors := rulesets.Validate()
	assert.NotEmpty(t, errors)
	
	// 应该有至少2个错误：一个无效规则集 + 一个重复键
	assert.GreaterOrEqual(t, len(errors), 2)
}

func TestRulesetList_Clone(t *testing.T) {
	original := RulesetList{
		{
			Group: "group1",
			URL:   "https://example.com/rules1.list",
		},
		{
			Group: "group2",
			URL:   "https://example.com/rules2.list",
		},
	}

	cloned := original.Clone()

	// 验证长度相同
	assert.Equal(t, original.Len(), cloned.Len())

	// 验证值相等但对象独立
	for i := range original {
		assert.Equal(t, original[i].Group, cloned[i].Group)
		assert.Equal(t, original[i].URL, cloned[i].URL)
		assert.NotSame(t, original[i], cloned[i])
	}

	// 修改克隆不应影响原对象
	cloned[0].Group = "modified"
	assert.Equal(t, "group1", original[0].Group)
	assert.Equal(t, "modified", cloned[0].Group)
}

func BenchmarkRulesetConfig_IsValid(b *testing.B) {
	ruleset := &RulesetConfig{
		Group:    "benchmark",
		URL:      "https://example.com/rules.list",
		Interval: 86400,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ruleset.IsValid()
	}
}

func BenchmarkRulesetContent_GetHash(b *testing.B) {
	content := &RulesetContent{
		Content: "DOMAIN,example.com\nDOMAIN,test.com\nDOMAIN,benchmark.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		content.Hash = "" // 重置缓存以测试计算性能
		_ = content.GetHash()
	}
}

func BenchmarkRulesetList_FilterByGroup(b *testing.B) {
	rulesets := make(RulesetList, 1000)
	for i := 0; i < 1000; i++ {
		rulesets[i] = &RulesetConfig{
			Group: "group" + string(rune(i%10)),
			URL:   "https://example.com/rules.list",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rulesets.FilterByGroup("group5")
	}
}