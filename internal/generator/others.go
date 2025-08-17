// internal/generator/others.go
package generator

import (
	"fmt"
	"strings"

	"subconverter-go/pkg/models"
)

// SurgeGenerator Surge配置生成器
type SurgeGenerator struct{}

func (g *SurgeGenerator) GetTarget() string { return "surge" }
func (g *SurgeGenerator) GetFormat() string { return "conf" }

func (g *SurgeGenerator) SupportsProxyType(proxyType models.ProxyType) bool {
	switch proxyType {
	case models.ProxyTypeShadowsocks, models.ProxyTypeVMess, models.ProxyTypeTrojan, models.ProxyTypeHTTP, models.ProxyTypeSOCKS5:
		return true
	default:
		return false
	}
}

func (g *SurgeGenerator) Validate(proxy *models.Proxy) error {
	if proxy == nil {
		return fmt.Errorf("proxy is nil")
	}
	if !g.SupportsProxyType(proxy.Type) {
		return fmt.Errorf("unsupported proxy type for surge: %s", proxy.Type.String())
	}
	return nil
}

func (g *SurgeGenerator) Generate(proxies []*models.Proxy, config *GenerateConfig) ([]byte, error) {
	var lines []string
	
	lines = append(lines, "[General]")
	lines = append(lines, "loglevel = notify")
	lines = append(lines, "")
	
	lines = append(lines, "[Proxy]")
	lines = append(lines, "DIRECT = direct")
	for _, proxy := range proxies {
		if line := g.convertProxy(proxy); line != "" {
			lines = append(lines, line)
		}
	}
	lines = append(lines, "")
	
	lines = append(lines, "[Proxy Group]")
	lines = append(lines, "PROXY = select, DIRECT")
	lines = append(lines, "")
	
	if config.EnableRule {
		lines = append(lines, "[Rule]")
		lines = append(lines, "DOMAIN-SUFFIX,local,DIRECT")
		lines = append(lines, "GEOIP,CN,DIRECT")
		lines = append(lines, "FINAL,PROXY")
	}
	
	return []byte(strings.Join(lines, "\n")), nil
}

func (g *SurgeGenerator) convertProxy(proxy *models.Proxy) string {
	name := proxy.Remark
	if name == "" {
		name = fmt.Sprintf("%s-%s-%d", proxy.Type.String(), proxy.Hostname, proxy.Port)
	}
	
	switch proxy.Type {
	case models.ProxyTypeShadowsocks:
		return fmt.Sprintf("%s = ss, %s, %d, encrypt-method=%s, password=%s", 
			name, proxy.Hostname, proxy.Port, proxy.EncryptMethod, proxy.Password)
	case models.ProxyTypeVMess:
		return fmt.Sprintf("%s = vmess, %s, %d, username=%s", 
			name, proxy.Hostname, proxy.Port, proxy.UserID)
	case models.ProxyTypeTrojan:
		return fmt.Sprintf("%s = trojan, %s, %d, password=%s", 
			name, proxy.Hostname, proxy.Port, proxy.Password)
	case models.ProxyTypeHTTP:
		if proxy.Username != "" && proxy.Password != "" {
			return fmt.Sprintf("%s = http, %s, %d, %s, %s", 
				name, proxy.Hostname, proxy.Port, proxy.Username, proxy.Password)
		}
		return fmt.Sprintf("%s = http, %s, %d", name, proxy.Hostname, proxy.Port)
	case models.ProxyTypeSOCKS5:
		if proxy.Username != "" && proxy.Password != "" {
			return fmt.Sprintf("%s = socks5, %s, %d, %s, %s", 
				name, proxy.Hostname, proxy.Port, proxy.Username, proxy.Password)
		}
		return fmt.Sprintf("%s = socks5, %s, %d", name, proxy.Hostname, proxy.Port)
	}
	return ""
}

// QuantumultXGenerator QuantumultX配置生成器
type QuantumultXGenerator struct{}

func (g *QuantumultXGenerator) GetTarget() string { return "quantumultx" }
func (g *QuantumultXGenerator) GetFormat() string { return "conf" }

func (g *QuantumultXGenerator) SupportsProxyType(proxyType models.ProxyType) bool {
	return proxyType == models.ProxyTypeShadowsocks || proxyType == models.ProxyTypeVMess || proxyType == models.ProxyTypeTrojan
}

func (g *QuantumultXGenerator) Validate(proxy *models.Proxy) error {
	if proxy == nil {
		return fmt.Errorf("proxy is nil")
	}
	if !g.SupportsProxyType(proxy.Type) {
		return fmt.Errorf("unsupported proxy type for quantumultx: %s", proxy.Type.String())
	}
	return nil
}

func (g *QuantumultXGenerator) Generate(proxies []*models.Proxy, config *GenerateConfig) ([]byte, error) {
	var lines []string
	
	lines = append(lines, "[general]")
	lines = append(lines, "")
	
	lines = append(lines, "[server_local]")
	for _, proxy := range proxies {
		if line := g.convertProxy(proxy); line != "" {
			lines = append(lines, line)
		}
	}
	lines = append(lines, "")
	
	if config.EnableRule {
		lines = append(lines, "[filter_local]")
		lines = append(lines, "host-suffix, local, direct")
		lines = append(lines, "geoip, cn, direct")
		lines = append(lines, "final, proxy")
	}
	
	return []byte(strings.Join(lines, "\n")), nil
}

func (g *QuantumultXGenerator) convertProxy(proxy *models.Proxy) string {
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
	}
	return ""
}

// 其他简化的生成器
type LoonGenerator struct{}
func (g *LoonGenerator) GetTarget() string { return "loon" }
func (g *LoonGenerator) GetFormat() string { return "conf" }
func (g *LoonGenerator) SupportsProxyType(proxyType models.ProxyType) bool { return false }
func (g *LoonGenerator) Validate(proxy *models.Proxy) error {
	return fmt.Errorf("loon generator not implemented")
}
func (g *LoonGenerator) Generate(proxies []*models.Proxy, config *GenerateConfig) ([]byte, error) {
	return nil, fmt.Errorf("loon generator not implemented")
}

type SingBoxGenerator struct{}
func (g *SingBoxGenerator) GetTarget() string { return "singbox" }
func (g *SingBoxGenerator) GetFormat() string { return "json" }
func (g *SingBoxGenerator) SupportsProxyType(proxyType models.ProxyType) bool { return false }
func (g *SingBoxGenerator) Validate(proxy *models.Proxy) error {
	return fmt.Errorf("singbox generator not implemented")
}
func (g *SingBoxGenerator) Generate(proxies []*models.Proxy, config *GenerateConfig) ([]byte, error) {
	return nil, fmt.Errorf("singbox generator not implemented")
}

type V2rayGenerator struct{}
func (g *V2rayGenerator) GetTarget() string { return "v2ray" }
func (g *V2rayGenerator) GetFormat() string { return "json" }
func (g *V2rayGenerator) SupportsProxyType(proxyType models.ProxyType) bool { return false }
func (g *V2rayGenerator) Validate(proxy *models.Proxy) error {
	return fmt.Errorf("v2ray generator not implemented")
}
func (g *V2rayGenerator) Generate(proxies []*models.Proxy, config *GenerateConfig) ([]byte, error) {
	return nil, fmt.Errorf("v2ray generator not implemented")
}

type SSGenerator struct{}
func (g *SSGenerator) GetTarget() string { return "ss" }
func (g *SSGenerator) GetFormat() string { return "json" }
func (g *SSGenerator) SupportsProxyType(proxyType models.ProxyType) bool { 
	return proxyType == models.ProxyTypeShadowsocks 
}
func (g *SSGenerator) Validate(proxy *models.Proxy) error {
	if proxy.Type != models.ProxyTypeShadowsocks {
		return fmt.Errorf("ss generator only supports shadowsocks")
	}
	return nil
}
func (g *SSGenerator) Generate(proxies []*models.Proxy, config *GenerateConfig) ([]byte, error) {
	return nil, fmt.Errorf("ss generator not implemented")
}