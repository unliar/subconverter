# SubConverter Go 版本配置管理和模板系统设计

## 1. 系统概述

### 1.1 设计目标

- **兼容性**: 100%兼容现有 C++版本的配置和模板
- **灵活性**: 支持多种配置源和动态更新
- **性能**: 高效的模板渲染和配置解析
- **可维护性**: 清晰的配置结构和模板组织

### 1.2 核心组件

- **配置管理器**: 统一配置加载和管理
- **模板引擎**: Jinja2 兼容的模板渲染
- **规则引擎**: 节点过滤和转换规则
- **验证系统**: 配置和数据验证

## 2. 配置管理系统架构

### 2.1 配置层次结构

```
配置管理系统
├── 应用配置 (Application Config)
│   ├── 服务器配置
│   ├── 日志配置
│   └── 性能配置
├── 规则配置 (Rules Config)
│   ├── 节点过滤规则
│   ├── 重命名规则
│   └── 地区分组规则
├── 模板配置 (Template Config)
│   ├── 客户端模板
│   ├── 自定义模板
│   └── 模板变量
└── 外部配置 (External Config)
    ├── 远程配置
    ├── 环境变量
    └── 命令行参数
```

### 2.2 配置管理接口设计

```go
// pkg/config/manager.go
package config

import (
    "context"
    "time"
)

// ConfigManager 配置管理器接口
type ConfigManager interface {
    // 加载配置
    Load(ctx context.Context) error

    // 获取应用配置
    GetApp() *AppConfig

    // 获取规则配置
    GetRules() *RulesConfig

    // 获取模板配置
    GetTemplates() *TemplateConfig

    // 动态更新配置
    Reload(ctx context.Context) error

    // 监听配置变化
    Watch(ctx context.Context) <-chan ConfigEvent

    // 验证配置
    Validate() error
}

// ConfigEvent 配置变更事件
type ConfigEvent struct {
    Type      EventType `json:"type"`
    Source    string    `json:"source"`
    Timestamp time.Time `json:"timestamp"`
    Changes   []Change  `json:"changes"`
}

type EventType string

const (
    EventTypeReload EventType = "reload"
    EventTypeUpdate EventType = "update"
    EventTypeError  EventType = "error"
)
```

### 2.3 应用配置结构

```go
// pkg/config/app.go
package config

// AppConfig 应用主配置
type AppConfig struct {
    Server    *ServerConfig    `yaml:"server" validate:"required"`
    Log       *LogConfig       `yaml:"log" validate:"required"`
    Cache     *CacheConfig     `yaml:"cache"`
    Security  *SecurityConfig  `yaml:"security"`
    Monitor   *MonitorConfig   `yaml:"monitor"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
    Host         string        `yaml:"host" validate:"required" default:"0.0.0.0"`
    Port         int           `yaml:"port" validate:"min=1,max=65535" default:"25500"`
    ReadTimeout  time.Duration `yaml:"read_timeout" default:"30s"`
    WriteTimeout time.Duration `yaml:"write_timeout" default:"30s"`
    MaxConns     int           `yaml:"max_connections" default:"1000"`
    TLS          *TLSConfig    `yaml:"tls"`
}

// LogConfig 日志配置
type LogConfig struct {
    Level      string `yaml:"level" validate:"oneof=debug info warn error" default:"info"`
    Format     string `yaml:"format" validate:"oneof=text json" default:"text"`
    Output     string `yaml:"output" default:"stdout"`
    MaxSize    int    `yaml:"max_size" default:"100"`    // MB
    MaxBackups int    `yaml:"max_backups" default:"3"`
    MaxAge     int    `yaml:"max_age" default:"28"`      // days
}

// CacheConfig 缓存配置
type CacheConfig struct {
    Enable       bool          `yaml:"enable" default:"true"`
    DefaultTTL   time.Duration `yaml:"default_ttl" default:"10m"`
    MaxEntries   int           `yaml:"max_entries" default:"1000"`
    CleanupTime  time.Duration `yaml:"cleanup_time" default:"1m"`
}
```

## 3. 规则配置系统

### 3.1 规则配置结构

```go
// pkg/config/rules.go
package config

// RulesConfig 规则配置
type RulesConfig struct {
    NodeFilters  []NodeFilter  `yaml:"node_filters"`
    Rename       []RenameRule  `yaml:"rename"`
    RegionGroup  []RegionRule  `yaml:"region_group"`
    Custom       []CustomRule  `yaml:"custom"`
}

// NodeFilter 节点过滤规则
type NodeFilter struct {
    Name     string   `yaml:"name" validate:"required"`
    Type     string   `yaml:"type" validate:"oneof=include exclude"`
    Patterns []string `yaml:"patterns" validate:"required"`
    Regex    bool     `yaml:"regex" default:"false"`
    Enabled  bool     `yaml:"enabled" default:"true"`
}

// RenameRule 重命名规则
type RenameRule struct {
    Name        string `yaml:"name" validate:"required"`
    Pattern     string `yaml:"pattern" validate:"required"`
    Replacement string `yaml:"replacement" validate:"required"`
    Regex       bool   `yaml:"regex" default:"false"`
    Enabled     bool   `yaml:"enabled" default:"true"`
}

// RegionRule 地区分组规则
type RegionRule struct {
    Name     string   `yaml:"name" validate:"required"`
    Regions  []string `yaml:"regions" validate:"required"`
    Patterns []string `yaml:"patterns" validate:"required"`
    Emoji    bool     `yaml:"emoji" default:"true"`
    Enabled  bool     `yaml:"enabled" default:"true"`
}
```

### 3.2 规则引擎实现

```go
// pkg/rules/engine.go
package rules

import (
    "regexp"
    "strings"
    "github.com/subconverter-go/pkg/config"
    "github.com/subconverter-go/pkg/models"
)

// Engine 规则引擎
type Engine struct {
    config   *config.RulesConfig
    filters  map[string]*regexp.Regexp
    renames  map[string]*regexp.Regexp
}

// NewEngine 创建规则引擎
func NewEngine(cfg *config.RulesConfig) (*Engine, error) {
    engine := &Engine{
        config:  cfg,
        filters: make(map[string]*regexp.Regexp),
        renames: make(map[string]*regexp.Regexp),
    }

    return engine, engine.compile()
}

// ApplyFilters 应用过滤规则
func (e *Engine) ApplyFilters(proxies []models.Proxy) []models.Proxy {
    var result []models.Proxy

    for _, proxy := range proxies {
        if e.shouldInclude(proxy.Name) {
            result = append(result, proxy)
        }
    }

    return result
}

// ApplyRename 应用重命名规则
func (e *Engine) ApplyRename(proxies []models.Proxy) []models.Proxy {
    for i := range proxies {
        proxies[i].Name = e.rename(proxies[i].Name)
    }
    return proxies
}

// ApplyRegionGroup 应用地区分组
func (e *Engine) ApplyRegionGroup(proxies []models.Proxy) map[string][]models.Proxy {
    groups := make(map[string][]models.Proxy)

    for _, proxy := range proxies {
        region := e.getRegion(proxy.Name)
        groups[region] = append(groups[region], proxy)
    }

    return groups
}
```

## 4. 模板系统设计

### 4.1 模板配置结构

```go
// pkg/config/template.go
package config

// TemplateConfig 模板配置
type TemplateConfig struct {
    BaseDir     string            `yaml:"base_dir" default:"templates"`
    Extension   string            `yaml:"extension" default:".tpl"`
    AutoReload  bool              `yaml:"auto_reload" default:"true"`
    Variables   map[string]string `yaml:"variables"`
    Functions   []string          `yaml:"functions"`
    Clients     []ClientTemplate  `yaml:"clients"`
}

// ClientTemplate 客户端模板配置
type ClientTemplate struct {
    Name     string            `yaml:"name" validate:"required"`
    File     string            `yaml:"file" validate:"required"`
    Type     string            `yaml:"type" validate:"required"`
    Variables map[string]string `yaml:"variables"`
    Enabled  bool              `yaml:"enabled" default:"true"`
}
```

### 4.2 模板引擎实现

```go
// pkg/template/engine.go
package template

import (
    "bytes"
    "context"
    "path/filepath"
    "github.com/flosch/pongo2/v6"
    "github.com/subconverter-go/pkg/config"
    "github.com/subconverter-go/pkg/models"
)

// Engine 模板引擎
type Engine struct {
    config    *config.TemplateConfig
    templates map[string]*pongo2.Template
    functions pongo2.Context
}

// NewEngine 创建模板引擎
func NewEngine(cfg *config.TemplateConfig) (*Engine, error) {
    engine := &Engine{
        config:    cfg,
        templates: make(map[string]*pongo2.Template),
        functions: make(pongo2.Context),
    }

    return engine, engine.initialize()
}

// Render 渲染模板
func (e *Engine) Render(ctx context.Context, client string, data *RenderData) (string, error) {
    template, exists := e.templates[client]
    if !exists {
        return "", ErrTemplateNotFound
    }

    context := e.buildContext(data)

    var buf bytes.Buffer
    err := template.ExecuteWriter(context, &buf)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

// RenderData 模板渲染数据
type RenderData struct {
    Proxies     []models.Proxy    `json:"proxies"`
    Groups      []models.Group    `json:"groups"`
    Rules       []models.Rule     `json:"rules"`
    Variables   map[string]string `json:"variables"`
    Metadata    *Metadata         `json:"metadata"`
}

// Metadata 元数据
type Metadata struct {
    Title       string    `json:"title"`
    Description string    `json:"description"`
    UpdateTime  time.Time `json:"update_time"`
    Version     string    `json:"version"`
    Author      string    `json:"author"`
}
```

### 4.3 自定义模板函数

```go
// pkg/template/functions.go
package template

import (
    "encoding/base64"
    "net/url"
    "strings"
    "github.com/flosch/pongo2/v6"
)

// RegisterFunctions 注册自定义函数
func (e *Engine) RegisterFunctions() {
    // Base64编码
    pongo2.RegisterFilter("b64encode", filterBase64Encode)

    // URL编码
    pongo2.RegisterFilter("urlencode", filterURLEncode)

    // 字符串处理
    pongo2.RegisterFilter("trim_prefix", filterTrimPrefix)
    pongo2.RegisterFilter("trim_suffix", filterTrimSuffix)

    // 节点处理
    pongo2.RegisterFilter("filter_by_type", filterByType)
    pongo2.RegisterFilter("group_by_region", groupByRegion)

    // 自定义全局函数
    e.functions["get_proxy_count"] = getProxyCount
    e.functions["format_speed"] = formatSpeed
    e.functions["generate_uuid"] = generateUUID
}

// filterBase64Encode Base64编码过滤器
func filterBase64Encode(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
    encoded := base64.StdEncoding.EncodeToString([]byte(in.String()))
    return pongo2.AsValue(encoded), nil
}

// filterURLEncode URL编码过滤器
func filterURLEncode(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
    encoded := url.QueryEscape(in.String())
    return pongo2.AsValue(encoded), nil
}
```

## 5. 配置加载和验证

### 5.1 配置加载器

```go
// pkg/config/loader.go
package config

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "github.com/spf13/viper"
    "github.com/go-playground/validator/v10"
)

// Loader 配置加载器
type Loader struct {
    configDir  string
    validator  *validator.Validate
}

// NewLoader 创建配置加载器
func NewLoader(configDir string) *Loader {
    return &Loader{
        configDir: configDir,
        validator: validator.New(),
    }
}

// LoadConfig 加载完整配置
func (l *Loader) LoadConfig(ctx context.Context) (*Config, error) {
    config := &Config{}

    // 加载应用配置
    appConfig, err := l.loadAppConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to load app config: %w", err)
    }
    config.App = appConfig

    // 加载规则配置
    rulesConfig, err := l.loadRulesConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to load rules config: %w", err)
    }
    config.Rules = rulesConfig

    // 加载模板配置
    templateConfig, err := l.loadTemplateConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to load template config: %w", err)
    }
    config.Templates = templateConfig

    // 验证配置
    if err := l.validator.Struct(config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return config, nil
}

// loadAppConfig 加载应用配置
func (l *Loader) loadAppConfig() (*AppConfig, error) {
    v := viper.New()
    v.SetConfigName("app")
    v.SetConfigType("yaml")
    v.AddConfigPath(l.configDir)

    // 设置默认值
    l.setAppDefaults(v)

    // 读取环境变量
    v.AutomaticEnv()
    v.SetEnvPrefix("SC")

    if err := v.ReadInConfig(); err != nil {
        return nil, err
    }

    var config AppConfig
    if err := v.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

### 5.2 配置验证系统

```go
// pkg/config/validator.go
package config

import (
    "fmt"
    "net"
    "regexp"
    "github.com/go-playground/validator/v10"
)

// SetupValidator 设置配置验证器
func SetupValidator() *validator.Validate {
    validate := validator.New()

    // 注册自定义验证器
    validate.RegisterValidation("host", validateHost)
    validate.RegisterValidation("regex", validateRegex)
    validate.RegisterValidation("template_file", validateTemplateFile)

    return validate
}

// validateHost 验证主机地址
func validateHost(fl validator.FieldLevel) bool {
    host := fl.Field().String()
    if host == "" {
        return false
    }

    // 检查是否为有效IP地址
    if ip := net.ParseIP(host); ip != nil {
        return true
    }

    // 检查是否为有效域名
    if len(host) > 253 {
        return false
    }

    return regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`).MatchString(host)
}

// validateRegex 验证正则表达式
func validateRegex(fl validator.FieldLevel) bool {
    pattern := fl.Field().String()
    _, err := regexp.Compile(pattern)
    return err == nil
}

// validateTemplateFile 验证模板文件
func validateTemplateFile(fl validator.FieldLevel) bool {
    filename := fl.Field().String()
    if filename == "" {
        return false
    }

    // 检查文件扩展名
    return strings.HasSuffix(filename, ".tpl") || strings.HasSuffix(filename, ".j2")
}
```

## 6. 动态配置更新

### 6.1 配置监听器

```go
// pkg/config/watcher.go
package config

import (
    "context"
    "log"
    "path/filepath"
    "github.com/fsnotify/fsnotify"
)

// Watcher 配置文件监听器
type Watcher struct {
    configDir string
    watcher   *fsnotify.Watcher
    events    chan ConfigEvent
    manager   ConfigManager
}

// NewWatcher 创建配置监听器
func NewWatcher(configDir string, manager ConfigManager) (*Watcher, error) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return nil, err
    }

    w := &Watcher{
        configDir: configDir,
        watcher:   watcher,
        events:    make(chan ConfigEvent, 100),
        manager:   manager,
    }

    return w, nil
}

// Start 启动监听
func (w *Watcher) Start(ctx context.Context) error {
    // 添加配置目录监听
    if err := w.watcher.Add(w.configDir); err != nil {
        return err
    }

    go w.watchLoop(ctx)
    return nil
}

// watchLoop 监听循环
func (w *Watcher) watchLoop(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case event := <-w.watcher.Events:
            w.handleFileEvent(event)
        case err := <-w.watcher.Errors:
            log.Printf("Config watcher error: %v", err)
        }
    }
}

// handleFileEvent 处理文件事件
func (w *Watcher) handleFileEvent(event fsnotify.Event) {
    if !w.isConfigFile(event.Name) {
        return
    }

    switch event.Op {
    case fsnotify.Write, fsnotify.Create:
        w.triggerReload(event.Name)
    case fsnotify.Remove:
        log.Printf("Config file removed: %s", event.Name)
    }
}

// triggerReload 触发配置重载
func (w *Watcher) triggerReload(filename string) {
    ctx := context.Background()
    if err := w.manager.Reload(ctx); err != nil {
        log.Printf("Failed to reload config: %v", err)
        w.events <- ConfigEvent{
            Type:   EventTypeError,
            Source: filename,
        }
    } else {
        w.events <- ConfigEvent{
            Type:   EventTypeReload,
            Source: filename,
        }
    }
}
```

## 7. 客户端配置模板

### 7.1 Clash 配置模板

```yaml
# templates/clash.tpl
# Clash 配置模板
port: {{ port | default("7890") }}
socks-port: {{ socks_port | default("7891") }}
allow-lan: {{ allow_lan | default("false") }}
mode: {{ mode | default("rule") }}
log-level: {{ log_level | default("info") }}
external-controller: {{ external_controller | default("127.0.0.1:9090") }}

proxies:
{% for proxy in proxies %}
  - name: "{{ proxy.name }}"
    type: {{ proxy.type }}
    server: {{ proxy.server }}
    port: {{ proxy.port }}
    {% if proxy.type == "ss" %}
    cipher: {{ proxy.cipher }}
    password: "{{ proxy.password }}"
    {% elif proxy.type == "vmess" %}
    uuid: {{ proxy.uuid }}
    alterId: {{ proxy.alter_id | default("0") }}
    cipher: {{ proxy.cipher | default("auto") }}
    {% if proxy.network %}
    network: {{ proxy.network }}
    {% endif %}
    {% if proxy.tls %}
    tls: true
    {% endif %}
    {% endif %}
{% endfor %}

proxy-groups:
{% for group in groups %}
  - name: "{{ group.name }}"
    type: {{ group.type }}
    proxies:
    {% for proxy_name in group.proxies %}
      - "{{ proxy_name }}"
    {% endfor %}
    {% if group.url %}
    url: "{{ group.url }}"
    interval: {{ group.interval | default("300") }}
    {% endif %}
{% endfor %}

rules:
{% for rule in rules %}
  - {{ rule.type }},{{ rule.payload }},{{ rule.policy }}
{% endfor %}
```

### 7.2 Surge 配置模板

```ini
# templates/surge.tpl
# Surge 配置模板
[General]
loglevel = {{ log_level | default("notify") }}
dns-server = {{ dns_server | default("223.5.5.5, 114.114.114.114") }}
skip-proxy = {{ skip_proxy | default("127.0.0.1, 192.168.0.0/16, 10.0.0.0/8, 172.16.0.0/12, 100.64.0.0/10, localhost, *.local") }}
bypass-tun = {{ bypass_tun | default("192.168.0.0/16, 10.0.0.0/8, 172.16.0.0/12") }}

[Proxy]
{% for proxy in proxies %}
{{ proxy.name }} = {{ proxy.type }}, {{ proxy.server }}, {{ proxy.port }}{% if proxy.username %}, {{ proxy.username }}{% endif %}{% if proxy.password %}, {{ proxy.password }}{% endif %}{% if proxy.encrypt_method %}, encrypt-method={{ proxy.encrypt_method }}{% endif %}{% if proxy.obfs %}, obfs={{ proxy.obfs }}{% endif %}{% if proxy.obfs_host %}, obfs-host={{ proxy.obfs_host }}{% endif %}
{% endfor %}

[Proxy Group]
{% for group in groups %}
{{ group.name }} = {{ group.type }}{% for proxy_name in group.proxies %}, {{ proxy_name }}{% endfor %}{% if group.url %}, url={{ group.url }}{% endif %}{% if group.interval %}, interval={{ group.interval }}{% endif %}
{% endfor %}

[Rule]
{% for rule in rules %}
{{ rule.type }}, {{ rule.payload }}, {{ rule.policy }}
{% endfor %}
```

## 8. 性能优化策略

### 8.1 配置缓存

```go
// pkg/config/cache.go
package config

import (
    "sync"
    "time"
    "github.com/patrickmn/go-cache"
)

// ConfigCache 配置缓存
type ConfigCache struct {
    cache      *cache.Cache
    templates  *sync.Map
    lastUpdate time.Time
    mutex      sync.RWMutex
}

// NewConfigCache 创建配置缓存
func NewConfigCache(defaultTTL time.Duration) *ConfigCache {
    return &ConfigCache{
        cache:     cache.New(defaultTTL, time.Minute),
        templates: &sync.Map{},
    }
}

// GetTemplate 获取缓存的模板
func (c *ConfigCache) GetTemplate(name string) (*pongo2.Template, bool) {
    value, exists := c.templates.Load(name)
    if !exists {
        return nil, false
    }
    return value.(*pongo2.Template), true
}

// SetTemplate 设置模板缓存
func (c *ConfigCache) SetTemplate(name string, template *pongo2.Template) {
    c.templates.Store(name, template)
}
```

### 8.2 模板预编译

```go
// pkg/template/precompile.go
package template

import (
    "path/filepath"
    "github.com/flosch/pongo2/v6"
)

// PrecompileTemplates 预编译所有模板
func (e *Engine) PrecompileTemplates() error {
    for _, client := range e.config.Clients {
        if !client.Enabled {
            continue
        }

        templatePath := filepath.Join(e.config.BaseDir, client.File)

        template, err := pongo2.FromFile(templatePath)
        if err != nil {
            return fmt.Errorf("failed to compile template %s: %w", client.Name, err)
        }

        e.templates[client.Name] = template
    }

    return nil
}
```

## 9. 总结

这个配置管理和模板系统设计具有以下特点：

### 9.1 核心优势

- **完全兼容**: 支持现有所有配置格式和模板语法
- **高性能**: 模板预编译和配置缓存优化
- **灵活配置**: 支持多层次配置和动态更新
- **强验证**: 完整的配置验证和错误处理

### 9.2 技术特性

- **模块化设计**: 清晰的接口分离和组件划分
- **异步处理**: 配置监听和动态更新
- **缓存优化**: 模板和配置的智能缓存
- **扩展性**: 易于添加新的客户端模板和规则

### 9.3 迁移友好

- **API 兼容**: 与 C++版本配置接口完全兼容
- **模板兼容**: Pongo2 提供 Jinja2 完全兼容性
- **配置格式**: 支持 YAML/JSON/TOML 多种格式
- **增量迁移**: 支持逐步迁移配置文件

这个设计为 SubConverter Go 版本提供了强大而灵活的配置管理和模板系统基础。
