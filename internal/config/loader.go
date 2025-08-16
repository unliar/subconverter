// internal/config/loader.go
package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Loader 配置加载器
type Loader struct {
	configDir string
	validator *validator.Validate
	viper     *viper.Viper
}

// NewLoader 创建配置加载器
func NewLoader(configDir string) *Loader {
	return &Loader{
		configDir: configDir,
		validator: SetupValidator(),
		viper:     viper.New(),
	}
}

// LoadConfig 加载完整配置
func (l *Loader) LoadConfig(ctx context.Context) (*Config, error) {
	config := GetDefaultConfig()

	// 加载应用配置
	if err := l.loadAppConfig(config); err != nil {
		return nil, fmt.Errorf("failed to load app config: %w", err)
	}

	// 加载规则配置
	if err := l.loadRulesConfig(config); err != nil {
		return nil, fmt.Errorf("failed to load rules config: %w", err)
	}

	// 加载模板配置
	if err := l.loadTemplateConfig(config); err != nil {
		return nil, fmt.Errorf("failed to load template config: %w", err)
	}

	// 验证配置
	if err := l.ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// loadAppConfig 加载应用配置
func (l *Loader) loadAppConfig(config *Config) error {
	configFile := filepath.Join(l.configDir, "app.yaml")
	
	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// 使用默认配置
		return nil
	}

	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")

	// 设置默认值
	l.setAppDefaults(v)

	// 读取环境变量
	v.AutomaticEnv()
	v.SetEnvPrefix("SC")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := v.Unmarshal(&config.App); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// loadRulesConfig 加载规则配置
func (l *Loader) loadRulesConfig(config *Config) error {
	configFile := filepath.Join(l.configDir, "rules.yaml")
	
	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// 使用默认配置
		return nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read rules config: %w", err)
	}

	if err := yaml.Unmarshal(data, &config.Rules); err != nil {
		return fmt.Errorf("failed to unmarshal rules config: %w", err)
	}

	return nil
}

// loadTemplateConfig 加载模板配置
func (l *Loader) loadTemplateConfig(config *Config) error {
	configFile := filepath.Join(l.configDir, "templates.yaml")
	
	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// 使用默认配置
		return nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read template config: %w", err)
	}

	if err := yaml.Unmarshal(data, &config.Templates); err != nil {
		return fmt.Errorf("failed to unmarshal template config: %w", err)
	}

	return nil
}

// setAppDefaults 设置应用默认值
func (l *Loader) setAppDefaults(v *viper.Viper) {
	// 服务器配置默认值
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 25500)
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.max_connections", 1000)
	v.SetDefault("server.api_mode", false)

	// 日志配置默认值
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "text")
	v.SetDefault("log.output", "stdout")
	v.SetDefault("log.max_size", 100)
	v.SetDefault("log.max_backups", 3)
	v.SetDefault("log.max_age", 28)

	// 缓存配置默认值
	v.SetDefault("cache.enable", true)
	v.SetDefault("cache.default_ttl", "10m")
	v.SetDefault("cache.max_entries", 1000)
	v.SetDefault("cache.cleanup_time", "1m")

	// 安全配置默认值
	v.SetDefault("security.enable_auth", false)
	v.SetDefault("security.rate_limiting", true)
	v.SetDefault("security.max_req_per_min", 60)

	// 监控配置默认值
	v.SetDefault("monitor.enable_metrics", false)
	v.SetDefault("monitor.metrics_path", "/metrics")
	v.SetDefault("monitor.enable_pprof", false)
	v.SetDefault("monitor.pprof_path", "/debug/pprof")
}

// ValidateConfig 验证配置
func (l *Loader) ValidateConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}

	// 使用 validator 验证结构体
	if err := l.validator.Struct(config); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// 额外的业务逻辑验证
	if config.App != nil && config.App.Server != nil {
		if config.App.Server.Port < 1 || config.App.Server.Port > 65535 {
			return fmt.Errorf("invalid port number: %d", config.App.Server.Port)
		}
	}

	return nil
}

// LoadFromFile 从文件加载配置
func (l *Loader) LoadFromFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	config := GetDefaultConfig()

	// 根据文件扩展名选择解析器
	ext := filepath.Ext(filename)
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal json: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config format: %s", ext)
	}

	// 验证配置
	if err := l.ValidateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig 保存配置到文件
func (l *Loader) SaveConfig(config *Config, filename string) error {
	// 验证配置
	if err := l.ValidateConfig(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// 序列化配置
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// 确保目录存在
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}