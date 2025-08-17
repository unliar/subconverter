// internal/parser/shadowsocks.go
package parser

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"subconverter-go/pkg/models"
)

// ShadowsocksParser Shadowsocks解析器
type ShadowsocksParser struct{}

// CanParse 检查是否能解析指定的链接
func (p *ShadowsocksParser) CanParse(link string) bool {
	return strings.HasPrefix(link, "ss://")
}

// GetType 获取解析器类型
func (p *ShadowsocksParser) GetType() models.ProxyType {
	return models.ProxyTypeShadowsocks
}

// Parse 解析Shadowsocks链接
func (p *ShadowsocksParser) Parse(link string) (*models.Proxy, error) {
	if !p.CanParse(link) {
		return nil, fmt.Errorf("not a shadowsocks link")
	}

	// 移除前缀
	link = strings.TrimPrefix(link, "ss://")

	// 解析URL
	u, err := url.Parse("ss://" + link)
	if err != nil {
		return nil, fmt.Errorf("invalid shadowsocks URL: %w", err)
	}

	proxy := &models.Proxy{
		Type:  models.ProxyTypeShadowsocks,
		Group: "ShadowsocksProvider",
	}

	// 解析主机和端口
	host := u.Hostname()
	portStr := u.Port()
	if host == "" || portStr == "" {
		return nil, fmt.Errorf("invalid host or port")
	}

	proxy.Hostname = host
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}
	proxy.Port = uint16(port)

	// 解析用户信息（method:password）
	userInfo := u.User
	if userInfo != nil {
		// 从用户信息解析 Base64 编码的 method:password
		encodedAuth := userInfo.Username()
		
		// Base64解码
		decodedAuth, err := base64.URLEncoding.DecodeString(encodedAuth)
		if err != nil {
			// 尝试标准Base64解码
			decodedAuth, err = base64.StdEncoding.DecodeString(encodedAuth)
			if err != nil {
				return nil, fmt.Errorf("failed to decode auth info: %w", err)
			}
		}
		
		// 解析 method:password
		authParts := strings.SplitN(string(decodedAuth), ":", 2)
		if len(authParts) == 2 {
			proxy.EncryptMethod = authParts[0]
			proxy.Password = authParts[1]
		} else {
			return nil, fmt.Errorf("invalid auth format")
		}
	} else {
		// 尝试从URL路径解析Base64编码的用户信息
		if u.Path != "" {
			decodedInfo, err := base64.URLEncoding.DecodeString(strings.TrimPrefix(u.Path, "/"))
			if err != nil {
				// 尝试标准Base64解码
				decodedInfo, err = base64.StdEncoding.DecodeString(strings.TrimPrefix(u.Path, "/"))
				if err != nil {
					return nil, fmt.Errorf("failed to decode user info: %w", err)
				}
			}
			
			parts := strings.SplitN(string(decodedInfo), "@", 2)
			if len(parts) == 2 {
				// 格式: method:password@server:port
				authParts := strings.SplitN(parts[0], ":", 2)
				if len(authParts) == 2 {
					proxy.EncryptMethod = authParts[0]
					proxy.Password = authParts[1]
				}
				
				// 更新服务器信息
				serverParts := strings.SplitN(parts[1], ":", 2)
				if len(serverParts) == 2 {
					proxy.Hostname = serverParts[0]
					port, err := strconv.ParseUint(serverParts[1], 10, 16)
					if err == nil {
						proxy.Port = uint16(port)
					}
				}
			}
		}
	}

	// 解析查询参数
	query := u.Query()
	
	// 解析名称
	if name := query.Get("name"); name != "" {
		proxy.Remark = name
	} else if fragment := u.Fragment; fragment != "" {
		// URL fragment作为名称
		proxy.Remark = fragment
	}

	// 解析分组
	if group := query.Get("group"); group != "" {
		proxy.Group = group
	}

	// 解析插件
	if plugin := query.Get("plugin"); plugin != "" {
		pluginParts := strings.SplitN(plugin, ";", 2)
		proxy.Plugin = pluginParts[0]
		if len(pluginParts) > 1 {
			proxy.PluginOption = pluginParts[1]
		}
	}

	// 验证必需字段
	if proxy.EncryptMethod == "" {
		return nil, fmt.Errorf("missing encrypt method")
	}
	if proxy.Password == "" {
		return nil, fmt.Errorf("missing password")
	}

	return proxy, nil
}

// ShadowsocksRParser ShadowsocksR解析器
type ShadowsocksRParser struct{}

// CanParse 检查是否能解析指定的链接
func (p *ShadowsocksRParser) CanParse(link string) bool {
	return strings.HasPrefix(link, "ssr://")
}

// GetType 获取解析器类型
func (p *ShadowsocksRParser) GetType() models.ProxyType {
	return models.ProxyTypeShadowsocksR
}

// Parse 解析ShadowsocksR链接
func (p *ShadowsocksRParser) Parse(link string) (*models.Proxy, error) {
	if !p.CanParse(link) {
		return nil, fmt.Errorf("not a shadowsocksr link")
	}

	// 移除前缀
	link = strings.TrimPrefix(link, "ssr://")

	// Base64解码
	decoded, err := base64.URLEncoding.DecodeString(link)
	if err != nil {
		// 尝试标准Base64解码
		decoded, err = base64.StdEncoding.DecodeString(link)
		if err != nil {
			return nil, fmt.Errorf("failed to decode ssr link: %w", err)
		}
	}

	// 解析格式: server:port:protocol:method:obfs:password_base64/?params
	parts := strings.Split(string(decoded), "/?")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ssr format")
	}

	// 解析主要部分
	mainParts := strings.Split(parts[0], ":")
	if len(mainParts) != 6 {
		return nil, fmt.Errorf("invalid ssr main format")
	}

	proxy := &models.Proxy{
		Type:     models.ProxyTypeShadowsocksR,
		Group:    "ShadowsocksRProvider",
		Hostname: mainParts[0],
	}

	// 解析端口
	port, err := strconv.ParseUint(mainParts[1], 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}
	proxy.Port = uint16(port)

	// 解析协议、方法、混淆
	proxy.Protocol = mainParts[2]
	proxy.EncryptMethod = mainParts[3]
	proxy.OBFS = mainParts[4]

	// 解析密码（Base64编码）
	passwordDecoded, err := base64.URLEncoding.DecodeString(mainParts[5])
	if err != nil {
		// 尝试标准Base64解码
		passwordDecoded, err = base64.StdEncoding.DecodeString(mainParts[5])
		if err != nil {
			return nil, fmt.Errorf("failed to decode password: %w", err)
		}
	}
	proxy.Password = string(passwordDecoded)

	// 解析查询参数
	if len(parts) > 1 {
		query, err := url.ParseQuery(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse query: %w", err)
		}

		// 解析混淆参数
		if obfsparam := query.Get("obfsparam"); obfsparam != "" {
			decoded, err := base64.URLEncoding.DecodeString(obfsparam)
			if err == nil {
				proxy.OBFSParam = string(decoded)
			}
		}

		// 解析协议参数
		if protoparam := query.Get("protoparam"); protoparam != "" {
			decoded, err := base64.URLEncoding.DecodeString(protoparam)
			if err == nil {
				proxy.ProtocolParam = string(decoded)
			}
		}

		// 解析备注
		if remarks := query.Get("remarks"); remarks != "" {
			decoded, err := base64.URLEncoding.DecodeString(remarks)
			if err == nil {
				proxy.Remark = string(decoded)
			}
		}

		// 解析分组
		if group := query.Get("group"); group != "" {
			decoded, err := base64.URLEncoding.DecodeString(group)
			if err == nil {
				proxy.Group = string(decoded)
			}
		}
	}

	return proxy, nil
}