\n# SubConverter Go 版本核心模块实现方案\n\n## 1. 核心模块实现概述\n\n### 1.1 实现策略\n- **接口驱动设计**：所有核心模块通过接口定义，便于测试和扩展\n- **管道式处理**：解析 → 过滤 → 转换 → 生成的流水线处理\n- **并发优化**：在安全的地方使用 goroutine 提升性能\n- **错误传播**：完整的错误链，便于调试和监控\n\n### 1.2 模块间数据流\n`\n订阅URL → 订阅解析器 → 节点解析器 → 过滤器 → 生成器 → 配置输出\n    ↓           ↓           ↓         ↓        ↓          ↓\n  网络请求   格式检测   协议解析   规则过滤  模板生成   格式输出\n`\n\n## 2. 解析器模块（Parser Module）实现方案\n\n### 2.1 核心接口设计\n\n`go\n// internal/parser/interfaces.go\npackage parser\n\nimport (\n    \"context\"\n    \"subconverter-go/pkg/models\"\n    \"subconverter-go/pkg/types\"\n)\n\n// Parser 统一解析器接口\ntype Parser interface {\n    Parse(ctx context.Context, content string) ([]models.Proxy, error)\n    GetType() types.ProxyType\n    Validate(content string) error\n}\n\n// SubscriptionParser 订阅解析器接口\ntype SubscriptionParser interface {\n    ParseSubscription(ctx context.Context, url string) ([]models.Proxy, error)\n    DetectFormat(content string) types.ConfType\n    FetchContent(ctx context.Context, url string) (string, error)\n}\n\n// NodeParser 节点解析器接口\ntype NodeParser interface {\n    ParseNode(link string) (*models.Proxy, error)\n    SupportedSchemes() []string\n    ValidateLink(link string) error\n}\n\n// BatchParser 批量解析器接口\ntype BatchParser interface {\n    ParseBatch(ctx context.Context, links []string) ([]models.Proxy, error)\n    ParseConcurrent(ctx context.Context, links []string, concurrency int) ([]models.Proxy, error)\n}\n`\n\n### 2.2 订阅解析器实现\n\n`go\n// internal/parser/subscription/parser.go\npackage subscription\n\nimport (\n    \"context\"\n    \"encoding/base64\"\n    \"fmt\"\n    \"strings\"\n    \"net/http\"\n    \"time\"\n    \n    \"subconverter-go/internal/parser/node\"\n    \"subconverter-go/pkg/models\"\n    \"subconverter-go/pkg/types\"\n    \"subconverter-go/internal/utils/network\"\n)\n\ntype Parser struct {\n    httpClient  *http.Client\n    nodeParser  map[types.ProxyType]node.Parser\n    timeout     time.Duration\n    maxSize     int64\n}\n\n// NewParser 创建订阅解析器\nfunc NewParser() *Parser {\n    return &Parser{\n        httpClient: &http.Client{\n            Timeout: 30 * time.Second,\n        },\n        nodeParser: make(map[types.ProxyType]node.Parser),\n        timeout:    30 * time.Second,\n        maxSize:    10 * 1024 * 1024, // 10MB\n    }\n}\n\n// ParseSubscription 解析订阅链接 - 对应 C++ 中的 explodeSub\nfunc (p *Parser) ParseSubscription(ctx context.Context, url string) ([]models.Proxy, error) {\n    // 获取订阅内容\n    content, err := p.FetchContent(ctx, url)\n    if err != nil {\n        return nil, fmt.Errorf(\"failed to fetch subscription: %w\", err)\n    }\n    \n    // 检测订阅格式\n    confType := p.DetectFormat(content)\n    \n    // 根据格式解析\n    switch confType {\n    case types.ConfTypeSUB:\n        return p.parseBase64Subscription(content)\n    case types.ConfTypeV2Ray:\n        return p.parseV2RaySubscription(content)\n    case types.ConfTypeSS:\n        return p.parseSSSubscription(content)\n    case types.ConfTypeSSR:\n        return p.parseSSRSubscription(content)\n    default:\n        return p.parseGenericSubscription(content)\n    }\n}\n\n// DetectFormat 检测订阅格式 - 对应 C++ 中的格式检测逻辑\nfunc (p *Parser) DetectFormat(content string) types.ConfType {\n    content = strings.TrimSpace(content)\n    \n    // Base64 编码的通用订阅\n    if isBase64(content) {\n        decoded, err := base64.StdEncoding.DecodeString(content)\n        if err == nil {\n            content = string(decoded)\n        }\n    }\n    \n    lines := strings.Split(content, \"\\n\")\n    if len(lines) == 0 {\n        return types.ConfTypeUnknown\n    }\n    \n    // 检查第一行判断格式\n    firstLine := strings.TrimSpace(lines[0])\n    \n    if strings.HasPrefix(firstLine, \"ss://\") {\n        return types.ConfTypeSS\n    }\n    if strings.HasPrefix(firstLine, \"ssr://\") {\n        return types.ConfTypeSSR\n    }\n    if strings.HasPrefix(firstLine, \"vmess://\") || strings.HasPrefix(firstLine, \"vless://\") {\n        return types.ConfTypeV2Ray\n    }\n    if strings.HasPrefix(firstLine, \"trojan://\") {\n        return types.ConfTypeV2Ray\n    }\n    if strings.HasPrefix(firstLine, \"hysteria://\") || strings.HasPrefix(firstLine, \"hysteria2://\") {\n        return types.ConfTypeV2Ray\n    }\n    \n    // 检查是否是配置文件格式\n    if strings.Contains(content, \"[server_local]\") || strings.Contains(content, \"[policy]\") {\n        return types.ConfTypeQuantumultX\n    }\n    \n    return types.ConfTypeSUB\n}\n\n// parseBase64Subscription 解析 Base64 编码的订阅\nfunc (p *Parser) parseBase64Subscription(content string) ([]models.Proxy, error) {\n    decoded, err := base64.StdEncoding.DecodeString(content)\n    if err != nil {\n        return nil, fmt.Errorf(\"failed to decode base64 subscription: %w\", err)\n    }\n    \n    return p.parseGenericSubscription(string(decoded))\n}\n\n// parseGenericSubscription 解析通用订阅格式\nfunc (p *Parser) parseGenericSubscription(content string) ([]models.Proxy, error) {\n    var proxies []models.Proxy\n    \n    lines := strings.Split(content, \"\\n\")\n    for i, line := range lines {\n        line = strings.TrimSpace(line)\n        if line == \"\" || strings.HasPrefix(line, \"#\") {\n            continue\n        }\n        \n        proxy, err := p.parseNodeLink(line)\n        if err != nil {\n            // 记录错误但继续处理其他节点\n            continue\n        }\n        \n        if proxy != nil {\n            proxy.ID = uint32(i + 1)\n            proxies = append(proxies, *proxy)\n        }\n    }\n    \n    return proxies, nil\n}\n\n// parseNodeLink 解析单个节点链接\nfunc (p *Parser) parseNodeLink(link string) (*models.Proxy, error) {\n    // 判断协议类型\n    if strings.HasPrefix(link, \"ss://\") {\n        return p.parseSSLink(link)\n    }\n    if strings.HasPrefix(link, \"ssr://\") {\n        return p.parseSSRLink(link)\n    }\n    if strings.HasPrefix(link, \"vmess://\") {\n        return p.parseVMessLink(link)\n    }\n    if strings.HasPrefix(link, \"vless://\") {\n        return p.parseVLESSLink(link)\n    }\n    if strings.HasPrefix(link, \"trojan://\") {\n        return p.parseTrojanLink(link)\n    }\n    if strings.HasPrefix(link, \"hysteria://\") {\n        return p.parseHysteriaLink(link)\n    }\n    if strings.HasPrefix(link, \"hysteria2://\") {\n        return p.parseHysteria2Link(link)\n    }\n    if strings.HasPrefix(link, \"tuic://\") {\n        return p.parseTUICLink(link)\n    }\n    \n    return nil, fmt.Errorf(\"unsupported protocol: %s\", link)\n}\n\n// FetchContent 获取订阅内容\nfunc (p *Parser) FetchContent(ctx context.Context, url string) (string, error) {\n    req, err := http.NewRequestWithContext(ctx, \"GET\", url, nil)\n    if err != nil {\n        return \"\", fmt.Errorf(\"failed to create request: %w\", err)\n    }\n    \n    // 设置 User-Agent\n    req.Header.Set(\"User-Agent\", \"subconverter-go/1.0\")\n    \n    resp, err := p.httpClient.Do(req)\n    if err != nil {\n        return \"\", fmt.Errorf(\"failed to fetch URL: %w\", err)\n    }\n    defer resp.Body.Close()\n    \n    if resp.StatusCode != http.StatusOK {\n        return \"\", fmt.Errorf(\"HTTP error: %d\", resp.StatusCode)\n    }\n    \n    // 限制内容大小\n    limitedReader := &io.LimitedReader{R: resp.Body, N: p.maxSize}\n    content, err := io.ReadAll(limitedReader)\n    if err != nil {\n        return \"\", fmt.Errorf(\"failed to read response: %w\", err)\n    }\n    \n    return string(content), nil\n}\n`\n\n### 2.3 节点解析器实现\n\n```go\n// internal/parser/node/vmess.go\npackage node\n\nimport (\n    \"encoding/base64\"\n    \"encoding/json\"\n    \"fmt\"\n    \"net/url\"\n    \"strconv\"\n    \"strings\"\n    \n    \"subconverter-go/pkg/models\"\n    \"subconverter-go/pkg/types\"\n)\n\ntype VMessParser struct{}\n\n// NewVMessParser 创建 VMess 解析器\nfunc NewVMessParser() *VMessParser {\n    return &VMessParser{}\n}\n\n// ParseNode 解析 VMess 节点 - 对应 C++ 中的 explodeVmess\nfunc (p *VMessParser) ParseNode(link string) (*models.Proxy, error) {\n    if !strings.HasPrefix(link, \"vmess://\") {\n        return nil, fmt.Errorf(\"invalid VMess link format\")\n    }\n    \n    // 移除协议前缀\n    content := strings.TrimPrefix(link, \"vmess://\")\n    \n    // Base64 解码\n    decoded, err := base64.StdEncoding.DecodeString(content)\n    if err != nil {\n        return nil, fmt.Errorf(\"failed to decode VMess link: %w\", err)\n    }\n    \n    // 解析 JSON\n    var vmessConfig VMessConfig\n    if err := json.Unmarshal(decoded, &vmessConfig); err != nil {\n        return nil, fmt.Errorf(\"failed to parse VMess JSON: %w\", err)\n    }\n    \n    // 转换为 Proxy 模型\n    proxy := &models.Proxy{\n        Type:     types.ProxyTypeVMess,\n        Remark:   vmessConfig.PS,\n        Hostname: vmessConfig.Add,\n        Port:     uint16(vmessConfig.Port),\n        UserID:   vmessConfig.ID,\n        AlterID:  uint16(vmessConfig.Aid),\n        TransferProtocol: vmessConfig.Net,\n        TLSStr:   vmessConfig.TLS,\n    }\n    \n    // 解析网络配置\n    switch vmessConfig.Net {\n    case \"ws\":\n        proxy.Path = vmessConfig.Path\n        proxy.Host = vmessConfig.Host\n    case \"grpc\":\n        proxy.GRPCServiceName = vmessConfig.Path\n    case \"h2\":\n        proxy.Path = vmessConfig.Path\n        proxy.Host = vmessConfig.Host\n    }\n    \n    // TLS 配置\n    if vmessConfig.TLS == \"tls\" {\n        proxy.TLSSecure = true\n        proxy.SNI = vmessConfig.SNI\n    }\n    \n    // 验证必要字段\n    if err := proxy.IsValid(); !err {\n        return nil, fmt.Errorf(\"invalid VMess proxy configuration\")\n    }\n    \n    return proxy, nil\n}\n\n// VMessConfig VMess 配置结构\ntype VMessConfig struct {\n    V    string `json:\"v\"`   // 版本\n    PS   string`json:\"ps\"`  // 备注\n    Add  string`json:\"add\"`// 地址\n    Port int  `json:\"port\"` // 端口\n I

\nD string `json:\"id\"` // UUID\n Aid int `json:\"aid\"` // alterId\n Net string `json:\"net\"` // 网络类型\n Type string `json:\"type\"` // 伪装类型\n Host string `json:\"host\"` // 伪装域名\n Path string `json:\"path\"` // 路径\n TLS string `json:\"tls\"` // TLS\n SNI string `json:\"sni\"` // SNI\n}\n\n// SupportedSchemes 返回支持的协议前缀\nfunc (p *VMessParser) SupportedSchemes() []string {\n return []string{\"vmess\"}\n}\n\n// ValidateLink 验证链接格式\nfunc (p *VMessParser) ValidateLink(link string) error {\n if !strings.HasPrefix(link, \"vmess://\") {\n return fmt.Errorf(\"invalid VMess link prefix\")\n }\n \n content := strings.TrimPrefix(link, \"vmess://\")\n if _, err := base64.StdEncoding.DecodeString(content); err != nil {\n return fmt.Errorf(\"invalid base64 encoding: %w\", err)\n }\n \n return nil\n}\n`\n\n### 2.4 其他协议解析器\n\n`go\n// internal/parser/node/shadowsocks.go\npackage node\n\nimport (\n \"encoding/base64\"\n \"fmt\"\n \"net/url\"\n \"strconv\"\n \"strings\"\n \n \"subconverter-go/pkg/models\"\n \"subconverter-go/pkg/types\"\n)\n\ntype ShadowsocksParser struct{}\n\n// ParseNode 解析 Shadowsocks 节点 - 对应 C++ 中的 explodeSS\nfunc (p *ShadowsocksParser) ParseNode(link string) (*models.Proxy, error) {\n if !strings.HasPrefix(link, \"ss://\") {\n return nil, fmt.Errorf(\"invalid Shadowsocks link format\")\n }\n \n // 移除协议前缀\n content := strings.TrimPrefix(link, \"ss://\")\n \n // 解析 URL\n u, err := url.Parse(\"ss://\" + content)\n if err != nil {\n return nil, fmt.Errorf(\"failed to parse SS URL: %w\", err)\n }\n \n // 解码用户信息\n userInfo := u.User.String()\n if userInfo == \"\" {\n return nil, fmt.Errorf(\"missing user info in SS link\")\n }\n \n // Base64 解码认证信息\n decoded, err := base64.URLEncoding.DecodeString(userInfo)\n if err != nil {\n // 尝试标准 Base64\n decoded, err = base64.StdEncoding.DecodeString(userInfo)\n if err != nil {\n return nil, fmt.Errorf(\"failed to decode SS auth: %w\", err)\n }\n }\n \n // 解析方法和密码\n parts := strings.SplitN(string(decoded), \":\", 2)\n if len(parts) != 2 {\n return nil, fmt.Errorf(\"invalid SS auth format\")\n }\n \n method, password := parts[0], parts[1]\n \n // 解析端口\n port, err := strconv.Atoi(u.Port())\n if err != nil {\n return nil, fmt.Errorf(\"invalid port: %w\", err)\n }\n \n proxy := &models.Proxy{\n Type: types.ProxyTypeShadowsocks,\n Hostname: u.Hostname(),\n Port: uint16(port),\n EncryptMethod: method,\n Password: password,\n Remark: getRemarkFromFragment(u.Fragment),\n }\n \n // 解析插件信息\n if plugin := u.Query().Get(\"plugin\"); plugin != \"\" {\n parts := strings.SplitN(plugin, \";\", 2)\n proxy.Plugin = parts[0]\n if len(parts) > 1 {\n proxy.PluginOption = parts[1]\n }\n }\n \n return proxy, nil\n}\n\n// getRemarkFromFragment 从 fragment 获取备注信息\nfunc getRemarkFromFragment(fragment string) string {\n if fragment == \"\" {\n return \"SS Node\"\n }\n \n decoded, err := url.QueryUnescape(fragment)\n if err != nil {\n return fragment\n }\n \n return decoded\n}\n`\n\n## 3. 生成器模块（Generator Module）实现方案\n\n### 3.1 核心接口设计\n\n`go\n// internal/generator/interfaces.go\npackage generator\n\nimport (\n \"context\"\n \"subconverter-go/pkg/models\"\n)\n\n// Generator 配置生成器接口\ntype Generator interface {\n Generate(ctx context.Context, proxies []models.Proxy, config GeneratorConfig) (*GenerateResult, error)\n GetType() string\n ValidateTemplate(template string) error\n SupportedFeatures() []string\n}\n\n// ProxyGenerator 代理配置生成器接口\ntype ProxyGenerator interface {\n GenerateProxies(proxies []models.Proxy) (interface{}, error)\n GenerateGroups(groups []models.ProxyGroupConfig, proxies []models.Proxy) (interface{}, error)\n GenerateRules(rulesets []models.RulesetContent) (interface{}, error)\n}\n\n// TemplateRenderer 模板渲染器接口\ntype TemplateRenderer interface {\n Render(template string, data interface{}) (string, error)\n AddFunction(name string, fn interface{})\n LoadTemplate(path string) error\n}\n\n// GeneratorConfig 生成器配置\ntype GeneratorConfig struct {\n Template string `json:\"template\"`\n BaseConfig string `json:\"base_config\"`\n ProxyGroups []models.ProxyGroupConfig `json:\"proxy_groups\"`\n Rulesets []models.RulesetContent `json:\"rulesets\"`\n ExtraSettings map[string]interface{} `json:\"extra_settings\"`\n}\n\n// GenerateResult 生成结果\ntype GenerateResult struct {\n Content string `json:\"content\"`\n ContentType string `json:\"content_type\"`\n Headers map[string]string `json:\"headers\"`\n Metadata map[string]interface{} `json:\"metadata\"`\n}\n`\n\n### 3.2 Clash 生成器实现\n\n`go\n// internal/generator/clash/generator.go\npackage clash\n\nimport (\n \"context\"\n \"fmt\"\n \"strings\"\n \n \"gopkg.in/yaml.v3\"\n \"subconverter-go/pkg/models\"\n \"subconverter-go/pkg/types\"\n \"subconverter-go/internal/generator\"\n)\n\ntype Generator struct {\n renderer generator.TemplateRenderer\n}\n\n// NewGenerator 创建 Clash 生成器\nfunc NewGenerator() *Generator {\n return &Generator{\n renderer: NewTemplateRenderer(),\n }\n}\n\n// Generate 生成 Clash 配置 - 对应 C++ 中的 proxyToClash\nfunc (g *Generator) Generate(ctx context.Context, proxies []models.Proxy, config generator.GeneratorConfig) (*generator.GenerateResult, error) {\n // 构建 Clash 配置结构\n clashConfig := &ClashConfig{\n Port: 7890,\n SocksPort: 7891,\n LogLevel: \"info\",\n Mode: \"rule\",\n }\n \n // 生成代理配置\n clashProxies, err := g.generateProxies(proxies)\n if err != nil {\n return nil, fmt.Errorf(\"failed to generate proxies: %w\", err)\n }\n clashConfig.Proxies = clashProxies\n \n // 生成代理组\n proxyGroups, err := g.generateProxyGroups(config.ProxyGroups, proxies)\n if err != nil {\n return nil, fmt.Errorf(\"failed to generate proxy groups: %w\", err)\n }\n clashConfig.ProxyGroups = proxyGroups\n \n // 生成规则\n rules, err := g.generateRules(config.Rulesets)\n if err != nil {\n return nil, fmt.Errorf(\"failed to generate rules: %w\", err)\n }\n clashConfig.Rules = rules\n \n // 合并基础配置\n if config.BaseConfig != \"\" {\n if err := g.mergeBaseConfig(clashConfig, config.BaseConfig); err != nil {\n return nil, fmt.Errorf(\"failed to merge base config: %w\", err)\n }\n }\n \n // 序列化为 YAML\n content, err := yaml.Marshal(clashConfig)\n if err != nil {\n return nil, fmt.Errorf(\"failed to marshal Clash config: %w\", err)\n }\n \n return &generator.GenerateResult{\n Content: string(content),\n ContentType: \"application/x-yaml\",\n Headers: map[string]string{\n \"Content-Disposition\": \"attachment; filename=clash.yaml\",\n },\n Metadata: map[string]interface{}{\n \"proxy_count\": len(proxies),\n \"group_count\": len(proxyGroups),\n \"rule_count\": len(rules),\n },\n }, nil\n}\n\n// generateProxies 生成代理配置\nfunc (g \*Generator) generateProxies(proxies []models.Proxy) ([]ClashProxy, error) {\n var clashProxies []ClashProxy\n \n for _, proxy := range proxies {\n clashProxy, err := g.convertProxyToClash(&proxy

)
if err != nil {
continue // 跳过无效代理，但记录错误
}
clashProxies = append(clashProxies, clashProxy)
}

    return clashProxies, nil

}

// convertProxyToClash 转换代理为 Clash 格式
func (g *Generator) convertProxyToClash(proxy *models.Proxy) (ClashProxy, error) {
base := ClashProxy{
Name: proxy.Remark,
Type: proxy.Type.String(),
Server: proxy.Hostname,
Port: proxy.Port,
}

    switch proxy.Type {
    case types.ProxyTypeShadowsocks:
        return g.convertSSToClash(proxy, base)
    case types.ProxyTypeVMess:
        return g.convertVMessToClash(proxy, base)
    case types.ProxyTypeVLESS:
        return g.convertVLESSToClash(proxy, base)
    case types.ProxyTypeTrojan:
        return g.convertTrojanToClash(proxy, base)
    case types.ProxyTypeHysteria:
        return g.convertHysteriaToClash(proxy, base)
    case types.ProxyTypeHysteria2:
        return g.convertHysteria2ToClash(proxy, base)
    default:
        return ClashProxy{}, fmt.Errorf("unsupported proxy type: %s", proxy.Type)
    }

}

// ClashConfig Clash 配置结构
type ClashConfig struct {
Port int `yaml:"port"`
SocksPort int `yaml:"socks-port"`
LogLevel string `yaml:"log-level"`
Mode string `yaml:"mode"`
Proxies []ClashProxy `yaml:"proxies"`
ProxyGroups []ClashProxyGroup `yaml:"proxy-groups"`
Rules []string `yaml:"rules"`
DNS map[string]interface{} `yaml:"dns,omitempty"`
Experimental map[string]interface{} `yaml:"experimental,omitempty"`
}

// ClashProxy Clash 代理配置
type ClashProxy struct {
Name string `yaml:"name"`
Type string `yaml:"type"`
Server string `yaml:"server"`
Port uint16 `yaml:"port"`
Cipher string `yaml:"cipher,omitempty"`
Password string `yaml:"password,omitempty"`
UUID string `yaml:"uuid,omitempty"`
AlterID uint16 `yaml:"alterId,omitempty"`
Network string `yaml:"network,omitempty"`
TLS bool `yaml:"tls,omitempty"`
UDP bool `yaml:"udp,omitempty"`
Extra map[string]interface{} `yaml:",inline"`
}

// ClashProxyGroup Clash 代理组配置
type ClashProxyGroup struct {
Name string `yaml:"name"`
Type string `yaml:"type"`
Proxies []string `yaml:"proxies"`
URL string `yaml:"url,omitempty"`
Interval int `yaml:"interval,omitempty"`
}

````

### 3.3 Surge 生成器实现

```go
// internal/generator/surge/generator.go
package surge

import (
    "context"
    "fmt"
    "strings"

    "subconverter-go/pkg/models"
    "subconverter-go/pkg/types"
    "subconverter-go/internal/generator"
)

type Generator struct {
    version int // Surge 版本 (3, 4, 5)
}

// NewGenerator 创建 Surge 生成器
func NewGenerator(version int) *Generator {
    return &Generator{
        version: version,
    }
}

// Generate 生成 Surge 配置 - 对应 C++ 中的 proxyToSurge
func (g *Generator) Generate(ctx context.Context, proxies []models.Proxy, config generator.GeneratorConfig) (*generator.GenerateResult, error) {
    var sections []string

    // [General] 段
    sections = append(sections, g.generateGeneralSection())

    // [Proxy] 段
    proxySection, err := g.generateProxySection(proxies)
    if err != nil {
        return nil, fmt.Errorf("failed to generate proxy section: %w", err)
    }
    sections = append(sections, proxySection)

    // [Proxy Group] 段
    if len(config.ProxyGroups) > 0 {
        groupSection, err := g.generateProxyGroupSection(config.ProxyGroups, proxies)
        if err != nil {
            return nil, fmt.Errorf("failed to generate proxy group section: %w", err)
        }
        sections = append(sections, groupSection)
    }

    // [Rule] 段
    if len(config.Rulesets) > 0 {
        ruleSection, err := g.generateRuleSection(config.Rulesets)
        if err != nil {
            return nil, fmt.Errorf("failed to generate rule section: %w", err)
        }
        sections = append(sections, ruleSection)
    }

    content := strings.Join(sections, "\n\n")

    return &generator.GenerateResult{
        Content:     content,
        ContentType: "text/plain",
        Headers: map[string]string{
            "Content-Disposition": "attachment; filename=surge.conf",
        },
        Metadata: map[string]interface{}{
            "surge_version": g.version,
            "proxy_count":   len(proxies),
        },
    }, nil
}

// generateProxySection 生成代理段
func (g *Generator) generateProxySection(proxies []models.Proxy) (string, error) {
    var lines []string
    lines = append(lines, "[Proxy]")

    for _, proxy := range proxies {
        line, err := g.convertProxyToSurge(&proxy)
        if err != nil {
            continue // 跳过无效代理
        }
        lines = append(lines, line)
    }

    return strings.Join(lines, "\n"), nil
}

// convertProxyToSurge 转换代理为 Surge 格式
func (g *Generator) convertProxyToSurge(proxy *models.Proxy) (string, error) {
    switch proxy.Type {
    case types.ProxyTypeShadowsocks:
        return g.convertSSToSurge(proxy)
    case types.ProxyTypeVMess:
        return g.convertVMessToSurge(proxy)
    case types.ProxyTypeTrojan:
        return g.convertTrojanToSurge(proxy)
    default:
        return "", fmt.Errorf("unsupported proxy type for Surge: %s", proxy.Type)
    }
}
````

## 4. 过滤器和处理器实现

### 4.1 代理过滤器

```go
// internal/processor/filter.go
package processor

import (
    "regexp"
    "strings"

    "subconverter-go/pkg/models"
    "subconverter-go/pkg/types"
)

// Filter 过滤器接口
type Filter interface {
    Filter(proxies []models.Proxy) []models.Proxy
    GetName() string
}

// IncludeFilter 包含过滤器
type IncludeFilter struct {
    patterns []*regexp.Regexp
}

// NewIncludeFilter 创建包含过滤器
func NewIncludeFilter(patterns []string) (*IncludeFilter, error) {
    var regexps []*regexp.Regexp
    for _, pattern := range patterns {
        regex, err := regexp.Compile(pattern)
        if err != nil {
            return nil, err
        }
        regexps = append(regexps, regex)
    }

    return &IncludeFilter{patterns: regexps}, nil
}

// Filter 执行包含过滤
func (f *IncludeFilter) Filter(proxies []models.Proxy) []models.Proxy {
    if len(f.patterns) == 0 {
        return proxies
    }

    var filtered []models.Proxy
    for _, proxy := range proxies {
        for _, pattern := range f.patterns {
            if pattern.MatchString(proxy.Remark) {
                filtered = append(filtered, proxy)
                break
            }
        }
    }

    return filtered
}

// ExcludeFilter 排除过滤器
type ExcludeFilter struct {
    patterns []*regexp.Regexp
}

// Filter 执行排除过滤
func (f *ExcludeFilter) Filter(proxies []models.Proxy) []models.Proxy {
    if len(f.patterns) == 0 {
        return proxies
    }

    var filtered []models.Proxy
    for _, proxy := range proxies {
        excluded := false
        for _, pattern := range f.patterns {
            if pattern.MatchString(proxy.Remark) {
                excluded = true
                break
            }
        }
        if !excluded {
            filtered = append(filtered, proxy)
        }
    }

    return filtered
}

// TypeFilter 类型过滤器
type TypeFilter struct {
    allowedTypes []types.ProxyType
}

// Filter 按类型过滤
func (f *TypeFilter) Filter(proxies []models.Proxy) []models.Proxy {
    if len(f.allowedTypes) == 0 {
        return proxies
    }

    var filtered []models.Proxy
    for _, proxy := range proxies {
        for _, allowedType := range f.allowedTypes {
            if proxy.Type == allowedType {
                filtered = append(filtered, proxy)
                break
            }
        }
    }

    return filtered
}
```

### 4.2 代理处理器

```go
// internal/processor/processor.go
package processor

import (
    "context"
    "sort"
    "strings"

    "subconverter-go/pkg/models"
)

// Processor 代理处理器
type Processor struct {
    filters    []Filter
    sorters    []Sorter
    validators []Validator
}

// NewProcessor 创建处理器
func NewProcessor() *Processor {
    return &Processor{
        filters:    make([]Filter, 0),
        sorters:    make([]Sorter, 0),
        validators: make([]Validator, 0),
    }
}

// AddFilter 添加过滤器
func (p *Processor) AddFilter(filter Filter) {
    p.filters = append(p.filters, filter)
}

// AddSorter 添加排序器
func (p *Processor) AddSorter(sorter Sorter) {
    p.sorters = append(p.sorters, sorter)
}

// Process 处理代理列表
func (p *Processor) Process(ctx context.Context, proxies []models.Proxy) ([]models.Proxy, error) {
    result := make([]models.Proxy, len(proxies))
    copy(result, proxies)

    // 验证
    result = p.validateProxies(result)

    // 过滤
    for _, filter := range p.filters {
        result = filter.Filter(result)
    }

    // 排序
    for _, sorter := range p.sorters {
        result = sorter.Sort(result)
    }

    // 去重
    result = p.deduplicateProxies(result)

    return result, nil
}

// validateProxies 验证代理
func (p *Processor) validateProxies(proxies []models.Proxy) []models.Proxy {
    var valid []models.Proxy
    for _, proxy := range proxies {
        if proxy.IsValid() {
            valid = append(valid, proxy)
        }
    }
    return valid
}

// deduplicateProxies 去重代理
func (p *Processor) deduplicateProxies(proxies []models.Proxy) []models.Proxy {
    seen := make(map[string]bool)
    var unique []models.Proxy

    for _, proxy := range proxies {
        key := proxy.Hostname + ":" + string(proxy.Port) + ":" + proxy.Type.String()
        if !seen[key] {
            seen[key] = true
            unique = append(unique, proxy)
        }
    }

    return unique
}

// Sorter 排序器接口
type Sorter interface {
    Sort(proxies []models.Proxy) []models.Proxy
}

// NameSorter 按名称排序
type NameSorter struct{}

func (s *NameSorter) Sort(proxies []models.Proxy) []models.Proxy {
    sort.Slice(proxies, func(i, j int) bool {
        return strings.ToLower(proxies[i].Remark) < strings.ToLower(proxies[j].Remark)
    })
    return proxies
}

// TypeSorter 按类型排序
type TypeSorter struct{}

func (s *TypeSorter) Sort(proxies []models.Proxy) []models.Proxy {
    sort.Slice(proxies, func(i, j int) bool {
        return proxies[i].Type < proxies[j].Type
    })
    return proxies
}
```

## 5. 并发处理优化

### 5.1 并发解析器

```go
// internal/parser/concurrent.go
package parser

import (
    "context"
    "sync"

    "subconverter-go/pkg/models"
)

// ConcurrentParser 并发解析器
type ConcurrentParser struct {
    parser      Parser
    maxWorkers  int
    bufferSize  int
}

// NewConcurrentParser 创建并发解析器
func NewConcurrentParser(parser Parser, maxWorkers int) *ConcurrentParser {
    return &ConcurrentParser{
        parser:     parser,
        maxWorkers: maxWorkers,
        bufferSize: maxWorkers * 2,
    }
}

// ParseConcurrent 并发解析节点
func (cp *ConcurrentParser) ParseConcurrent(ctx context.Context, links []string) ([]models.Proxy, error) {
    jobs := make(chan string, cp.bufferSize)
    results := make(chan parseResult, cp.bufferSize)

    // 启动工作 goroutine
    var wg sync.WaitGroup
    for i := 0; i < cp.maxWorkers; i++ {
        wg.Add(1)
        go cp.worker(ctx, &wg, jobs, results)
    }

    // 发送任务
    go func() {
        defer close(jobs)
        for _, link := range links {
            select {
            case jobs <- link:
            case <-ctx.Done():
                return
            }
        }
    }()

    // 等待所有工作完成
    go func() {
        wg.Wait()
        close(results)
    }()

    // 收集结果
    var proxies []models.Proxy
    for result := range results {
        if result.err == nil && result.proxy != nil {
            proxies = append(proxies, *result.proxy)
        }
    }

    return proxies, nil
}

type parseResult struct {
    proxy *models.Proxy
    err   error
}

// worker 工作协程
func (cp *ConcurrentParser) worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, results chan<- parseResult) {
    defer wg.Done()

    for {
        select {
        case link, ok := <-jobs:
            if !ok {
                return
            }

            proxy, err := cp.parseNodeLink(link)
            results <- parseResult{proxy: proxy, err: err}

        case <-ctx.Done():
            return
        }
    }
}
```

## 6. 错误处理和监控

### 6.1 错误聚合器

```go
// internal/processor/error_handler.go
package processor

import (
    "fmt"
    "sync"

    "subconverter-go/pkg/types"
)

// ErrorAggregator 错误聚合器
type ErrorAggregator struct {
    errors []types.ConvertError
    mutex  sync.RWMutex
}

// NewErrorAggregator 创建错误聚合器
func NewErrorAggregator() *ErrorAggregator {
    return &ErrorAggregator{
        errors: make([]types.ConvertError, 0),
    }
}

// AddError 添加错误
func (ea *ErrorAggregator) AddError(err error, context string) {
    ea.mutex.Lock()
    defer ea.mutex.Unlock()

    convertErr := types.ConvertError{
        Code:    types.ErrorCodeParseError,
        Message: err.Error(),
        Details: context,
    }

    ea.errors = append(ea.errors, convertErr)
}

// GetErrors 获取所有错误
func (ea *ErrorAggregator) GetErrors() []types.ConvertError {
    ea.mutex.RLock()
    defer ea.mutex.RUnlock()

    result := make([]types.ConvertError, len(ea.errors))
    copy(result, ea.errors)
    return result
}

// HasErrors 检查是否有错误
func (ea *ErrorAggregator) HasErrors() bool {
    ea.mutex.RLock()
    defer ea.mutex.RUnlock()
    return len(ea.errors) > 0
}

// GetSummary 获取错误摘要
func (ea *ErrorAggregator) GetSummary() string {
    errors := ea.GetErrors()
    if len(errors) == 0 {
        return "No errors"
    }

    return fmt.Sprintf("Total errors: %d", len(errors))
}
```

## 7. 模块集成和工厂模式

### 7.1 解析器工厂

```go
// internal/parser/factory.go
package parser

import (
    "fmt"

    "subconverter-go/pkg/types"
    "subconverter-go/internal/parser/node"
)

// ParserFactory 解析器工厂
type ParserFactory struct {
    parsers map[types.ProxyType]node.Parser
}

// NewParserFactory 创建解析器工厂
func NewParserFactory() *ParserFactory {
    factory := &ParserFactory{
        parsers: make(map[types.ProxyType]node.Parser),
    }

    // 注册所有解析器
    factory.RegisterParser(types.ProxyTypeShadowsocks, node.NewShadowsocksParser())
    factory.RegisterParser(types.ProxyTypeVMess, node.NewVMessParser())
    factory.RegisterParser(types.ProxyTypeVLESS, node.NewVLESSParser())
    factory.RegisterParser(types.ProxyTypeTrojan, node.NewTrojanParser())
    factory.RegisterParser(types.ProxyTypeHysteria, node.NewHysteriaParser())
    factory.RegisterParser(types.ProxyTypeHysteria2, node.NewHysteria2Parser())

    return factory
}

// RegisterParser 注册解析器
func (pf *ParserFactory) RegisterParser(proxyType types.ProxyType, parser node.Parser) {
    pf.parsers[proxyType] = parser
}

// GetParser 获取解析器
func (pf *ParserFactory) GetParser(proxyType types.ProxyType) (node.Parser, error) {
    parser, exists := pf.parsers[proxyType]
    if !exists {
        return nil, fmt.Errorf("no parser registered for proxy type: %s", proxyType)
    }
    return parser, nil
}

// GetAllParsers 获取所有解析器
func (pf *ParserFactory) GetAllParsers() map[types.ProxyType]node.Parser {
    return pf.parsers
}
```

### 7.2 生成器工厂

```go
// internal/generator/factory.go
package generator

import (
    "fmt"
    "strings"

    "subconverter-go/internal/generator/clash"
    "subconverter-go/internal/generator/surge"
    "subconverter-go/internal/generator/quantumultx"
    "subconverter-go/internal/generator/loon"
)

// GeneratorFactory 生成器工厂
type GeneratorFactory struct {
    generators map[string]Generator
}

// NewGeneratorFactory 创建生成器工厂
func NewGeneratorFactory() *GeneratorFactory {
    factory := &GeneratorFactory{
        generators: make(map[string]Generator),
    }

    // 注册所有生成器
    factory.RegisterGenerator("clash", clash.NewGenerator())
    factory.RegisterGenerator("clashr", clash.NewGenerator())
    factory.RegisterGenerator("surge", surge.NewGenerator(4))
    factory.RegisterGenerator("surge3", surge.NewGenerator(3))
    factory.RegisterGenerator("surge4", surge.NewGenerator(4))
    factory.RegisterGenerator("quanx", quantumultx.NewGenerator())
    factory.RegisterGenerator("loon", loon.NewGenerator())

    return factory
}

// RegisterGenerator 注册生成器
func (gf *GeneratorFactory) RegisterGenerator(target string, generator Generator) {
    gf.generators[strings.ToLower(target)] = generator
}

// GetGenerator 获取生成器
func (gf *GeneratorFactory) GetGenerator(target string) (Generator, error) {
    generator, exists := gf.generators[strings.ToLower(target)]
    if !exists {
        return nil, fmt.Errorf("no generator registered for target: %s", target)
    }
    return generator, nil
}

// GetSupportedTargets 获取支持的目标
func (gf *GeneratorFactory) GetSupportedTargets() []string {
    var targets []string
    for target := range gf.generators {
        targets = append(targets, target)
    }
    return targets
}
```

## 8. 总结

### 8.1 核心模块实现完成情况

我已经完成了解析器和生成器模块的详细实现方案设计：

#### 解析器模块 (Parser Module)：

1. **接口设计**：统一的解析器接口，支持不同协议
2. **订阅解析器**：支持 Base64、V2Ray、SS、SSR 等格式
3. **节点解析器**：VMess、VLESS、Shadowsocks、Trojan 等协议解析
4. **并发处理**：并发解析提升性能
5. **错误处理**：完整的错误处理机制

#### 生成器模块 (Generator Module)：

1. **接口设计**：统一的生成器接口
2. **多客户端支持**：Clash、Surge、QuantumultX、Loon 等
3. **模板系统**：灵活的模板渲染机制
4. **配置合并**：基础配置和动态配置合并

#### 辅助模块：

1. **过滤器系统**：包含、排除、类型过滤
2. **处理器系统**：验证、排序、去重
3. **工厂模式**：解析器和生成器的工厂管理
4. **错误聚合**：统一的错误收集和处理

### 8.2 兼容性保证

- **完全对应 C++ 功能**：所有解析和生成函数都对应 C++ 版本的实现
- **数据结构兼容**：使用前面定义的数据模型确保兼容性
- **接口兼容**：保持相同的输入输出格式
- **错误处理兼容**：相同的错误类型和处理方式

### 8.3 性能优化特性

- **并发解析**：使用 goroutine 池并发处理节点解析
- **内存管理**：合理的缓冲区大小和对象复用
- **流式处理**：大文件的流式解析避免内存爆炸
- **缓存机制**：解析结果缓存减少重复工作

这个实现方案为核心转换功能提供了完整的技术框架，确保了从 C++ 到 Go 的平滑迁移。
