---
name: bailian-rag-retrieve
description: 检索阿里云百炼RAG知识库，获取精准文档片段。使用百炼 Go SDK 进行语义检索，支持向量召回、关键词召回、重排和多轮对话上下文。当用户需要查询知识库内容、检索文档、或从RAG知识库获取信息时使用此技能。使用 Go 语言实现检索功能。
---

# 百炼 RAG 知识库检索

## 认证配置

使用阿里云 AccessKey 进行认证。确保环境变量已设置：

```bash
export ALIBABA_CLOUD_ACCESS_KEY_ID="<your-access-key-id>"
export ALIBABA_CLOUD_ACCESS_KEY_SECRET="<your-access-key-secret>"
```

RAM 子账号需授予 `AliyunBailianDataFullAccess` 权限策略并加入业务空间。

## Go SDK

### 安装

```bash
go get github.com/alibabacloud-go/bailian-20231229/v2
```

### 客户端初始化

```go
import (
    openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
    bailian "github.com/alibabacloud-go/bailian-20231229/v2/client"
)

func CreateBailianClient() (*bailian.Client, error) {
    config := &openapi.Config{
        AccessKeyId:     tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")),
        AccessKeySecret: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")),
        Endpoint:        tea.String("bailian.cn-beijing.aliyuncs.com"),
    }
    return bailian.NewClient(config)
}
```

### Retrieve API

```go
// 简单调用
func (client *Client) Retrieve(WorkspaceId *string, request *RetrieveRequest) (*RetrieveResponse, error)

// 带运行时选项调用（可设置超时等）
func (client *Client) RetrieveWithOptions(WorkspaceId *string, tmpReq *RetrieveRequest, headers map[string]*string, runtime *dara.RuntimeOptions) (*RetrieveResponse, error)
```

### RetrieveRequest 字段

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| IndexId | *string | 是 | - | 知识库ID |
| Query | *string | 否 | - | 检索文本 |
| DenseSimilarityTopK | *int32 | 否 | 100 | 向量召回数量 (0-100) |
| SparseSimilarityTopK | *int32 | 否 | 100 | 关键词召回数量 (0-100) |
| EnableReranking | *bool | 否 | true | 是否启用重排 |
| EnableRewrite | *bool | 否 | false | 是否启用多轮改写 |
| RerankMinScore | *float32 | 否 | - | 重排相似度阈值 (0.01-1.00) |
| RerankTopN | *int32 | 否 | 5 | 重排后返回数量 (1-20) |
| Rerank | []*RetrieveRequestRerank | 否 | - | 重排配置 (ModelName, RerankMode, RerankInstruct) |
| Rewrite | []*RetrieveRequestRewrite | 否 | - | 改写配置 (ModelName) |
| QueryHistory | []*RetrieveRequestQueryHistory | 否 | - | 对话上下文 (Role, Content) |
| SaveRetrieverHistory | *bool | 否 | false | 是否保存检索历史 |
| SearchFilters | []map[string]*string | 否 | - | 过滤条件 |
| Images | []*string | 否 | - | 图片URL列表 |

注意：DenseSimilarityTopK + SparseSimilarityTopK 之和须 <= 200。

### RetrieveResponse 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| Body.Code | *string | 状态码 |
| Body.Success | *bool | 是否成功 |
| Body.Message | *string | 错误信息 |
| Body.Data.Nodes | []*RetrieveResponseBodyDataNodes | 检索结果列表 |
| Body.Data.Nodes[].Text | *string | 命中的文档片段 |
| Body.Data.Nodes[].Score | *float64 | 匹配度分数 (0-1) |
| Body.Data.Nodes[].Metadata | interface{} | 文档元信息 |

## 检索流程

### 步骤 1: 获取必要参数

从用户处获取或确认：
- `WorkspaceId`: 业务空间标识
- `IndexId`: 知识库标识
- `Query`: 检索内容

### 步骤 2: 创建临时 Go 项目

在 `/tmp` 目录下创建一个独立的 Go 模块目录，**源文件和 go.mod 必须放在同一个目录内**。

**CRITICAL - 常见错误规避：**
1. **文件名绝不能包含 `_test.go` 后缀**，Go 会将其视为测试文件，`go run` 拒绝执行。使用 `main.go` 即可。
2. **源文件必须放在 go.mod 所在目录内**，否则 `go run` 无法解析模块依赖。
3. **go.mod 中不要写 `require` 和版本号**，让 `go mod tidy` 根据源文件 import 自动解析。
4. **必须先创建源文件再执行 `go mod tidy`**，否则 tidy 无法确定依赖。
5. **执行 `go run` 前必须 `cd` 到模块目录**，使用 `go run .` 而非指定外部文件路径。

按以下顺序操作：

```bash
# 1. 创建模块目录
mkdir -p /tmp/bailian_retrieve_mod

# 2. 创建 go.mod（不含 require，让 tidy 自动解析）
cat > /tmp/bailian_retrieve_mod/go.mod << 'GOMOD'
module bailian_retrieve

go 1.24
GOMOD

# 3. 创建源文件（在同一目录内，文件名为 main.go）
cat > /tmp/bailian_retrieve_mod/main.go << 'GOSRC'
package main

import (
	"fmt"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	bailian "github.com/alibabacloud-go/bailian-20231229/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

func main() {
	workspaceID := os.Getenv("WORKSPACE_ID")
	indexID := os.Getenv("INDEX_ID")
	query := os.Getenv("QUERY")

	// 初始化客户端
	client, err := bailian.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")),
		AccessKeySecret: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")),
		Endpoint:        tea.String("bailian.cn-beijing.aliyuncs.com"),
	})
	if err != nil {
		fmt.Printf("创建客户端失败: %v\n", err)
		os.Exit(1)
	}

	// 构建请求
	req := &bailian.RetrieveRequest{}
	req.SetIndexId(indexID)
	req.SetQuery(query)
	req.SetDenseSimilarityTopK(50)
	req.SetEnableReranking(true)
	req.SetRerankTopN(10)

	// 执行检索
	resp, err := client.Retrieve(tea.String(workspaceID), req)
	if err != nil {
		fmt.Printf("检索失败: %v\n", err)
		os.Exit(1)
	}

	// 处理响应
	if resp.Body != nil && resp.Body.Data != nil && len(resp.Body.Data.Nodes) > 0 {
		fmt.Printf("===== 检索结果 (共 %d 条) =====\n\n", len(resp.Body.Data.Nodes))
		for i, node := range resp.Body.Data.Nodes {
			fmt.Printf("--- 结果 %d (匹配度: %.3f) ---\n", i+1, tea.Float64Value(node.Score))
			fmt.Printf("内容:\n%s\n\n", tea.StringValue(node.Text))
		}
	} else {
		code := ""
		msg := ""
		if resp.Body != nil {
			code = tea.StringValue(resp.Body.Code)
			msg = tea.StringValue(resp.Body.Message)
		}
		fmt.Printf("未找到相关结果 (Code: %s, Message: %s)\n", code, msg)
	}
}
GOSRC
```

### 步骤 3: 安装依赖并执行检索

**必须 cd 到模块目录后执行**，使用 `go run .` 运行：

```bash
cd /tmp/bailian_retrieve_mod && \
  GOPROXY=https://mirrors.aliyun.com/goproxy/,direct go mod tidy && \
  WORKSPACE_ID="<workspace_id>" INDEX_ID="<index_id>" QUERY="<query>" \
  go run .
```

**关键要点：**
- `go mod tidy` 必须在源文件已存在后执行，才能根据 import 解析依赖
- 使用 `GOPROXY=https://mirrors.aliyun.com/goproxy/,direct` 加速依赖下载
- 使用 `go run .` 而非 `go run ../外部文件.go`，确保模块依赖正确解析

### 步骤 4: 处理响应

1. 检查 `Body.Success` 是否为 `true`
2. 遍历 `Body.Data.Nodes` 获取检索结果
3. 按 `Score` 排序，提取最相关的内容
4. 将 `Text` 字段内容作为上下文提供给大模型

### 步骤 5: 错误处理

常见错误：
- `Index.InvalidParameter`: 参数缺失或无效
- `Index.NoWorkspacePermissions`: AccessKey 没有业务空间权限
- 认证失败: 检查 AccessKey 配置
- 超时: 检索计算复杂，使用 `RetrieveWithOptions` 设置超时

## 使用示例

### 基础检索

```go
req := &bailian.RetrieveRequest{}
req.SetIndexId("51x6qp66p2")
req.SetQuery("DND 5e 战斗中掩护规则")
resp, err := client.Retrieve(tea.String("llm-xxx"), req)
```

### 带重排阈值和 TopN

```go
req := &bailian.RetrieveRequest{}
req.SetIndexId("51x6qp66p2")
req.SetQuery("通义千问模型能力")
req.SetDenseSimilarityTopK(50)
req.SetEnableReranking(true)
req.SetRerankMinScore(0.3)
req.SetRerankTopN(5)
resp, err := client.Retrieve(tea.String("llm-xxx"), req)
```

### 多轮对话检索

```go
req := &bailian.RetrieveRequest{}
req.SetIndexId("51x6qp66p2")
req.SetQuery("它的主要功能是什么")
req.SetEnableRewrite(true)
req.SetQueryHistory([]*bailian.RetrieveRequestQueryHistory{
    {Role: tea.String("user"), Content: tea.String("什么是阿里云百炼")},
    {Role: tea.String("assistant"), Content: tea.String("阿里云百炼是...")},
})
resp, err := client.Retrieve(tea.String("llm-xxx"), req)
```

### 带超时的检索

```go
import "github.com/alibabacloud-go/tea/dara"

runtime := &dara.RuntimeOptions{
    ReadTimeout:  tea.Int(60000),
    ConnectTimeout: tea.Int(30000),
}
headers := make(map[string]*string)
resp, err := client.RetrieveWithOptions(tea.String("llm-xxx"), req, headers, runtime)
```

## 最佳实践

1. **TopK 设置**: 首次召回设 50-100，重排后取 Top 5-10
2. **重排阈值**: 使用 `RerankMinScore` 设 0.3 可过滤低质量结果
3. **多轮对话**: 开启 `EnableRewrite` 并传入 `QueryHistory` 提升上下文理解
4. **超时策略**: 检索计算复杂，建议通过 `RetrieveWithOptions` 设置 30-60 秒超时
5. **结果处理**: 优先使用 `Score` 高的片段，合并相关内容
6. **临时文件**: 所有临时 Go 脚本和 go.mod 放置在 `/tmp` 目录，执行完成后可删除

## Go 工程化注意事项（CRITICAL）

以下是基于实际执行经验总结的常见错误，务必规避：

### 1. 文件名不能包含 `_test.go`

Go 将 `_test.go` 后缀的文件视为测试文件，`go run` 会拒绝执行并报错。源文件统一命名为 `main.go`。

### 2. 源文件与 go.mod 必须在同一目录

Go 模块依赖解析基于 `go.mod` 所在目录。如果 `.go` 文件在模块目录外，`go run /外部路径/main.go` 会报 `no required module provides package` 错误。**所有文件必须放在同一个目录下。**

### 3. 先写源文件再 go mod tidy

`go mod tidy` 根据目录内的 `.go` 文件 import 语句来解析依赖。如果目录内还没有源文件，tidy 无法确定需要哪些依赖，会产生 `matched no packages` 警告。正确顺序：
1. 创建模块目录
2. 创建 `go.mod`（不含 require）
3. 创建 `main.go`（含 import）
4. 执行 `go mod tidy`
5. 执行 `go run .`

### 4. go.mod 不写 require 和版本号

`require github.com/alibabacloud-go/bailian-20231229/v2 latest` 中的 `latest` 不是有效的 Go module 版本语法。只需声明 `module` 和 `go` 版本，让 `go mod tidy` 自动解析所有依赖及版本。

### 5. 使用 go run . 而非指定外部文件

在模块目录内执行 `go run .`，而不是 `go run ../外部文件.go`。后者会导致 Go 在文件所在目录（而非模块目录）查找 `go.mod`，从而找不到依赖。

### 6. 执行命令时必须 cd 到模块目录

Bash 工具的默认工作目录是项目根目录，不会自动切换。所有 `go mod tidy` 和 `go run` 命令必须通过 `cd /tmp/bailian_retrieve_mod && ...` 的方式确保在正确目录执行。

### 推荐的一键执行模板

```bash
# 创建模块目录和文件
mkdir -p /tmp/bailian_retrieve_mod && \
cat > /tmp/bailian_retrieve_mod/go.mod << 'EOF'
module bailian_retrieve
go 1.24
EOF
# (写入 main.go 省略，见步骤2)
# 然后一步到位：
cd /tmp/bailian_retrieve_mod && \
  GOPROXY=https://mirrors.aliyun.com/goproxy/,direct go mod tidy && \
  WORKSPACE_ID="<workspace_id>" INDEX_ID="<index_id>" QUERY="<query>" \
  go run .
```
