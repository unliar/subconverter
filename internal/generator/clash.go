// internal/generator/clash.go
package generator

import (
	"fmt"
	"strings"

	"subconverter-go/pkg/models"

	"gopkg.in/yaml.v3"
)

// ClashGenerator Clash配置生成器
type ClashGenerator struct{}

// GetTarget 获取目标客户端类型
func (g *ClashGenerator) GetTarget() string {
	return "clash"
}

// GetFormat 获取输出格式
func (g *ClashGenerator) GetFormat() string {
	return "yaml"
}

// SupportsProxyType 检查是否支持指定的代理类型
func (g *ClashGenerator) SupportsProxyType(proxyType models.ProxyType) bool {
	switch proxyType {
	case models.ProxyTypeShadowsocks,
		models.ProxyTypeShadowsocksR,
		models.ProxyTypeVMess,
		models.ProxyTypeVLESS,
		models.ProxyTypeTrojan,
		models.ProxyTypeHTTP,
		models.ProxyTypeHTTPS,
		models.ProxyTypeSOCKS5:
		return true
	default:
		return false
	}
}

// Validate 验证代理节点是否兼容该生成器
func (g *ClashGenerator) Validate(proxy *models.Proxy) error {
	if proxy == nil {
		return fmt.Errorf("proxy is nil")
	}

	if !g.SupportsProxyType(proxy.Type) {
		return fmt.Errorf("unsupported proxy type for clash: %s", proxy.Type.String())
	}

	if proxy.Hostname == "" {
		return fmt.Errorf("hostname is required")
	}

	if proxy.Port == 0 {
		return fmt.Errorf("port is required")
	}

	return nil
}

// ClashConfig Clash配置结构
type ClashConfig struct {
	Port               int                      `yaml:"port"`
	SocksPort          int                      `yaml:"socks-port"`
	AllowLan           bool                     `yaml:"allow-lan"`
	Mode               string                   `yaml:"mode"`
	LogLevel           string                   `yaml:"log-level"`
	ExternalController string                   `yaml:"external-controller"`
	Proxies            []map[string]interface{} `yaml:"proxies"`
	ProxyGroups        []map[string]interface{} `yaml:"proxy-groups,omitempty"`
	Rules              []string                 `yaml:"rules,omitempty"`
}

// Generate 生成Clash配置
func (g *ClashGenerator) Generate(proxies []*models.Proxy, config *GenerateConfig) ([]byte, error) {
	if len(proxies) == 0 {
		return nil, fmt.Errorf("no proxies provided")
	}

	clashConfig := &ClashConfig{
		Port:               7890,
		SocksPort:          7891,
		AllowLan:           false,
		Mode:               "rule",
		LogLevel:           "info",
		ExternalController: "127.0.0.1:9090",
		Proxies:            make([]map[string]interface{}, 0, len(proxies)),
	}

	// 转换代理节点
	var proxyNames []string
	for _, proxy := range proxies {
		clashProxy, err := g.convertProxy(proxy, config)
		if err != nil {
			continue // 跳过转换失败的节点
		}
		clashConfig.Proxies = append(clashConfig.Proxies, clashProxy)
		proxyNames = append(proxyNames, clashProxy["name"].(string))
	}

	// 添加策略组
	if len(proxyNames) > 0 {
		clashConfig.ProxyGroups = g.generateProxyGroups(proxyNames)
	}

	// 添加基础规则
	if config.EnableRule {
		clashConfig.Rules = g.generateRules()
	}

	// 序列化为YAML
	return yaml.Marshal(clashConfig)
}

// convertProxy 转换代理节点为Clash格式
func (g *ClashGenerator) convertProxy(proxy *models.Proxy, config *GenerateConfig) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"name":   proxy.Remark,
		"server": proxy.Hostname,
		"port":   proxy.Port,
	}

	// 如果名称为空，使用默认格式
	if proxy.Remark == "" {
		result["name"] = fmt.Sprintf("%s-%s-%d", proxy.Type.String(), proxy.Hostname, proxy.Port)
	}

	switch proxy.Type {
	case models.ProxyTypeShadowsocks:
		result["type"] = "ss"
		result["cipher"] = proxy.EncryptMethod
		result["password"] = proxy.Password

	case models.ProxyTypeShadowsocksR:
		result["type"] = "ssr"
		result["cipher"] = proxy.EncryptMethod
		result["password"] = proxy.Password
		result["protocol"] = proxy.Protocol
		result["obfs"] = proxy.OBFS
		if proxy.ProtocolParam != "" {
			result["protocol-param"] = proxy.ProtocolParam
		}
		if proxy.OBFSParam != "" {
			result["obfs-param"] = proxy.OBFSParam
		}

	case models.ProxyTypeVMess:
		result["type"] = "vmess"
		result["uuid"] = proxy.UserID
		result["alterId"] = proxy.AlterID
		result["cipher"] = proxy.EncryptMethod
		if proxy.EncryptMethod == "" {
			result["cipher"] = "auto"
		}

		// 添加传输层配置
		if proxy.TransferProtocol != "" && proxy.TransferProtocol != "tcp" {
			result["network"] = proxy.TransferProtocol
		}

		if proxy.TLS == "tls" {
			result["tls"] = true
			if proxy.SNI != "" {
				result["servername"] = proxy.SNI
			}
		}

	case models.ProxyTypeTrojan:
		result["type"] = "trojan"
		result["password"] = proxy.Password
		result["sni"] = proxy.SNI
		if proxy.SNI == "" {
			result["sni"] = proxy.Hostname
		}

	case models.ProxyTypeHTTP:
		result["type"] = "http"
		if proxy.Username != "" {
			result["username"] = proxy.Username
		}
		if proxy.Password != "" {
			result["password"] = proxy.Password
		}

	case models.ProxyTypeSOCKS5:
		result["type"] = "socks5"
		if proxy.Username != "" {
			result["username"] = proxy.Username
		}
		if proxy.Password != "" {
			result["password"] = proxy.Password
		}

	default:
		return nil, fmt.Errorf("unsupported proxy type: %s", proxy.Type.String())
	}

	// 添加通用选项
	if config.UDP {
		result["udp"] = true
	}

	return result, nil
}

// generateProxyGroups 生成策略组
func (g *ClashGenerator) generateProxyGroups(proxyNames []string) []map[string]interface{} {
	var groups []map[string]interface{}

	// 手动选择组
	selectGroup := map[string]interface{}{
		"name":    "PROXY",
		"type":    "select",
		"proxies": append([]string{"Auto"}, proxyNames...),
	}
	groups = append(groups, selectGroup)

	// 自动选择组
	autoGroup := map[string]interface{}{
		"name":     "Auto",
		"type":     "url-test",
		"proxies":  proxyNames,
		"url":      "http://www.gstatic.com/generate_204",
		"interval": 300,
	}
	groups = append(groups, autoGroup)

	return groups
}

// generateRules 生成基础规则
func (g *ClashGenerator) generateRules() []string {
	return []string{
		"DOMAIN-SUFFIX,local,DIRECT",
		"IP-CIDR,127.0.0.0/8,DIRECT",
		"IP-CIDR,172.16.0.0/12,DIRECT",
		"IP-CIDR,192.168.0.0/16,DIRECT",
		"IP-CIDR,10.0.0.0/8,DIRECT",
		"IP-CIDR,17.0.0.0/8,DIRECT",
		"IP-CIDR,100.64.0.0/10,DIRECT",
		"GEOIP,CN,DIRECT",
		"MATCH,PROXY",
	}
}

// parsePluginOptions 解析插件选项（简化版）
func parsePluginOptions(options string) map[string]interface{} {
	result := make(map[string]interface{})
	parts := strings.Split(options, ";")
	for _, part := range parts {
		if kv := strings.SplitN(part, "=", 2); len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	return result
}