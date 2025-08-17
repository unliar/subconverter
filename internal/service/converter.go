package service

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"subconverter-go/internal/config"
	"subconverter-go/internal/generator"
	"subconverter-go/internal/parser"
	"subconverter-go/pkg/models"

	"github.com/sirupsen/logrus"
)

// ConverterService 转换服务
type ConverterService struct {
	config           *config.Manager
	logger           *logrus.Logger
	parserManager    *parser.Manager
	generatorManager *generator.Manager
	httpClient       *http.Client
}

// NewConverterService 创建转换服务
func NewConverterService(configManager *config.Manager, logger *logrus.Logger) *ConverterService {
	return &ConverterService{
		config:           configManager,
		logger:           logger,
		parserManager:    parser.NewManager(),
		generatorManager: generator.NewManager(),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ConvertSubscription 转换订阅
func (cs *ConverterService) ConvertSubscription(subscriptionURL, target string, options map[string]string) ([]byte, error) {
	cs.logger.WithFields(logrus.Fields{
		"url":    subscriptionURL,
		"target": target,
	}).Info("开始转换订阅")

	// 1. 获取订阅内容
	content, err := cs.fetchSubscription(subscriptionURL)
	if err != nil {
		return nil, fmt.Errorf("获取订阅失败: %w", err)
	}

	// 2. 解析代理节点
	proxies, err := cs.parserManager.ParseSubscription(string(content))
	if err != nil {
		return nil, fmt.Errorf("解析代理节点失败: %w", err)
	}

	if len(proxies) == 0 {
		return nil, fmt.Errorf("没有找到有效的代理节点")
	}

	cs.logger.WithField("count", len(proxies)).Info("成功解析代理节点")

	// 3. 构建生成配置
	generateConfig, err := cs.buildGenerateConfig(target, options)
	if err != nil {
		return nil, fmt.Errorf("构建生成配置失败: %w", err)
	}

	// 4. 生成配置
	result, err := cs.generatorManager.Generate(proxies, generateConfig)
	if err != nil {
		return nil, fmt.Errorf("生成配置失败: %w", err)
	}

	cs.logger.WithField("size", len(result)).Info("成功生成配置")
	return result, nil
}

// ConvertNodes 转换代理节点
func (cs *ConverterService) ConvertNodes(nodeList []string, target string, options map[string]string) ([]byte, error) {
	cs.logger.WithFields(logrus.Fields{
		"nodeCount": len(nodeList),
		"target":    target,
	}).Info("开始转换代理节点")

	if len(nodeList) == 0 {
		return nil, fmt.Errorf("代理节点列表为空")
	}

	// 解析所有代理节点
	var allProxies []*models.Proxy
	for _, nodeData := range nodeList {
		proxies, err := cs.parserManager.ParseSubscription(nodeData)
		if err != nil {
			cs.logger.WithError(err).Warn("解析代理节点失败，跳过")
			continue
		}
		allProxies = append(allProxies, proxies...)
	}

	if len(allProxies) == 0 {
		return nil, fmt.Errorf("没有找到有效的代理节点")
	}

	// 构建生成配置
	generateConfig, err := cs.buildGenerateConfig(target, options)
	if err != nil {
		return nil, fmt.Errorf("构建生成配置失败: %w", err)
	}

	// 生成配置
	result, err := cs.generatorManager.Generate(allProxies, generateConfig)
	if err != nil {
		return nil, fmt.Errorf("生成配置失败: %w", err)
	}

	cs.logger.WithField("size", len(result)).Info("成功生成配置")
	return result, nil
}

// GetSupportedTargets 获取支持的目标客户端
func (cs *ConverterService) GetSupportedTargets() []string {
	return cs.generatorManager.GetSupportedTargets()
}

// ValidateSubscription 验证订阅链接
func (cs *ConverterService) ValidateSubscription(subscriptionURL string) error {
	if subscriptionURL == "" {
		return fmt.Errorf("订阅链接不能为空")
	}

	if !strings.HasPrefix(subscriptionURL, "http://") && !strings.HasPrefix(subscriptionURL, "https://") {
		return fmt.Errorf("订阅链接格式无效")
	}

	return nil
}

func (cs *ConverterService) fetchSubscription(url string) ([]byte, error) {
	cs.logger.WithField("url", url).Debug("获取订阅内容")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("User-Agent", "SubConverter-Go/1.0")

	resp, err := cs.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP错误: %d", resp.StatusCode)
	}

	content := make([]byte, 0)
	buffer := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			content = append(content, buffer[:n]...)
		}
		if err != nil {
			break
		}
	}

	if len(content) == 0 {
		return nil, fmt.Errorf("订阅内容为空")
	}

	return content, nil
}

func (cs *ConverterService) buildGenerateConfig(target string, options map[string]string) (*generator.GenerateConfig, error) {
	config := &generator.GenerateConfig{
		Target: target,
	}

	for key, value := range options {
		switch key {
		case "format":
			config.Format = value
		case "include":
			config.Include = value
		case "exclude":
			config.Exclude = value
		case "sort":
			config.Sort = parseBool(value)
		case "udp":
			config.UDP = parseBool(value)
		case "enable_rule":
			config.EnableRule = parseBool(value)
		}
	}

	return config, nil
}

func parseBool(value string) bool {
	if value == "" {
		return false
	}
	result, _ := strconv.ParseBool(value)
	return result
}
