package template

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"subconverter-go/pkg/models"

	"github.com/sirupsen/logrus"
)

// Engine 模板引擎
type Engine struct {
	templates map[string]*template.Template
	logger    *logrus.Logger
	baseDir   string
}

// NewEngine 创建模板引擎
func NewEngine(baseDir string, logger *logrus.Logger) *Engine {
	return &Engine{
		templates: make(map[string]*template.Template),
		logger:    logger,
		baseDir:   baseDir,
	}
}

// LoadTemplates 加载所有模板
func (e *Engine) LoadTemplates() error {
	templateDirs := []string{"clash", "surge", "quantumultx", "loon", "singbox"}
	
	for _, dir := range templateDirs {
		templateDir := filepath.Join(e.baseDir, dir)
		err := e.loadTemplatesFromDir(templateDir, dir)
		if err != nil {
			e.logger.WithError(err).Warnf("Failed to load templates from %s", templateDir)
		}
	}
	
	return nil
}

// loadTemplatesFromDir 从目录加载模板
func (e *Engine) loadTemplatesFromDir(dir, prefix string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".tpl") {
			continue
		}

		templatePath := filepath.Join(dir, file.Name())
		templateName := fmt.Sprintf("%s/%s", prefix, strings.TrimSuffix(file.Name(), ".tpl"))
		
		content, err := ioutil.ReadFile(templatePath)
		if err != nil {
			e.logger.WithError(err).Warnf("Failed to read template file %s", templatePath)
			continue
		}

		tmpl, err := template.New(templateName).Funcs(e.getFuncMap()).Parse(string(content))
		if err != nil {
			e.logger.WithError(err).Warnf("Failed to parse template %s", templateName)
			continue
		}

		e.templates[templateName] = tmpl
		e.logger.Debugf("Loaded template: %s", templateName)
	}

	return nil
}

// Render 渲染模板
func (e *Engine) Render(templateName string, data interface{}) (string, error) {
	tmpl, exists := e.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template not found: %s", templateName)
	}

	var builder strings.Builder
	err := tmpl.Execute(&builder, data)
	if err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", templateName, err)
	}

	return builder.String(), nil
}

// HasTemplate 检查模板是否存在
func (e *Engine) HasTemplate(templateName string) bool {
	_, exists := e.templates[templateName]
	return exists
}

// GetTemplateNames 获取所有模板名称
func (e *Engine) GetTemplateNames() []string {
	var names []string
	for name := range e.templates {
		names = append(names, name)
	}
	return names
}

// getFuncMap 获取模板函数映射
func (e *Engine) getFuncMap() template.FuncMap {
	return template.FuncMap{
		"upper":    strings.ToUpper,
		"lower":    strings.ToLower,
		"title":    strings.Title,
		"join":     strings.Join,
		"contains": strings.Contains,
		"replace":  strings.Replace,
		"split":    strings.Split,
		"trim":     strings.TrimSpace,
		"quote":    func(s string) string { return fmt.Sprintf("\"%s\"", s) },
		"add":      func(a, b int) int { return a + b },
		"sub":      func(a, b int) int { return a - b },
		"mul":      func(a, b int) int { return a * b },
		"div":      func(a, b int) int { return a / b },
		"mod":      func(a, b int) int { return a % b },
		"eq":       func(a, b interface{}) bool { return a == b },
		"ne":       func(a, b interface{}) bool { return a != b },
		"lt":       func(a, b int) bool { return a < b },
		"le":       func(a, b int) bool { return a <= b },
		"gt":       func(a, b int) bool { return a > b },
		"ge":       func(a, b int) bool { return a >= b },
		"default":  func(defaultValue, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},
		"proxyType": func(proxy *models.Proxy) string {
			return proxy.Type.String()
		},
		"formatProxy": e.formatProxy,
	}
}

// formatProxy 格式化代理节点
func (e *Engine) formatProxy(proxy *models.Proxy, format string) string {
	switch format {
	case "clash":
		return e.formatClashProxy(proxy)
	case "surge":
		return e.formatSurgeProxy(proxy)
	case "quantumultx":
		return e.formatQuantumultXProxy(proxy)
	default:
		return fmt.Sprintf("%s://%s:%d", proxy.Type.String(), proxy.Hostname, proxy.Port)
	}
}

// formatClashProxy 格式化Clash代理
func (e *Engine) formatClashProxy(proxy *models.Proxy) string {
	switch proxy.Type {
	case models.ProxyTypeShadowsocks:
		return fmt.Sprintf("- {name: %s, type: ss, server: %s, port: %d, cipher: %s, password: %s}",
			proxy.Remark, proxy.Hostname, proxy.Port, proxy.EncryptMethod, proxy.Password)
	case models.ProxyTypeVMess:
		return fmt.Sprintf("- {name: %s, type: vmess, server: %s, port: %d, uuid: %s, alterId: %d, cipher: %s}",
			proxy.Remark, proxy.Hostname, proxy.Port, proxy.UserID, proxy.AlterID, proxy.EncryptMethod)
	case models.ProxyTypeTrojan:
		return fmt.Sprintf("- {name: %s, type: trojan, server: %s, port: %d, password: %s}",
			proxy.Remark, proxy.Hostname, proxy.Port, proxy.Password)
	default:
		return fmt.Sprintf("# Unsupported proxy type: %s", proxy.Type.String())
	}
}

// formatSurgeProxy 格式化Surge代理
func (e *Engine) formatSurgeProxy(proxy *models.Proxy) string {
	switch proxy.Type {
	case models.ProxyTypeShadowsocks:
		return fmt.Sprintf("%s = ss, %s, %d, encrypt-method=%s, password=%s",
			proxy.Remark, proxy.Hostname, proxy.Port, proxy.EncryptMethod, proxy.Password)
	case models.ProxyTypeVMess:
		return fmt.Sprintf("%s = vmess, %s, %d, username=%s",
			proxy.Remark, proxy.Hostname, proxy.Port, proxy.UserID)
	case models.ProxyTypeTrojan:
		return fmt.Sprintf("%s = trojan, %s, %d, password=%s",
			proxy.Remark, proxy.Hostname, proxy.Port, proxy.Password)
	default:
		return fmt.Sprintf("# Unsupported proxy type: %s", proxy.Type.String())
	}
}

// formatQuantumultXProxy 格式化QuantumultX代理
func (e *Engine) formatQuantumultXProxy(proxy *models.Proxy) string {
	switch proxy.Type {
	case models.ProxyTypeShadowsocks:
		return fmt.Sprintf("shadowsocks=%s:%d, method=%s, password=%s, tag=%s",
			proxy.Hostname, proxy.Port, proxy.EncryptMethod, proxy.Password, proxy.Remark)
	case models.ProxyTypeVMess:
		return fmt.Sprintf("vmess=%s:%d, method=chacha20-poly1305, password=%s, tag=%s",
			proxy.Hostname, proxy.Port, proxy.UserID, proxy.Remark)
	case models.ProxyTypeTrojan:
		return fmt.Sprintf("trojan=%s:%d, password=%s, tag=%s",
			proxy.Hostname, proxy.Port, proxy.Password, proxy.Remark)
	default:
		return fmt.Sprintf("# Unsupported proxy type: %s", proxy.Type.String())
	}
}

// TemplateData 模板数据结构
type TemplateData struct {
	Proxies []models.Proxy
	Config  map[string]interface{}
	Meta    map[string]interface{}
}

// NewTemplateData 创建模板数据
func NewTemplateData(proxies []*models.Proxy, config map[string]interface{}) *TemplateData {
	// 转换代理列表
	var proxyList []models.Proxy
	for _, proxy := range proxies {
		if proxy != nil {
			proxyList = append(proxyList, *proxy)
		}
	}

	return &TemplateData{
		Proxies: proxyList,
		Config:  config,
		Meta: map[string]interface{}{
			"ProxyCount": len(proxyList),
			"Generated":  true,
		},
	}
}
