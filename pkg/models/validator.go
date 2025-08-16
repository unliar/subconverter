// Package models 定义验证器相关的功能\npackage models\n\nimport (\n\t\"fmt\"\n\t\"reflect\"\n\t\"regexp\"\n\t\"strings\"\n\n\t\"github.com/go-playground/validator/v10\"\n\t\"subconverter-go/pkg/constants\"\n\t\"subconverter-go/pkg/types\"\n)\n\nvar validate *validator.Validate\n\n// 预编译的正则表达式\nvar (\n\turlRegex    *regexp.Regexp\n\tipv4Regex   *regexp.Regexp\n\tipv6Regex   *regexp.Regexp\n\tdomainRegex *regexp.Regexp\n\tportRegex   *regexp.Regexp\n\tuuidRegex   *regexp.Regexp\n\tbase64Regex *regexp.Regexp\n)\n\nfunc init() {\n\tvalidate = validator.New()\n\n\t// 编译正则表达式\n\turlRegex = regexp.MustCompile(constants.URLPattern)\n\tipv4Regex = regexp.MustCompile(constants.IPv4Pattern)\n\tipv6Regex = regexp.MustCompile(constants.IPv6Pattern)\n\tdomainRegex = regexp.MustCompile(constants.DomainPattern)\n\tportRegex = regexp.MustCompile(constants.PortPattern)\n\tuuidRegex = regexp.MustCompile(constants.UUIDPattern)\n\tbase64Regex = regexp.MustCompile(constants.Base64Pattern)\n\n\t// 注册自定义验证器\n\tvalidate.RegisterValidation(\"proxytype\", validateProxyType)\n\tvalidate.RegisterValidation(\"target\", validateTarget)\n\tvalidate.RegisterValidation(\"hostname\", validateHostname)\n\tvalidate.RegisterValidation(\"uuid\", validateUUID)\n\tvalidate.RegisterValidation(\"base64\", validateBase64)\n\tvalidate.RegisterValidation(\"cipher\", validateCipher)\n\tvalidate.RegisterValidation(\"protocol\", validateProtocol)\n\tvalidate.RegisterValidation(\"obfs\", validateObfs)\n\n\t// 注册字段名映射\n\tvalidate.RegisterTagNameFunc(func(fld reflect.StructField) string {\n\t\tname := strings.SplitN(fld.Tag.Get(\"json\"), \",\", 2)[0]\n\t\tif name == \"-\" {\n\t\t\treturn \"\"\n\t\t}\n\t\treturn name\n\t})\n}\n\n// ValidateStruct 验证结构体\nfunc ValidateStruct(s interface{}) error {\n\treturn validate.Struct(s)\n}\n\n// ValidateVar 验证单个变量\nfunc ValidateVar(field interface{}, tag string) error {\n\treturn validate.Var(field, tag)\n}\n\n// GetValidationErrors 获取详细的验证错误信息\nfunc GetValidationErrors(err error) types.ValidationErrors {\n\tvar validationErrors types.ValidationErrors\n\n\tif err == nil {\n\t\treturn validationErrors\n\t}\n\n\tif validatorErrors, ok := err.(validator.ValidationErrors); ok {\n\t\tfor _, fieldError := range validatorErrors {\n\t\t\tvalidationErrors.Add(\n\t\t\t\tfieldError.Field(),\n\t\t\t\tfmt.Sprintf(\"%v\", fieldError.Value()),\n\t\t\t\tgetErrorMessage(fieldError),\n\t\t\t\tfieldError.Tag(),\n\t\t\t)\n\t\t}\n\t}\n\n\treturn validationErrors\n}\n\n// getErrorMessage 获取友好的错误消息\nfunc getErrorMessage(fe validator.FieldError) string {\n\tswitch fe.Tag() {\n\tcase \"required\":\n\t\treturn \"该字段为必填项\"\n\tcase \"url\":\n\t\treturn \"无效的URL格式\"\n\tcase \"hostname\":\n\t\treturn \"无效的主机名或IP地址\"\n\tcase \"min\":\n\t\treturn \"值太小\"\n\tcase \"max\":\n\t\treturn \"值太大\"\n\tcase \"uuid\":\n\t\treturn \"无效的UUID格式\"\n\tcase \"base64\":\n\t\treturn \"无效的Base64格式\"\n\tcase \"proxytype\":\n\t\treturn \"不支持的代理类型\"\n\tcase \"target\":\n\t\treturn \"不支持的目标客户端类型\"\n\tcase \"cipher\":\n\t\treturn \"不支持的加密方法\"\n\tcase \"protocol\":\n\t\treturn \"不支持的协议\"\n\tcase \"obfs\":\n\t\treturn \"不支持的混淆方法\"\n\tdefault:\n\t\treturn \"验证失败\"\n\t}\n}\n\n// validateProxyType 验证代理类型\nfunc validateProxyType(fl validator.FieldLevel) bool {\n\tproxyType, ok := fl.Field().Interface().(types.ProxyType)\n\tif !ok {\n\t\treturn false\n\t}\n\treturn proxyType.IsValid()\n}\n\n// validateTarget 验证目标类型\nfunc validateTarget(fl validator.FieldLevel) bool {\n\ttarget := fl.Field().String()\n\tfor _, validTarget := range constants.SupportedTargets {\n\t\tif target == validTarget {\n\t\t\treturn true\n\t\t}\n\t}\n\treturn false\n}\n\n// validateHostname 验证主机名（域名或IP地址）\nfunc validateHostname(fl validator.FieldLevel) bool {\n\thostname := fl.Field().String()\n\tif hostname == \"\" {\n\t\treturn false\n\t}\n\n\t// 检查是否为IPv4地址\n\tif ipv4Regex.MatchString(hostname) {\n\t\treturn true\n\t}\n\n\t// 检查是否为IPv6地址\n\tif ipv6Regex.MatchString(hostname) {\n\t\treturn true\n\t}\n\n\t// 检查是否为域名\n\treturn domainRegex.MatchString(hostname)\n}\n\n// validateUUID 验证UUID格式\nfunc validateUUID(fl validator.FieldLevel) bool {\n\tuuid := fl.Field().String()\n\treturn uuidRegex.MatchString(uuid)\n}\n\n// validateBase64 验证Base64格式\nfunc validateBase64(fl validator.FieldLevel) bool {\n\tbase64Str := fl.Field().String()\n\treturn base64Regex.MatchString(base64Str)\n}\n\n// validateCipher 验证加密方法\nfunc validateCipher(fl validator.FieldLevel) bool {\n\tcipher := fl.Field().String()\n\tif cipher == \"\" {\n\t\treturn true // 允许空值\n\t}\n\n\t// 检查所有支持的加密方法\n\tfor _, ciphers := range constants.SupportedCiphers {\n\t\tfor _, supportedCipher := range ciphers {\n\t\t\tif cipher == supportedCipher {\n\t\t\t\treturn true\n\t\t\t}\n\t\t}\n\t}\n\treturn false\n}\n\n// validateProtocol 验证协议\nfunc validateProtocol(fl validator.FieldLevel) bool {\n\tprotocol := fl.Field().String()\n\tif protocol == \"\" {\n\t\treturn true // 允许空值\n\t}\n\n\t// 检查所有支持的协议\n\tfor _, protocols := range constants.SupportedProtocols {\n\t\tfor _, supportedProtocol := range protocols {\n\t\t\tif protocol == supportedProtocol {\n\t\t\t\treturn true\n\t\t\t}\n\t\t}\n\t}\n\treturn false\n}\n\n// validateObfs 验证混淆方法\nfunc validateObfs(fl validator.FieldLevel) bool {\n\tobfs := fl.Field().String()\n\tif obfs == \"\" {\n\t\treturn true // 允许空值\n\t}\n\n\t// 检查所有支持的混淆方法\n\tfor _, obfsList := range constants.SupportedObfs {\n\t\tfor _, supportedObfs := range obfsList {\n\t\t\tif obfs == supportedObfs {\n\t\t\t\treturn true\n\t\t\t}\n\t\t}\n\t}\n\treturn false\n}\n\n// ValidateProxy 验证代理配置\nfunc ValidateProxy(proxy *Proxy) error {\n\tif proxy == nil {\n\t\treturn types.NewConvertError(types.ErrorCodeValidationError, \"proxy is nil\")\n\t}\n\n\t// 基础验证\n\tif err :
// Package models 定义验证器相关的功能
package models

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"subconverter-go/pkg/constants"
	"subconverter-go/pkg/types"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// 预编译的正则表达式
var (
	urlRegex    *regexp.Regexp
	ipv4Regex   *regexp.Regexp
	ipv6Regex   *regexp.Regexp
	domainRegex *regexp.Regexp
	portRegex   *regexp.Regexp
	uuidRegex   *regexp.Regexp
	base64Regex *regexp.Regexp
)

func init() {
	validate = validator.New()

	// 编译正则表达式
	urlRegex = regexp.MustCompile(constants.URLPattern)
	ipv4Regex = regexp.MustCompile(constants.IPv4Pattern)
	ipv6Regex = regexp.MustCompile(constants.IPv6Pattern)
	domainRegex = regexp.MustCompile(constants.DomainPattern)
	portRegex = regexp.MustCompile(constants.PortPattern)
	uuidRegex = regexp.MustCompile(constants.UUIDPattern)
	base64Regex = regexp.MustCompile(constants.Base64Pattern)

	// 注册自定义验证器
	validate.RegisterValidation("proxytype", validateProxyType)
	validate.RegisterValidation("target", validateTarget)
	validate.RegisterValidation("hostname", validateHostname)
	validate.RegisterValidation("uuid", validateUUID)
	validate.RegisterValidation("base64", validateBase64)
	validate.RegisterValidation("cipher", validateCipher)
	validate.RegisterValidation("protocol", validateProtocol)
	validate.RegisterValidation("obfs", validateObfs)

	// 注册字段名映射
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// ValidateVar 验证单个变量
func ValidateVar(field interface{}, tag string) error {
	return validate.Var(field, tag)
}

// GetValidationErrors 获取详细的验证错误信息
func GetValidationErrors(err error) types.ValidationErrors {
	var validationErrors types.ValidationErrors

	if err == nil {
		return validationErrors
	}

	if validatorErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validatorErrors {
			validationErrors.Add(
				fieldError.Field(),
				fmt.Sprintf("%v", fieldError.Value()),
				getErrorMessage(fieldError),
				fieldError.Tag(),
			)
		}
	}

	return validationErrors
}

// getErrorMessage 获取友好的错误消息
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "该字段为必填项"
	case "url":
		return "无效的URL格式"
	case "hostname":
		return "无效的主机名或IP地址"
	case "min":
		return "值太小"
	case "max":
		return "值太大"
	case "uuid":
		return "无效的UUID格式"
	case "base64":
		return "无效的Base64格式"
	case "proxytype":
		return "不支持的代理类型"
	case "target":
		return "不支持的目标客户端类型"
	case "cipher":
		return "不支持的加密方法"
	case "protocol":
		return "不支持的协议"
	case "obfs":
		return "不支持的混淆方法"
	default:
		return "验证失败"
	}
}

// validateProxyType 验证代理类型
func validateProxyType(fl validator.FieldLevel) bool {
	proxyType, ok := fl.Field().Interface().(types.ProxyType)
	if !ok {
		return false
	}
	return proxyType.IsValid()
}

// validateTarget 验证目标类型
func validateTarget(fl validator.FieldLevel) bool {
	target := fl.Field().String()
	for _, validTarget := range constants.SupportedTargets {
		if target == validTarget {
			return true
		}
	}
	return false
}

// validateHostname 验证主机名（域名或IP地址）
func validateHostname(fl validator.FieldLevel) bool {
	hostname := fl.Field().String()
	if hostname == "" {
		return false
	}

	// 检查是否为IPv4地址
	if ipv4Regex.MatchString(hostname) {
		return true
	}

	// 检查是否为IPv6地址
	if ipv6Regex.MatchString(hostname) {
		return true
	}

	// 检查是否为域名
	return domainRegex.MatchString(hostname)
}

// validateUUID 验证UUID格式
func validateUUID(fl validator.FieldLevel) bool {
	uuid := fl.Field().String()
	return uuidRegex.MatchString(uuid)
}

// validateBase64 验证Base64格式
func validateBase64(fl validator.FieldLevel) bool {
	base64Str := fl.Field().String()
	return base64Regex.MatchString(base64Str)
}

// validateCipher 验证加密方法
func validateCipher(fl validator.FieldLevel) bool {
	cipher := fl.Field().String()
	if cipher == "" {
		return true // 允许空值
	}

	// 检查所有支持的加密方法
	for _, ciphers := range constants.SupportedCiphers {
		for _, supportedCipher := range ciphers {
			if cipher == supportedCipher {
				return true
			}
		}
	}
	return false
}

// validateProtocol 验证协议
func validateProtocol(fl validator.FieldLevel) bool {
	protocol := fl.Field().String()
	if protocol == "" {
		return true // 允许空值
	}

	// 检查所有支持的协议
	for _, protocols := range constants.SupportedProtocols {
		for _, supportedProtocol := range protocols {
			if protocol == supportedProtocol {
				return true
			}
		}
	}
	return false
}

// validateObfs 验证混淆方法
func validateObfs(fl validator.FieldLevel) bool {
	obfs := fl.Field().String()
	if obfs == "" {
		return true // 允许空值
	}

	// 检查所有支持的混淆方法
	for _, obfsList := range constants.SupportedObfs {
		for _, supportedObfs := range obfsList {
			if obfs == supportedObfs {
				return true
			}
		}
	}
	return false
}

// ValidateProxy 验证代理配置
func ValidateProxy(proxy *Proxy) error {
	if proxy == nil {
		return types.NewConvertError(types.ErrorCodeValidationError, "proxy is nil")
	}

	// 基础验证
	if err := ValidateStruct(proxy); err != nil {
		return types.WrapError(err, types.ErrorCodeValidationError, "proxy validation failed")
	}

	// 自定义业务逻辑验证
	if !proxy.IsValid() {
		return types.NewConvertError(types.ErrorCodeValidationError, "proxy configuration is invalid")
	}

	return nil
}

// ValidateProxyGroup 验证代理组配置
func ValidateProxyGroup(group *ProxyGroupConfig) error {
	if group == nil {
		return types.NewConvertError(types.ErrorCodeValidationError, "proxy group is nil")
	}

	if err := ValidateStruct(group); err != nil {
		return types.WrapError(err, types.ErrorCodeValidationError, "proxy group validation failed")
	}

	if !group.IsValid() {
		return types.NewConvertError(types.ErrorCodeValidationError, "proxy group configuration is invalid")
	}

	return nil
}

// ValidateRuleset 验证规则集配置
func ValidateRuleset(ruleset *RulesetConfig) error {
	if ruleset == nil {
		return types.NewConvertError(types.ErrorCodeValidationError, "ruleset is nil")
	}

	if err := ValidateStruct(ruleset); err != nil {
		return types.WrapError(err, types.ErrorCodeValidationError, "ruleset validation failed")
	}

	if !ruleset.IsValid() {
		return types.NewConvertError(types.ErrorCodeValidationError, "ruleset configuration is invalid")
	}

	return nil
}

// ValidateConvertRequest 验证转换请求
func ValidateConvertRequest(request *ConvertRequest) error {
	if request == nil {
		return types.NewConvertError(types.ErrorCodeValidationError, "convert request is nil")
	}

	if err := ValidateStruct(request); err != nil {
		return types.WrapError(err, types.ErrorCodeValidationError, "convert request validation failed")
	}

	if !request.IsValid() {
		return types.NewConvertError(types.ErrorCodeValidationError, "convert request is invalid")
	}

	return nil
}