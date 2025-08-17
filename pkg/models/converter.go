// Package models 定义数据转换器相关的功能
package models

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v3"
)

// ProxyConverter 代理转换器
type ProxyConverter struct{}

// NewProxyConverter 创建新的代理转换器
func NewProxyConverter() *ProxyConverter {
	return &ProxyConverter{}
}

// ToMap 将代理转换为 map，便于模板处理
func (pc *ProxyConverter) ToMap(proxy *Proxy) map[string]interface{} {
	if proxy == nil {
		return nil
	}

	result := make(map[string]interface{})

	// 基础字段
	result["type"] = proxy.Type.String()
	result["name"] = proxy.Remark
	result["server"] = proxy.Hostname
	result["port"] = proxy.Port

	// 根据代理类型添加特定字段
	switch proxy.Type {
	case ProxyTypeShadowsocks:
		pc.addShadowsocksFields(result, proxy)
	case ProxyTypeShadowsocksR:
		pc.addShadowsocksRFields(result, proxy)
	case ProxyTypeVMess, ProxyTypeVLESS:
		pc.addVMessFields(result, proxy)
	case ProxyTypeTrojan:
		pc.addTrojanFields(result, proxy)
	case ProxyTypeHTTP, ProxyTypeHTTPS:
		pc.addHTTPFields(result, proxy)
	case ProxyTypeSOCKS5:
		pc.addSOCKS5Fields(result, proxy)
	}

	// 通用选项
	pc.addCommonFields(result, proxy)

	return result
}

// addShadowsocksFields 添加 Shadowsocks 相关字段
func (pc *ProxyConverter) addShadowsocksFields(result map[string]interface{}, proxy *Proxy) {
	result["cipher"] = proxy.EncryptMethod
	result["password"] = proxy.Password

	if proxy.Plugin != "" {
		result["plugin"] = proxy.Plugin
		if proxy.PluginOption != "" {
			result["plugin-opts"] = pc.parsePluginOptions(proxy.PluginOption)
		}
	}
}

// addShadowsocksRFields 添加 ShadowsocksR 相关字段
func (pc *ProxyConverter) addShadowsocksRFields(result map[string]interface{}, proxy *Proxy) {
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
}

// addVMessFields 添加 VMess/VLESS 相关字段
func (pc *ProxyConverter) addVMessFields(result map[string]interface{}, proxy *Proxy) {
	result["uuid"] = proxy.UserID
	result["alterId"] = proxy.AlterID
	result["cipher"] = proxy.EncryptMethod

	// 网络配置
	if proxy.TransferProtocol != "" {
		result["network"] = proxy.TransferProtocol

		switch proxy.TransferProtocol {
		case "ws":
			if proxy.Path != "" {
				result["ws-path"] = proxy.Path
			}
			if proxy.Host != "" {
				result["ws-headers"] = map[string]string{"Host": proxy.Host}
			}
		case "grpc":
			if proxy.GRPCServiceName != "" {
				result["grpc-service-name"] = proxy.GRPCServiceName
			}
		}
	}

	// TLS 配置
	if proxy.TLS == "tls" {
		result["tls"] = true
		if proxy.SNI != "" {
			result["servername"] = proxy.SNI
		}
	}

	// VLESS 特有字段
	if proxy.Type == ProxyTypeVLESS {
		if proxy.Flow != "" {
			result["flow"] = proxy.Flow
		}
		// VLESS specific fields can be added here
	}
}

// addTrojanFields 添加 Trojan 相关字段
func (pc *ProxyConverter) addTrojanFields(result map[string]interface{}, proxy *Proxy) {
	result["password"] = proxy.Password

	if proxy.SNI != "" {
		result["sni"] = proxy.SNI
	}

	// 网络配置
	if proxy.TransferProtocol != "" {
		result["network"] = proxy.TransferProtocol

		if proxy.TransferProtocol == "grpc" && proxy.GRPCServiceName != "" {
			result["grpc-service-name"] = proxy.GRPCServiceName
		}
	}
}

// addHTTPFields 添加 HTTP/HTTPS 相关字段
func (pc *ProxyConverter) addHTTPFields(result map[string]interface{}, proxy *Proxy) {
	if proxy.Username != "" {
		result["username"] = proxy.Username
	}
	if proxy.Password != "" {
		result["password"] = proxy.Password
	}

	// HTTPS 特有配置
	if proxy.Type == ProxyTypeHTTPS {
		result["tls"] = true
		if proxy.SNI != "" {
			result["sni"] = proxy.SNI
		}
	}
}

// addSOCKS5Fields 添加 SOCKS5 相关字段
func (pc *ProxyConverter) addSOCKS5Fields(result map[string]interface{}, proxy *Proxy) {
	if proxy.Username != "" {
		result["username"] = proxy.Username
	}
	if proxy.Password != "" {
		result["password"] = proxy.Password
	}
}

// addCommonFields 添加通用字段
func (pc *ProxyConverter) addCommonFields(result map[string]interface{}, proxy *Proxy) {
	if proxy.UDP != nil {
		result["udp"] = *proxy.UDP
	}
	if proxy.TCPFastOpen != nil {
		result["tcp-fast-open"] = *proxy.TCPFastOpen
	}
	if proxy.AllowInsecure != nil {
		result["skip-cert-verify"] = *proxy.AllowInsecure
	}
}

// parsePluginOptions 解析插件选项
func (pc *ProxyConverter) parsePluginOptions(opts string) map[string]interface{} {
	result := make(map[string]interface{})
	
	// 简化处理，实际应该根据插件类型来解析
	pairs := strings.Split(opts, ";")
	for _, pair := range pairs {
		if kv := strings.SplitN(pair, "=", 2); len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	
	return result
}

// ToJSON 将代理转换为 JSON
func (pc *ProxyConverter) ToJSON(proxy *Proxy) ([]byte, error) {
	return json.Marshal(proxy)
}

// ToYAML 将代理转换为 YAML
func (pc *ProxyConverter) ToYAML(proxy *Proxy) ([]byte, error) {
	return yaml.Marshal(proxy)
}

// FromJSON 从 JSON 创建代理
func (pc *ProxyConverter) FromJSON(data []byte) (*Proxy, error) {
	var proxy Proxy
	err := json.Unmarshal(data, &proxy)
	return &proxy, err
}

// FromYAML 从 YAML 创建代理
func (pc *ProxyConverter) FromYAML(data []byte) (*Proxy, error) {
	var proxy Proxy
	err := yaml.Unmarshal(data, &proxy)
	return &proxy, err
}

// ProxyListToMaps 将代理列表转换为 map 列表
func (pc *ProxyConverter) ProxyListToMaps(proxies ProxyList) []map[string]interface{} {
	result := make([]map[string]interface{}, len(proxies))
	for i, proxy := range proxies {
		result[i] = pc.ToMap(proxy)
	}
	return result
}

// FormatForTarget 根据目标客户端格式化代理
func (pc *ProxyConverter) FormatForTarget(proxy *Proxy, target string) map[string]interface{} {
	baseMap := pc.ToMap(proxy)
	if baseMap == nil {
		return nil
	}

	// 根据不同目标客户端调整字段名和格式
	switch target {
	case "clash", "clashr":
		return pc.formatForClash(baseMap, proxy)
	case "surge":
		return pc.formatForSurge(baseMap, proxy)
	case "quan", "quanx":
		return pc.formatForQuantumultX(baseMap, proxy)
	case "loon":
		return pc.formatForLoon(baseMap, proxy)
	default:
		return baseMap
	}
}

// formatForClash 格式化为 Clash 格式
func (pc *ProxyConverter) formatForClash(baseMap map[string]interface{}, proxy *Proxy) map[string]interface{} {
	// Clash 格式特殊处理
	result := make(map[string]interface{})
	for k, v := range baseMap {
		result[k] = v
	}
	return result
}

// formatForSurge 格式化为 Surge 格式
func (pc *ProxyConverter) formatForSurge(baseMap map[string]interface{}, proxy *Proxy) map[string]interface{} {
	// Surge 格式特殊处理
	result := make(map[string]interface{})
	for k, v := range baseMap {
		result[k] = v
	}
	return result
}

// formatForQuantumultX 格式化为 QuantumultX 格式
func (pc *ProxyConverter) formatForQuantumultX(baseMap map[string]interface{}, proxy *Proxy) map[string]interface{} {
	// QuantumultX 格式特殊处理
	result := make(map[string]interface{})
	for k, v := range baseMap {
		result[k] = v
	}
	return result
}

// formatForLoon 格式化为 Loon 格式
func (pc *ProxyConverter) formatForLoon(baseMap map[string]interface{}, proxy *Proxy) map[string]interface{} {
	// Loon 格式特殊处理
	result := make(map[string]interface{})
	for k, v := range baseMap {
		result[k] = v
	}
	return result
}