package template

import (
	"testing"

	"subconverter-go/pkg/models"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngine(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel) // 减少测试输出

	engine := NewEngine("../../templates", logger)

	// 测试加载模板
	err := engine.LoadTemplates()
	assert.NoError(t, err)

	// 测试获取模板名称
	names := engine.GetTemplateNames()
	assert.NotEmpty(t, names)

	// 测试模板存在性检查
	if len(names) > 0 {
		assert.True(t, engine.HasTemplate(names[0]))
	}
	assert.False(t, engine.HasTemplate("non-existent-template"))
}

func TestTemplateData(t *testing.T) {
	// 创建测试代理
	proxy := &models.Proxy{
		Type:          models.ProxyTypeShadowsocks,
		Hostname:      "127.0.0.1",
		Port:          8388,
		EncryptMethod: "aes-256-gcm",
		Password:      "test",
		Remark:        "Test SS",
	}

	config := map[string]interface{}{
		"target": "clash",
		"sort":   true,
	}

	data := NewTemplateData([]*models.Proxy{proxy}, config)

	require.NotNil(t, data)
	assert.Len(t, data.Proxies, 1)
	assert.Equal(t, config, data.Config)
	assert.Equal(t, 1, data.Meta["ProxyCount"])
	assert.Equal(t, true, data.Meta["Generated"])
}

func TestFormatProxyFunctions(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	engine := NewEngine("../../templates", logger)

	proxy := &models.Proxy{
		Type:          models.ProxyTypeShadowsocks,
		Hostname:      "127.0.0.1",
		Port:          8388,
		EncryptMethod: "aes-256-gcm",
		Password:      "test",
		Remark:        "Test SS",
	}

	// 测试Clash格式化
	clashResult := engine.formatClashProxy(proxy)
	assert.Contains(t, clashResult, "Test SS")
	assert.Contains(t, clashResult, "127.0.0.1")
	assert.Contains(t, clashResult, "8388")

	// 测试Surge格式化
	surgeResult := engine.formatSurgeProxy(proxy)
	assert.Contains(t, surgeResult, "Test SS")
	assert.Contains(t, surgeResult, "127.0.0.1")
	assert.Contains(t, surgeResult, "8388")

	// 测试QuantumultX格式化
	quantumultxResult := engine.formatQuantumultXProxy(proxy)
	assert.Contains(t, quantumultxResult, "Test SS")
	assert.Contains(t, quantumultxResult, "127.0.0.1")
	assert.Contains(t, quantumultxResult, "8388")
}

func TestFuncMap(t *testing.T) {
	logger := logrus.New()
	engine := NewEngine("../../templates", logger)
	funcMap := engine.getFuncMap()

	// 测试字符串函数
	upperFunc := funcMap["upper"].(func(string) string)
	assert.Equal(t, "HELLO", upperFunc("hello"))

	lowerFunc := funcMap["lower"].(func(string) string)
	assert.Equal(t, "hello", lowerFunc("HELLO"))

	// 测试数学函数
	addFunc := funcMap["add"].(func(int, int) int)
	assert.Equal(t, 5, addFunc(2, 3))

	// 测试比较函数
	eqFunc := funcMap["eq"].(func(interface{}, interface{}) bool)
	assert.True(t, eqFunc("test", "test"))
	assert.False(t, eqFunc("test", "other"))

	// 测试默认值函数
	defaultFunc := funcMap["default"].(func(interface{}, interface{}) interface{})
	assert.Equal(t, "default", defaultFunc("default", ""))
	assert.Equal(t, "value", defaultFunc("default", "value"))
}
