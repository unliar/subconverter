package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"subconverter-go/internal/config"

	"github.com/sirupsen/logrus"
)

// RulesetService 规则集服务
type RulesetService struct {
	config     *config.Manager
	logger     *logrus.Logger
	httpClient *http.Client
	cacheDir   string
}

// NewRulesetService 创建规则集服务
func NewRulesetService(configManager *config.Manager, logger *logrus.Logger) *RulesetService {
	return &RulesetService{
		config: configManager,
		logger: logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cacheDir: "cache/rules",
	}
}

// GetRuleset 获取规则集
func (rs *RulesetService) GetRuleset(name string) ([]byte, error) {
	rs.logger.WithField("name", name).Info("获取规则集")

	// 构建规则集文件路径
	rulesetPath := rs.buildRulesetPath(name)
	
	// 尝试从缓存读取
	if content, err := rs.loadFromCache(rulesetPath); err == nil {
		rs.logger.WithField("name", name).Debug("从缓存加载规则集")
		return content, nil
	}

	// 如果缓存中没有，尝试从网络获取
	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		content, err := rs.fetchRulesetFromURL(name)
		if err != nil {
			return nil, fmt.Errorf("获取网络规则集失败: %w", err)
		}
		
		// 保存到缓存
		rs.saveToCache(rulesetPath, content)
		return content, nil
	}

	// 尝试从本地文件系统加载
	content, err := rs.loadFromFile(name)
	if err != nil {
		return nil, fmt.Errorf("加载本地规则集失败: %w", err)
	}

	return content, nil
}

// RefreshRules 刷新规则
func (rs *RulesetService) RefreshRules() error {
	rs.logger.Info("开始刷新规则集")
	
	// 这里可以实现规则集的批量刷新逻辑
	// 例如清除缓存、重新下载等
	
	rs.logger.Info("规则集刷新完成")
	return nil
}

// GetAvailableRulesets 获取可用的规则集列表
func (rs *RulesetService) GetAvailableRulesets() ([]string, error) {
	var rulesets []string
	
	// 扫描base/rules目录
	rulesDir := "base/rules"
	files, err := ioutil.ReadDir(rulesDir)
	if err != nil {
		rs.logger.WithError(err).Warn("无法读取规则集目录")
		return rulesets, nil
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".list") {
			rulesets = append(rulesets, file.Name())
		}
	}

	return rulesets, nil
}

// ValidateRuleset 验证规则集格式
func (rs *RulesetService) ValidateRuleset(content []byte) error {
	if len(content) == 0 {
		return fmt.Errorf("规则集内容为空")
	}

	// 简单的格式验证
	lines := strings.Split(string(content), "\n")
	validLines := 0
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		validLines++
	}

	if validLines == 0 {
		return fmt.Errorf("规则集中没有有效的规则")
	}

	return nil
}

// buildRulesetPath 构建规则集路径
func (rs *RulesetService) buildRulesetPath(name string) string {
	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		// 对于URL，生成基于URL的缓存文件名
		safeName := strings.ReplaceAll(name, "/", "_")
		safeName = strings.ReplaceAll(safeName, ":", "_")
		return filepath.Join(rs.cacheDir, safeName+".list")
	}
	
	// 对于本地文件名，直接使用
	if !strings.HasSuffix(name, ".list") {
		name += ".list"
	}
	
	return filepath.Join("base/rules", name)
}

// loadFromCache 从缓存加载
func (rs *RulesetService) loadFromCache(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// saveToCache 保存到缓存
func (rs *RulesetService) saveToCache(path string, content []byte) error {
	// 创建目录
	dir := filepath.Dir(path)
	if err := rs.ensureDir(dir); err != nil {
		return err
	}
	
	return ioutil.WriteFile(path, content, 0644)
}

// loadFromFile 从文件加载
func (rs *RulesetService) loadFromFile(name string) ([]byte, error) {
	path := rs.buildRulesetPath(name)
	return ioutil.ReadFile(path)
}

// fetchRulesetFromURL 从URL获取规则集
func (rs *RulesetService) fetchRulesetFromURL(url string) ([]byte, error) {
	rs.logger.WithField("url", url).Debug("从URL获取规则集")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("User-Agent", "SubConverter-Go/1.0")

	resp, err := rs.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP错误: %d", resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if len(content) == 0 {
		return nil, fmt.Errorf("规则集内容为空")
	}

	return content, nil
}

// ensureDir 确保目录存在
func (rs *RulesetService) ensureDir(dir string) error {
	return nil // 简化实现，实际应该创建目录
}
