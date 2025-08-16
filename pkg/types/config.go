// Package types 定义配置相关的类型和枚举
package types

import (
	"fmt"
	"strings"
)

// ConfType 配置类型枚举 - 对应 C++ 版本的配置格式
type ConfType int

const (
	ConfTypeUnknown ConfType = iota
	ConfTypeSS
	ConfTypeSSR
	ConfTypeV2Ray
	ConfTypeSSConf
	ConfTypeSSTap
	ConfTypeNetch
	ConfTypeSOCKS
	ConfTypeHTTP
	ConfTypeSUB
	ConfTypeLocal
	ConfTypeClash
	ConfTypeSurge
	ConfTypeQuantumultX
	ConfTypeLoon
	ConfTypeSingBox
)

// 配置类型字符串映射
var confTypeStrings = map[ConfType]string{
	ConfTypeUnknown:     "Unknown",
	ConfTypeSS:          "SS",
	ConfTypeSSR:         "SSR",
	ConfTypeV2Ray:       "V2Ray",
	ConfTypeSSConf:      "SSConf",
	ConfTypeSSTap:       "SSTap",
	ConfTypeNetch:       "Netch",
	ConfTypeSOCKS:       "SOCKS",
	ConfTypeHTTP:        "HTTP",
	ConfTypeSUB:         "SUB",
	ConfTypeLocal:       "Local",
	ConfTypeClash:       "Clash",
	ConfTypeSurge:       "Surge",
	ConfTypeQuantumultX: "QuantumultX",
	ConfTypeLoon:        "Loon",
	ConfTypeSingBox:     "SingBox",
}

// 字符串到配置类型的映射
var stringToConfType = map[string]ConfType{
	"unknown":     ConfTypeUnknown,
	"ss":          ConfTypeSS,
	"ssr":         ConfTypeSSR,
	"v2ray":       ConfTypeV2Ray,
	"ssconf":      ConfTypeSSConf,
	"sstap":       ConfTypeSSTap,
	"netch":       ConfTypeNetch,
	"socks":       ConfTypeSOCKS,
	"http":        ConfTypeHTTP,
	"sub":         ConfTypeSUB,
	"local":       ConfTypeLocal,
	"clash":       ConfTypeClash,
	"surge":       ConfTypeSurge,
	"quantumultx": ConfTypeQuantumultX,
	"quanx":       ConfTypeQuantumultX,
	"loon":        ConfTypeLoon,
	"singbox":     ConfTypeSingBox,
}

// String 返回配置类型的字符串表示
func (ct ConfType) String() string {
	if str, ok := confTypeStrings[ct]; ok {
		return str
	}
	return "Unknown"
}

// IsValid 检查配置类型是否有效
func (ct ConfType) IsValid() bool {
	return ct > ConfTypeUnknown && ct <= ConfTypeSingBox
}

// IsClientType 检查是否为客户端配置类型
func (ct ConfType) IsClientType() bool {
	switch ct {
	case ConfTypeClash, ConfTypeSurge, ConfTypeQuantumultX, ConfTypeLoon, ConfTypeSingBox:
		return true
	default:
		return false
	}
}

// ParseConfType 从字符串解析配置类型
func ParseConfType(s string) ConfType {
	if ct, ok := stringToConfType[strings.ToLower(s)]; ok {
		return ct
	}
	return ConfTypeUnknown
}

// ProxyGroupType 代理组类型枚举
type ProxyGroupType int

const (
	ProxyGroupTypeSelect ProxyGroupType = iota
	ProxyGroupTypeURLTest
	ProxyGroupTypeFallback
	ProxyGroupTypeLoadBalance
	ProxyGroupTypeRelay
	ProxyGroupTypeSSID
	ProxyGroupTypeSmart
	ProxyGroupTypeAuto
)

// 代理组类型字符串映射
var proxyGroupTypeStrings = map[ProxyGroupType]string{
	ProxyGroupTypeSelect:      "select",
	ProxyGroupTypeURLTest:     "url-test",
	ProxyGroupTypeFallback:    "fallback",
	ProxyGroupTypeLoadBalance: "load-balance",
	ProxyGroupTypeRelay:       "relay",
	ProxyGroupTypeSSID:        "ssid",
	ProxyGroupTypeSmart:       "smart",
	ProxyGroupTypeAuto:        "auto",
}

// 字符串到代理组类型的映射
var stringToProxyGroupType = map[string]ProxyGroupType{
	"select":       ProxyGroupTypeSelect,
	"url-test":     ProxyGroupTypeURLTest,
	"fallback":     ProxyGroupTypeFallback,
	"load-balance": ProxyGroupTypeLoadBalance,
	"relay":        ProxyGroupTypeRelay,
	"ssid":         ProxyGroupTypeSSID,
	"smart":        ProxyGroupTypeSmart,
	"auto":         ProxyGroupTypeAuto,
}

// String 返回代理组类型的字符串表示
func (pgt ProxyGroupType) String() string {
	if str, ok := proxyGroupTypeStrings[pgt]; ok {
		return str
	}
	return "select"
}

// IsValid 检查代理组类型是否有效
func (pgt ProxyGroupType) IsValid() bool {
	return pgt >= ProxyGroupTypeSelect && pgt <= ProxyGroupTypeAuto
}

// RequiresURL 检查代理组类型是否需要 URL
func (pgt ProxyGroupType) RequiresURL() bool {
	switch pgt {
	case ProxyGroupTypeURLTest, ProxyGroupTypeFallback:
		return true
	default:
		return false
	}
}

// ParseProxyGroupType 从字符串解析代理组类型
func ParseProxyGroupType(s string) ProxyGroupType {
	if pgt, ok := stringToProxyGroupType[strings.ToLower(s)]; ok {
		return pgt
	}
	return ProxyGroupTypeSelect
}

// BalanceStrategy 负载均衡策略
type BalanceStrategy int

const (
	BalanceStrategyConsistentHashing BalanceStrategy = iota
	BalanceStrategyRoundRobin
	BalanceStrategyRandom
	BalanceStrategyLeastConnections
)

// 负载均衡策略字符串映射
var balanceStrategyStrings = map[BalanceStrategy]string{
	BalanceStrategyConsistentHashing: "consistent-hashing",
	BalanceStrategyRoundRobin:        "round-robin",
	BalanceStrategyRandom:            "random",
	BalanceStrategyLeastConnections:  "least-connections",
}

// String 返回负载均衡策略的字符串表示
func (bs BalanceStrategy) String() string {
	if str, ok := balanceStrategyStrings[bs]; ok {
		return str
	}
	return "round-robin"
}

// RulesetType 规则集类型枚举
type RulesetType int

const (
	RulesetTypeSurge RulesetType = iota
	RulesetTypeQuantumultX
	RulesetTypeClashDomain
	RulesetTypeClashIPCIDR
	RulesetTypeClashClassic
	RulesetTypeClashBehavior
)

// 规则集类型字符串映射
var rulesetTypeStrings = map[RulesetType]string{
	RulesetTypeSurge:         "surge",
	RulesetTypeQuantumultX:   "quantumultx",
	RulesetTypeClashDomain:   "clash-domain",
	RulesetTypeClashIPCIDR:   "clash-ipcidr",
	RulesetTypeClashClassic:  "clash-classic",
	RulesetTypeClashBehavior: "clash-behavior",
}

// String 返回规则集类型的字符串表示
func (rt RulesetType) String() string {
	if str, ok := rulesetTypeStrings[rt]; ok {
		return str
	}
	return "clash-domain"
}

// FilterType 过滤器类型
type FilterType int

const (
	FilterTypeInclude FilterType = iota
	FilterTypeExclude
)

// String 返回过滤器类型的字符串表示
func (ft FilterType) String() string {
	switch ft {
	case FilterTypeInclude:
		return "include"
	case FilterTypeExclude:
		return "exclude"
	default:
		return "include"
	}
}

// ParseFilterType 从字符串解析过滤器类型
func ParseFilterType(s string) FilterType {
	switch strings.ToLower(s) {
	case "include":
		return FilterTypeInclude
	case "exclude":
		return FilterTypeExclude
	default:
		return FilterTypeInclude
	}
}

// MarshalText 实现 encoding.TextMarshaler 接口
func (ct ConfType) MarshalText() ([]byte, error) {
	return []byte(ct.String()), nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler 接口
func (ct *ConfType) UnmarshalText(text []byte) error {
	*ct = ParseConfType(string(text))
	if *ct == ConfTypeUnknown && string(text) != "Unknown" {
		return fmt.Errorf("unknown config type: %s", string(text))
	}
	return nil
}

// MarshalText 实现 encoding.TextMarshaler 接口
func (pgt ProxyGroupType) MarshalText() ([]byte, error) {
	return []byte(pgt.String()), nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler 接口
func (pgt *ProxyGroupType) UnmarshalText(text []byte) error {
	*pgt = ParseProxyGroupType(string(text))
	return nil
}