// internal/config/validator.go
package config

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// SetupValidator 设置配置验证器
func SetupValidator() *validator.Validate {
	validate := validator.New()

	// 注册自定义验证器
	validate.RegisterValidation("host", validateHost)
	validate.RegisterValidation("regex", validateRegex)
	validate.RegisterValidation("template_file", validateTemplateFile)
	validate.RegisterValidation("config_file", validateConfigFile)
	validate.RegisterValidation("port", validatePort)

	return validate
}

// validateHost 验证主机地址
func validateHost(fl validator.FieldLevel) bool {
	host := fl.Field().String()
	if host == "" {
		return false
	}

	// 检查是否为有效IP地址
	if ip := net.ParseIP(host); ip != nil {
		return true
	}

	// 检查是否为有效域名
	if len(host) > 253 {
		return false
	}

	domainPattern := `^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`
	return regexp.MustCompile(domainPattern).MatchString(host)
}

// validateRegex 验证正则表达式
func validateRegex(fl validator.FieldLevel) bool {
	pattern := fl.Field().String()
	_, err := regexp.Compile(pattern)
	return err == nil
}

// validateTemplateFile 验证模板文件
func validateTemplateFile(fl validator.FieldLevel) bool {
	filename := fl.Field().String()
	if filename == "" {
		return false
	}

	// 检查文件扩展名
	validExts := []string{".tpl", ".j2", ".tmpl", ".template"}
	for _, ext := range validExts {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

// validateConfigFile 验证配置文件
func validateConfigFile(fl validator.FieldLevel) bool {
	filename := fl.Field().String()
	if filename == "" {
		return true // 允许空值
	}

	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	// 检查文件扩展名
	validExts := []string{".yaml", ".yml", ".json", ".toml", ".ini"}
	for _, ext := range validExts {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

// validatePort 验证端口号
func validatePort(fl validator.FieldLevel) bool {
	port := fl.Field().Int()
	return port >= 1 && port <= 65535
}

// ValidateNodeFilter 验证节点过滤器
func ValidateNodeFilter(filter *NodeFilter) error {
	if filter == nil {
		return fmt.Errorf("filter is nil")
	}

	if filter.Name == "" {
		return fmt.Errorf("filter name is required")
	}

	if filter.Type != "include" && filter.Type != "exclude" {
		return fmt.Errorf("filter type must be 'include' or 'exclude'")
	}

	if len(filter.Patterns) == 0 {
		return fmt.Errorf("filter patterns is empty")
	}

	// 如果启用正则表达式，验证每个模式
	if filter.Regex {
		for _, pattern := range filter.Patterns {
			if _, err := regexp.Compile(pattern); err != nil {
				return fmt.Errorf("invalid regex pattern '%s': %w", pattern, err)
			}
		}
	}

	return nil
}

// ValidateRenameRule 验证重命名规则
func ValidateRenameRule(rule *RenameRule) error {
	if rule == nil {
		return fmt.Errorf("rename rule is nil")
	}

	if rule.Name == "" {
		return fmt.Errorf("rename rule name is required")
	}

	if rule.Pattern == "" {
		return fmt.Errorf("rename rule pattern is required")
	}

	if rule.Replacement == "" {
		return fmt.Errorf("rename rule replacement is required")
	}

	// 如果启用正则表达式，验证模式
	if rule.Regex {
		if _, err := regexp.Compile(rule.Pattern); err != nil {
			return fmt.Errorf("invalid regex pattern '%s': %w", rule.Pattern, err)
		}
	}

	return nil
}

// ValidateRegionRule 验证地区分组规则
func ValidateRegionRule(rule *RegionRule) error {
	if rule == nil {
		return fmt.Errorf("region rule is nil")
	}

	if rule.Name == "" {
		return fmt.Errorf("region rule name is required")
	}

	if len(rule.Regions) == 0 {
		return fmt.Errorf("region rule regions is empty")
	}

	if len(rule.Patterns) == 0 {
		return fmt.Errorf("region rule patterns is empty")
	}

	return nil
}

// ValidateClientTemplate 验证客户端模板
func ValidateClientTemplate(template *ClientTemplate) error {
	if template == nil {
		return fmt.Errorf("client template is nil")
	}

	if template.Name == "" {
		return fmt.Errorf("client template name is required")
	}

	if template.File == "" {
		return fmt.Errorf("client template file is required")
	}

	if template.Type == "" {
		return fmt.Errorf("client template type is required")
	}

	// 验证客户端类型
	validTypes := []string{"clash", "surge", "quantumultx", "loon", "singbox", "v2ray", "trojan", "ss", "ssr"}
	validType := false
	for _, t := range validTypes {
		if template.Type == t {
			validType = true
			break
		}
	}

	if !validType {
		return fmt.Errorf("invalid client template type: %s", template.Type)
	}

	return nil
}