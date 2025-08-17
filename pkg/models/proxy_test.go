package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyCreation(t *testing.T) {
	t.Run("Create Basic Proxy", func(t *testing.T) {
		proxy := &Proxy{
			Type:          ProxyTypeShadowsocks,
			Remark:        "test-ss",
			Hostname:      "example.com",
			Port:          443,
			Password:      "password123",
			EncryptMethod: "aes-256-gcm",
		}
		
		assert.NotNil(t, proxy)
		assert.Equal(t, ProxyTypeShadowsocks, proxy.Type)
		assert.Equal(t, "test-ss", proxy.Remark)
		assert.Equal(t, "example.com", proxy.Hostname)
		assert.Equal(t, uint16(443), proxy.Port)
		assert.Equal(t, "password123", proxy.Password)
		assert.Equal(t, "aes-256-gcm", proxy.EncryptMethod)
	})
}

func TestProxyTypeConstants(t *testing.T) {
	t.Run("Test Proxy Type Constants", func(t *testing.T) {
		assert.Equal(t, ProxyType(1), ProxyTypeShadowsocks)
		assert.Equal(t, ProxyType(2), ProxyTypeShadowsocksR)
		assert.Equal(t, ProxyType(3), ProxyTypeVMess)
		assert.Equal(t, ProxyType(4), ProxyTypeTrojan)
		assert.Equal(t, ProxyType(6), ProxyTypeHTTP)
		assert.Equal(t, ProxyType(8), ProxyTypeSOCKS5)
		assert.Equal(t, ProxyType(10), ProxyTypeVLESS)
	})
}

func TestProxyAllTypes(t *testing.T) {
	t.Run("Test All Proxy Types", func(t *testing.T) {
		tests := []struct {
			name      string
			proxyType ProxyType
		}{
			{"Shadowsocks", ProxyTypeShadowsocks},
			{"ShadowsocksR", ProxyTypeShadowsocksR},
			{"VMess", ProxyTypeVMess},
			{"VLESS", ProxyTypeVLESS},
			{"Trojan", ProxyTypeTrojan},
			{"HTTP", ProxyTypeHTTP},
			{"SOCKS5", ProxyTypeSOCKS5},
		}
		
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				proxy := &Proxy{
					Type:     tt.proxyType,
					Remark:   "test-" + tt.name,
					Hostname: "example.com",
					Port:     443,
				}
				
				assert.Equal(t, tt.proxyType, proxy.Type)
				assert.Equal(t, "test-"+tt.name, proxy.Remark)
				assert.Equal(t, "example.com", proxy.Hostname)
				assert.Equal(t, uint16(443), proxy.Port)
			})
		}
	})
}

func TestProxyListOperations(t *testing.T) {
	t.Run("Test ProxyList Operations", func(t *testing.T) {
		proxy1 := &Proxy{
			Type:     ProxyTypeShadowsocks,
			Remark:   "ss-1",
			Hostname: "ss.example.com",
			Port:     443,
		}
		
		proxy2 := &Proxy{
			Type:     ProxyTypeVMess,
			Remark:   "vmess-1",
			Hostname: "vmess.example.com",
			Port:     443,
		}
		
		proxy3 := &Proxy{
			Type:     ProxyTypeShadowsocks,
			Remark:   "ss-2",
			Hostname: "ss2.example.com",
			Port:     443,
		}
		
		proxies := ProxyList{proxy1, proxy2, proxy3}
		
		// 测试长度
		assert.Equal(t, 3, len(proxies))
		
		// 测试访问元素
		assert.Equal(t, proxy1, proxies[0])
		assert.Equal(t, proxy2, proxies[1])
		assert.Equal(t, proxy3, proxies[2])
		
		// 测试代理类型分布
		ssCount := 0
		vmessCount := 0
		for _, proxy := range proxies {
			switch proxy.Type {
			case ProxyTypeShadowsocks:
				ssCount++
			case ProxyTypeVMess:
				vmessCount++
			}
		}
		assert.Equal(t, 2, ssCount)
		assert.Equal(t, 1, vmessCount)
	})
}

func TestProxyFieldValidation(t *testing.T) {
	t.Run("Test Required Fields", func(t *testing.T) {
		// 测试基本字段设置
		proxy := &Proxy{}
		
		// 设置类型
		proxy.Type = ProxyTypeShadowsocks
		assert.Equal(t, ProxyTypeShadowsocks, proxy.Type)
		
		// 设置主机名
		proxy.Hostname = "test.example.com"
		assert.Equal(t, "test.example.com", proxy.Hostname)
		
		// 设置端口
		proxy.Port = 8080
		assert.Equal(t, uint16(8080), proxy.Port)
		
		// 设置备注
		proxy.Remark = "test-proxy"
		assert.Equal(t, "test-proxy", proxy.Remark)
		
		// 测试可选字段
		proxy.Group = "MyGroup"
		assert.Equal(t, "MyGroup", proxy.Group)
		
		proxy.Password = "secret"
		assert.Equal(t, "secret", proxy.Password)
		
		proxy.EncryptMethod = "aes-256-gcm"
		assert.Equal(t, "aes-256-gcm", proxy.EncryptMethod)
	})
}

func TestProxyAdvancedFields(t *testing.T) {
	t.Run("Test Optional Fields", func(t *testing.T) {
		proxy := &Proxy{
			Type:     ProxyTypeVMess,
			Hostname: "example.com",
			Port:     443,
			Remark:   "test-vmess",
		}
		
		// VMess特定字段
		proxy.UserID = "12345678-1234-1234-1234-123456789abc"
		assert.Equal(t, "12345678-1234-1234-1234-123456789abc", proxy.UserID)
		
		proxy.AlterID = 0
		assert.Equal(t, uint16(0), proxy.AlterID)
		
		proxy.Path = "/path"
		assert.Equal(t, "/path", proxy.Path)
		
		proxy.Host = "host.example.com"
		assert.Equal(t, "host.example.com", proxy.Host)
		
		// TLS相关字段
		proxy.TLSSecure = true
		assert.True(t, proxy.TLSSecure)
		
		proxy.SNI = "sni.example.com"
		assert.Equal(t, "sni.example.com", proxy.SNI)
	})
}

func TestProxyGroupListOperations(t *testing.T) {
	t.Run("Test ProxyGroupList", func(t *testing.T) {
		group1 := &ProxyGroupConfig{
			Name:    "group1",
			Proxies: []string{"proxy1", "proxy2"},
		}
		
		group2 := &ProxyGroupConfig{
			Name:    "group2", 
			Proxies: []string{"proxy3"},
		}
		
		groups := ProxyGroupList{group1, group2}
		
		assert.Equal(t, 2, len(groups))
		assert.Equal(t, "group1", groups[0].Name)
		assert.Equal(t, "group2", groups[1].Name)
		assert.Len(t, groups[0].Proxies, 2)
		assert.Len(t, groups[1].Proxies, 1)
	})
}

func TestNilHandling(t *testing.T) {
	t.Run("Test Nil Proxy Handling", func(t *testing.T) {
		// 测试nil代理不会导致panic
		var nilProxy *Proxy
		assert.Nil(t, nilProxy)
		
		// 测试空代理列表
		var emptyList ProxyList
		assert.Equal(t, 0, len(emptyList))
		
		// 测试空组列表
		var emptyGroups ProxyGroupList
		assert.Equal(t, 0, len(emptyGroups))
	})
}

func BenchmarkProxyCreation(b *testing.B) {
	b.Run("Create Proxy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = &Proxy{
				Type:          ProxyTypeShadowsocks,
				Remark:        "benchmark-proxy",
				Hostname:      "example.com",
				Port:          443,
				Password:      "password123",
				EncryptMethod: "aes-256-gcm",
			}
		}
	})
}

func BenchmarkProxyListOperations(b *testing.B) {
	// 创建测试数据
	proxies := make(ProxyList, 1000)
	for i := 0; i < 1000; i++ {
		proxies[i] = &Proxy{
			Type:     ProxyTypeShadowsocks,
			Remark:   "proxy-" + string(rune(i)),
			Hostname: "example.com",
			Port:     443,
		}
	}
	
	b.Run("Access ProxyList", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = proxies[i%len(proxies)]
		}
	})
	
	b.Run("Iterate ProxyList", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, proxy := range proxies {
				_ = proxy.Type
			}
		}
	})
}