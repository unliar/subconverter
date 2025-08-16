\n# SubConverter Go 版本第三方依赖选择方案\n\n## 1. 核心依赖分类\n\n### 1.1 Web 框架层\n- **主框架**: [Gin](https://github.com/gin-gonic/gin)\n - 版本: v1.9+\n - 理由: 高性能、轻量级、生态丰富、中间件支持完善\n - 用途: API 路由、中间件管理、HTTP 请求处理\n\n### 1.2 配置管理\n- **配置解析**: [Viper](https://github.com/spf13/viper)\n - 版本: v1.16+\n - 理由: 支持多格式配置文件、环境变量、远程配置\n - 用途: 应用配置、规则配置、模板配置管理\n\n### 1.3 模板引擎\n- **模板渲染**: [Pongo2](https://github.com/flosch/pongo2)\n - 版本: v6.0+\n - 理由: Django/Jinja2 兼容语法，与现有模板完全兼容\n - 备选: [Jet](https://github.com/CloudyKit/jet) (更高性能但需要模板修改)\n - 用途: 配置文件模板渲染、客户端配置生成\n\n## 2. 数据处理依赖\n\n### 2.1 序列化/反序列化\n- **JSON 处理**: 标准库 `encoding/json` + [jsoniter](https://github.com/json-iterator/go)\n - 版本: jsoniter v1.1.12+\n - 理由: 性能优异，API 兼容标准库\n - 用途: API 数据交换、配置解析\n\n- **YAML 处理**: [go-yaml](https://github.com/go-yaml/yaml)\n - 版本: v3.0+\n - 理由: 标准 YAML 库，功能完整\n - 用途: 配置文件解析、规则文件处理\n\n- **TOML 处理**: [go-toml](https://github.com/pelletier/go-toml)\n - 版本: v2.0+\n - 理由: 高性能、完整 TOML 支持\n - 用途: 配置文件格式支持\n\n### 2.2 字符串处理\n- **正则表达式**: 标准库 `regexp`\n- **字符串工具**: [go-humanize](https://github.com/dustin/go-humanize)\n - 版本: v1.0+\n - 理由: 人性化字符串格式化\n - 用途: 日志输出、错误信息格式化\n\n## 3. 网络和加密依赖\n\n### 3.1 HTTP 客户端\n- **HTTP 客户端**: [Resty](https://github.com/go-resty/resty)\n - 版本: v2.7+\n - 理由: 简单易用、功能丰富、支持重试和中间件\n - 用途: 订阅 URL 获取、远程规则下载\n\n### 3.2 加密和哈希\n- **加密库**: 标准库 `crypto/*` + [x/crypto](https://golang.org/x/crypto)\n - 理由: 官方扩展库，安全可靠\n - 用途: 协议加密、密钥生成、证书处理\n\n- **UUID 生成**: [google/uuid](https://github.com/google/uuid)\n - 版本: v1.3+\n - 理由: 标准 UUID 实现\n - 用途: 节点 ID 生成、会话管理\n\n### 3.3 网络工具\n- **IP 地址处理**: [netip](https://golang.org/x/net/netip) (Go 1.18+)\n - 理由: 高性能 IP 地址解析和处理\n - 用途: 代理节点 IP 验证、CIDR 处理\n\n- **域名解析**: 标准库 `net` + [miekg/dns](https://github.com/miekg/dns)\n - 版本: dns v1.1.50+\n - 理由: 功能完整的 DNS 库\n - 用途: 域名解析、DNS 查询\n\n## 4. 存储和缓存\n\n### 4.1 内存缓存\n- **本地缓存**: [go-cache](https://github.com/patrickmn/go-cache)\n - 版本: v2.1+\n - 理由: 轻量级、支持 TTL、线程安全\n - 用途: 解析结果缓存、配置缓存\n\n- **LRU 缓存**: [groupcache/lru](https://github.com/golang/groupcache/tree/master/lru)\n - 理由: 高效 LRU 实现\n - 用途: 大容量数据缓存\n\n### 4.2 持久化存储 (可选)\n- **SQLite**: [modernc.org/sqlite](https://gitlab.com/cznic/sqlite)\n - 版本: v1.20+\n - 理由: 纯 Go 实现、无 CGO 依赖\n - 用途: 本地数据存储、统计信息\n\n## 5. 日志和监控\n\n### 5.1 日志系统\n- **结构化日志**: [logrus](https://github.com/sirupsen/logrus)\n - 版本: v1.9+\n - 理由: 功能丰富、插件生态完善、JSON 格式支持\n - 备选: [zap](https://github.com/uber-go/zap) (更高性能)\n - 用途: 应用日志、错误跟踪、调试信息\n\n### 5.2 性能监控\n- **指标收集**: [prometheus/client_golang](https://github.com/prometheus/client_golang)\n - 版本: v1.14+\n - 理由: 标准指标收集库\n - 用途: 性能监控、API 统计\n\n## 6. 测试和开发工具\n\n### 6.1 测试框架\n- **单元测试**: 标准库 `testing` + [testify](https://github.com/stretchr/testify)\n - 版本: testify v1.8+\n - 理由: 丰富的断言和 Mock 功能\n - 用途: 单元测试、集成测试\n\n- **HTTP 测试**: [httptest](https://golang.org/pkg/net/http/httptest/) (标准库)\n - 理由: 官方 HTTP 测试工具\n - 用途: API 端点测试\n\n### 6.2 代码质量\n- **代码检查**: [golangci-lint](https://github.com/golangci/golangci-lint)\n - 版本: v1.50+\n - 理由: 集成多种 linter、配置灵活\n - 用途: 代码质量检查、CI/CD 集成\n\n## 7. 构建和部署\n\n### 7.1 构建工具\n- **构建管理**: [Taskfile](https://taskfile.dev/)\n - 版本: v3.20+\n - 理由: 现代化的 Make 替代品、跨平台\n - 用途: 构建脚本、开发任务自动化\n\n### 7.2 容器化\n- **Docker**: 多阶段构建\n - 基础镜像: `golang:1.20-alpine`

- 运行镜像: `alpine:latest`
- 用途: 生产环境部署、开发环境统一

## 8. 专用协议依赖

### 8.1 代理协议支持

- **V2Ray 协议**: [v2ray-core](https://github.com/v2fly/v2ray-core) (作为参考)

  - 注意: 仅作为协议规范参考，不直接依赖
  - 用途: VMess/VLESS 协议实现参考

- **Shadowsocks**: 自实现
  - 理由: 协议相对简单，减少外部依赖
  - 用途: SS/SSR 协议支持

### 8.2 工具库

- **Base64 编码**: 标准库 `encoding/base64`
- **URL 处理**: 标准库 `net/url`
- **时间处理**: 标准库 `time`

## 9. 完整依赖列表

### 9.1 go.mod 主要依赖

```go
module subconverter-go

go 1.20

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/spf13/viper v1.16.0
    github.com/flosch/pongo2/v6 v6.0.0
    github.com/json-iterator/go v1.1.12
    gopkg.in/yaml.v3 v3.0.1
    github.com/pelletier/go-toml/v2 v2.0.8
    github.com/go-resty/resty/v2 v2.7.0
    github.com/google/uuid v1.3.0
    github.com/patrickmn/go-cache v2.1.0+incompatible
    github.com/sirupsen/logrus v1.9.3
    github.com/prometheus/client_golang v1.16.0
    github.com/stretchr/testify v1.8.4
    github.com/miekg/dns v1.1.55
    github.com/dustin/go-humanize v1.0.1
    golang.org/x/crypto v0.10.0
    golang.org/x/net v0.10.0
)
```

### 9.2 开发依赖

```yaml
# 开发工具
- golangci-lint: v1.53.3
- taskfile: v3.28.0
- air: v1.44.0 (热重载开发)
- gotests: v1.6.0 (测试生成)
- mockery: v2.32.0 (Mock生成)
```

## 10. 依赖管理策略

### 10.1 版本管理

- **语义化版本**: 严格遵循 SemVer 规范
- **依赖锁定**: 使用 go.sum 锁定依赖版本
- **定期更新**: 每月检查依赖更新
- **安全扫描**: 使用 govulncheck 检查安全漏洞

### 10.2 依赖最小化原则

- **标准库优先**: 能用标准库的不使用第三方库
- **功能聚合**: 选择功能完整的库，避免功能重复
- **社区活跃度**: 优选维护活跃、社区支持好的库
- **性能考虑**: 关键路径优选高性能库

### 10.3 依赖替换方案

- **模板引擎**: Pongo2 → Jet (性能优化)
- **日志系统**: Logrus → Zap (性能优化)
- **JSON 处理**: jsoniter → sonic (极致性能)
- **HTTP 框架**: Gin → Fiber (内存优化)

## 11. 兼容性考虑

### 11.1 Go 版本兼容

- **最低版本**: Go 1.19 (泛型支持)
- **推荐版本**: Go 1.20+ (性能优化)
- **兼容策略**: 支持最新 3 个大版本

### 11.2 平台兼容

- **操作系统**: Linux, Windows, macOS
- **架构支持**: amd64, arm64
- **容器支持**: Docker, Podman
- **云平台**: 支持主流云服务商部署

## 12. 性能优化依赖

### 12.1 内存优化

- **对象池**: sync.Pool (标准库)
- **内存分析**: [pprof](https://golang.org/pkg/net/http/pprof/) (标准库)
- **内存监控**: prometheus 指标

### 12.2 并发优化

- **工作池**: [ants](https://github.com/panjf2000/ants)
  - 版本: v2.8+
  - 理由: 高性能 goroutine 池
  - 用途: 控制并发数量、减少 goroutine 创建开销

### 12.3 IO 优化

- **缓冲 IO**: 标准库 bufio
- **零拷贝**: [fasthttp](https://github.com/valyala/fasthttp) (可选)
  - 用途: 极高性能 HTTP 处理

## 13. 安全依赖

### 13.1 输入验证

- **数据验证**: [validator](https://github.com/go-playground/validator)
  - 版本: v10.14+
  - 理由: 功能完整的数据验证库
  - 用途: API 参数验证、配置验证

### 13.2 安全工具

- **安全扫描**: [gosec](https://github.com/securecodewarrior/gosec)
- **依赖检查**: [nancy](https://github.com/sonatypecommunity/nancy)
- **许可证检查**: [fossa](https://fossa.com/)

## 14. 总结

这个依赖选择方案具有以下特点:

### 14.1 核心优势

- **高性能**: 选择性能优异的库
- **稳定性**: 优选维护活跃的成熟库
- **兼容性**: 确保与现有 C++版本 100%兼容
- **可维护性**: 依赖关系清晰，升级路径明确

### 14.2 风险控制

- **依赖最小化**: 减少外部依赖数量
- **版本锁定**: 避免意外更新导致的问题
- **替换方案**: 为关键依赖准备替换选项
- **安全监控**: 持续监控依赖安全性

### 14.3 迁移友好

- **API 兼容**: 确保与现有接口完全兼容
- **模板兼容**: Pongo2 提供 Jinja2 兼容性
- **配置兼容**: 支持现有配置文件格式
- **部署兼容**: 支持现有部署方式

这个依赖选择方案为 SubConverter 的 Go 版本迁移提供了坚实的技术基础，确保了高性能、高稳定性和完全兼容性。
