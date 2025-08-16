\n// internal/config/config.go\npackage config\n\nimport (\n\t\"context\"\n\t\"time\"\n)\n\n// Config 主配置结构\ntype Config struct {\n\tApp       *AppConfig       `yaml:\"app\" json:\"app\" validate:\"required\"`\n\tRules     *RulesConfig     `yaml:\"rules\" json:\"rules\"`\n\tTemplates *TemplatesConfig `yaml:\"templates\" json:\"templates\"`\n}\n\n// AppConfig 应用配置\ntype AppConfig struct {\n\tServer   *ServerConfig   `yaml:\"server\" json:\"server\" validate:\"required\"`\n\tLog      *LogConfig      `yaml:\"log\" json:\"log\" validate:\"required\"`\n\tCache    *CacheConfig    `yaml:\"cache\" json:\"cache\"`\n\tSecurity *SecurityConfig `yaml:\"security\" json:\"security\"`\n\tMonitor  *MonitorConfig  `yaml:\"monitor\" json:\"monitor\"`\n}\n\n// ServerConfig 服务器配置\ntype ServerConfig struct {\n\tHost           string        `yaml:\"host\" json:\"host\" validate:\"required,host\"`\n\tPort           int           `yaml:\"port\" json:\"port\" validate:\"required,port\"`\n\tReadTimeout    time.Duration `yaml:\"read_timeout\" json:\"read_timeout\"`\n\tWriteTimeout   time.Duration `yaml:\"write_timeout\" json:\"write_timeout\"`\n\tMaxConnections int           `yaml:\"max_connections\" json:\"max_connections\" validate:\"min=1\"`\n\tAPIMode        bool          `yaml:\"api_mode\" json:\"api_mode\"`\n}\n\n// LogConfig 日志配置\ntype LogConfig struct {\n\tLevel      string `yaml:\"level\" json:\"level\" validate:\"required,oneof=debug info warn error\"`\n\tFormat     string `yaml:\"format\" json:\"format\" validate:\"oneof=text json\"`\n\tOutput     string `yaml:\"output\" json:\"output\"`\n\tMaxSize    int    `yaml:\"max_size\" json:\"max_size\" validate:\"min=1\"`\n\tMaxBackups int    `yaml:\"max_backups\" json:\"max_backups\" validate:\"min=1\"`\n\tMaxAge     int    `yaml:\"max_age\" json:\"max_age\" validate:\"min=1\"`\n}\n\n// CacheConfig 缓存配置\ntype CacheConfig struct {\n\tEnable      bool          `yaml:\"enable\" json:\"enable\"`\n\tDefaultTTL  time.Duration `yaml:\"default_ttl\" json:\"default_ttl\"`\n\tMaxEntries  int           `yaml:\"max_entries\" json:\"max_entries\" validate:\"min=1\"`\n\tCleanupTime time.Duration `yaml:\"cleanup_time\" json:\"cleanup_time\"`\n}\n\n// SecurityConfig 安全配置\ntype SecurityConfig struct {\n\tEnableAuth     bool     `yaml:\"enable_auth\" json:\"enable_auth\"`\n\tRateLimiting   bool     `yaml:\"rate_limiting\" json:\"rate_limiting\"`\n\tMaxReqPerMin   int      `yaml:\"max_req_per_min\" json:\"max_req_per_min\" validate:\"min=1\"`\n\tAllowedOrigins []string `yaml:\"allowed_origins\" json:\"allowed_origins\"`\n}\n\n// MonitorConfig 监控配置\ntype MonitorConfig struct {\n\tEnableMetrics bool   `yaml:\"enable_metrics\" json:\"enable_metrics\"`\n\tMetricsPath   string `yaml:\"metrics_path\" json:\"metrics_path\"`\n\tEnablePprof   bool   `yaml:\"enable_pprof\" json:\"enable_pprof\"`\n\tPprofPath     string `yaml:\"pprof_path\" json:\"pprof_path\"`\n}\n\n// RulesConfig 规则配置\ntype RulesConfig struct {\n\tNodeFilters  []*NodeFilter  `yaml:\"node_filters\" json:\"node_filters\"`\n\tRenameRules  []*RenameRule  `yaml:\"rename_rules\" json:\"rename_rules\"`\n\tRegionRules  []*RegionRule  `yaml:\"region_rules\" json:\"region_rules\"`\n\tCustomRules  []*CustomRule  `yaml:\"custom_rules\" json:\"custom_rules\"`\n\tDefaultRules *DefaultRules  `yaml:\"default_rules\" json:\"default_rules\"`\n}\n\n// NodeFilter 节点过滤器\ntype NodeFilter struct {\n\tName     string   `yaml:\"name\" json:\"name\" validate:\"required\"`\n\tType     string   `yaml:\"type\" json:\"type\" validate:\"required,oneof=include exclude\"`\n\tPatterns []string `yaml:\"patterns\" json:\"patterns\" validate:\"required,min=1\"`\n\tRegex    bool     `yaml:\"regex\" json:\"regex\"`\n\tEnabled  bool     `yaml:\"enabled\" json:\"enabled\"`\n}\n\n// RenameRule 重命名规则\ntype RenameRule struct {\n\tName        string `yaml:\"name\" json:\"name\" validate:\"required\"`\n\tPattern     string `yaml:\"pattern\" json:\"pattern\" validate:\"required\"`\n\tReplacement string `yaml:\"replacement\" json:\"replacement\" validate:\"required\"`\n\tRegex       bool   `yaml:\"regex\" json:\"regex\"`\n\tEnabled     bool   `yaml:\"enabled\" json:\"enabled\"`\n}\n\n// RegionRule 地区分组规则\ntype RegionRule struct {\n\tName     string   `yaml:\"name\" json:\"name\" validate:\"required\"`\n\tRegions  []string `yaml:\"regions\" json:\"regions\" validate:\"required,min=1\"`\n\tPatterns []string `yaml:\"patterns\" json:\"patterns\" validate:\"required,min=1\"`\n\tRegex    bool     `yaml:\"regex\" json:\"regex\"`\n\tEnabled  bool     `yaml:\"enabled\" json:\"enabled\"`\n}\n\n// CustomRule 自定义规则\ntype CustomRule struct {\n\tName        string            `yaml:\"name\" json:\"name\" validate:\"required\"`\n\tType        string            `yaml:\"type\" json:\"type\" validate:\"required\"`\n\tParameters  map[string]string `yaml:\"parameters\" json:\"parameters\"`\n\tEnabled     bool              `yaml:\"enabled\" json:\"enabled\"`\n\tDescription string            `yaml:\"description\" json:\"description\"`\n}\n\n// DefaultRules 默认规则配置\ntype DefaultRules struct {\n\tEnableNodeFilter bool `yaml:\"enable_node_filter\" json:\"enable_node_filter\"`\n\tEnableRename     bool `yaml:\"enable_rename\" json:\"enable_rename\"`\n\tEnableRegion     bool `yaml:\"enable_region\" json:\"enable_region\"`\n\tSortNodes        bool `yaml:\"sort_nodes\" json:\"sort_nodes\"`\n\tUDPSupport       bool `yaml:\"udp_support\" json:\"udp_support\"`\n}\n\n// TemplatesConfig 模板配置\ntype TemplatesConfig struct {\n\tClientTemplates []*ClientTemplate `yaml:\"client_templates\" json:\"client_templates\"`\n\tDefaultTemplate string            `yaml:\"default_template\" json:\"default_template\"`\n\tTemplateDir     string            `yaml:\"template_dir\" json:\"template_dir\"`\n\tCacheTemplates  bool              `yaml:\"cache_templates\" json:\"cache_templates\"`\n}\n\n// ClientTemplate 客户端模板\ntype ClientTemplate struct {\n\tName        string            `yaml:\"name\" json:\"name\" validate:\"required\"`\n\tType        string            `yaml:\"type\" json:\"type\" validate:\"required\"`\n\tFile        string            `yaml:\"file\" json:\"file\" validate:\"required,template_file\"`\n\tDescription string            `yaml:\"description\" json:\"description\"`\n\tEnabled     bool              `yaml:\"enabled\" json:\"enabled\"`\n\tOptions     map[string]string `yaml:\"options\" json:\"options\"`\n}\n\n// Manager 配置管理器\ntype Manager struct {\n\tconfig *Config\n\tloader *Loader\n}\n\n// NewManager 创建配置管理器\nfunc NewManager(configDir string) *Manager {\n\treturn &Manager{\n\t\tloader: NewLoader(configDir),\n\t}\n}\n\n// LoadConfig 加载配置\nfunc (m *Manager) LoadConfig() error {\n\tconfig, err := m.loader.LoadConfig(context.Background())\n\tif err != nil {\n\t\treturn err\n\t}\n\tm.config = config\n\treturn nil\n}\n\n// GetConfig 获取配置\nfunc (m *Manager) GetConfig() *Config {\n\treturn m.config\n}\n\n// GetAppConfig 获取应用配置\nfunc (m *Manager) GetAppConfig() *AppConfig {\n\tif m.config == nil {\n\t\treturn nil\n\t}\n\treturn m.config.App\n}\n\n// GetRulesConfig 获取规则配置\nfunc (m *Manager) GetRulesConfig() *RulesConfig {\n\tif m.config == nil {\n\t\treturn nil\n\t}\n\treturn m.config.Rules\n}\n\n// GetTemplatesConfig 获取模板配置\nfunc (m *Manager) GetTemplatesConfig() *TemplatesConfig {\n\tif m.config == nil {\n\t\treturn nil\n\t}\n\treturn m.config.Templates\n}\n\n// GetDefaultConfig 获取默认配置\nfunc GetDefaultConfig() *Config {\n\treturn &Config{\n\t\tApp: &AppConfig{\n\t\t\tServer: &ServerConfig{\n\t\t\t\tHost:           \"0.0.0.0\",\n\t\t\t\tPort:           25500,\n\t\t\t\tReadTimeout:    30 * time.Second,\n\t\t\t\tWriteTimeout:   30 * time.Second,\n\t\t\t\tMaxConnections: 1000,\n\t\t\t\tAPIMode:        f
alse,
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