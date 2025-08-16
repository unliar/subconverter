// internal/api/handlers.go
package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"subconverter-go/pkg/models"
)

// handleSubscriptionConvert 处理订阅转换 - 兼容 C++ 版本的 /sub 接口
func (s *Server) handleSubscriptionConvert(c *gin.Context) {
	// 解析请求参数
	req := &models.ConvertRequest{
		URL:              c.Query("url"),
		Target:           c.Query("target"),
		Config:           c.Query("config"),
		Filename:         c.Query("filename"),
		Interval:         s.parseIntParam(c.Query("interval"), 0),
		Strict:           s.parseBoolParam(c.Query("strict"), false),
		Sort:             s.parseBoolParam(c.Query("sort"), false),
		FilterDeprecated: s.parseBoolParam(c.Query("fdn"), false),
		AppendType:       s.parseBoolParam(c.Query("append_type"), false),
		List:             s.parseBoolParam(c.Query("list"), false),
	}

	// 处理包含和排除过滤器
	if include := c.Query("include"); include != "" {
		req.IncludeFilters = strings.Split(include, "|")
	}
	if exclude := c.Query("exclude"); exclude != "" {
		req.ExcludeFilters = strings.Split(exclude, "|")
	}

	// 处理三态布尔参数
	req.UDP = s.parseTriStateBool(c.Query("udp"))
	req.TFO = s.parseTriStateBool(c.Query("tfo"))
	req.SkipCertVerify = s.parseTriStateBool(c.Query("scv"))
	req.TLS13 = s.parseTriStateBool(c.Query("tls13"))
	req.Emoji = s.parseTriStateBool(c.Query("emoji"))

	// 添加元数据
	req.RequestID = c.GetString("request_id")
	req.CreatedAt = time.Now()
	req.ClientIP = c.ClientIP()
	req.UserAgent = c.GetHeader("User-Agent")

	// 验证请求
	if !req.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters",
		})
		return
	}

	// 调用服务层处理转换
	result, err := s.converterSvc.ConvertSubscription(c.Request.Context(), req)
	if err != nil {
		s.logger.Errorf("Failed to convert subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 设置响应头
	for k, v := range result.Headers {
		c.Header(k, v)
	}

	// 根据目标类型设置 Content-Type
	contentType := s.getContentTypeForTarget(req.Target)
	c.Header("Content-Type", contentType)

	// 如果指定了文件名，设置下载头
	if req.Filename != "" {
		c.Header("Content-Disposition", "attachment; filename="+req.Filename)
	}

	// 返回结果
	c.String(http.StatusOK, result.Content)
}

// handleGetRuleset 处理获取规则集请求
func (s *Server) handleGetRuleset(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ruleset name is required",
		})
		return
	}

	// 调用服务层获取规则集
	ruleset, err := s.rulesetSvc.GetRuleset(c.Request.Context(), name)
	if err != nil {
		s.logger.Errorf("Failed to get ruleset: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Ruleset not found",
		})
		return
	}

	// 返回规则集内容
	c.String(http.StatusOK, ruleset.Content)
}

// handleRefreshRules 处理刷新规则请求
func (s *Server) handleRefreshRules(c *gin.Context) {
	// 调用服务层刷新所有规则集
	err := s.rulesetSvc.RefreshAllRulesets(c.Request.Context())
	if err != nil {
		s.logger.Errorf("Failed to refresh rulesets: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to refresh rulesets",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Rulesets refreshed successfully",
	})
}

// handleGetProfile 处理获取配置档案请求
func (s *Server) handleGetProfile(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Profile name is required",
		})
		return
	}

	// TODO: 实现配置档案获取逻辑
	c.JSON(http.StatusOK, gin.H{
		"name": name,
		"content": "# Profile content",
	})
}

// handleUpdateConfig 处理更新配置请求
func (s *Server) handleUpdateConfig(c *gin.Context) {
	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	// TODO: 实现配置更新逻辑
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Configuration updated",
	})
}

// handleReadConfig 处理读取配置请求
func (s *Server) handleReadConfig(c *gin.Context) {
	// 重新加载配置
	err := s.config.Reload(c.Request.Context())
	if err != nil {
		s.logger.Errorf("Failed to reload config: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to reload configuration",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Configuration reloaded",
	})
}

// handleShortURL 处理短链接获取
func (s *Server) handleShortURL(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Short URL ID is required",
		})
		return
	}

	// TODO: 实现短链接获取逻辑
	c.JSON(http.StatusOK, gin.H{
		"id": id,
		"url": "https://example.com/subscription",
	})
}

// handleCreateShortURL 处理创建短链接
func (s *Server) handleCreateShortURL(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// TODO: 实现短链接创建逻辑
	c.JSON(http.StatusOK, gin.H{
		"id": "abc123",
		"short_url": "https://short.link/abc123",
	})
}

// 辅助方法

// parseIntParam 解析整数参数
func (s *Server) parseIntParam(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}

// parseBoolParam 解析布尔参数
func (s *Server) parseBoolParam(value string, defaultValue bool) bool {
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}

// parseTriStateBool 解析三态布尔参数
func (s *Server) parseTriStateBool(value string) *bool {
	if value == "" {
		return nil
	}
	result := s.parseBoolParam(value, false)
	return &result
}

// getContentTypeForTarget 根据目标类型获取内容类型
func (s *Server) getContentTypeForTarget(target string) string {
	switch strings.ToLower(target) {
	case "clash", "clashr":
		return "application/x-yaml"
	case "surge", "loon", "quan", "quanx":
		return "text/plain"
	case "v2ray", "singbox":
		return "application/json"
	default:
		return "text/plain"
	}
}