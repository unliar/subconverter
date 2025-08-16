package models

import (
	"testing"
	"time"

	"subconverter-go/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxy_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		proxy   *Proxy
		want    bool
		wantErr bool
	}{
		{
			name: "valid shadowsocks proxy",
			proxy: &Proxy{
				Type:          types.ProxyTypeShadowsocks,
				Remark:        "test-ss",
				Hostname:      "example.com",
				Port:          443,
				Password:      "password123",
				EncryptMethod: "aes-256-gcm",
			},
			want: true,
		},
		{
			name: "valid vmess proxy",
			proxy: &Proxy{
				Type:     types.ProxyTypeVMess,
				Remark:   "test-vmess",
				Hostname: "example.com",
				Port:     443,
				UserID:   "12345678-1234-1234-1234-123456789abc",
			},
			want: true,
		},
		{
			name: "invalid proxy - missing hostname",
			proxy: &Proxy{
				Type:   types.ProxyTypeShadowsocks,
				Remark: "test-ss",
				Port:   443,
			},
			want: false,
		},
		{
			name: "invalid proxy - missing password for shadowsocks",
			proxy: &Proxy{
				Type:     types.ProxyTypeShadowsocks,
				Remark:   "test-ss",
				Hostname: "example.com",
				Port:     443,
			},
			want: false,
		},
		{
			name:  "nil proxy",
			proxy: nil,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.proxy.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProxy_Clone(t *testing.T) {
	udp := true
	tfo := false
	original := &Proxy{
		Type:          types.ProxyTypeShadowsocks,
		Remark:        "test-proxy",
		Hostname:      "example.com",
		Port:          443,
		Password:      "password123",
		EncryptMethod: "aes-256-gcm",
		UDP:           &udp,
		TCPFastOpen:   &tfo,
		CreatedAt:     time.Now(),
	}

	cloned := original.Clone()

	// 验证深拷贝成功
	assert.Equal(t, original.Type, cloned.Type)
	assert.Equal(t, original.Remark, cloned.Remark)
	assert.Equal(t, original.Hostname, cloned.Hostname)
	assert.Equal(t, original.Port, cloned.Port)
	assert.Equal(t, *original.UDP, *cloned.UDP)
	assert.Equal(t, *original.TCPFastOpen, *cloned.TCPFastOpen)

	// 验证指针是独立的
	assert.NotSame(t, original.UDP, cloned.UDP)
	assert.NotSame(t, original.TCPFastOpen, cloned.TCPFastOpen)

	// 修改克隆对象不应影响原对象
	*cloned.UDP = false
	assert.True(t, *original.UDP)
	assert.False(t, *cloned.UDP)
}

func TestProxy_SetDefaults(t *testing.T) {
	proxy := &Proxy{
		Type:     types.ProxyTypeShadowsocks,
		Hostname: "example.com",
		Remark:   "test",
	}

	proxy.SetDefaults()

	assert.Equal(t, uint16(443), proxy.Port)
	assert.NotEmpty(t, proxy.Group)
	assert.False(t, proxy.CreatedAt.IsZero())
	assert.False(t, proxy.UpdatedAt.IsZero())
}

func TestProxy_JSON(t *testing.T) {
	proxy := &Proxy{
		Type:          types.ProxyTypeShadowsocks,
		Remark:        "test-proxy",
		Hostname:      "example.com",
		Port:          443,
		Password:      "password123",
		EncryptMethod: "aes-256-gcm",
	}

	// 测试序列化
	data, err := proxy.ToJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// 测试反序列化
	decoded, err := FromJSON(data)
	require.NoError(t, err)
	assert.Equal(t, proxy.Type, decoded.Type)
	assert.Equal(t, proxy.Remark, decoded.Remark)
	assert.Equal(t, proxy.Hostname, decoded.Hostname)
	assert.Equal(t, proxy.Port, decoded.Port)
}

func TestProxyList_Methods(t *testing.T) {
	proxies := ProxyList{
		{
			Type:     types.ProxyTypeShadowsocks,
			Remark:   "ss-1",
			Hostname: "ss.example.com",
			Port:     443,
		},
		{
			Type:     types.ProxyTypeVMess,
			Remark:   "vmess-1",
			Hostname: "vmess.example.com",
			Port:     443,
		},
		{
			Type:     types.ProxyTypeShadowsocks,
			Remark:   "ss-2",
			Hostname: "ss2.example.com",
			Port:     443,
		},
	}

	// 测试长度
	assert.Equal(t, 3, proxies.Len())

	// 测试按类型过滤
	ssProxies := proxies.FilterByType(types.ProxyTypeShadowsocks)
	assert.Equal(t, 2, ssProxies.Len())

	vmessProxies := proxies.FilterByType(types.ProxyTypeVMess)
	assert.Equal(t, 1, vmessProxies.Len())

	// 测试按类型分组
	groups := proxies.GroupByType()
	assert.Len(t, groups, 2)
	assert.Len(t, groups[types.ProxyTypeShadowsocks], 2)
	assert.Len(t, groups[types.ProxyTypeVMess], 1)

	// 测试获取备注
	remarks := proxies.GetRemarks()
	expected := []string{"ss-1", "vmess-1", "ss-2"}
	assert.Equal(t, expected, remarks)
}

func TestProxy_GetDefaultGroup(t *testing.T) {
	tests := []struct {
		name      string
		proxyType types.ProxyType
		expected  string
	}{
		{
			name:      "shadowsocks",
			proxyType: types.ProxyTypeShadowsocks,
			expected:  "SSProvider",
		},
		{
			name:      "vmess",
			proxyType: types.ProxyTypeVMess,
			expected:  "V2RayProvider",
		},
		{
			name:      "trojan",
			proxyType: types.ProxyTypeTrojan,
			expected:  "TrojanProvider",
		},
		{
			name:      "unknown type",
			proxyType: types.ProxyTypeUnknown,
			expected:  "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proxy := &Proxy{Type: tt.proxyType}
			assert.Equal(t, tt.expected, proxy.GetDefaultGroup())
		})
	}
}

func TestProxy_HasFeature(t *testing.T) {
	proxy := &Proxy{Type: types.ProxyTypeShadowsocks}

	tests := []struct {
		feature  string
		expected bool
	}{
		{"udp", true},  // SS supports UDP
		{"tls", false}, // SS doesn't support TLS by default
		{"xyz", false}, // Unknown feature
	}

	for _, tt := range tests {
		t.Run(tt.feature, func(t *testing.T) {
			assert.Equal(t, tt.expected, proxy.HasFeature(tt.feature))
		})
	}
}

func TestProxy_GetKey(t *testing.T) {
	proxy := &Proxy{
		Type:     types.ProxyTypeShadowsocks,
		Hostname: "example.com",
		Port:     443,
	}

	key := proxy.GetKey()
	expected := "example.com:443:SS"
	assert.Equal(t, expected, key)
}

func TestProxy_IsSecure(t *testing.T) {
	tests := []struct {
		name     string
		proxy    *Proxy
		expected bool
	}{
		{
			name: "TLS enabled",
			proxy: &Proxy{
				TLSSecure: true,
			},
			expected: true,
		},
		{
			name: "TLS string set",
			proxy: &Proxy{
				TLSStr: "tls",
			},
			expected: true,
		},
		{
			name: "VMess (supports TLS)",
			proxy: &Proxy{
				Type: types.ProxyTypeVMess,
			},
			expected: true,
		},
		{
			name: "Shadowsocks (no TLS)",
			proxy: &Proxy{
				Type: types.ProxyTypeShadowsocks,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.proxy.IsSecure())
		})
	}
}

func BenchmarkProxy_Clone(b *testing.B) {
	udp := true
	proxy := &Proxy{
		Type:          types.ProxyTypeShadowsocks,
		Remark:        "test-proxy",
		Hostname:      "example.com",
		Port:          443,
		Password:      "password123",
		EncryptMethod: "aes-256-gcm",
		UDP:           &udp,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = proxy.Clone()
	}
}

func BenchmarkProxy_IsValid(b *testing.B) {
	proxy := &Proxy{
		Type:          types.ProxyTypeShadowsocks,
		Remark:        "test-proxy",
		Hostname:      "example.com",
		Port:          443,
		Password:      "password123",
		EncryptMethod: "aes-256-gcm",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = proxy.IsValid()
	}
}