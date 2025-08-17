package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// handleSubscriptionConvert 处理订阅转换
func (s *Server) handleSubscriptionConvert(c *gin.Context) {
	// 获取基本参数
	subscriptionURL := c.Query("url")
	target := c.Query("target")

	if subscriptionURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Subscription URL is required",
		})
		return
	}

	if target == "" {
		target = "clash" // 默认值
	}

	// 构建选项映射
	options := make(map[string]string)
	
	// 解析所有查询参数作为选项
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			options[key] = values[0]
		}
	}

	// 验证订阅链接
	if err := s.converterSvc.ValidateSubscription(subscriptionURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 调用服务层处理转换
	result, err := s.converterSvc.ConvertSubscription(subscriptionURL, target, options)
	if err != nil {
		s.logger.WithError(err).Error("Failed to convert subscription")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 设置响应头
	contentType := s.getContentTypeForTarget(target)
	c.Header("Content-Type", contentType)

	// 如果指定了文件名，设置下载头
	if filename := c.Query("filename"); filename != "" {
		c.Header("Content-Disposition", "attachment; filename="+filename)
	}

	// 返回结果
	c.Data(http.StatusOK, contentType, result)
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
	ruleset, err := s.rulesetSvc.GetRuleset(name)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get ruleset")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Ruleset not found",
		})
		return
	}

	// 返回规则集内容
	c.Data(http.StatusOK, "text/plain", ruleset)
}

// handleRefreshRules 处理刷新规则请求
func (s *Server) handleRefreshRules(c *gin.Context) {
	// 调用服务层刷新规则集
	err := s.rulesetSvc.RefreshRules()
	if err != nil {
		s.logger.WithError(err).Error("Failed to refresh rulesets")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to refresh rulesets",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
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
		"name":    name,
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
		"status":  "success",
		"message": "Configuration updated",
	})
}

// handleReadConfig 处理读取配置请求
func (s *Server) handleReadConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
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
		"id":  id,
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
		"id":        "abc123",
		"short_url": "https://short.link/abc123",
	})
}

// handleGetSupportedTargets 处理获取支持的目标客户端
func (s *Server) handleGetSupportedTargets(c *gin.Context) {
	targets := s.converterSvc.GetSupportedTargets()
	c.JSON(http.StatusOK, gin.H{
		"targets": targets,
	})
}

// handleGetAvailableRulesets 处理获取可用规则集列表
func (s *Server) handleGetAvailableRulesets(c *gin.Context) {
	rulesets, err := s.rulesetSvc.GetAvailableRulesets()
	if err != nil {
		s.logger.WithError(err).Error("Failed to get available rulesets")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get rulesets",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rulesets": rulesets,
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
