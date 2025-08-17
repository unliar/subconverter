package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelsCreation(t *testing.T) {
	t.Run("Create Basic Models", func(t *testing.T) {
		// 测试所有核心数据模型是否可以正常创建
		
		// 1. 测试 Proxy 模型
		proxy := &Proxy{}
		assert.NotNil(t, proxy)
		
		// 2. 测试 ProxyGroupConfig 模型
		group := &ProxyGroupConfig{}
		assert.NotNil(t, group)
		
		// 3. 测试 RulesetConfig 模型
		ruleset := &RulesetConfig{}
		assert.NotNil(t, ruleset)
		
		// 4. 测试 ConvertRequest 模型
		request := &ConvertRequest{}
		assert.NotNil(t, request)
		
		// 5. 测试 ConvertResponse 模型
		response := &ConvertResponse{}
		assert.NotNil(t, response)
		
		// 6. 测试 ServerConfig 模型
		serverConfig := &ServerConfig{}
		assert.NotNil(t, serverConfig)
		
		// 7. 测试 ConverterConfig 模型
		converterConfig := &ConverterConfig{}
		assert.NotNil(t, converterConfig)
		
		// 8. 测试 TemplateConfig 模型
		templateConfig := &TemplateConfig{}
		assert.NotNil(t, templateConfig)
		
		// 9. 测试 ApplicationConfig 模型
		appConfig := &ApplicationConfig{}
		assert.NotNil(t, appConfig)
	})
}

func TestProxyTypes(t *testing.T) {
	t.Run("Proxy Type Constants", func(t *testing.T) {
		// 测试代理类型常量
		assert.Equal(t, ProxyType(1), ProxyTypeShadowsocks)
		assert.Equal(t, ProxyType(2), ProxyTypeShadowsocksR)
		assert.Equal(t, ProxyType(3), ProxyTypeVMess)
		assert.Equal(t, ProxyType(4), ProxyTypeTrojan)
		assert.Equal(t, ProxyType(6), ProxyTypeHTTP)
		assert.Equal(t, ProxyType(8), ProxyTypeSOCKS5)
		assert.Equal(t, ProxyType(10), ProxyTypeVLESS)
	})
}

func TestProxyBasicOperations(t *testing.T) {
	t.Run("Proxy Creation and Basic Properties", func(t *testing.T) {
		proxy := &Proxy{
			Type:     ProxyTypeShadowsocks,
			Remark:   "test-proxy",
			Hostname: "example.com",
			Port:     443,
		}
		
		assert.Equal(t, ProxyTypeShadowsocks, proxy.Type)
		assert.Equal(t, "test-proxy", proxy.Remark)
		assert.Equal(t, "example.com", proxy.Hostname)
		assert.Equal(t, uint16(443), proxy.Port)
	})
}

func TestProxyList(t *testing.T) {
	t.Run("ProxyList Operations", func(t *testing.T) {
		proxy1 := &Proxy{
			Type:     ProxyTypeShadowsocks,
			Remark:   "proxy1",
			Hostname: "example1.com",
			Port:     443,
		}
		
		proxy2 := &Proxy{
			Type:     ProxyTypeVMess,
			Remark:   "proxy2", 
			Hostname: "example2.com",
			Port:     8080,
		}
		
		proxies := ProxyList{proxy1, proxy2}
		assert.Len(t, proxies, 2)
		assert.Equal(t, proxy1, proxies[0])
		assert.Equal(t, proxy2, proxies[1])
	})
}

func TestValidationFunctions(t *testing.T) {
	t.Run("Validation Functions Exist", func(t *testing.T) {
		// 测试验证函数是否存在且不会panic
		
		proxy := &Proxy{
			Type:     ProxyTypeShadowsocks,
			Remark:   "test",
			Hostname: "example.com",
			Port:     443,
		}
		
		assert.NotPanics(t, func() {
			ValidateStruct(proxy)
		})
		
		assert.NotPanics(t, func() {
			ValidateProxy(proxy)
		})
		
		group := &ProxyGroupConfig{
			Name:    "test-group",
			Proxies: []string{"proxy1"},
		}
		assert.NotPanics(t, func() {
			ValidateProxyGroup(group)
		})
		
		ruleset := &RulesetConfig{
			Group: "test",
			URL:   "https://example.com/rules.list",
		}
		assert.NotPanics(t, func() {
			ValidateRuleset(ruleset)
		})
		
		request := &ConvertRequest{
			URL:    "https://example.com/subscription",
			Target: "clash",
		}
		assert.NotPanics(t, func() {
			ValidateConvertRequest(request)
		})
	})
}

func TestConverterOperations(t *testing.T) {
	t.Run("ProxyConverter Operations", func(t *testing.T) {
		converter := NewProxyConverter()
		require.NotNil(t, converter)
		
		proxy := &Proxy{
			Type:     ProxyTypeShadowsocks,
			Remark:   "test-proxy",
			Hostname: "example.com",
			Port:     443,
		}
		
		// 测试转换为 Map
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
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("Nil Object Handling", func(t *testing.T) {
		// 测试nil对象不会导致panic
		
		assert.NotPanics(t, func() {
			var nilProxy *Proxy
			_ = nilProxy
		})
		
		assert.NotPanics(t, func() {
			var nilGroup *ProxyGroupConfig
			_ = nilGroup
		})
		
		assert.NotPanics(t, func() {
			var nilRuleset *RulesetConfig
			_ = nilRuleset
		})
		
		// 测试空切片操作
		var emptyProxyList ProxyList
		assert.Equal(t, 0, len(emptyProxyList))
		
		var emptyGroupList ProxyGroupList
		assert.Equal(t, 0, len(emptyGroupList))
	})
}

func BenchmarkBasicOperations(b *testing.B) {
	proxy := &Proxy{
		Type:     ProxyTypeShadowsocks,
		Remark:   "benchmark-proxy",
		Hostname: "example.com",
		Port:     443,
	}
	
	converter := NewProxyConverter()
	
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
}