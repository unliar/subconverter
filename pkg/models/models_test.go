package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelsIntegration(t *testing.T) {
	t.Run("Data Model Completeness", func(t *testing.T) {
		// 测试所有核心数据模型是否可以正常创建和使用
		testDataModelCompleteness(t)
	})

	t.Run("Validation System", func(t *testing.T) {
		// 测试验证系统的完整性
		testValidationSystem(t)
	})

	t.Run("Conversion System", func(t *testing.T) {
		// 测试转换器系统的完整性
		testConversionSystem(t)
	})

	t.Run("Error Handling", func(t *testing.T) {
		// 测试错误处理系统
		testErrorHandling(t)
	})
}

func testDataModelCompleteness(t *testing.T) {
	// 验证所有核心模型都可以创建和基本操作
	
	// 1. 测试 Proxy 模型
	proxy := &Proxy{}
	assert.NotNil(t, proxy)
	assert.NotNil(t, proxy.Clone())
	
	// 2. 测试 ProxyGroupConfig 模型
	group := &ProxyGroupConfig{}
	assert.NotNil(t, group)
	assert.NotNil(t, group.Clone())
	
	// 3. 测试 RulesetConfig 模型
	ruleset := &RulesetConfig{}
	assert.NotNil(t, ruleset)
	assert.NotNil(t, ruleset.Clone())
	
	// 4. 测试 ConvertRequest 模型
	request := &ConvertRequest{}
	assert.NotNil(t, request)
	assert.NotNil(t, request.Clone())
	
	// 5. 测试 ConvertResponse 模型
	response := &ConvertResponse{}
	assert.NotNil(t, response)
	assert.NotNil(t, response.Clone())
	
	// 6. 测试 ServerConfig 模型
	serverConfig := &ServerConfig{}
	assert.NotNil(t, serverConfig)
	assert.NotNil(t, serverConfig.Clone())
	
	// 7. 测试 ConverterConfig 模型
	converterConfig := &ConverterConfig{}
	assert.NotNil(t, converterConfig)
	assert.NotNil(t, converterConfig.Clone())
	
	// 8. 测试 TemplateConfig 模型
	templateConfig := &TemplateConfig{}
	assert.NotNil(t, templateConfig)
	assert.NotNil(t, templateConfig.Clone())
	
	// 9. 测试 ApplicationConfig 模型
	appConfig := &ApplicationConfig{}
	assert.NotNil(t, appConfig)
	assert.NotNil(t, appConfig.Clone())
}

func testValidationSystem(t *testing.T) {
	// 测试验证器系统的各个组件
	
	// 1. 测试 ValidateStruct 函数存在
	proxy := &Proxy{}
	err := ValidateStruct(proxy)
	// 应该有错误，因为 proxy 是空的
	assert.Error(t, err)
	
	// 2. 测试具体验证方法
	validProxy := &Proxy{
		Type:     1, // types.ProxyTypeShadowsocks
		Remark:   "test",
		Hostname: "example.com",
		Port:     443,
	}
	err = ValidateProxy(validProxy)
	// 可能有错误，但不应该panic
	assert.NotPanics(t, func() {
		ValidateProxy(validProxy)
	})
	
	// 3. 测试代理组验证
	validGroup := &ProxyGroupConfig{
		Name:    "test-group",
		Proxies: []string{"proxy1"},
	}
	assert.NotPanics(t, func() {
		ValidateProxyGroup(validGroup)
	})
	
	// 4. 测试规则集验证
	validRuleset := &RulesetConfig{
		Group: "test",
		URL:   "https://example.com/rules.list",
	}
	assert.NotPanics(t, func() {
		ValidateRuleset(validRuleset)
	})
	
	// 5. 测试请求验证
	validRequest := &ConvertRequest{
		URL:    "https://example.com/subscription",
		Target: "clash",
	}
	assert.NotPanics(t, func() {
		ValidateConvertRequest(validRequest)
	})
}

func testConversionSystem(t *testing.T) {
	// 测试转换器系统
	
	converter := NewProxyConverter()
	assert.NotNil(t, converter)
	
	// 测试代理转换为 Map
	proxy := &Proxy{
		Type:     1, // types.ProxyTypeShadowsocks
		Remark:   "test-proxy",
		Hostname: "example.com",
		Port:     443,
	}
	
	proxyMap := converter.ToMap(proxy)
	assert.NotNil(t, proxyMap)
	assert.Contains(t, proxyMap, "name")
	assert.Contains(t, proxyMap, "server")
	assert.Contains(t, proxyMap, "port")
	
	// 测试 JSON 序列化
	jsonData, err := converter.ToJSON(proxy)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)
	
	// 测试 JSON 反序列化
	decodedProxy, err := converter.FromJSON(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, proxy.Remark, decodedProxy.Remark)
	
	// 测试 YAML 序列化
	yamlData, err := converter.ToYAML(proxy)
	assert.NoError(t, err)
	assert.NotEmpty(t, yamlData)
	
	// 测试代理列表转换
	proxies := ProxyList{proxy}
	proxyMaps := converter.ProxyListToMaps(proxies)
	assert.Len(t, proxyMaps, 1)
	assert.Equal(t, proxyMap, proxyMaps[0])
}

func testErrorHandling(t *testing.T) {
	// 测试错误处理系统
	
	// 1. 测试 nil 对象验证
	assert.NotPanics(t, func() {
		var nilProxy *Proxy
		nilProxy.IsValid()
	})
	
	assert.NotPanics(t, func() {
		var nilGroup *ProxyGroupConfig
		nilGroup.IsValid()
	})
	
	assert.NotPanics(t, func() {
		var nilRuleset *RulesetConfig
		nilRuleset.IsValid()
	})
	
	assert.NotPanics(t, func() {
		var nilRequest *ConvertRequest
		nilRequest.IsValid()
	})
	
	assert.NotPanics(t, func() {
		var nilResponse *ConvertResponse
		nilResponse.IsSuccess()
	})
	
	// 2. 测试 Clone nil 对象
	assert.NotPanics(t, func() {
		var nilProxy *Proxy
		cloned := nilProxy.Clone()
		assert.Nil(t, cloned)
	})
	
	assert.NotPanics(t, func() {
		var nilGroup *ProxyGroupConfig
		cloned := nilGroup.Clone()
		assert.Nil(t, cloned)
	})
	
	// 3. 测试空切片操作
	var emptyProxyList ProxyList
	assert.Equal(t, 0, emptyProxyList.Len())
	assert.Empty(t, emptyProxyList.GetRemarks())
	
	var emptyGroupList ProxyGroupList
	assert.Equal(t, 0, emptyGroupList.Len())
	assert.Empty(t, emptyGroupList.GetNames())
	
	var emptyRulesetList RulesetList
	assert.Equal(t, 0, emptyRulesetList.Len())
	assert.Empty(t, emptyRulesetList.GetURLs())
}

func TestDataModelConsistency(t *testing.T) {
	// 测试数据模型的一致性
	
	t.Run("Default Values Consistency", func(t *testing.T) {
		// 测试默认值设置的一致性
		
		proxy := &Proxy{}
		proxy.SetDefaults()
		assert.False(t, proxy.CreatedAt.IsZero())
		assert.False(t, proxy.UpdatedAt.IsZero())
		
		group := &ProxyGroupConfig{}
		group.SetDefaults()
		assert.False(t, group.CreatedAt.IsZero())
		assert.False(t, group.UpdatedAt.IsZero())
		
		ruleset := &RulesetConfig{}
		ruleset.SetDefaults()
		assert.False(t, ruleset.CreatedAt.IsZero())
		assert.False(t, ruleset.UpdatedAt.IsZero())
		
		request := &ConvertRequest{}
		request.SetDefaults()
		assert.False(t, request.CreatedAt.IsZero())
		
		serverConfig := &ServerConfig{}
		serverConfig.SetDefaults()
		assert.False(t, serverConfig.CreatedAt.IsZero())
		assert.False(t, serverConfig.UpdatedAt.IsZero())
		
		converterConfig := &ConverterConfig{}
		converterConfig.SetDefaults()
		assert.False(t, converterConfig.CreatedAt.IsZero())
		assert.False(t, converterConfig.UpdatedAt.IsZero())
		
		templateConfig := &TemplateConfig{}
		templateConfig.SetDefaults()
		assert.False(t, templateConfig.CreatedAt.IsZero())
		assert.False(t, templateConfig.UpdatedAt.IsZero())
		
		appConfig := &ApplicationConfig{}
		appConfig.SetDefaults()
		assert.False(t, appConfig.CreatedAt.IsZero())
		assert.False(t, appConfig.UpdatedAt.IsZero())
	})
	
	t.Run("Clone Independence", func(t *testing.T) {
		// 测试克隆对象的独立性
		
		// 创建一个复杂的代理对象
		udp := true
		original := &Proxy{
			Type:     1, // types.ProxyTypeShadowsocks
			Remark:   "original",
			Hostname: "example.com",
			Port:     443,
			UDP:      &udp,
		}
		
		cloned := original.Clone()
		require.NotNil(t, cloned)
		
		// 修改克隆对象
		cloned.Remark = "cloned"
		*cloned.UDP = false
		
		// 原对象应该保持不变
		assert.Equal(t, "original", original.Remark)
		assert.True(t, *original.UDP)
		
		// 克隆对象应该有新值
		assert.Equal(t, "cloned", cloned.Remark)
		assert.False(t, *cloned.UDP)
	})
}

func BenchmarkDataModels(b *testing.B) {
	// 性能基准测试
	
	b.Run("Proxy Operations", func(b *testing.B) {
		proxy := &Proxy{
			Type:     1, // types.ProxyTypeShadowsocks
			Remark:   "benchmark-proxy",
			Hostname: "example.com",
			Port:     443,
		}
		
		b.Run("IsValid", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				proxy.IsValid()
			}
		})
		
		b.Run("Clone", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				proxy.Clone()
			}
		})
		
		b.Run("SetDefaults", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				testProxy := &Proxy{}
				testProxy.SetDefaults()
			}
		})
	})
	
	b.Run("ProxyList Operations", func(b *testing.B) {
		// 创建一个包含1000个代理的列表
		proxies := make(ProxyList, 1000)
		for i := 0; i < 1000; i++ {
			proxies[i] = &Proxy{
				Type:     1, // types.ProxyTypeShadowsocks
				Remark:   "proxy-" + string(rune(i)),
				Hostname: "example.com",
				Port:     443,
			}
		}
		
		b.Run("FilterByType", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				proxies.FilterByType(1) // types.ProxyTypeShadowsocks
			}
		})
		
		b.Run("GetRemarks", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				proxies.GetRemarks()
			}
		})
		
		b.Run("GroupByType", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				proxies.GroupByType()
			}
		})
	})
	
	b.Run("Converter Operations", func(b *testing.B) {
		converter := NewProxyConverter()
		proxy := &Proxy{
			Type:     1, // types.ProxyTypeShadowsocks
			Remark:   "benchmark-proxy",
			Hostname: "example.com",
			Port:     443,
		}
		
		b.Run("ToMap", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				converter.ToMap(proxy)
			}
		})
		
		b.Run("ToJSON", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				converter.ToJSON(proxy)
			}
		})
		
		b.Run("ToYAML", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				converter.ToYAML(proxy)
			}
		})
	})
}