package models

import (
	"testing"

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
		Name:          "test-group",
		Type:          types.ProxyGroupTypeURLTest,
		Proxies:       []string{"proxy1", "proxy2", "proxy3"},
		UsingProvider: []string{"provider1"},
		URL:           "http://www.google.com/generate_204",
		Interval:      300,
		Timeout:       5000,
		Lazy:          &lazy,
	}

	cloned := original.Clone()

	// 验证值相等
	assert.Equal(t, original.Name, cloned.Name)
	assert.Equal(t, original.Type, cloned.Type)
	assert.Equal(t, original.URL, cloned.URL)
	assert.Equal(t, original.Interval, cloned.Interval)
	assert.Equal(t, *original.Lazy, *cloned.Lazy)

	// 验证切片内容相等
	assert.Equal(t, original.Proxies, cloned.Proxies)
	assert.Equal(t, original.UsingProvider, cloned.UsingProvider)

	// 验证是独立的对象
	assert.NotSame(t, original, cloned)

	// 修改克隆对象不应影响原对象
	if len(cloned.Proxies) > 0 {
		cloned.Proxies[0] = "modified"
		assert.Equal(t, "proxy1", original.Proxies[0])
		assert.Equal(t, "modified", cloned.Proxies[0])
	}

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

func TestProxyGroupConfig_BasicFields(t *testing.T) {
	t.Run("Test Basic Field Operations", func(t *testing.T) {
		group := &ProxyGroupConfig{}
		
		// 设置基本字段
		group.Name = "test-group"
		assert.Equal(t, "test-group", group.Name)
		
		group.Type = types.ProxyGroupTypeSelect
		assert.Equal(t, types.ProxyGroupTypeSelect, group.Type)
		
		group.Proxies = []string{"proxy1", "proxy2"}
		assert.Equal(t, []string{"proxy1", "proxy2"}, group.Proxies)
		
		group.URL = "http://www.google.com/generate_204"
		assert.Equal(t, "http://www.google.com/generate_204", group.URL)
		
		group.Interval = 300
		assert.Equal(t, 300, group.Interval)
		
		group.Timeout = 5000
		assert.Equal(t, 5000, group.Timeout)
		
		group.Tolerance = 50
		assert.Equal(t, 50, group.Tolerance)
	})
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
	assert.Equal(t, 3, len(groups))

	// 测试按类型计数
	selectCount := 0
	urlTestCount := 0
	for _, group := range groups {
		switch group.Type {
		case types.ProxyGroupTypeSelect:
			selectCount++
		case types.ProxyGroupTypeURLTest:
			urlTestCount++
		}
	}
	assert.Equal(t, 2, selectCount)
	assert.Equal(t, 1, urlTestCount)

	// 测试按名称查找
	var found *ProxyGroupConfig
	for _, group := range groups {
		if group.Name == "auto-1" {
			found = group
			break
		}
	}
	require.NotNil(t, found)
	assert.Equal(t, "auto-1", found.Name)

	// 测试查找不存在的组
	var notFound *ProxyGroupConfig
	for _, group := range groups {
		if group.Name == "nonexistent" {
			notFound = group
			break
		}
	}
	assert.Nil(t, notFound)
}

func TestProxyGroupList_Validation(t *testing.T) {
	t.Run("Test Group List Validation", func(t *testing.T) {
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
		}

		// 手动验证
		validCount := 0
		invalidCount := 0
		for _, group := range groups {
			if group.IsValid() {
				validCount++
			} else {
				invalidCount++
			}
		}
		
		assert.Equal(t, 1, validCount)
		assert.Equal(t, 1, invalidCount)
	})
}

func TestProxyGroupConfig_TypeValidation(t *testing.T) {
	t.Run("Test Different Group Types", func(t *testing.T) {
		groupTypes := []types.ProxyGroupType{
			types.ProxyGroupTypeSelect,
			types.ProxyGroupTypeURLTest,
			types.ProxyGroupTypeFallback,
			types.ProxyGroupTypeLoadBalance,
		}
		
		for _, groupType := range groupTypes {
			t.Run(groupType.String(), func(t *testing.T) {
				group := &ProxyGroupConfig{
					Name:    "test-group",
					Type:    groupType,
					Proxies: []string{"proxy1", "proxy2"},
				}
				
				// url-test 和 fallback 需要URL
				if groupType.RequiresURL() {
					group.URL = "http://www.google.com/generate_204"
				}
				
				assert.True(t, group.IsValid(), "Group type %s should be valid", groupType.String())
			})
		}
	})
}

func TestProxyGroupConfig_NilHandling(t *testing.T) {
	t.Run("Test Nil Handling", func(t *testing.T) {
		// 测试nil组不会导致panic
		var nilGroup *ProxyGroupConfig
		assert.Nil(t, nilGroup)
		assert.False(t, nilGroup.IsValid())
		
		// 测试nil克隆
		cloned := nilGroup.Clone()
		assert.Nil(t, cloned)
		
		// 测试空组列表
		var emptyGroups ProxyGroupList
		assert.Equal(t, 0, len(emptyGroups))
	})
}

func BenchmarkProxyGroupConfig_Creation(b *testing.B) {
	b.Run("Create ProxyGroupConfig", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = &ProxyGroupConfig{
				Name:    "benchmark-group",
				Type:    types.ProxyGroupTypeSelect,
				Proxies: []string{"proxy1", "proxy2", "proxy3"},
			}
		}
	})
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