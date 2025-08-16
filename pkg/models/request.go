// Package models 定义请求响应相关的数据模型
package models

import (
	"subconverter-go/pkg/constants"
	"subconverter-go/pkg/types"
	"time"
)

// ConvertRequest 转换请求 - 完全兼容 C++ 版本的请求参数
type ConvertRequest struct {
	// 基础参数
	URL      string `json:"url" form:"url" validate:"required,url"`
	Target   string `json:"target" form:"target" validate:"required"`
	Config   string `json:"config,omitempty" form:"config,omitempty"`
	Filename string `json:"filename,omitempty" form:"filename,omitempty"`
	Interval int    `json:"interval,omitempty" form:"interval,omitempty"`
	Strict   bool   `json:"strict,omitempty" form:"strict,omitempty"`

	// 过滤参数
	IncludeFilters []string `json:"include,omitempty" form:"include,omitempty"`
	ExcludeFilters []string `json:"exclude,omitempty" form:"exclude,omitempty"`

	// 功能开关
	Sort             bool `json:"sort,omitempty" form:"sort,omitempty"`
	FilterDeprecated bool `json:"fdn,omitempty" form:"fdn,omitempty"`
	AppendType       bool `json:"append_type,omitempty" form:"append_type,omitempty"`
	List             bool `json:"list,omitempty" form:"list,omitempty"`

	// 特性开关 - 使用指针实现三态逻辑
	UDP            *bool `json:"udp,omitempty" form:"udp,omitempty"`
	TFO            *bool `json:"tfo,omitempty" form:"tfo,omitempty"`
	SkipCertVerify *bool `json:"scv,omitempty" form:"scv,omitempty"`
	TLS13          *bool `json:"tls13,omitempty" form:"tls13,omitempty"`
	Emoji          *bool `json:"emoji,omitempty" form:"emoji,omitempty"`

	// 元数据
	RequestID  string    `json:"request_id,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	ClientIP   string    `json:"client_ip,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
	AccessToken string   `json:"access_token,omitempty"`
}

// IsValid 验证转换请求是否有效
func (cr *ConvertRequest) IsValid() bool {
	if cr == nil || cr.URL == "" || cr.Target == "" {
		return false
	}

	// 验证目标类型
	for _, target := range constants.SupportedTargets {
		if cr.Target == target {
			return true
		}
	}

	return false
}

// GetSafeTarget 获取安全的目标类型
func (cr *ConvertRequest) GetSafeTarget() string {
	if cr.IsValid() {
		return cr.Target
	}
	return "clash" // 默认使用 clash
}

// SetDefaults 设置默认值
func (cr *ConvertRequest) SetDefaults() {
	if cr.Interval <= 0 {
		cr.Interval = constants.DefaultRulesetInterval
	}
	if cr.CreatedAt.IsZero() {
		cr.CreatedAt = time.Now()
	}
}

// GetContentType 根据目标类型获取内容类型
func (cr *ConvertRequest) GetContentType() string {
	switch cr.Target {
	case "clash", "clashr":
		return constants.ContentTypeYAML
	case "surge", "surfboard":
		return constants.ContentTypeConf
	case "quan", "quanx", "loon":
		return constants.ContentTypeConf
	case "ss", "ssr", "v2ray", "trojan":
		return constants.ContentTypeText
	case "singbox":
		return constants.ContentTypeJSON
	default:
		return constants.ContentTypeYAML
	}
}

// GetFilename 获取默认文件名
func (cr *ConvertRequest) GetFilename() string {
	if cr.Filename != "" {
		return cr.Filename
	}

	switch cr.Target {
	case "clash", "clashr":
		return "config.yaml"
	case "surge":
		return "surge.conf"
	case "quan":
		return "quan.conf"
	case "quanx":
		return "quantumultx.conf"
	case "loon":
		return "loon.conf"
	case "surfboard":
		return "surfboard.conf"
	case "ss":
		return "shadowsocks.txt"
	case "ssr":
		return "shadowsocksr.txt"
	case "v2ray":
		return "v2ray.txt"
	case "trojan":
		return "trojan.txt"
	case "singbox":
		return "singbox.json"
	default:
		return "config.yaml"
	}
}

// Clone 深拷贝请求对象
func (cr *ConvertRequest) Clone() *ConvertRequest {
	if cr == nil {
		return nil
	}

	clone := *cr

	// 深拷贝切片
	if cr.IncludeFilters != nil {
		clone.IncludeFilters = make([]string, len(cr.IncludeFilters))
		copy(clone.IncludeFilters, cr.IncludeFilters)
	}
	if cr.ExcludeFilters != nil {
		clone.ExcludeFilters = make([]string, len(cr.ExcludeFilters))
		copy(clone.ExcludeFilters, cr.ExcludeFilters)
	}

	// 深拷贝指针字段
	if cr.UDP != nil {
		udp := *cr.UDP
		clone.UDP = &udp
	}
	if cr.TFO != nil {
		tfo := *cr.TFO
		clone.TFO = &tfo
	}
	if cr.SkipCertVerify != nil {
		scv := *cr.SkipCertVerify
		clone.SkipCertVerify = &scv
	}
	if cr.TLS13 != nil {
		tls13 := *cr.TLS13
		clone.TLS13 = &tls13
	}
	if cr.Emoji != nil {
		emoji := *cr.Emoji
		clone.Emoji = &emoji
	}

	return &clone
}

// ConvertResponse 转换响应
type ConvertResponse struct {
	// 响应内容
	Content     string            `json:"content,omitempty"`
	ContentType string            `json:"content_type,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`

	// 状态信息
	StatusCode int    `json:"status_code"`
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`

	// 统计信息
	ProxyCount int `json:"proxy_count,omitempty"`
	GroupCount int `json:"group_count,omitempty"`
	RuleCount  int `json:"rule_count,omitempty"`

	// 元数据
	RequestID   string        `json:"request_id,omitempty"`
	ProcessTime time.Duration `json:"process_time,omitempty"`
	GeneratedAt time.Time     `json:"generated_at,omitempty"`
}

// IsSuccess 检查响应是否成功
func (cr *ConvertResponse) IsSuccess() bool {
	return cr != nil && cr.Success && cr.StatusCode >= 200 && cr.StatusCode < 300
}

// SetHeaders 设置响应头
func (cr *ConvertResponse) SetHeaders(headers map[string]string) {
	if cr.Headers == nil {
		cr.Headers = make(map[string]string)
	}
	for k, v := range headers {
		cr.Headers[k] = v
	}
}

// AddHeader 添加响应头
func (cr *ConvertResponse) AddHeader(key, value string) {
	if cr.Headers == nil {
		cr.Headers = make(map[string]string)
	}
	cr.Headers[key] = value
}

// SetError 设置错误响应
func (cr *ConvertResponse) SetError(err *types.ConvertError) {
	cr.Success = false
	cr.StatusCode = err.Code.HTTPStatus()
	cr.Message = err.Message
	if cr.RequestID == "" && err.RequestID != "" {
		cr.RequestID = err.RequestID
	}
}

// SetSuccess 设置成功响应
func (cr *ConvertResponse) SetSuccess(content string, stats map[string]int) {
	cr.Success = true
	cr.StatusCode = 200
	cr.Content = content
	cr.GeneratedAt = time.Now()

	if stats != nil {
		if count, ok := stats["proxy"]; ok {
			cr.ProxyCount = count
		}
		if count, ok := stats["group"]; ok {
			cr.GroupCount = count
		}
		if count, ok := stats["rule"]; ok {
			cr.RuleCount = count
		}
	}
}

// GetSize 获取响应内容大小
func (cr *ConvertResponse) GetSize() int {
	return len(cr.Content)
}

// Clone 深拷贝响应对象
func (cr *ConvertResponse) Clone() *ConvertResponse {
	if cr == nil {
		return nil
	}

	clone := *cr

	// 深拷贝 map
	if cr.Headers != nil {
		clone.Headers = make(map[string]string)
		for k, v := range cr.Headers {
			clone.Headers[k] = v
		}
	}

	return &clone
}

// RequestStats 请求统计
type RequestStats struct {
	TotalRequests   int64         `json:"total_requests"`
	SuccessRequests int64         `json:"success_requests"`
	ErrorRequests   int64         `json:"error_requests"`
	AverageTime     time.Duration `json:"average_time"`
	LastRequest     time.Time     `json:"last_request"`
}

// SuccessRate 获取成功率
func (rs *RequestStats) SuccessRate() float64 {
	if rs.TotalRequests == 0 {
		return 0
	}
	return float64(rs.SuccessRequests) / float64(rs.TotalRequests) * 100
}

// ErrorRate 获取错误率
func (rs *RequestStats) ErrorRate() float64 {
	if rs.TotalRequests == 0 {
		return 0
	}
	return float64(rs.ErrorRequests) / float64(rs.TotalRequests) * 100
}

// AddRequest 添加请求统计
func (rs *RequestStats) AddRequest(success bool, duration time.Duration) {
	rs.TotalRequests++
	if success {
		rs.SuccessRequests++
	} else {
		rs.ErrorRequests++
	}
	
	// 计算平均时间（简化版本）
	if rs.TotalRequests == 1 {
		rs.AverageTime = duration
	} else {
		rs.AverageTime = (rs.AverageTime*time.Duration(rs.TotalRequests-1) + duration) / time.Duration(rs.TotalRequests)
	}
	
	rs.LastRequest = time.Now()
}