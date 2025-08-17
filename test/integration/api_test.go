package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"subconverter-go/internal/api"
	"subconverter-go/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIIntegration(t *testing.T) {
	// 创建配置管理器
	configManager := config.NewManager("../../configs")
	
	// 创建API服务器
	server, err := api.NewServer(configManager)
	require.NoError(t, err)

	// 测试健康检查
	t.Run("HealthCheck", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "healthy", response["status"])
	})

	// 测试版本信息
	t.Run("Version", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/version", nil)
		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "name")
		assert.Contains(t, response, "version")
	})

	// 测试获取支持的目标客户端
	t.Run("GetSupportedTargets", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/targets", nil)
		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "targets")
		
		targets, ok := response["targets"].([]interface{})
		require.True(t, ok)
		assert.NotEmpty(t, targets)
	})

	// 测试订阅转换API（无效参数）
	t.Run("ConvertSubscription_InvalidParams", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/sub", nil)
		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "error")
	})

	// 测试获取规则集（无效参数）
	t.Run("GetRuleset_InvalidParams", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/getruleset", nil)
		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "error")
	})

	// 测试刷新规则
	t.Run("RefreshRules", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/refreshrules", nil)
		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "success", response["status"])
	})
}

func TestCORSHeaders(t *testing.T) {
	configManager := config.NewManager("../../configs")
	server, err := api.NewServer(configManager)
	require.NoError(t, err)

	// 测试OPTIONS请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/health", nil)
	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
}

func TestJSONPayload(t *testing.T) {
	configManager := config.NewManager("../../configs")
	server, err := api.NewServer(configManager)
	require.NoError(t, err)

	// 测试POST请求（更新配置）
	payload := map[string]interface{}{
		"test": "value",
	}
	
	jsonData, _ := json.Marshal(payload)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/updateconf", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "success", response["status"])
}
