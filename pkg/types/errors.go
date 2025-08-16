// Package types 定义错误相关的类型和枚举
package types

import (
	"fmt"
)

// ErrorCode 错误码枚举
type ErrorCode int

const (
	ErrorCodeUnknown ErrorCode = iota
	ErrorCodeInvalidRequest
	ErrorCodeInvalidURL
	ErrorCodeInvalidTarget
	ErrorCodeInvalidConfig
	ErrorCodeParseError
	ErrorCodeGenerateError
	ErrorCodeNetworkError
	ErrorCodeTimeoutError
	ErrorCodeAuthError
	ErrorCodeInternalError
	ErrorCodeValidationError
	ErrorCodeNotFound
	ErrorCodeTooManyRequests
	ErrorCodeServiceUnavailable
)

// 错误码字符串映射
var errorCodeStrings = map[ErrorCode]string{
	ErrorCodeUnknown:            "UNKNOWN_ERROR",
	ErrorCodeInvalidRequest:     "INVALID_REQUEST",
	ErrorCodeInvalidURL:         "INVALID_URL",
	ErrorCodeInvalidTarget:      "INVALID_TARGET",
	ErrorCodeInvalidConfig:      "INVALID_CONFIG",
	ErrorCodeParseError:         "PARSE_ERROR",
	ErrorCodeGenerateError:      "GENERATE_ERROR",
	ErrorCodeNetworkError:       "NETWORK_ERROR",
	ErrorCodeTimeoutError:       "TIMEOUT_ERROR",
	ErrorCodeAuthError:          "AUTH_ERROR",
	ErrorCodeInternalError:      "INTERNAL_ERROR",
	ErrorCodeValidationError:    "VALIDATION_ERROR",
	ErrorCodeNotFound:           "NOT_FOUND",
	ErrorCodeTooManyRequests:    "TOO_MANY_REQUESTS",
	ErrorCodeServiceUnavailable: "SERVICE_UNAVAILABLE",
}

// 错误码到 HTTP 状态码的映射
var errorCodeToHTTPStatus = map[ErrorCode]int{
	ErrorCodeUnknown:            500, // Internal Server Error
	ErrorCodeInvalidRequest:     400, // Bad Request
	ErrorCodeInvalidURL:         400, // Bad Request
	ErrorCodeInvalidTarget:      400, // Bad Request
	ErrorCodeInvalidConfig:      400, // Bad Request
	ErrorCodeParseError:         400, // Bad Request
	ErrorCodeGenerateError:      500, // Internal Server Error
	ErrorCodeNetworkError:       502, // Bad Gateway
	ErrorCodeTimeoutError:       504, // Gateway Timeout
	ErrorCodeAuthError:          401, // Unauthorized
	ErrorCodeInternalError:      500, // Internal Server Error
	ErrorCodeValidationError:    422, // Unprocessable Entity
	ErrorCodeNotFound:           404, // Not Found
	ErrorCodeTooManyRequests:    429, // Too Many Requests
	ErrorCodeServiceUnavailable: 503, // Service Unavailable
}

// String 返回错误码的字符串表示
func (ec ErrorCode) String() string {
	if str, ok := errorCodeStrings[ec]; ok {
		return str
	}
	return "UNKNOWN_ERROR"
}

// HTTPStatus 返回错误码对应的 HTTP 状态码
func (ec ErrorCode) HTTPStatus() int {
	if status, ok := errorCodeToHTTPStatus[ec]; ok {
		return status
	}
	return 500
}

// IsClientError 检查是否为客户端错误（4xx）
func (ec ErrorCode) IsClientError() bool {
	status := ec.HTTPStatus()
	return status >= 400 && status < 500
}

// IsServerError 检查是否为服务端错误（5xx）
func (ec ErrorCode) IsServerError() bool {
	status := ec.HTTPStatus()
	return status >= 500 && status < 600
}

// ConvertError 自定义错误类型
type ConvertError struct {
	Code      ErrorCode `json:"code"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
	Cause     error     `json:"-"`
	RequestID string    `json:"request_id,omitempty"`
}

// Error 实现 error 接口
func (ce *ConvertError) Error() string {
	if ce.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", ce.Code.String(), ce.Message, ce.Details)
	}
	return fmt.Sprintf("[%s] %s", ce.Code.String(), ce.Message)
}

// Unwrap 支持错误链
func (ce *ConvertError) Unwrap() error {
	return ce.Cause
}

// WithRequestID 添加请求 ID
func (ce *ConvertError) WithRequestID(requestID string) *ConvertError {
	ce.RequestID = requestID
	return ce
}

// WithDetails 添加详细信息
func (ce *ConvertError) WithDetails(details string) *ConvertError {
	ce.Details = details
	return ce
}

// WithCause 添加原因错误
func (ce *ConvertError) WithCause(cause error) *ConvertError {
	ce.Cause = cause
	return ce
}

// NewConvertError 创建新的转换错误
func NewConvertError(code ErrorCode, message string, details ...string) *ConvertError {
	err := &ConvertError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

// NewConvertErrorWithCause 创建带原因的转换错误
func NewConvertErrorWithCause(code ErrorCode, message string, cause error) *ConvertError {
	return &ConvertError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// 预定义的常用错误
var (
	// 请求相关错误
	ErrInvalidURL         = NewConvertError(ErrorCodeInvalidURL, "Invalid URL provided")
	ErrInvalidTarget      = NewConvertError(ErrorCodeInvalidTarget, "Invalid target client type")
	ErrMissingURL         = NewConvertError(ErrorCodeInvalidRequest, "Subscription URL is required")
	ErrMissingTarget      = NewConvertError(ErrorCodeInvalidRequest, "Target client type is required")
	ErrInvalidConfig      = NewConvertError(ErrorCodeInvalidConfig, "Invalid configuration")

	// 解析相关错误
	ErrParseSubscription  = NewConvertError(ErrorCodeParseError, "Failed to parse subscription")
	ErrParseProxy         = NewConvertError(ErrorCodeParseError, "Failed to parse proxy configuration")
	ErrParseConfig        = NewConvertError(ErrorCodeParseError, "Failed to parse configuration file")
	ErrUnsupportedFormat  = NewConvertError(ErrorCodeParseError, "Unsupported subscription format")

	// 生成相关错误
	ErrGenerateConfig     = NewConvertError(ErrorCodeGenerateError, "Failed to generate configuration")
	ErrTemplateRender     = NewConvertError(ErrorCodeGenerateError, "Failed to render template")
	ErrUnsupportedClient  = NewConvertError(ErrorCodeGenerateError, "Unsupported client type")

	// 网络相关错误
	ErrNetworkTimeout     = NewConvertError(ErrorCodeTimeoutError, "Network request timeout")
	ErrNetworkFailure     = NewConvertError(ErrorCodeNetworkError, "Network request failed")
	ErrSubscriptionFetch  = NewConvertError(ErrorCodeNetworkError, "Failed to fetch subscription")

	// 认证相关错误
	ErrInvalidToken       = NewConvertError(ErrorCodeAuthError, "Invalid access token")
	ErrMissingToken       = NewConvertError(ErrorCodeAuthError, "Access token required")
	ErrTokenExpired       = NewConvertError(ErrorCodeAuthError, "Access token expired")

	// 验证相关错误
	ErrValidationFailed   = NewConvertError(ErrorCodeValidationError, "Validation failed")
	ErrInvalidProxyData   = NewConvertError(ErrorCodeValidationError, "Invalid proxy data")

	// 资源相关错误
	ErrConfigNotFound     = NewConvertError(ErrorCodeNotFound, "Configuration not found")
	ErrTemplateNotFound   = NewConvertError(ErrorCodeNotFound, "Template not found")
	ErrRulesetNotFound    = NewConvertError(ErrorCodeNotFound, "Ruleset not found")

	// 限流相关错误
	ErrRateLimitExceeded  = NewConvertError(ErrorCodeTooManyRequests, "Rate limit exceeded")

	// 服务相关错误
	ErrServiceUnavailable = NewConvertError(ErrorCodeServiceUnavailable, "Service temporarily unavailable")
	ErrInternalError      = NewConvertError(ErrorCodeInternalError, "Internal server error")
)

// IsConvertError 检查错误是否为 ConvertError 类型
func IsConvertError(err error) bool {
	_, ok := err.(*ConvertError)
	return ok
}

// GetConvertError 从错误中提取 ConvertError
func GetConvertError(err error) *ConvertError {
	if ce, ok := err.(*ConvertError); ok {
		return ce
	}
	return nil
}

// WrapError 包装普通错误为 ConvertError
func WrapError(err error, code ErrorCode, message string) *ConvertError {
	if err == nil {
		return nil
	}
	
	// 如果已经是 ConvertError，直接返回
	if ce, ok := err.(*ConvertError); ok {
		return ce
	}
	
	return &ConvertError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

// ValidationError 验证错误，包含具体的字段错误信息
type ValidationError struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
	Tag     string `json:"tag,omitempty"`
}

// Error 实现 error 接口
func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", ve.Field, ve.Message)
}

// ValidationErrors 验证错误集合
type ValidationErrors []*ValidationError

// Error 实现 error 接口
func (ves ValidationErrors) Error() string {
	if len(ves) == 0 {
		return "validation errors"
	}
	if len(ves) == 1 {
		return ves[0].Error()
	}
	return fmt.Sprintf("validation failed for %d fields", len(ves))
}

// Add 添加验证错误
func (ves *ValidationErrors) Add(field, value, message, tag string) {
	*ves = append(*ves, &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
		Tag:     tag,
	})
}

// HasErrors 检查是否有验证错误
func (ves ValidationErrors) HasErrors() bool {
	return len(ves) > 0
}

// GetFieldErrors 获取指定字段的错误
func (ves ValidationErrors) GetFieldErrors(field string) []*ValidationError {
	var errors []*ValidationError
	for _, err := range ves {
		if err.Field == field {
			errors = append(errors, err)
		}
	}
	return errors
}