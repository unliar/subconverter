package generator

import (
	"testing"

	"subconverter-go/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClashGenerator(t *testing.T) {
	generator := &ClashGenerator{}

	// 创建测试代理
	proxy := &models.Proxy{
		Type:          models.ProxyTypeShadowsocks,
		Hostname:      "127.0.0.1",
		Port:          8388,
		EncryptMethod: "aes-256-gcm",
		Password:      "test",
		Remark:        "Test SS",
	}

	// 测试代理类型支持
	assert.True(t, generator.SupportsProxyType(models.ProxyTypeShadowsocks))
	assert.True(t, generator.SupportsProxyType(models.ProxyTypeVMess))
	assert.True(t, generator.SupportsProxyType(models.ProxyTypeTrojan))
	assert.False(t, generator.SupportsProxyType(models.ProxyType(99)))

	// 测试验证
	err := generator.Validate(proxy)
	assert.NoError(t, err)

	// 测试生成配置
	config := &GenerateConfig{
		Target:     "clash",
		UDP:        true,
		EnableRule: true,
	}

	result, err := generator.Generate([]*models.Proxy{proxy}, config)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// 验证生成的YAML内容
	resultStr := string(result)
	assert.Contains(t, resultStr, "port: 7890")
	assert.Contains(t, resultStr, "proxies:")
	assert.Contains(t, resultStr, "Test SS")
	assert.Contains(t, resultStr, "127.0.0.1")
}

func TestSurgeGenerator(t *testing.T) {
	generator := &SurgeGenerator{}

	proxy := &models.Proxy{
		Type:          models.ProxyTypeShadowsocks,
		Hostname:      "127.0.0.1",
		Port:          8388,
		EncryptMethod: "aes-256-gcm",
		Password:      "test",
		Remark:        "Test SS",
	}

	// 测试代理类型支持
	assert.True(t, generator.SupportsProxyType(models.ProxyTypeShadowsocks))
	assert.False(t, generator.SupportsProxyType(models.ProxyTypeVLESS))

	// 测试验证
	err := generator.Validate(proxy)
	assert.NoError(t, err)

	// 测试生成配置
	config := &GenerateConfig{
		Target:     "surge",
		EnableRule: true,
	}

	result, err := generator.Generate([]*models.Proxy{proxy}, config)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// 验证生成的内容
	resultStr := string(result)
	assert.Contains(t, resultStr, "[General]")
	assert.Contains(t, resultStr, "[Proxy]")
	assert.Contains(t, resultStr, "Test SS")
}

func TestQuantumultXGenerator(t *testing.T) {
	generator := &QuantumultXGenerator{}

	proxy := &models.Proxy{
		Type:          models.ProxyTypeShadowsocks,
		Hostname:      "127.0.0.1",
		Port:          8388,
		EncryptMethod: "aes-256-gcm",
		Password:      "test",
		Remark:        "Test SS",
	}

	// 测试验证
	err := generator.Validate(proxy)
	assert.NoError(t, err)

	// 测试生成配置
	config := &GenerateConfig{
		Target:     "quantumultx",
		EnableRule: true,
	}

	result, err := generator.Generate([]*models.Proxy{proxy}, config)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// 验证生成的内容
	resultStr := string(result)
	assert.Contains(t, resultStr, "[general]")
	assert.Contains(t, resultStr, "[server_local]")
	assert.Contains(t, resultStr, "shadowsocks=")
}

func TestManager(t *testing.T) {
	manager := NewManager()

	// 测试获取支持的目标
	targets := manager.GetSupportedTargets()
	assert.NotEmpty(t, targets)
	assert.Contains(t, targets, "clash")
	assert.Contains(t, targets, "surge")
	assert.Contains(t, targets, "quantumultx")

	// 创建测试代理
	proxy := &models.Proxy{
		Type:          models.ProxyTypeShadowsocks,
		Hostname:      "127.0.0.1",
		Port:          8388,
		EncryptMethod: "aes-256-gcm",
		Password:      "test",
		Remark:        "Test SS",
	}

	// 测试生成配置
	config := &GenerateConfig{
		Target:     "clash",
		EnableRule: true,
		Sort:       true,
	}

	result, err := manager.Generate([]*models.Proxy{proxy}, config)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// 测试无效目标
	invalidConfig := &GenerateConfig{
		Target: "invalid",
	}

	_, err = manager.Generate([]*models.Proxy{proxy}, invalidConfig)
	assert.Error(t, err)
}

func TestManagerFiltering(t *testing.T) {
	manager := NewManager()

	// 创建多个测试代理
	proxies := []*models.Proxy{
		{
			Type:          models.ProxyTypeShadowsocks,
			Hostname:      "127.0.0.1",
			Port:          8388,
			EncryptMethod: "aes-256-gcm",
			Password:      "test",
			Remark:        "HK Test SS",
		},
		{
			Type:          models.ProxyTypeVMess,
			Hostname:      "127.0.0.2",
			Port:          8080,
			UserID:        "test-uuid",
			Remark:        "US Test VMess",
		},
		{
			Type:          models.ProxyTypeTrojan,
			Hostname:      "127.0.0.3",
			Port:          443,
			Password:      "test",
			Remark:        "SG Test Trojan",
		},
	}

	// 测试包含过滤
	config := &GenerateConfig{
		Target:         "clash",
		IncludeRemarks: []string{"HK", "US"},
		Sort:           true,
	}

	result, err := manager.Generate(proxies, config)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// 测试排除过滤
	excludeConfig := &GenerateConfig{
		Target:         "clash",
		ExcludeRemarks: []string{"SG"},
	}

	result, err = manager.Generate(proxies, excludeConfig)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestValidateConfig(t *testing.T) {
	manager := NewManager()

	// 测试有效配置
	validConfig := &GenerateConfig{
		Target: "clash",
	}
	err := manager.ValidateConfig(validConfig)
	assert.NoError(t, err)

	// 测试无效配置
	invalidConfig := &GenerateConfig{
		Target: "invalid",
	}
	err = manager.ValidateConfig(invalidConfig)
	assert.Error(t, err)

	// 测试空配置
	err = manager.ValidateConfig(nil)
	assert.Error(t, err)

	// 测试空目标
	emptyTargetConfig := &GenerateConfig{}
	err = manager.ValidateConfig(emptyTargetConfig)
	assert.Error(t, err)
}
