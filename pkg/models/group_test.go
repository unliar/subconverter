package models

import (
	"testing"
	"time"

	"subconverter-go/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxyGroupConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		group    *ProxyGroupConfig
		expected bool
	}{
		{
			name: "valid select group",
			group: &ProxyGroupConfig{
				Name:    "test-group",
				Type:    types.ProxyGroupTypeSelect,
				Proxies: []string{"proxy1", "proxy2"},
			},
			expected: true,
		},
		{
			name: "valid url-test group",
			group: &ProxyGroupConfig{
				Name:    "auto-group",
				Type:    types.ProxyGroupTypeURLTest,
				Proxies: []string{"proxy1", "proxy2"},
				URL:     "http://www.google.com/generate_204",
			},
			expected: true,
		},
		{
			name: "invalid - missing name",
			group: &ProxyGroupConfig{
				Type:    types.ProxyGroupTypeSelect,
				Proxies: []string{"proxy1"},
			},
			expected: false,
		},
		{
			name: "invalid - url-test without URL",
			group: &ProxyGroupConfig{
				Name:    "auto-group",
				Type:    types.ProxyGroupTypeURLTest,
				Proxies: []string{"proxy1"},
			},
			expected: false,
		},
		{
			name: "invalid - no proxies or providers",
			group: &ProxyGroupConfig{
				Name: "empty-group",
				Type: types.ProxyGroupTypeSelect,
			},
			expected: false,
		},
		{
			name:     "nil group",
			group:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.group.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProxyGroupConfig_Clone(t *testing.T) {
	lazy := true
	original := &ProxyGroupConfig{
		Name:              "test-group",
		Type:              types.ProxyGroupTypeURLTest,
		Proxies:           []string{"proxy1", "proxy2", "proxy3"},
		UsingProvider:     []string{"provider1"},
		URL:               "http://www.google.com/generate_204",
		Interval:          300,
		Timeout:           5000,
		Lazy:              &lazy,
		CreatedAt:         time.Now(),
	}

	cloned := original.Clone()

	// 验证值相等
	assert.Equal(t, original.Name, cloned.Name)
	assert.Equal(t, original.Type, cloned.Type)
	assert.Equal(t, original.URL, cloned.URL)
	assert.Equal(t, original.Interval, cloned.Interval)
	assert.Equal(t, *original.Lazy, *cloned.Lazy)

	// 验证切片是独立的
	assert.NotSame(t, original.Proxies, cloned.Proxies)
	assert.NotSame(t, original.UsingProvider, cloned.UsingProvider)
	assert.NotSame(t, original.Lazy, cloned.Lazy)

	// 修改克隆对象不应影响原对象
	cloned.Proxies[0] = "modified"
	assert.Equal(t, "proxy1", original.Proxies[0])
	assert.Equal(t, "modified", cloned.Proxies[0])

	*cloned.Lazy = false
	assert.True(t, *original.Lazy)
	assert.False(t, *cloned.Lazy)
}

func TestProxyGroupConfig_SetDefaults(t *testing.T) {
	group := &ProxyGroupConfig{
		Name: "test-group",
		Type: types.ProxyGroupTypeURLTest,
	}

	group.SetDefaults()

	assert.Equal(t, 300, group.Interval)
	assert.Equal(t, 5000, group.Timeout)
	assert.Equal(t, 50, group.Tolerance)
	assert.False(t, group.CreatedAt.IsZero())
	assert.False(t, group.UpdatedAt.IsZero())
}

func TestProxyGroupConfig_ProxyManagement(t *testing.T) {
	group := &ProxyGroupConfig{
		Name:    "test-group",
		Type:    types.ProxyGroupTypeSelect,
		Proxies: []string{"proxy1", "proxy2"},
	}

	// 测试检查代理是否存在
	assert.True(t, group.HasProxy("proxy1"))
	assert.False(t, group.HasProxy("proxy3"))

	// 测试添加代理
	group.AddProxy("proxy3")
	assert.True(t, group.HasProxy("proxy3"))
	assert.Len(t, group.Proxies, 3)

	// 测试添加重复代理（不应增加）
	group.AddProxy("proxy1")
	assert.Len(t, group.Proxies, 3)

	// 测试删除代理
	group.RemoveProxy("proxy2")
	assert.False(t, group.HasProxy("proxy2"))
	assert.Len(t, group.Proxies, 2)
}

func TestProxyGroupList_Methods(t *testing.T) {
	groups := ProxyGroupList{
		{
			Name: "select-1",
			Type: types.ProxyGroupTypeSelect,
		},
		{
			Name: "auto-1",
			Type: types.ProxyGroupTypeURLTest,
		},
		{
			Name: "select-2",
			Type: types.ProxyGroupTypeSelect,
		},
	}

	// 测试长度
	assert.Equal(t, 3, groups.Len())

	// 测试按类型过滤
	selectGroups := groups.FilterByType(types.ProxyGroupTypeSelect)
	assert.Equal(t, 2, selectGroups.Len())

	autoGroups := groups.FilterByType(types.ProxyGroupTypeURLTest)
	assert.Equal(t, 1, autoGroups.Len())

	// 测试获取名称
	names := groups.GetNames()
	expected := []string{"select-1", "auto-1", "select-2"}
	assert.Equal(t, expected, names)

	// 测试按名称查找
	found := groups.FindByName("auto-1")
	require.NotNil(t, found)
	assert.Equal(t, "auto-1", found.Name)

	notFound := groups.FindByName("nonexistent")
	assert.Nil(t, notFound)
}

func TestProxyGroupList_Validate(t *testing.T) {
	groups := ProxyGroupList{
		{
			Name:    "valid-group",
			Type:    types.ProxyGroupTypeSelect,
			Proxies: []string{"proxy1"},
		},
		{
			Name: "invalid-group",
			Type: types.ProxyGroupTypeURLTest,
			// 缺少 URL 和 Proxies
		},
		{
			Name:    "duplicate-name",
			Type:    types.ProxyGroupTypeSelect,
			Proxies: []string{"proxy1"},
		},
		{
			Name:    "duplicate-name", // 重复名称
			Type:    types.ProxyGroupTypeFallback,
			Proxies: []string{"proxy1"},
			URL:     "http://www.google.com/generate_204",
		},
	}

	errors := groups.Validate()
	assert.NotEmpty(t, errors)
	
	// 应该有至少2个错误：一个无效组 + 一个重复名称
	assert.GreaterOrEqual(t, len(errors), 2)
}

func TestProxyGroupConfig_GetTypeString(t *testing.T) {
	group := &ProxyGroupConfig{
		Type: types.ProxyGroupTypeURLTest,
	}

	assert.Equal(t, "url-test", group.GetTypeString())
}

func TestProxyGroupConfig_GetStrategyString(t *testing.T) {
	group := &ProxyGroupConfig{
		Strategy: types.BalanceStrategyRoundRobin,
	}

	assert.Equal(t, "round-robin", group.GetStrategyString())
}

func BenchmarkProxyGroupConfig_Clone(b *testing.B) {
	group := &ProxyGroupConfig{
		Name:          "benchmark-group",
		Type:          types.ProxyGroupTypeSelect,
		Proxies:       []string{"proxy1", "proxy2", "proxy3", "proxy4", "proxy5"},
		UsingProvider: []string{"provider1", "provider2"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = group.Clone()
	}
}

func BenchmarkProxyGroupConfig_IsValid(b *testing.B) {
	group := &ProxyGroupConfig{
		Name:    "benchmark-group",
		Type:    types.ProxyGroupTypeSelect,
		Proxies: []string{"proxy1", "proxy2"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = group.IsValid()
	}
}