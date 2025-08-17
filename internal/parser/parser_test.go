package parser

import (
	"testing"

	"subconverter-go/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShadowsocksParser(t *testing.T) {
	parser := &ShadowsocksParser{}
	
	tests := []struct {
		name     string
		link     string
		expected *models.Proxy
		wantErr  bool
	}{
		{
			name: "valid shadowsocks link",
			link: "ss://YWVzLTI1Ni1nY206dGVzdA==@127.0.0.1:8388#Test%20SS",
			expected: &models.Proxy{
				Type:          models.ProxyTypeShadowsocks,
				Hostname:      "127.0.0.1",
				Port:          uint16(8388),
				EncryptMethod: "aes-256-gcm",
				Password:      "test",
				Remark:        "Test SS",
			},
			wantErr: false,
		},
		{
			name:     "invalid shadowsocks link",
			link:     "ss://invalid",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !parser.CanParse(tt.link) && !tt.wantErr {
				t.Errorf("CanParse() should return true for valid link")
				return
			}

			result, err := parser.Parse(tt.link)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.expected.Type, result.Type)
				assert.Equal(t, tt.expected.Hostname, result.Hostname)
				assert.Equal(t, tt.expected.Port, result.Port)
				assert.Equal(t, tt.expected.EncryptMethod, result.EncryptMethod)
				assert.Equal(t, tt.expected.Password, result.Password)
				assert.Equal(t, tt.expected.Remark, result.Remark)
			}
		})
	}
}

func TestVMeSSParser(t *testing.T) {
	parser := &VMeSSParser{}
	
	// 测试有效的VMess链接
	link := "vmess://eyJ2IjoiMiIsInBzIjoidGVzdCIsImFkZCI6IjEyNy4wLjAuMSIsInBvcnQiOiI4MDgwIiwiaWQiOiI1NWY3NGExNi1mNjc2LTQxNTktOGEzYy0yN2YwZTEzNGY1YjAiLCJhaWQiOiIwIiwibmV0IjoidGNwIiwidHlwZSI6Im5vbmUiLCJob3N0IjoiIiwicGF0aCI6IiIsInRscyI6IiJ9"
	
	assert.True(t, parser.CanParse(link))
	
	result, err := parser.Parse(link)
	require.NoError(t, err)
	require.NotNil(t, result)
	
	assert.Equal(t, models.ProxyTypeVMess, result.Type)
	assert.Equal(t, "127.0.0.1", result.Hostname)
	assert.Equal(t, uint16(8080), result.Port)
	assert.Equal(t, "55f74a16-f676-4159-8a3c-27f0e134f5b0", result.UserID)
}

func TestTrojanParser(t *testing.T) {
	parser := &TrojanParser{}
	
	link := "trojan://password@127.0.0.1:443#Test%20Trojan"
	
	assert.True(t, parser.CanParse(link))
	
	result, err := parser.Parse(link)
	require.NoError(t, err)
	require.NotNil(t, result)
	
	assert.Equal(t, models.ProxyTypeTrojan, result.Type)
	assert.Equal(t, "127.0.0.1", result.Hostname)
	assert.Equal(t, uint16(443), result.Port)
	assert.Equal(t, "password", result.Password)
	assert.Equal(t, "Test Trojan", result.Remark)
}

func TestManager(t *testing.T) {
	manager := NewManager()
	
	// 测试解析单个链接
	link := "ss://YWVzLTI1Ni1nY206dGVzdA==@127.0.0.1:8388#Test%20SS"
	result, err := manager.Parse(link)
	
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, models.ProxyTypeShadowsocks, result.Type)
	
	// 测试解析多个链接
	links := []string{
		"ss://YWVzLTI1Ni1nY206dGVzdA==@127.0.0.1:8388#Test%20SS",
		"trojan://password@127.0.0.1:443#Test%20Trojan",
	}
	
	results, err := manager.ParseMultiple(links)
	require.NoError(t, err)
	assert.Len(t, results, 2)
	
	// 测试解析订阅内容
	subscription := `ss://YWVzLTI1Ni1nY206dGVzdA==@127.0.0.1:8388#Test%20SS
trojan://password@127.0.0.1:443#Test%20Trojan`
	
	subResults, err := manager.ParseSubscription(subscription)
	require.NoError(t, err)
	assert.Len(t, subResults, 2)
}

func TestManagerWithInvalidInput(t *testing.T) {
	manager := NewManager()
	
	// 测试无效链接
	_, err := manager.Parse("invalid://link")
	assert.Error(t, err)
	
	// 测试空订阅
	_, err = manager.ParseSubscription("")
	assert.Error(t, err)
	
	// 测试空链接列表
	results, err := manager.ParseMultiple([]string{})
	assert.NoError(t, err)
	assert.Empty(t, results)
}
