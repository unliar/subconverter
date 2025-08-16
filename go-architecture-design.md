\n# SubConverter Go 版本项目架构设计\n\n## 1. 架构设计原则\n\n### 1.1 兼容性优先原则\n- **API 完全兼容**：保持与现有 C++ 版本 100% 的 API 兼容性\n- **配置格式兼容**：支持现有的所有配置文件格式（YAML、TOML、INI）\n- **模板系统兼容**：确保现有模板文件无需修改即可使用\n- **数据结构兼容**：保持与现有输出格式的完全一致性\n\n### 1.2 模块化设计原则\n- **清晰的职责分离**：每个模块负责单一功能领域\n- **松耦合设计**：模块间通过接口交互，降低依赖关系\n- **高内聚**：相关功能集中在同一模块内\n\n### 1.3 性能优化原则\n- **并发处理**：利用 Go 的 goroutine 优化并发请求处理\n- **内存管理**：优化大文件和批量节点处理的内存使用\n- **缓存策略**：合理使用缓存减少重复计算\n\n## 2. 整体架构图\n\n`\n┌─────────────────────────────────────────────────────────────────┐\n│                        HTTP API Layer                          │\n├─────────────────────────────────────────────────────────────────┤\n│  GET /sub │ GET /version │ GET /getruleset │ POST /updateconf   │\n└─────────────────────────────────────────────────────────────────┘\n                                    │\n                                    ▼\n┌─────────────────────────────────────────────────────────────────┐\n│                     Service Layer (Business Logic)             │\n├─────────────────────────────────────────────────────────────────┤\n│ ConverterService │ RulesetService │ ConfigService │ ProfileService│\n└─────────────────────────────────────────────────────────────────┘\n                                    │\n                                    ▼\n┌─────────────────────────────────────────────────────────────────┐\n│                        Core Processing Layer                   │\n├──────────────────┬─────────────────┬─────────────────────────────┤\n│   Parser Module  │ Generator Module│      Utility Module         │\n│                  │                 │                             │\n│ • SubscriptionParser │ • ClashGenerator │ • NetworkUtils        │\n│ • NodeParser         │ • SurgeGenerator │ • StringUtils         │\n│ • ConfigParser       │ • QuanXGenerator │ • CryptoUtils         │\n│ • RulesetParser      │ • LoonGenerator  │ • FileUtils           │\n└──────────────────┴─────────────────┴─────────────────────────────┘\n                                    │\n                                    ▼\n┌─────────────────────────────────────────────────────────────────┐\n│                       Data Model Layer                         │\n├─────────────────────────────────────────────────────────────────┤\n│   Proxy Models   │  Ruleset Models │ Config Models │Template Models│\n└─────────────────────────────────────────────────────────────────┘\n                                    │\n                                    ▼\n┌─────────────────────────────────────────────────────────────────┐\n│                    Infrastructure Layer                        │\n├─────────────────────────────────────────────────────────────────┤\n│   HTTP Client    │    File System  │   Template Engine │  Logger │\n└─────────────────────────────────────────────────────────────────┘\n`\n\n## 3. 核心模块设计\n\n### 3.1 HTTP API Layer - Web 服务层\n\n#### 接口兼容性保证\n- **完全兼容现有 API 端点**：\n - `GET /sub` - 订阅转换核心接口\n - `GET /version` - 版本信息接口\n - `GET /getruleset` - 规则集获取接口\n - `GET /getprofile` - 配置档案接口\n - `POST /updateconf` - 配置更新接口\n - `GET /refreshrules` - 规则刷新接口\n - `GET /readconf` - 配置重载接口\n\n#### 技术实现方案\n`go\n// HTTP 路由器设计\ntype Router struct {\n    gin    *gin.Engine\n    config *config.ServerConfig\n}\n\n// API 处理器接口\ntype Handler interface {\n    Handle(c *gin.Context) error\n}\n\n// 主要处理器\ntype ConverterHandler struct {\n    service ConverterService\n}\n`\n\n### 3.2 Service Layer - 业务逻辑层\n\n#### 3.2.1 ConverterService - 转换服务\n 负责整个订阅转换的业务逻辑流程：\n\n```go\ntype ConverterService interface {\n ConvertSubscription(req *ConvertRequest) (*ConvertResponse, error)\n ValidateRequest(req \*ConvertRequest) error\n ProcessFilters(nodes []Proxy, filters []Filter) []Proxy\n}\n\ntype ConvertReq

\nuest struct {\n URL string `json:\"url\" form:\"url\"`\n Target string `json:\"target\" form:\"target\"`\n Config string `json:\"config\" form:\"config\"`\n Filename string `json:\"filename\" form:\"filename\"`\n Interval int `json:\"interval\" form:\"interval\"`\n Strict bool `json:\"strict\" form:\"strict\"`\n IncludeFilters []string `json:\"include\" form:\"include\"`\n ExcludeFilters []string `json:\"exclude\" form:\"exclude\"`\n Sort bool `json:\"sort\" form:\"sort\"`\n FilterDeprecated bool `json:\"fdn\" form:\"fdn\"`\n AppendType bool `json:\"append_type\" form:\"append_type\"`\n List bool `json:\"list\" form:\"list\"`\n UDP *bool `json:\"udp\" form:\"udp\"`\n TFO *bool `json:\"tfo\" form:\"tfo\"`\n SkipCertVerify *bool `json:\"scv\" form:\"scv\"`\n Emoji *bool `json:\"emoji\" form:\"emoji\"`\n}\n`\n\n#### 3.2.2 RulesetService - 规则集服务\n管理规则集的获取、缓存和更新：\n\n`go\ntype RulesetService interface {\n GetRuleset(name string) (*RulesetContent, error)\n RefreshRulesets() error\n ValidateRuleset(content string) error\n}\n`\n\n#### 3.2.3 ConfigService - 配置服务\n处理应用配置的读取、更新和验证：\n\n`go\ntype ConfigService interface {\n ReadConfig() error\n UpdateConfig(content string) error\n GetServerConfig() *ServerConfig\n GetConverterConfig() _ConverterConfig\n}\n`\n\n### 3.3 Parser Module - 解析器模块\n\n#### 3.3.1 接口设计保证兼容性\n`go\n// 统一解析器接口\ntype Parser interface {\n Parse(content string) ([]Proxy, error)\n GetType() ProxyType\n}\n\n// 订阅解析器\ntype SubscriptionParser interface {\n ParseSubscription(url string) ([]Proxy, error)\n DetectFormat(content string) ConfType\n}\n\n// 节点解析器 - 对应 C++ 中的 explode_ 函数\ntype NodeParser struct{}\n\nfunc (p *NodeParser) ParseVMess(link string) (*Proxy, error) // 对应 explodeVmess\nfunc (p *NodeParser) ParseVless(link string) (*Proxy, error) // 对应 explodeVless \nfunc (p *NodeParser) ParseSS(link string) (*Proxy, error) // 对应 explodeSS\nfunc (p *NodeParser) ParseSSR(link string) (*Proxy, error) // 对应 explodeSSR\nfunc (p *NodeParser) ParseTrojan(link string) (*Proxy, error) // 对应 explodeTrojan\nfunc (p *NodeParser) ParseHysteria(link string) (*Proxy, error) // 对应 explodeHysteria\nfunc (p *NodeParser) ParseHysteria2(link string) (*Proxy, error)// 对应 explodeHysteria2\nfunc (p *NodeParser) ParseTUIC(link string) (*Proxy, error) // 对应 explodeTUIC\n`\n\n### 3.4 Generator Module - 生成器模块\n\n#### 3.4.1 生成器接口设计\n`go\n// 配置生成器接口\ntype Generator interface {\n Generate(nodes []Proxy, config GeneratorConfig) (string, error)\n GetType() string\n ValidateTemplate(template string) error\n}\n\n// 具体生成器实现 - 对应 C++ 中的 proxyTo* 函数\ntype ClashGenerator struct{} // 对应 proxyToClash\ntype SurgeGenerator struct{} // 对应 proxyToSurge\ntype QuanXGenerator struct{} // 对应 proxyToQuanX\ntype LoonGenerator struct{} // 对应 proxyToLoon\ntype SingBoxGenerator struct{} // 对应 proxyToSingBox\n`\n\n## 4. 数据模型设计\n\n### 4.1 Proxy 数据结构 - 完全兼容 C++ 版本\n\n`go\n// ProxyType 枚举 - 完全对应 C++ 版本\ntype ProxyType int\n\nconst (\n ProxyTypeUnknown ProxyType = iota\n ProxyTypeShadowsocks\n ProxyTypeShadowsocksR\n ProxyTypeVMess\n ProxyTypeTrojan\n ProxyTypeSnell\n ProxyTypeHTTP\n ProxyTypeHTTPS\n ProxyTypeSOCKS5\n ProxyTypeWireGuard\n ProxyTypeVLESS\n ProxyTypeHysteria\n ProxyTypeHysteria2\n ProxyTypeTUIC\n ProxyTypeAnyTLS\n ProxyTypeMieru\n)\n\n// Proxy 结构体 - 对应 C++ 中的 Proxy struct\ntype Proxy struct {\n Type ProxyType `json:\"type\"`\n ID uint32 `json:\"id\"`\n GroupID uint32 `json:\"group_id\"`\n Group string `json:\"group\"`\n Remark string `json:\"remark\"`\n Hostname string `json:\"hostname\"`\n Port uint16 `json:\"port\"`\n CongestionControl string `json:\"congestion_control\"`\n Username string `json:\"username\"`\n Password string `json:\"password\"`\n EncryptMethod string `json:\"encrypt_method\"`\n Plugin string `json:\"plugin\"`\n PluginOption string `json:\"plugin_option\"`\n Protocol string `json:\"protocol\"`\n ProtocolParam string `json:\"protocol_param\"`\n OBFS string `json:\"obfs\"`\n OBFSParam string `json:\"obfs_param\"`\n UserID string `json:\"user_id\"`\n AlterID uint16 `json:\"alter_id\"`\n TransferProtocol string `json:\"transfer_protocol\"`\n FakeType string `json:\"fake_type\"`\n AuthStr string `json:\"auth_str\"`\n \n // TLS 相关\n TLSStr string `json:\"tls_str\"`\n TLSSecure bool `json:\"tls_secure\"`\n \n // 网络相关\n Host string `json:\"host\"`\n Path string `json:\"path\"`\n Edge string `json:\"edge\"`\n \n // QUIC 相关\n QUICSecure string `json:\"quic_secure\"`\n QUICSecret string `json:\"quic_secret\"`\n \n // 特性开关 - 使用指针实现三态逻辑\n UDP *bool `json:\"udp,omitempty\"`\n XUDP *bool `json:\"xudp,omitempty\"`\n TCPFastOpen *bool `json:\"tfo,omitempty\"`\n AllowInsecure *bool `json:\"allow_insecure,omitempty\"`\n TLS13 *bool `json:\"tls13,omitempty\"`\n \n // WireGuard 专用字段\n SelfIP string `jso

\nn\":\"self*ip\"`\n    SelfIPv6            string    `json:\"self_ipv6\"`\n    PublicKey           string    `json:\"public_key\"`\n    PrivateKey          string    `json:\"private_key\"`\n    PreSharedKey        string    `json:\"pre_shared_key\"`\n    DNSServers          []string  `json:\"dns_servers\"`\n    MTU                 uint16    `json:\"mtu\"`\n    AllowedIPs          string    `json:\"allowed_ips\"`\n    KeepAlive           uint16    `json:\"keep_alive\"`\n    \n    // Hysteria/TUIC 专用字段\n    TestURL             string    `json:\"test_url\"`\n    ClientID            string    `json:\"client_id\"`\n    Ports               string    `json:\"ports\"`\n    Auth                string    `json:\"auth\"`\n    ALPN                string    `json:\"alpn\"`\n    UpMbps              string    `json:\"up_mbps\"`\n    DownMbps            string    `json:\"down_mbps\"`\n    Insecure            string    `json:\"insecure\"`\n    Fingerprint         string    `json:\"fingerprint\"`\n    OBFSPassword        string    `json:\"obfs_password\"`\n    UDPRelayMode        string    `json:\"udp_relay_mode\"`\n    RequestTimeout      uint16    `json:\"request_timeout\"`\n    Token               string    `json:\"token\"`\n    \n    // 其他扩展字段\n    UnderlyingProxy     string    `json:\"underlying_proxy\"`\n    ALPNList            []string  `json:\"alpn_list\"`\n    PacketEncoding      string    `json:\"packet_encoding\"`\n    Multiplexing        string    `json:\"multiplexing\"`\n    V2rayHTTPUpgrade    *bool     `json:\"v2ray_http_upgrade,omitempty\"` \n}\n```\n\n### 4.2 配置和规则集模型\n\n```go\n// 服务器配置\ntype ServerConfig struct {\n    ListenAddress     string  `yaml:\"listen_address\"`\n    ListenPort        int    `yaml:\"listen_port\"`\n    MaxPendingConns   int    `yaml:\"max_pending_conns\"`\n    MaxConcurThreads  int    `yaml:\"max_concur_threads\"`\n    APIMode           bool   `yaml:\"api_mode\"`\n    AccessToken       string `yaml:\"access_token\"`\n}\n\n// 转换器配置\ntype ConverterConfig struct {\n    DefaultConfig     string                `yaml:\"default_config\"`\n    EnableRuleGenerator bool               `yaml:\"enable_rule_generator\"`\n    OverwriteOriginalRules bool            `yaml:\"overwrite_original_rules\"`\n    CustomRulesets    []RulesetConfig       `yaml:\"custom_rulesets\"`\n    ProxyGroups       []ProxyGroupConfig    `yaml:\"proxy_groups\"`\n}\n\n// 规则集配置\ntype RulesetConfig struct {\n    Group    string `yaml:\"group\"`\n    URL      string `yaml:\"url\"`\n    Interval int    `yaml:\"interval\"`\n}\n\n// 代理组配置\ntype ProxyGroupConfig struct {\n    Name     string   `yaml:\"name\"`\n    Type     string   `yaml:\"type\"`\n    Proxies  []string `yaml:\"proxies\"`\n    URL      string   `yaml:\"url\"`\n    Interval int      `yaml:\"interval\"`\n    Timeout  int      `yaml:\"timeout\"` \n}\n```\n\n## 5. 兼容性实现策略\n\n### 5.1 API 兼容性\n\n#### 请求参数映射\n```go\n// 完全兼容 C++ 版本的请求参数\nvar CompatibleParams = map[string]string{\n    \"url\":         \"订阅链接\",\n    \"target\":      \"目标客户端类型\", \n    \"config\":      \"配置文件链接\",\n    \"filename\":    \"生成的配置文件名\",\n    \"interval\":    \"更新间隔\",\n    \"strict\":      \"严格模式\",\n    \"include\":     \"包含过滤器\",\n    \"exclude\":     \"排除过滤器\",\n    \"sort\":        \"节点排序\",\n    \"fdn\":         \"过滤废弃节点\",\n    \"append_type\": \"添加节点类型\",\n    \"list\":        \"生成节点列表\",\n    \"udp\":         \"UDP 支持\",\n    \"tfo\":         \"TCP Fast Open\",\n    \"scv\":         \"跳过证书验证\",\n    \"emoji\":       \"添加 emoji\",\n}\n```\n\n#### 响应格式兼容\n```go\n// 保持与 C++ 版本完全一致的响应格式\ntype ConvertResponse struct {\n    Content     string             `json:\"-\"`             // 主要内容\n    Headers     map[string]string`json:\"-\"`            // HTTP 头部\n    StatusCode  int             `json:\"-\"`            // 状态码\n    ContentType string          `json:\"-\"` // 内容类型\n}\n`\n\n### 5.2 配置文件兼容性\n\n#### 多格式支持\n`go\ntype ConfigReader interface {\n ReadConfig(path string) (*Config, error)\n GetFormat() string\n}\n\n// YAML 配置读取器\ntype YAMLConfigReader struct{}\nfunc (r *YAMLConfigReader) ReadConfig(path string) (*Config, error) {\n // 使用 gopkg.in/yaml.v3 保持与 C++ yaml-cpp 的兼容性\n}\n\n// TOML 配置读取器 \ntype TOMLConfigReader struct{}\nfunc (r *TOMLConfigReader) ReadConfig(path string) (*Config, error) {\n // 使用 github.com/BurntSushi/toml 保持与 C++ toml11 的兼容性\n}\n\n// INI 配置读取器\ntype INIConfigReader struct{}\nfunc (r *INIConfigReader) ReadConfig(path string) (*Config, error) {\n // 保持与 C++ INI 格式的兼容性\n}\n`\n\n### 5.3 模板系统兼容性\n\n#### 模板引擎适配\n`go\n// 模板引擎接口 - 兼容 Jinja2 语法\ntype TemplateEngine interface {\n Render(template string, data interface{}) (string, error)\n ParseTemplate(content string) (\*Template, error)\n}\n\n// Go 模板引擎适配器\ntype GoTemplateEngine struct {\n funcMap template.FuncMap\n}\n\n// 实现 Jinja2 兼容的函数映射\nfunc (e \_GoTemplateEngine) buildCompatibleFuncMap() template.FuncMap {\n return template.FuncMap{\n // 兼容 Jinja2 的常用过滤器\n \"length\": func(v interface{}) int { /* 实现 _/ },\n \"upper\": strings.ToUpper,\n \"lower\": strings.ToLower,\n \"default\": func(v, d interface{}) interface{} { /_ 实现 _/ },\n \"join\": func(sep string, items []string) string { return strings.Join(items, sep) },\n \"split\": func(sep, s string) []string { return strings.Split(s, sep) },\n \"replace\": func(old, new, s string) string { return strings.ReplaceAll(s, old, new) },\n // 自定义函数保持兼容\n \"webget\": func(url string) string { /_ HTTP 请求实现 _/ },\n \"parseHostname\": func(url string) string { /_ 解析主机名实现 */ },\n }\n}\n`\n\n## 6. 性能优化设计\n\n### 6.1 并发处理架构\n\n`go\n// 并发请求处理器\ntype ConcurrentProcessor struct {\n workerPool *WorkerPool\n resultCache *Cache\n rateLimiter *RateLimiter\n}\n\n// 工作池设计\ntype WorkerPool struct {\n workers chan chan ConvertJob\n jobQueue chan ConvertJob\n maxWorkers int\n}\n\n// 并发安全的缓存\ntype Cache struct {\n data sync.Map\n ttl time.Duration\n janitor *time.Ticker\n}\n`\n\n### 6.2 内存优化策略\n\n`go\n// 流式处理大型订阅\ntype StreamProcessor struct {\n bufferSize int\n parser Parser\n}\n\nfunc (p *StreamProcessor) ProcessLargeSubscription(reader io.Reader) <-chan Proxy {\n resultChan := make(chan Proxy, p.bufferSize)\n go func() {\n defer close(resultChan)\n // 流式解析和处理\n scanner := bufio.NewScanner(reader)\n for scanner.Scan() {\n if proxy, err := p.parser.ParseLine(scanner.Text()); err == nil {\n resultChan <- proxy\n }\n }\n }()\n return resultChan\n}\n`\n\n## 7. 第三方依赖选择\n\n### 7.1 核心依赖库\n\n| 功能领域 | Go 库 | C++ 对应库 | 选择理由 |\n|---------|-------|------------|----------|\n| Web 框架 | github.com/gin-gonic/gin | httplib | 高性能、简洁 API |\n| YAML 处理 | gopkg.in/yaml.v3 | yaml-cpp | 完全兼容、广泛使用 |\n| JSON 处理 | encoding/json (标准库) | rapidjson | 标准库、性能优秀 |\n| 正则表达式 | regexp (标准库) | PCRE2 | 标准库、功能完备 |\n| HTTP 客户端 | net/http (标准库) | libcurl | 标准库、功能齐全 |\n| 模板引擎 | text/template (标准库) | inja | 标准库、可扩展 |\n| 配置管理 | github.com/spf13/viper | 自实现 | 功能强大、多格式支持 |\n| 日志记录 | github.com/sirupsen/logrus | 自实现 | 结构化日志、级别控制 |\n| Base64 编码 | encoding/base64 (标准库) | 自实现 | 标准库、性能优秀 |\n| URL 编码 | net/url (标准库) | 自实现 | 标准库、完全兼容 |\n\n### 7.2 依赖版本管理\n\n`go\n// go.mod 文件示例\nmodule github.com/user/subconverter-go\n\ngo 1.21\n\nrequire (\n github.com/gin-gonic/gin v1.9.1\n github.com/spf13/viper v1.16.0\n github.com/sirupsen/logrus v1.9.3\n gopkg.in/yaml.v3 v3.0.1\n)\n`\n\n## 8. 部署和运维考虑\n\n### 8.1 容器化支持\n\n`dockerfile\n# Dockerfile\nFROM golang:1.21-alpine AS builder\nWORKDIR /app\nCOPY . .\nRUN go mod download\nRUN CGO_ENA

BLED=0 go build -ldflags '-s -w' -o subconverter-go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/subconverter-go .
COPY --from=builder /app/base ./base
EXPOSE 25500
CMD ["./subconverter-go"]

````

### 8.2 配置管理
```go
// 支持环境变量覆盖配置
type EnvConfig struct {
    ListenPort    int    `env:"PORT" envDefault:"25500"`
    ListenAddress string `env:"LISTEN_ADDRESS" envDefault:"0.0.0.0"`
    APIMode       bool   `env:"API_MODE" envDefault:"false"`
    AccessToken   string `env:"ACCESS_TOKEN"`
}
````

## 9. 测试策略

### 9.1 单元测试覆盖率目标：95%+

```go
// 解析器测试示例
func TestVMessParser(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected *Proxy
        wantErr  bool
    }{
        {
            name: "标准VMess链接",
            input: "vmess://eyJ2IjoiMiIsInBzIjoid...",
            expected: &Proxy{
                Type: ProxyTypeVMess,
                Remark: "测试节点",
                // ...
            },
            wantErr: false,
        },
    }

    parser := &NodeParser{}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := parser.ParseVMess(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 9.2 集成测试

```go
// API 兼容性测试
func TestAPICompatibility(t *testing.T) {
    // 启动测试服务器
    server := setupTestServer()
    defer server.Close()

    // 测试与 C++ 版本的完全兼容性
    testCases := []struct {
        endpoint string
        params   map[string]string
        expected string
    }{
        {
            endpoint: "/sub",
            params: map[string]string{
                "target": "clash",
                "url": "https://example.com/sub",
            },
            expected: "clash配置内容",
        },
    }

    for _, tc := range testCases {
        // 执行兼容性测试
    }
}
```

### 9.3 性能基准测试

```go
func BenchmarkConvertSubscription(b *testing.B) {
    service := setupConverterService()
    req := &ConvertRequest{
        URL: "https://example.com/large-subscription",
        Target: "clash",
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := service.ConvertSubscription(req)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## 10. 迁移验证策略

### 10.1 对比测试框架

```go
// C++ 版本输出对比测试
type CompatibilityTester struct {
    cppEndpoint string  // C++ 版本的端点
    goEndpoint  string  // Go 版本的端点
}

func (t *CompatibilityTester) CompareOutputs(params map[string]string) error {
    cppResult := t.callCppVersion(params)
    goResult := t.callGoVersion(params)

    // 逐字节对比输出
    if !bytes.Equal(cppResult, goResult) {
        return fmt.Errorf("输出不匹配")
    }
    return nil
}
```

### 10.2 渐进式部署

```yaml
# 部署策略配置
deployment:
  strategy: blue-green
  validation:
    - compatibility_tests
    - performance_benchmarks
    - stress_tests
  rollback_triggers:
    - error_rate > 1%
    - response_time > 500ms
```

## 11. 总结

### 11.1 架构优势

1. **完全兼容性**：API、配置、模板系统与 C++ 版本 100% 兼容
2. **性能优化**：利用 Go 并发特性，处理性能显著提升
3. **代码质量**：现代化的 Go 语言特性，提高可维护性
4. **部署简化**：单文件部署，容器化支持，运维便利

### 11.2 关键技术特点

- **数据结构映射**：Go struct 完全对应 C++ struct，保证数据兼容性
- **接口抽象**：清晰的接口设计，便于扩展新的代理协议和客户端格式
- **并发处理**：goroutine 池优化大并发请求处理
- **模板兼容**：Jinja2 语法适配，现有模板无需修改

### 11.3 实施路径

1. **阶段一**：核心数据模型和解析器实现
2. **阶段二**：生成器模块和模板系统
3. **阶段三**：HTTP API 层和服务层
4. **阶段四**：性能优化和兼容性验证
5. **阶段五**：部署和运维工具完善

这个架构设计确保了 Go 版本与现有 C++ 版本的平滑迁移，同时获得更好的性能和可维护性。
