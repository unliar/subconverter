# SubConverter Go 版本测试策略和持续集成方案

## 1. 测试策略概述

### 1.1 测试目标

- **功能正确性**: 确保所有核心功能与 C++版本完全一致
- **性能达标**: 满足或超越现有性能指标
- **稳定性保证**: 长期运行稳定，异常处理完善
- **兼容性验证**: 全面的向后兼容性测试

### 1.2 测试层次

```
测试金字塔
├── E2E测试 (5%)
│   ├── 完整订阅转换流程
│   ├── 多客户端兼容性
│   └── 性能基准测试
├── 集成测试 (25%)
│   ├── API端点测试
│   ├── 模板渲染测试
│   └── 配置加载测试
└── 单元测试 (70%)
    ├── 解析器测试
    ├── 生成器测试
    ├── 规则引擎测试
    └── 工具函数测试
```

## 2. 单元测试策略

### 2.1 测试框架选择

- **主框架**: 标准库 `testing` + `testify/suite`
- **断言库**: `testify/assert` + `testify/require`
- **Mock 工具**: `testify/mock` + `mockery`
- **测试数据**: `testify/suite` + `gofakeit`

### 2.2 核心模块测试

#### 2.2.1 解析器模块测试

```go
// pkg/parser/ss_test.go
package parser

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type SSParserTestSuite struct {
    suite.Suite
    parser *SSParser
}

func (s *SSParserTestSuite) SetupTest() {
    s.parser = NewSSParser()
}

func (s *SSParserTestSuite) TestParseValidURL() {
    tests := []struct {
        name     string
        input    string
        expected *models.Proxy
        wantErr  bool
    }{
        {
            name:  "标准SS链接",
            input: "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ@example.com:8388#测试节点",
            expected: &models.Proxy{
                Type:     "ss",
                Name:     "测试节点",
                Server:   "example.com",
                Port:     8388,
                Cipher:   "aes-256-gcm",
                Password: "password",
            },
            wantErr: false,
        },
        {
            name:    "无效SS链接",
            input:   "invalid://url",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        s.Run(tt.name, func() {
            result, err := s.parser.Parse(tt.input)

            if tt.wantErr {
                s.Error(err)
                return
            }

            s.NoError(err)
            s.Equal(tt.expected.Type, result.Type)
            s.Equal(tt.expected.Name, result.Name)
            s.Equal(tt.expected.Server, result.Server)
            s.Equal(tt.expected.Port, result.Port)
        })
    }
}

func (s *SSParserTestSuite) TestParseBatch() {
    urls := []string{
        "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ@example1.com:8388#节点1",
        "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ@example2.com:8389#节点2",
    }

    results, err := s.parser.ParseBatch(urls)
    s.NoError(err)
    s.Len(results, 2)
    s.Equal("节点1", results[0].Name)
    s.Equal("节点2", results[1].Name)
}

func TestSSParserTestSuite(t *testing.T) {
    suite.Run(t, new(SSParserTestSuite))
}
```

#### 2.2.2 生成器模块测试

```go
// pkg/generator/clash_test.go
package generator

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/subconverter-go/pkg/models"
)

func TestClashGenerator_Generate(t *testing.T) {
    generator := NewClashGenerator()

    proxies := []models.Proxy{
        {
            Type:     "ss",
            Name:     "测试节点",
            Server:   "example.com",
            Port:     8388,
            Cipher:   "aes-256-gcm",
            Password: "password",
        },
    }

    config := &models.ClientConfig{
        Proxies: proxies,
        Groups: []models.Group{
            {
                Name:    "自动选择",
                Type:    "url-test",
                Proxies: []string{"测试节点"},
                URL:     "http://www.gstatic.com/generate_204",
            },
        },
    }

    result, err := generator.Generate(config)
    assert.NoError(t, err)
    assert.Contains(t, result, "测试节点")
    assert.Contains(t, result, "自动选择")
    assert.Contains(t, result, "aes-256-gcm")
}

func TestClashGenerator_Performance(t *testing.T) {
    generator := NewClashGenerator()

    // 生成大量测试数据
    proxies := make([]models.Proxy, 1000)
    for i := 0; i < 1000; i++ {
        proxies[i] = models.Proxy{
            Type:     "ss",
            Name:     fmt.Sprintf("节点%d", i),
            Server:   fmt.Sprintf("example%d.com", i),
            Port:     8388 + i,
            Cipher:   "aes-256-gcm",
            Password: "password",
        }
    }

    config := &models.ClientConfig{Proxies: proxies}

    start := time.Now()
    _, err := generator.Generate(config)
    duration := time.Since(start)

    assert.NoError(t, err)
    assert.Less(t, duration, time.Second, "生成1000个节点应在1秒内完成")
}
```

#### 2.2.3 规则引擎测试

```go
// pkg/rules/engine_test.go
package rules

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/subconverter-go/pkg/config"
    "github.com/subconverter-go/pkg/models"
)

func TestEngine_ApplyFilters(t *testing.T) {
    rulesConfig := &config.RulesConfig{
        NodeFilters: []config.NodeFilter{
            {
                Name:     "包含香港",
                Type:     "include",
                Patterns: []string{"香港", "HK"},
                Enabled:  true,
            },
            {
                Name:     "排除测试",
                Type:     "exclude",
                Patterns: []string{"测试", "test"},
                Enabled:  true,
            },
        },
    }

    engine, err := NewEngine(rulesConfig)
    assert.NoError(t, err)

    proxies := []models.Proxy{
        {Name: "香港节点1"},
        {Name: "香港节点2"},
        {Name: "美国节点"},
        {Name: "测试节点"},
        {Name: "HK Premium"},
    }

    filtered := engine.ApplyFilters(proxies)

    assert.Len(t, filtered, 3) // 香港节点1, 香港节点2, HK Premium

    names := make([]string, len(filtered))
    for i, proxy := range filtered {
        names[i] = proxy.Name
    }

    assert.Contains(t, names, "香港节点1")
    assert.Contains(t, names, "香港节点2")
    assert.Contains(t, names, "HK Premium")
    assert.NotContains(t, names, "美国节点")
    assert.NotContains(t, names, "测试节点")
}
```

### 2.3 测试覆盖率要求

- **总体覆盖率**: ≥ 85%
- **核心模块**: ≥ 90% (parser, generator, rules)
- **工具函数**: ≥ 95%
- **API 层**: ≥ 80%

## 3. 集成测试策略

### 3.1 API 集成测试

```go
// tests/integration/api_test.go
package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type APITestSuite struct {
    suite.Suite
    router *gin.Engine
    server *httptest.Server
}

func (s *APITestSuite) SetupSuite() {
    gin.SetMode(gin.TestMode)
    s.router = setupTestRouter()
    s.server = httptest.NewServer(s.router)
}

func (s *APITestSuite) TearDownSuite() {
    s.server.Close()
}

func (s *APITestSuite) TestConvertSubscription() {
    // 测试订阅转换API
    requestBody := map[string]interface{}{
        "url":    "https://example.com/subscription",
        "target": "clash",
        "config": "https://example.com/config.ini",
    }

    jsonBody, _ := json.Marshal(requestBody)

    resp, err := http.Post(
        s.server.URL+"/sub",
        "application/json",
        bytes.NewBuffer(jsonBody),
    )

    s.NoError(err)
    defer resp.Body.Close()

    s.Equal(http.StatusOK, resp.Status)
    s.Equal("application/yaml", resp.Header.Get("Content-Type"))

    // 验证响应内容
    var responseBody bytes.Buffer
    _, err = responseBody.ReadFrom(resp.Body)
    s.NoError(err)

    content := responseBody.String()
    s.Contains(content, "proxies:")
    s.Contains(content, "proxy-groups:")
    s.Contains(content, "rules:")
}

func (s *APITestSuite) TestConvertWithInvalidURL() {
    requestBody := map[string]interface{}{
        "url":    "invalid-url",
        "target": "clash",
    }

    jsonBody, _ := json.Marshal(requestBody)

    resp, err := http.Post(
        s.server.URL+"/sub",
        "application/json",
        bytes.NewBuffer(jsonBody),
    )

    s.NoError(err)
    defer resp.Body.Close()

    s.Equal(http.StatusBadRequest, resp.Status)
}

func TestAPITestSuite(t *testing.T) {
    suite.Run(t, new(APITestSuite))
}
```

### 3.2 配置加载集成测试

```go
// tests/integration/config_test.go
package integration

import (
    "context"
    "io/ioutil"
    "os"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/subconverter-go/pkg/config"
)

func TestConfigLoader_Integration(t *testing.T) {
    // 创建临时配置目录
    tempDir, err := ioutil.TempDir("", "config_test")
    assert.NoError(t, err)
    defer os.RemoveAll(tempDir)

    // 写入测试配置文件
    appConfig := `
server:
  host: "127.0.0.1"
  port: 25500
log:
  level: "info"
  format: "json"
cache:
  enable: true
  default_ttl: "10m"
`

    err = ioutil.WriteFile(filepath.Join(tempDir, "app.yaml"), []byte(appConfig), 0644)
    assert.NoError(t, err)

    rulesConfig := `
node_filters:
  - name: "test_filter"
    type: "include"
    patterns: ["test"]
    enabled: true
`

    err = ioutil.WriteFile(filepath.Join(tempDir, "rules.yaml"), []byte(rulesConfig), 0644)
    assert.NoError(t, err)

    // 测试配置加载
    loader := config.NewLoader(tempDir)
    cfg, err := loader.LoadConfig(context.Background())

    assert.NoError(t, err)
    assert.NotNil(t, cfg)
    assert.Equal(t, "127.0.0.1", cfg.App.Server.Host)
    assert.Equal(t, 25500, cfg.App.Server.Port)
    assert.Equal(t, "info", cfg.App.Log.Level)
    assert.Len(t, cfg.Rules.NodeFilters, 1)
}
```

## 4. 端到端测试

### 4.1 完整流程测试

```go
// tests/e2e/subscription_test.go
package e2e

import (
    "context"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
    suite.Suite
    app            *App
    subscriptionServer *httptest.Server
}

func (s *E2ETestSuite) SetupSuite() {
    // 启动模拟订阅服务器
    s.subscriptionServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        subscription := `
ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ@example1.com:8388#节点1
ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ@example2.com:8389#节点2
vmess://eyJ2IjoiMiIsInBzIjoi6IqC54K5MyIsImFkZCI6ImV4YW1wbGUzLmNvbSIsInBvcnQiOiI0NDMiLCJ0eXBlIjoibm9uZSIsImlkIjoiMTIzNDU2NzgtMTIzNC0xMjM0LTEyMzQtMTIzNDU2Nzg5YWJjIiwiYWlkIjoiMCIsIm5ldCI6IndzIiwicGF0aCI6Ii8iLCJob3N0IjoiZXhhbXBsZTMuY29tIiwidGxzIjoidGxzIn0=
`
        w.Header().Set("Content-Type", "text/plain")
        w.Write([]byte(subscription))
    }))

    // 启动应用
    s.app = NewTestApp()
    go s.app.Start()

    // 等待应用启动
    time.Sleep(100 * time.Millisecond)
}

func (s *E2ETestSuite) TearDownSuite() {
    s.subscriptionServer.Close()
    s.app.Stop()
}

func (s *E2ETestSuite) TestCompleteConversionFlow() {
    // 测试完整的订阅转换流程
    testCases := []struct {
        target   string
        expected []string
    }{
        {
            target: "clash",
            expected: []string{
                "proxies:",
                "proxy-groups:",
                "rules:",
                "节点1",
                "节点2",
                "节点3",
            },
        },
        {
            target: "surge",
            expected: []string{
                "[Proxy]",
                "[Proxy Group]",
                "[Rule]",
                "节点1 = ss",
                "节点2 = ss",
                "节点3 = vmess",
            },
        },
    }

    for _, tc := range testCases {
        s.Run(fmt.Sprintf("Convert to %s", tc.target), func() {
            url := fmt.Sprintf("http://localhost:25500/sub?target=%s&url=%s",
                tc.target, s.subscriptionServer.URL)

            resp, err := http.Get(url)
            s.NoError(err)
            defer resp.Body.Close()

            s.Equal(http.StatusOK, resp.Status)

            body, err := io.ReadAll(resp.Body)
            s.NoError(err)

            content := string(body)
            for _, expected := range tc.expected {
                s.Contains(content, expected)
            }
        })
    }
}

func (s *E2ETestSuite) TestPerformanceBenchmark() {
    // 性能基准测试
    url := fmt.Sprintf("http://localhost:25500/sub?target=clash&url=%s",
        s.subscriptionServer.URL)

    // 预热
    for i := 0; i < 5; i++ {
        resp, _ := http.Get(url)
        resp.Body.Close()
    }

    // 基准测试
    start := time.Now()
    for i := 0; i < 100; i++ {
        resp, err := http.Get(url)
        s.NoError(err)
        resp.Body.Close()
    }
    duration := time.Since(start)

    avgResponseTime := duration / 100
    s.Less(avgResponseTime, 100*time.Millisecond, "平均响应时间应小于100ms")
}

func TestE2ETestSuite(t *testing.T) {
    suite.Run(t, new(E2ETestSuite))
}
```

## 5. 性能测试

### 5.1 基准测试

```go
// tests/benchmark/parser_benchmark_test.go
package benchmark

import (
    "testing"
    "github.com/subconverter-go/pkg/parser"
)

func BenchmarkSSParser_Parse(b *testing.B) {
    p := parser.NewSSParser()
    url := "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ@example.com:8388#测试节点"

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := p.Parse(url)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkClashGenerator_Generate(b *testing.B) {
    generator := generator.NewClashGenerator()

    // 准备测试数据
    proxies := make([]models.Proxy, 100)
    for i := 0; i < 100; i++ {
        proxies[i] = models.Proxy{
            Type:     "ss",
            Name:     fmt.Sprintf("节点%d", i),
            Server:   "example.com",
            Port:     8388,
            Cipher:   "aes-256-gcm",
            Password: "password",
        }
    }

    config := &models.ClientConfig{Proxies: proxies}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := generator.Generate(config)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkMemoryUsage(b *testing.B) {
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        // 执行内存密集型操作
        data := make([]byte, 1024*1024) // 1MB
        _ = data
    }
}
```

### 5.2 压力测试

```go
// tests/stress/concurrent_test.go
package stress

import (
    "context"
    "net/http"
    "sync"
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
)

func TestConcurrentRequests(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过压力测试")
    }

    const (
        numWorkers = 50
        numRequests = 1000
    )

    url := "http://localhost:25500/sub?target=clash&url=https://example.com/sub"

    var wg sync.WaitGroup
    requests := make(chan struct{}, numRequests)

    // 启动工作协程
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            client := &http.Client{Timeout: 10 * time.Second}

            for range requests {
                resp, err := client.Get(url)
                if err != nil {
                    t.Errorf("请求失败: %v", err)
                    continue
                }
                resp.Body.Close()

                if resp.Status != http.StatusOK {
                    t.Errorf("期望状态码200，实际: %d", resp.Status)
                }
            }
        }()
    }

    // 发送请求
    start := time.Now()
    for i := 0; i < numRequests; i++ {
        requests <- struct{}{}
    }
    close(requests)

    wg.Wait()
    duration := time.Since(start)

    qps := float64(numRequests) / duration.Seconds()
    t.Logf("完成 %d 个请求，耗时 %v，QPS: %.2f", numRequests, duration, qps)

    assert.Greater(t, qps, 100.0, "QPS应大于100")
}
```

## 6. 兼容性测试

### 6.1 向后兼容性测试

```go
// tests/compatibility/api_compatibility_test.go
package compatibility

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestAPICompatibility(t *testing.T) {
    // 测试与C++版本API的兼容性
    testCases := []struct {
        name        string
        endpoint    string
        params      map[string]string
        expectJSON  bool
    }{
        {
            name:     "基本订阅转换",
            endpoint: "/sub",
            params: map[string]string{
                "target": "clash",
                "url":    "https://example.com/sub",
            },
        },
        {
            name:     "短链接生成",
            endpoint: "/short",
            params: map[string]string{
                "data": "target=clash&url=https://example.com/sub",
            },
            expectJSON: true,
        },
        {
            name:     "版本信息",
            endpoint: "/version",
            params:   map[string]string{},
            expectJSON: true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // 构建请求URL
            url := buildTestURL(tc.endpoint, tc.params)

            // 发送请求
            resp, err := http.Get(url)
            assert.NoError(t, err)
            defer resp.Body.Close()

            // 验证响应
            assert.Equal(t, http.StatusOK, resp.Status)

            if tc.expectJSON {
                assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
            }
        })
    }
}
```

### 6.2 客户端兼容性测试

```go
// tests/compatibility/client_compatibility_test.go
package compatibility

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestClientConfigCompatibility(t *testing.T) {
    // 测试生成的配置文件与各客户端的兼容性
    clients := []string{"clash", "surge", "quantumult-x", "loon", "singbox"}

    for _, client := range clients {
        t.Run(client, func(t *testing.T) {
            config, err := generateTestConfig(client)
            assert.NoError(t, err)

            // 验证配置格式
            err = validateClientConfig(client, config)
            assert.NoError(t, err, "生成的%s配置应该有效", client)
        })
    }
}
```

## 7. 持续集成配置

### 7.1 GitHub Actions 配置

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.19, 1.20, 1.21]

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: ${{ matrix.go-version }}

      - name: Cache dependencies
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Install dependencies
        run: go mod download

      - name: Run linter
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
          args: --timeout 5m

      - name: Run unit tests
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Run integration tests
        run: go test -v -tags=integration ./tests/integration/...

      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out

      - name: Build binary
        run: go build -v ./cmd/subconverter

  e2e:
    runs-on: ubuntu-latest
    needs: test
    if: github.event_name == 'push'

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.20

      - name: Build application
        run: go build -o subconverter ./cmd/subconverter

      - name: Start application
        run: |
          ./subconverter &
          sleep 5

      - name: Run E2E tests
        run: go test -v -tags=e2e ./tests/e2e/...

      - name: Run performance tests
        run: go test -v -bench=. -benchmem ./tests/benchmark/...

  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Run Gosec Security Scanner
        uses: securecodewarrior/github-action-gosec@master
        with:
          args: "./..."

      - name: Run Nancy vulnerability check
        run: |
          go list -json -m all | nancy sleuth
```

### 7.2 质量门禁配置

```yaml
# .github/workflows/quality-gate.yml
name: Quality Gate

on:
  pull_request:
    branches: [main]

jobs:
  quality-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.20

      - name: Run tests with coverage
        run: |
          go test -v -race -coverprofile=coverage.out ./...
          go tool cover -func=coverage.out

      - name: Check coverage threshold
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "总覆盖率: ${COVERAGE}%"
          if (( $(echo "$COVERAGE < 85" | bc -l) )); then
            echo "错误: 代码覆盖率 ${COVERAGE}% 低于要求的 85%"
            exit 1
          fi

      - name: Check cyclomatic complexity
        run: |
          go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
          gocyclo -over 10 .

      - name: Check code duplication
        run: |
          go install github.com/mibk/dupl@latest
          dupl -threshold 100 .

      - name: Performance regression check
        run: |
          go test -bench=. -benchmem ./tests/benchmark/... > current.bench
          # 与基准性能对比 (需要存储历史基准数据)
```

## 8. 测试数据管理

### 8.1 测试数据生成

```go
// tests/testdata/generator.go
package testdata

import (
    "fmt"
    "github.com/brianvoe/gofakeit/v6"
    "github.com/subconverter-go/pkg/models"
)

// GenerateProxies 生成测试代理数据
func GenerateProxies(count int) []models.Proxy {
    proxies := make([]models.Proxy, count)

    for i := 0; i < count; i++ {
        proxies[i] = models.Proxy{
            Type:     randomProxyType(),
            Name:     fmt.Sprintf("测试节点%d", i+1),
            Server:   gofakeit.IPv4Address(),
            Port:     gofakeit.Number(1000, 65535),
            Cipher:   randomCipher(),
            Password: gofakeit.Password(true, true, true, false, false, 16),
        }
    }

    return proxies
}

// GenerateSubscription 生成测试订阅内容
func GenerateSubscription(count int) string {
    proxies := GenerateProxies(count)
    var lines []string

    for _, proxy := range proxies {
        switch proxy.Type {
        case "ss":
            lines = append(lines, generateSSURL(proxy))
        case "vmess":
            lines = append(lines, generateVMessURL(proxy))
        }
    }

    return strings.Join(lines, "\n")
}
```

### 8.2 固定测试数据

```yaml
# tests/testdata/subscriptions/sample.yaml
name: "示例订阅数据"
description: "用于测试的标准订阅数据"
proxies:
  - type: "ss"
    name: "香港节点1"
    server: "hk1.example.com"
    port: 8388
    cipher: "aes-256-gcm"
    password: "password123"

  - type: "vmess"
    name: "美国节点1"
    server: "us1.example.com"
    port: 443
    uuid: "12345678-1234-1234-1234-123456789abc"
    alterId: 0
    cipher: "auto"
    network: "ws"
    path: "/"
    tls: true

rules:
  - "DOMAIN-SUFFIX,google.com,Proxy"
  - "DOMAIN-KEYWORD,youtube,Proxy"
  - "GEOIP,CN,DIRECT"
  - "FINAL,Proxy"
```

## 9. 测试报告和监控

### 9.1 测试报告生成

```go
// tests/report/generator.go
package report

import (
    "encoding/json"
    "html/template"
    "time"
)

type TestReport struct {
    Timestamp    time.Time     `json:"timestamp"`
    Duration     time.Duration `json:"duration"`
    Coverage     float64       `json:"coverage"`
    PassedTests  int           `json:"passed_tests"`
    FailedTests  int           `json:"failed_tests"`
    Performance  Performance   `json:"performance"`
}

type Performance struct {
    QPS              float64 `json:"qps"`
    AvgResponseTime  int64   `json:"avg_response_time_ms"`
    MemoryUsage      int64   `json:"memory_usage_mb"`
}

func GenerateHTMLReport(report *TestReport) error {
    tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>SubConverter 测试报告</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .metric { margin: 10px 0; }
        .pass { color: green; }
        .fail { color: red; }
    </style>
</head>
<body>
    <h1>SubConverter Go 版本测试报告</h1>
    <div class="metric">测试时间: {{.Timestamp.Format "2006-01-02 15:04:05"}}</div>
    <div class="metric">执行时长: {{.Duration}}</div>
    <div class="metric">代码覆盖率: {{printf "%.2f" .Coverage}}%</div>
    <div class="metric {{if gt .FailedTests 0}}fail{{else}}pass{{end}}">
        测试结果: {{.PassedTests}} 通过, {{.FailedTests}} 失败
    </div>
    <h2>性能指标</h2>
    <div class="metric">QPS: {{printf "%.2f" .Performance.QPS}}</div>
    <div class="metric">平均响应时间: {{.Performance.AvgResponseTime}}ms</div>
    <div class="metric">内存使用: {{.Performance.MemoryUsage}}MB</div>
</body>
</html>
`

    t, err := template.New("report").Parse(tmpl)
    if err != nil {
        return err
    }

    file, err := os.Create("test-report.html")
    if err != nil {
        return err
    }
    defer file.Close()

    return t.Execute(file, report)
}
```

## 10. 测试最佳实践

### 10.1 测试编写规范

1. **命名规范**: 测试函数以 `Test` 开头，基准测试以 `Benchmark` 开头
2. **测试结构**: 使用 AAA 模式 (Arrange, Act, Assert)
3. **错误处理**: 使用 `testify/require` 处理致命错误，`testify/assert` 处理非致命错误
4. **测试隔离**: 每个测试函数应该独立，不依赖其他测试
5. **数据驱动**: 使用表格驱动测试处理多个测试用例

### 10.2 性能测试指导

1. **基准测试**: 使用 `b.ResetTimer()` 排除准备时间
2. **内存分析**: 使用 `b.ReportAllocs()` 报告内存分配
3. **并发测试**: 使用 `b.RunParallel()` 进行并发基准测试
4. **性能回归**: 建立性能基线，监控性能变化

### 10.3 Mock 和 Stub 策略

1. **接口 Mock**: 使用 `testify/mock` 创建接口 Mock
2. **HTTP Mock**: 使用 `httptest` 创建测试服务器
3. **数据库 Mock**: 使用 `go-sqlmock` 模拟数据库操作
4. **文件系统 Mock**: 使用 `afero` 模拟文件系统操作

## 11. 总结

这个测试策略和持续集成方案提供了：

### 11.1 全面的测试覆盖

- **多层次测试**: 从单元测试到端到端测试的完整覆盖
- **性能保障**: 基准测试和压力测试确保性能达标
- **兼容性验证**: 确保与现有系统完全兼容
- **安全检查**: 代码安全扫描和漏洞检测

### 11.2 高效的 CI/CD 流程

- **自动化测试**: 全流程自动化测试执行
- **质量门禁**: 严格的代码质量和覆盖率要求
- **快速反馈**: 快速发现和报告问题
- **持续监控**: 性能和质量指标持续监控

### 11.3 可维护的测试体系

- **清晰的测试结构**: 良好的测试组织和命名规范
- **可复用的测试工具**: 测试数据生成和工具函数
- **详细的测试报告**: 可视化的测试结果和性能报告
- **文档化的测试流程**: 完整的测试指导和最佳实践

这个方案确保了 SubConverter Go 版本的高质量交付，为项目的长期维护和发展提供了坚实的质量保障基础。
