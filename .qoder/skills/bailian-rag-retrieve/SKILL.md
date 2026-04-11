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

### 步骤 2: 创建临时 Go 脚本

在 `/tmp` 目录创建临时 Go 文件。使用以下模板：

```go
// /tmp/bailian_retrieve_{timestamp}.go
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
	req.SetRerankTopN(5)

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
```

同时创建 `go.mod`：

```
// /tmp/bailian_retrieve_{timestamp}_mod/go.mod
module bailian_retrieve

go 1.24

```

### 步骤 3: 执行检索脚本

```bash
cd /tmp/bailian_retrieve_{timestamp}_mod && \
  go mod tidy && \
  WORKSPACE_ID="<workspace_id>" INDEX_ID="<index_id>" QUERY="<query>" \
  go run ../bailian_retrieve_{timestamp}.go
```

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
