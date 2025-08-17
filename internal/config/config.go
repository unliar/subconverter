// internal/config/config.go
package config

import (
	"context"
	"os"
	"path/filepath"
	"time"
)

// Config 主配置结构
type Config struct {
	App       *AppConfig       `yaml:"app" json:"app" validate:"required"`
	Rules     *RulesConfig     `yaml:"rules" json:"rules"`
	Templates *TemplatesConfig `yaml:"templates" json:"templates"`
}

// AppConfig 应用配置
type AppConfig struct {
	Server   *ServerConfig   `yaml:"server" json:"server" validate:"required"`
	Log      *LogConfig      `yaml:"log" json:"log" validate:"required"`
	Cache    *CacheConfig    `yaml:"cache" json:"cache"`
	Security *SecurityConfig `yaml:"security" json:"security"`
	Monitor  *MonitorConfig  `yaml:"monitor" json:"monitor"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host           string        `yaml:"host" json:"host" validate:"required,host"`
	Port           int           `yaml:"port" json:"port" validate:"required,port"`
	ReadTimeout    time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout" json:"write_timeout"`
	MaxConnections int           `yaml:"max_connections" json:"max_connections" validate:"min=1"`
	APIMode        bool          `yaml:"api_mode" json:"api_mode"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level" json:"level" validate:"required,oneof=debug info warn error"`
	Format     string `yaml:"format" json:"format" validate:"oneof=text json"`
	Output     string `yaml:"output" json:"output"`
	MaxSize    int    `yaml:"max_size" json:"max_size" validate:"min=1"`
	MaxBackups int    `yaml:"max_backups" json:"max_backups" validate:"min=1"`
	MaxAge     int    `yaml:"max_age" json:"max_age" validate:"min=1"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enable      bool          `yaml:"enable" json:"enable"`
	DefaultTTL  time.Duration `yaml:"default_ttl" json:"default_ttl"`
	MaxEntries  int           `yaml:"max_entries" json:"max_entries" validate:"min=1"`
	CleanupTime time.Duration `yaml:"cleanup_time" json:"cleanup_time"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	EnableAuth     bool     `yaml:"enable_auth" json:"enable_auth"`
	RateLimiting   bool     `yaml:"rate_limiting" json:"rate_limiting"`
	MaxReqPerMin   int      `yaml:"max_req_per_min" json:"max_req_per_min" validate:"min=1"`
	AllowedOrigins []string `yaml:"allowed_origins" json:"allowed_origins"`
}

// MonitorConfig 监控配置
type MonitorConfig struct {
	EnableMetrics bool   `yaml:"enable_metrics" json:"enable_metrics"`
	MetricsPath   string `yaml:"metrics_path" json:"metrics_path"`
	EnablePprof   bool   `yaml:"enable_pprof" json:"enable_pprof"`
	PprofPath     string `yaml:"pprof_path" json:"pprof_path"`
}

// RulesConfig 规则配置
type RulesConfig struct {
	NodeFilters  []*NodeFilter  `yaml:"node_filters" json:"node_filters"`
	RenameRules  []*RenameRule  `yaml:"rename_rules" json:"rename_rules"`
	RegionRules  []*RegionRule  `yaml:"region_rules" json:"region_rules"`
	CustomRules  []*CustomRule  `yaml:"custom_rules" json:"custom_rules"`
	DefaultRules *DefaultRules  `yaml:"default_rules" json:"default_rules"`
}

// NodeFilter 节点过滤器
type NodeFilter struct {
	Name     string   `yaml:"name" json:"name" validate:"required"`
	Type     string   `yaml:"type" json:"type" validate:"required,oneof=include exclude"`
	Patterns []string `yaml:"patterns" json:"patterns" validate:"required,min=1"`
	Regex    bool     `yaml:"regex" json:"regex"`
	Enabled  bool     `yaml:"enabled" json:"enabled"`
}

// RenameRule 重命名规则
type RenameRule struct {
	Name        string `yaml:"name" json:"name" validate:"required"`
	Pattern     string `yaml:"pattern" json:"pattern" validate:"required"`
	Replacement string `yaml:"replacement" json:"replacement" validate:"required"`
	Regex       bool   `yaml:"regex" json:"regex"`
	Enabled     bool   `yaml:"enabled" json:"enabled"`
}

// RegionRule 地区分组规则
type RegionRule struct {
	Name     string   `yaml:"name" json:"name" validate:"required"`
	Regions  []string `yaml:"regions" json:"regions" validate:"required,min=1"`
	Patterns []string `yaml:"patterns" json:"patterns" validate:"required,min=1"`
	Regex    bool     `yaml:"regex" json:"regex"`
	Enabled  bool     `yaml:"enabled" json:"enabled"`
}

// CustomRule 自定义规则
type CustomRule struct {
	Name        string            `yaml:"name" json:"name" validate:"required"`
	Type        string            `yaml:"type" json:"type" validate:"required"`
	Parameters  map[string]string `yaml:"parameters" json:"parameters"`
	Enabled     bool              `yaml:"enabled" json:"enabled"`
	Description string            `yaml:"description" json:"description"`
}

// DefaultRules 默认规则配置
type DefaultRules struct {
	EnableNodeFilter bool `yaml:"enable_node_filter" json:"enable_node_filter"`
	EnableRename     bool `yaml:"enable_rename" json:"enable_rename"`
	EnableRegion     bool `yaml:"enable_region" json:"enable_region"`
	SortNodes        bool `yaml:"sort_nodes" json:"sort_nodes"`
	UDPSupport       bool `yaml:"udp_support" json:"udp_support"`
}

// TemplatesConfig 模板配置
type TemplatesConfig struct {
	ClientTemplates []*ClientTemplate `yaml:"client_templates" json:"client_templates"`
	DefaultTemplate string            `yaml:"default_template" json:"default_template"`
	TemplateDir     string            `yaml:"template_dir" json:"template_dir"`
	CacheTemplates  bool              `yaml:"cache_templates" json:"cache_templates"`
}

// ClientTemplate 客户端模板
type ClientTemplate struct {
	Name        string            `yaml:"name" json:"name" validate:"required"`
	Type        string            `yaml:"type" json:"type" validate:"required"`
	File        string            `yaml:"file" json:"file" validate:"required,template_file"`
	Description string            `yaml:"description" json:"description"`
	Enabled     bool              `yaml:"enabled" json:"enabled"`
	Options     map[string]string `yaml:"options" json:"options"`
}

// Manager 配置管理器
type Manager struct {
	config *Config
	loader *Loader
}

// NewManager 创建配置管理器
func NewManager(configDir string) *Manager {
	return &Manager{
		loader: NewLoader(configDir),
	}
}

// LoadConfig 加载配置
func (m *Manager) LoadConfig() error {
	// 首先尝试加载统一配置文件
	configFile := filepath.Join(m.loader.configDir, "config.yaml")
	if _, err := os.Stat(configFile); err == nil {
		config, err := m.loader.LoadFromFile(configFile)
		if err != nil {
			return err
		}
		m.config = config
		return nil
	}

	// 如果统一配置文件不存在，则使用分离的配置文件
	config, err := m.loader.LoadConfig(context.Background())
	if err != nil {
		return err
	}
	m.config = config
	return nil
}

// GetConfig 获取配置
func (m *Manager) GetConfig() *Config {
	return m.config
}

// GetAppConfig 获取应用配置
func (m *Manager) GetAppConfig() *AppConfig {
	if m.config == nil {
		return nil
	}
	return m.config.App
}

// GetRulesConfig 获取规则配置
func (m *Manager) GetRulesConfig() *RulesConfig {
	if m.config == nil {
		return nil
	}
	return m.config.Rules
}

// GetTemplatesConfig 获取模板配置
func (m *Manager) GetTemplatesConfig() *TemplatesConfig {
	if m.config == nil {
		return nil
	}
	return m.config.Templates
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	return &Config{
		App: &AppConfig{
			Server: &ServerConfig{
				Host:           "0.0.0.0",
				Port:           25500,
				ReadTimeout:    30 * time.Second,
				WriteTimeout:   30 * time.Second,
				MaxConnections: 1000,
				APIMode:        false,
			},
			Log: &LogConfig{
				Level:      "info",
				Format:     "text",
				Output:     "stdout",
				MaxSize:    100,
				MaxBackups: 3,
				MaxAge:     28,
			},
			Cache: &CacheConfig{
				Enable:      true,
				DefaultTTL:  10 * time.Minute,
				MaxEntries:  1000,
				CleanupTime: 1 * time.Minute,
			},
			Security: &SecurityConfig{
				EnableAuth:     false,
				RateLimiting:   true,
				MaxReqPerMin:   60,
				AllowedOrigins: []string{"*"},
			},
			Monitor: &MonitorConfig{
				EnableMetrics: false,
				MetricsPath:   "/metrics",
				EnablePprof:   false,
				PprofPath:     "/debug/pprof",
			},
		},
		Rules: &RulesConfig{
			NodeFilters:  []*NodeFilter{},
			RenameRules:  []*RenameRule{},
			RegionRules:  []*RegionRule{},
			CustomRules:  []*CustomRule{},
			DefaultRules: &DefaultRules{
				EnableNodeFilter: true,
				EnableRename:     true,
				EnableRegion:     true,
				SortNodes:        true,
				UDPSupport:       true,
			},
		},
		Templates: &TemplatesConfig{
			ClientTemplates: []*ClientTemplate{},
			DefaultTemplate: "clash",
			TemplateDir:     "templates",
			CacheTemplates:  true,
		},
	}
}