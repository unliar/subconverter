# SubConverter Go 版本项目目录结构设计

## 1. 项目根目录结构

```
subconverter-go/
├── cmd/                          # 应用程序入口
│   └── subconverter/
│       └── main.go              # 主程序入口
├── internal/                     # 内部包（不对外暴露）
│   ├── api/                     # HTTP API 层
│   ├── service/                 # 业务逻辑层
│   ├── parser/                  # 解析器模块
│   ├── generator/               # 生成器模块
│   ├── config/                  # 配置管理
│   ├── cache/                   # 缓存实现
│   └── utils/                   # 工具函数
├── pkg/                         # 可对外暴露的包
│   ├── models/                  # 数据模型
│   ├── types/                   # 类型定义
│   └── constants/               # 常量定义
├── configs/                     # 配置文件
│   ├── config.yaml              # 默认配置
│   ├── config.example.yaml      # 配置示例
│   └── docker.yaml              # Docker 配置
├── templates/                   # 模板文件（兼容原有模板）
│   ├── clash/                   # Clash 模板
│   ├── surge/                   # Surge 模板
│   ├── quantumultx/             # QuantumultX 模板
│   ├── loon/                    # Loon 模板
│   └── singbox/                 # SingBox 模板
├── assets/                      # 静态资源文件
│   ├── rules/                   # 规则文件（兼容原有规则）
│   ├── profiles/                # 配置档案
│   └── snippets/                # 代码片段
├── test/                        # 测试文件
│   ├── integration/             # 集成测试
│   ├── compatibility/           # 兼容性测试
│   └── testdata/                # 测试数据
├── scripts/                     # 构建和部署脚本
│   ├── build.sh                # 构建脚本
│   ├── docker-build.sh         # Docker 构建
│   └── deploy.sh               # 部署脚本
├── docs/                        # 项目文档
│   ├── api.md                  # API 文档
│   ├── deployment.md           # 部署文档
│   └── migration.md            # 迁移指南
├── .github/                     # GitHub 配置
│   └── workflows/               # CI/CD 工作流
├── docker/                      # Docker 相关文件
│   ├── Dockerfile              # 容器构建文件
│   └── docker-compose.yml      # 容器编排
├── go.mod                       # Go 模块定义
├── go.sum                       # 依赖校验文件
├── Makefile                     # 构建规则
├── README.md                    # 项目说明
└── LICENSE                      # 许可证文件
```

## 2. 详细目录说明

### 2.1 cmd/ - 应用程序入口

```
cmd/
└── subconverter/
    └── main.go                  # 主程序：命令行参数解析、服务启动
```

**职责**：

- 解析命令行参数
- 初始化配置
- 启动 HTTP 服务器
- 优雅关闭处理

### 2.2 internal/ - 内部核心模块

#### 2.2.1 api/ - HTTP API 层

```
internal/api/
├── handlers/                    # HTTP 处理器
│   ├── converter.go            # 转换接口处理器
│   ├── ruleset.go              # 规则集接口处理器
│   ├── config.go               # 配置接口处理器
│   ├── profile.go              # 配置档案处理器
│   └── template.go             # 模板渲染处理器
├── middleware/                  # 中间件
│   ├── auth.go                 # 认证中间件
│   ├── ratelimit.go            # 限流中间件
│   ├── cors.go                 # CORS 处理
│   └── logger.go               # 日志中间件
├── routes/                      # 路由定义
│   └── routes.go               # 路由配置
└── server.go                    # HTTP 服务器配置
```

#### 2.2.2 service/ - 业务逻辑层

```
internal/service/
├── converter/                   # 转换服务
│   ├── service.go              # 转换服务主逻辑
│   ├── filter.go               # 节点过滤逻辑
│   └── validator.go            # 请求验证逻辑
├── ruleset/                     # 规则集服务
│   ├── service.go              # 规则集服务
│   ├── cache.go                # 规则集缓存
│   └── fetcher.go              # 规则集获取器
├── config/                      # 配置服务
│   ├── service.go              # 配置服务
│   ├── loader.go               # 配置加载器
│   └── watcher.go              # 配置监控器
└── profile/                     # 配置档案服务
    ├── service.go              # 档案服务
    └── manager.go              # 档案管理器
```

#### 2.2.3 parser/ - 解析器模块

```
internal/parser/
├── subscription/                # 订阅解析器
│   ├── parser.go               # 订阅解析主逻辑
│   ├── detector.go             # 格式检测器
│   └── fetcher.go              # 订阅获取器
├── node/                        # 节点解析器
│   ├── vmess.go                # VMess 解析器
│   ├── vless.go                # VLESS 解析器
│   ├── shadowsocks.go          # Shadowsocks 解析器
│   ├── shadowsocksr.go         # ShadowsocksR 解析器
│   ├── trojan.go               # Trojan 解析器
│   ├── hysteria.go             # Hysteria 解析器
│   ├── hysteria2.go            # Hysteria2 解析器
│   ├── tuic.go                 # TUIC 解析器
│   ├── snell.go                # Snell 解析器
│   └── socks.go                # SOCKS 解析器
├── config/                      # 配置解析器
│   ├── yaml.go                 # YAML 解析器
│   ├── toml.go                 # TOML 解析器
│   └── ini.go                  # INI 解析器
└── base.go                      # 解析器基类和接口
```

#### 2.2.4 generator/ - 生成器模块

```
internal/generator/
├── clash/                       # Clash 生成器
│   ├── generator.go            # Clash 配置生成器
│   ├── proxy.go                # 代理配置生成
│   ├── rules.go                # 规则配置生成
│   └── groups.go               # 代理组配置生成
├── surge/                       # Surge 生成器
│   ├── generator.go            # Surge 配置生成器
│   ├── proxy.go                # 代理配置生成
│   └── rules.go                # 规则配置生成
├── quantumultx/                 # QuantumultX 生成器
│   ├── generator.go            # QuanX 配置生成器
│   └── proxy.go                # 代理配置生成
├── loon/                        # Loon 生成器
│   ├── generator.go            # Loon 配置生成器
│   └── proxy.go                # 代理配置生成
├── singbox/                     # SingBox 生成器
│   ├── generator.go            # SingBox 配置生成器
│   └── proxy.go                # 代理配置生成
├── template/                    # 模板引擎
│   ├── engine.go               # 模板引擎接口
│   ├── functions.go            # 模板函数库
│   └── compat.go               # Jinja2 兼容层
└── base.go                      # 生成器基类和接口
```

#### 2.2.5 config/ - 配置管理

```
internal/config/
├── config.go                    # 配置结构定义
├── loader.go                    # 配置加载器
├── validator.go                 # 配置验证器
├── env.go                       # 环境变量处理
└── defaults.go                  # 默认配置
```

#### 2.2.6 cache/ - 缓存实现

```
internal/cache/
├── cache.go                     # 缓存接口定义
├── memory.go                    # 内存缓存实现
├── redis.go                     # Redis 缓存实现
└── manager.go                   # 缓存管理器
```

#### 2.2.7 utils/ - 工具函数

```
internal/utils/
├── network/                     # 网络工具
│   ├── http.go                 # HTTP 客户端
│   ├── dns.go                  # DNS 解析
│   └── proxy.go                # 代理测试
├── crypto/                      # 加密工具
│   ├── base64.go               # Base64 编解码
│   ├── md5.go                  # MD5 哈希
│   └── uuid.go                 # UUID 生成
├── string/                      # 字符串工具
│   ├── convert.go              # 字符串转换
│   ├── validate.go             # 字符串验证
│   └── format.go               # 字符串格式化
├── file/                        # 文件工具
│   ├── io.go                   # 文件读写
│   ├── path.go                 # 路径处理
│   └── watch.go                # 文件监控
└── log/                         # 日志工具
    ├── logger.go               # 日志接口
    └── logrus.go               # Logrus 实现
```

### 2.3 pkg/ - 对外暴露包

#### 2.3.1 models/ - 数据模型

```
pkg/models/
├── proxy.go                     # 代理节点模型
├── ruleset.go                   # 规则集模型
├── config.go                    # 配置模型
├── group.go                     # 代理组模型
└── request.go                   # 请求响应模型
```

#### 2.3.2 types/ - 类型定义

```
pkg/types/
├── proxy.go                     # 代理类型枚举
├── generator.go                 # 生成器类型
├── errors.go                    # 错误类型
└── constants.go                 # 常量定义
```

### 2.4 configs/ - 配置文件

```
configs/
├── config.yaml                  # 默认配置文件
├── config.example.yaml          # 配置示例文件
├── docker.yaml                  # Docker 环境配置
├── test.yaml                    # 测试环境配置
└── prod.yaml                    # 生产环境配置
```

### 2.5 templates/ - 模板文件（完全兼容原有结构）

```
templates/
├── clash/
│   ├── base.yaml               # Clash 基础模板
│   ├── config.yaml             # Clash 配置模板
│   └── providers.yaml          # Clash 提供商模板
├── surge/
│   ├── base.conf               # Surge 基础模板
│   └── config.conf             # Surge 配置模板
├── quantumultx/
│   ├── base.conf               # QuanX 基础模板
│   └── config.conf             # QuanX 配置模板
├── loon/
│   ├── base.conf               # Loon 基础模板
│   └── config.conf             # Loon 配置模板
└── singbox/
    ├── base.json               # SingBox 基础模板
    └── config.json             # SingBox 配置模板
```

### 2.6 assets/ - 静态资源（完全兼容原有结构）

```
assets/
├── rules/                       # 规则文件目录
│   ├── ACL4SSR/                # ACL4SSR 规则集
│   ├── DivineEngine/           # DivineEngine 规则集
│   ├── lhie1/                  # lhie1 规则集
│   └── NobyDa/                 # NobyDa 规则集
├── profiles/                    # 配置档案目录
│   └── example_profile.ini     # 示例配置档案
└── snippets/                    # 代码片段目录
    ├── emoji.toml              # Emoji 配置
    ├── groups.toml             # 代理组配置
    ├── rename_node.toml        # 节点重命名配置
    └── rulesets.toml           # 规则集配置
```

### 2.7 test/ - 测试文件

```
test/
├── integration/                 # 集成测试
│   ├── api_test.go             # API 集成测试
│   ├── converter_test.go       # 转换器集成测试
│   └── compatibility_test.go   # 兼容性集成测试
├── compatibility/               # 兼容性测试
│   ├── cpp_comparison_test.go  # C++ 版本对比测试
│   ├── template_test.go        # 模板兼容性测试
│   └── config_test.go          # 配置兼容性测试
├── testdata/                    # 测试数据
│   ├── subscriptions/          # 测试订阅
│   ├── configs/                # 测试配置
│   ├── templates/              # 测试模板
│   └── expected/               # 期望输出
└── benchmark/                   # 性能测试
    ├── parser_bench_test.go    # 解析器性能测试
    └── generator_bench_test.go # 生成器性能测试
```

## 3. 关键设计原则

### 3.1 兼容性优先

- **目录映射**：`templates/` 和 `assets/` 目录完全对应原有 C++ 版本的 `base/` 目录结构
- **配置兼容**：支持原有的 YAML、TOML、INI 配置格式
- **模板兼容**：直接使用原有模板文件，无需修改

### 3.2 清晰的分层架构

- **cmd/**：应用程序入口，负责启动和命令行处理
- **internal/api/**：HTTP API 层，处理 Web 请求
- **internal/service/**：业务逻辑层，核心功能实现
- **internal/parser/ & internal/generator/**：核心处理层
- **pkg/**：对外暴露的数据模型和类型

### 3.3 可维护性设计

- **模块化**：每个功能模块独立，便于开发和测试
- **接口驱动**：通过接口实现松耦合
- **测试覆盖**：完整的单元测试、集成测试和兼容性测试

### 3.4 部署友好

- **容器化支持**：完整的 Docker 配置
- **配置外化**：支持环境变量和配置文件
- **静态资源分离**：便于更新和维护

## 4. 文件组织策略

### 4.1 包命名规范

- 使用小写字母和下划线
- 包名简洁明了，体现功能
- 避免与标准库冲突

### 4.2 文件命名规范

- 使用小写字母和下划线
- 文件名体现具体功能
- 测试文件以 `_test.go` 结尾

### 4.3 导入路径

```go
// 内部包导入
import (
    "subconverter-go/internal/service/converter"
    "subconverter-go/internal/parser/node"
    "subconverter-go/pkg/models"
)

// 外部依赖导入
import (
    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
    "gopkg.in/yaml.v3"
)
```

## 5. 构建和部署

### 5.1 Makefile 主要目标

```makefile
# 构建二进制文件
build:
	go build -o bin/subconverter ./cmd/subconverter

# 运行测试
test:
	go test ./...

# 运行集成测试
test-integration:
	go test ./test/integration/...

# 构建 Docker 镜像
docker-build:
	docker build -t subconverter-go .

# 代码检查
lint:
	golangci-lint run

# 清理构建文件
clean:
	rm -rf bin/
```

### 5.2 Docker 多阶段构建

```dockerfile
# 构建阶段
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o subconverter ./cmd/subconverter

# 运行阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/subconverter .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/configs ./configs
EXPOSE 25500
CMD ["./subconverter"]
```

这个目录结构设计确保了：

1. **完全兼容**：与现有 C++ 版本的文件结构保持一致
2. **清晰分层**：职责明确，便于开发和维护
3. **易于测试**：完整的测试策略支持
4. **部署友好**：支持容器化和传统部署方式
5. **扩展性强**：便于添加新的代理协议和客户端支持
