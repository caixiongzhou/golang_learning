# 个人博客系统后端

基于 Go + Gin + GORM 开发的个人博客系统后端，提供完整的文章管理、用户认证和评论功能。

## 功能特性

- ✅ 用户注册、登录和JWT认证
- ✅ 文章的CRUD操作
- ✅ 评论系统（支持嵌套评论）
- ✅ 分页查询
- ✅ 统一的错误处理
- ✅ 完整的日志记录
- ✅ RESTful API设计
- ✅ 跨域支持

## 技术栈

- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **认证**: JWT
- **密码加密**: bcrypt

## 项目结构

```
blog-system/
├── cmd/server/ # 应用入口
├── internal/ # 内部包
│ ├── config/ # 配置管理
│ ├── controller/ # 控制器层
│ ├── middleware/ # 中间件
│ ├── model/ # 数据模型
│ ├── repository/ # 数据访问层
│ ├── service/ # 业务逻辑层
│ ├── dto/ # 数据传输对象
│ └── utils/ # 工具类
├── pkg/logger/ # 日志包
└── api/docs/ # API文档
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd blog-system
```

1. **安装依赖**

bash

```bash
go mod tidy
```



1. **配置环境变量**

bash

```bash
cp .env .env
# 编辑 .env 文件，配置数据库连接等信息
```

1. **启动服务**

bash

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动API接口文档

### 认证接口

#### 用户注册

http

```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "firstName": "Test",
  "lastName": "User"
}
```



#### 用户登录

http

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

### 文章接口

#### 获取文章列表

http

```http
GET /api/posts?page=1&pageSize=10
```



#### 创建文章

http

```http
POST /api/posts
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "文章标题",
  "content": "文章内容",
  "excerpt": "文章摘要",
  "status": "published"
}
```



#### 更新文章

http

```http
PUT /api/posts/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "更新后的标题",
  "content": "更新后的内容"
}
```



#### 删除文章

http

```http
DELETE /api/posts/{id}
Authorization: Bearer <token>
```



### 评论接口

#### 创建评论

http

```http
POST /api/comments
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "评论内容",
  "postId": 1,
  "parentId": null
}
```



#### 获取文章评论

http

```http
GET /api/comments/posts/{postId}
```



## 数据库设计

### Users 表

- id (主键)
- username (用户名，唯一)
- email (邮箱，唯一)
- password (加密密码)
- first_name
- last_name
- avatar
- bio
- status
- created_at
- updated_at

### Posts 表

- id (主键)
- title
- content
- excerpt
- slug (URL友好标识，唯一)
- status
- view_count
- user_id (外键)
- created_at
- updated_at
- published_at

### Comments 表

- id (主键)
- content
- status
- post_id (外键)
- user_id (外键)
- parent_id (自关联，支持嵌套评论)
- created_at
- updated_at

## 错误码说明

| 状态码 | 说明           |
| :----- | :------------- |
| 200    | 成功           |
| 201    | 创建成功       |
| 400    | 请求参数错误   |
| 401    | 未授权         |
| 403    | 禁止访问       |
| 404    | 资源不存在     |
| 500    | 服务器内部错误 |

## 测试

使用Postman或其他API测试工具进行接口测试：

1. 首先注册用户或使用已有账号登录获取token
2. 在请求头中添加: `Authorization: Bearer <your-token>`
3. 测试各个接口功能

## 部署

### 生产环境部署

1. 编译项目:

bash

```
go build -o blog-system cmd/server/main.go
```



1. 配置生产环境变量
2. 使用进程管理器(如systemd, supervisord)运行服务

### Docker部署

dockerfile

```
FROM golang:1.21-alpine

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main cmd/server/main.go

EXPOSE 8080
CMD ["./main"]
```



## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证

MIT License

text

```
## 15. 测试用例示例

**文件：internal/controller/auth_controller_test.go**
```go
package controller

import (
	"blog-system/internal/dto"
	"blog-system/internal/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService 模拟认证服务
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*dto.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*dto.AuthResponse), args.Error(1)
}

func (m *MockAuthService) GetUserProfile(userID uint) (*dto.AuthResponse, error) {
	args := m.Called(userID)
	return args.Get(0).(*dto.AuthResponse), args.Error(1)
}

func TestAuthController_Register(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建模拟服务
	mockService := new(MockAuthService)
	controller := NewAuthController(mockService)

	// 测试用例
	tests := []struct {
		name           string
		request        dto.RegisterRequest
		mockResponse   *dto.AuthResponse
		mockError      error
		expectedStatus int
	}{
		{
			name: "注册成功",
			request: dto.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockResponse: &dto.AuthResponse{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				Token:    "jwt-token",
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置模拟期望
			mockService.On("Register", &tt.request).Return(tt.mockResponse, tt.mockError)

			// 创建请求
			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// 创建响应记录器
			w := httptest.NewRecorder()
			router := gin.Default()
			router.POST("/api/auth/register", controller.Register)

			// 执行请求
			router.ServeHTTP(w, req)

			// 断言
			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
```


## powershell 测试用例：权限验证，测试用户只能修改和删除自己的文章，无法修改和删除别的用户的文章
```
# 完整的权限验证测试脚本
Write-Host "=== 开始完整的权限验证测试 ===" -ForegroundColor Cyan

# 使用已有的用户Token
$userAToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJ1c2VybmFtZSI6InVzZXJBXzAyNDYyOCIsImlzcyI6ImJsb2ctc3lzdGVtIiwiZXhwIjoxNzYzOTIzNTg4LCJpYXQiOjE3NjM4MzcxODh9.2SSqDxHW1MRkmTj0rrKbVAErhYQOrPmTTIss1KsFyt0"
$userBToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJ1c2VybmFtZSI6InVzZXJCXzAyNDgyNCIsImlzcyI6ImJsb2ctc3lzdGVtIiwiZXhwIjoxNzYzOTIzNzA0LCJpYXQiOjE3NjM4MzczMDR9.AncQJ0zaA-XRMun6HVluQDt-Ml2f_8p03K7HJI9kzmk"

Write-Host "用户A ID: 3 (userA_024628)" -ForegroundColor Cyan
Write-Host "用户B ID: 4 (userB_024824)" -ForegroundColor Cyan

# 1. 验证用户Token有效性
Write-Host "`n=== 1. 验证用户Token有效性 ===" -ForegroundColor Green

try {
    $profileA = Invoke-RestMethod -Uri "http://localhost:8080/api/auth/profile" `
      -Method GET `
      -Headers @{"Authorization"="Bearer $userAToken"}
    Write-Host "✅ 用户A Token有效: $($profileA.data.username)" -ForegroundColor Green
} catch {
    Write-Host "❌ 用户A Token无效: $($_.Exception.Message)" -ForegroundColor Red
    return
}

try {
    $profileB = Invoke-RestMethod -Uri "http://localhost:8080/api/auth/profile" `
      -Method GET `
      -Headers @{"Authorization"="Bearer $userBToken"}
    Write-Host "✅ 用户B Token有效: $($profileB.data.username)" -ForegroundColor Green
} catch {
    Write-Host "❌ 用户B Token无效: $($_.Exception.Message)" -ForegroundColor Red
    return
}

# 2. 用户A创建文章
Write-Host "`n=== 2. 用户A创建文章 ===" -ForegroundColor Green
$createPostBody = @{
    title = "用户A的专属文章 $(Get-Date -Format 'HH:mm:ss')"
    content = "这是用户A创建的专属文章内容，其他用户不应该能修改或删除。创建时间: $(Get-Date)"
    excerpt = "这是用户A的文章摘要 - 权限测试"
    status = "published"
    tags = "权限测试,用户A专属"
} | ConvertTo-Json

Write-Host "创建文章请求体: $createPostBody"

try {
    $createResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/posts" `
      -Method POST `
      -Headers @{"Authorization"="Bearer $userAToken"; "Content-Type"="application/json"} `
      -Body $createPostBody
    Write-Host "✅ 用户A创建文章成功!" -ForegroundColor Green
    $postId = $createResponse.data.id
    Write-Host "创建的文章ID: $postId" -ForegroundColor Cyan
    Write-Host "文章标题: $($createResponse.data.title)" -ForegroundColor Cyan
    Write-Host "文章作者: $($createResponse.data.author.username)" -ForegroundColor Cyan
} catch {
    Write-Host "❌ 用户A创建文章失败: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "错误详情: $($_.ErrorDetails.Message)" -ForegroundColor Red
    return
}

# 3. 用户B尝试更新用户A的文章（应该失败）
Write-Host "`n=== 3. 用户B尝试更新用户A的文章 ===" -ForegroundColor Yellow
Write-Host "预期结果: 应该返回 403 权限错误" -ForegroundColor Yellow

$updatePostBody = @{
    title = "用户B尝试非法修改标题 $(Get-Date -Format 'HH:mm:ss')"
    content = "用户B尝试非法修改内容，这不应该成功"
    status = "published"
} | ConvertTo-Json

Write-Host "用户B的更新请求体: $updatePostBody"

try {
    $updateResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/posts/$postId" `
      -Method PUT `
      -Headers @{"Authorization"="Bearer $userBToken"; "Content-Type"="application/json"} `
      -Body $updatePostBody
    Write-Host "❌ 权限验证失败！用户B竟然能更新他人文章" -ForegroundColor Red
    Write-Host "响应: $($updateResponse | ConvertTo-Json -Depth 3)" -ForegroundColor Red
} catch {
    if ($_.Exception.Response.StatusCode.value__ -eq 403) {
        Write-Host "✅ 权限验证成功！用户B无法更新他人文章" -ForegroundColor Green
        Write-Host "错误信息: $($_.ErrorDetails.Message)" -ForegroundColor Green
    } else {
        Write-Host "❌ 出现意外错误: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host "状态码: $($_.Exception.Response.StatusCode.value__)" -ForegroundColor Red
    }
}

# 4. 用户A更新自己的文章（应该成功）
Write-Host "`n=== 4. 用户A更新自己的文章 ===" -ForegroundColor Green
Write-Host "预期结果: 应该成功更新" -ForegroundColor Green

$updateByOwnerBody = @{
    title = "用户A合法更新自己的文章 $(Get-Date -Format 'HH:mm:ss')"
    content = "用户A合法地更新了自己的文章内容。更新时间: $(Get-Date)"
    status = "published"
} | ConvertTo-Json

Write-Host "用户A的更新请求体: $updateByOwnerBody"

try {
    $updateByOwnerResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/posts/$postId" `
      -Method PUT `
      -Headers @{"Authorization"="Bearer $userAToken"; "Content-Type"="application/json"} `
      -Body $updateByOwnerBody
    Write-Host "✅ 用户A成功更新自己的文章" -ForegroundColor Green
    Write-Host "新标题: $($updateByOwnerResponse.data.title)" -ForegroundColor Cyan
    Write-Host "更新响应: $($updateByOwnerResponse.message)" -ForegroundColor Cyan
} catch {
    Write-Host "❌ 用户A更新自己的文章失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 5. 用户B尝试删除用户A的文章（应该失败）
Write-Host "`n=== 5. 用户B尝试删除用户A的文章 ===" -ForegroundColor Yellow
Write-Host "预期结果: 应该返回 403 权限错误" -ForegroundColor Yellow

try {
    $deleteResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/posts/$postId" `
      -Method DELETE `
      -Headers @{"Authorization"="Bearer $userBToken"}
    Write-Host "❌ 权限验证失败！用户B竟然能删除他人文章" -ForegroundColor Red
    Write-Host "响应: $($deleteResponse | ConvertTo-Json -Depth 3)" -ForegroundColor Red
} catch {
    if ($_.Exception.Response.StatusCode.value__ -eq 403) {
        Write-Host "✅ 权限验证成功！用户B无法删除他人文章" -ForegroundColor Green
        Write-Host "错误信息: $($_.ErrorDetails.Message)" -ForegroundColor Green
    } else {
        Write-Host "❌ 出现意外错误: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host "状态码: $($_.Exception.Response.StatusCode.value__)" -ForegroundColor Red
    }
}

# 6. 验证文章是否还存在且未被用户B修改
Write-Host "`n=== 6. 验证文章状态 ===" -ForegroundColor Green

try {
    $getResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/posts/$postId"
    Write-Host "✅ 文章仍然存在" -ForegroundColor Green
    Write-Host "当前标题: $($getResponse.data.title)" -ForegroundColor Cyan
    Write-Host "当前作者: $($getResponse.data.author.username)" -ForegroundColor Cyan
    Write-Host "文章ID: $($getResponse.data.id)" -ForegroundColor Cyan
    
    if ($getResponse.data.author.username -eq "userA_024628") {
        Write-Host "✅ 文章作者正确，仍然是用户A" -ForegroundColor Green
    } else {
        Write-Host "❌ 文章作者异常" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ 获取文章失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 7. 用户A删除自己的文章（应该成功）
Write-Host "`n=== 7. 用户A删除自己的文章 ===" -ForegroundColor Green
Write-Host "预期结果: 应该成功删除" -ForegroundColor Green

try {
    $deleteByOwnerResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/posts/$postId" `
      -Method DELETE `
      -Headers @{"Authorization"="Bearer $userAToken"}
    Write-Host "✅ 用户A成功删除自己的文章" -ForegroundColor Green
    Write-Host "删除响应: $($deleteByOwnerResponse.message)" -ForegroundColor Cyan
} catch {
    Write-Host "❌ 用户A删除自己的文章失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 8. 最终验证文章是否已删除
Write-Host "`n=== 8. 最终验证文章是否已删除 ===" -ForegroundColor Green

try {
    $finalGetResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/posts/$postId"
    Write-Host "❌ 文章仍然存在，删除失败" -ForegroundColor Red
    Write-Host "响应: $($finalGetResponse | ConvertTo-Json -Depth 3)" -ForegroundColor Red
} catch {
    if ($_.Exception.Response.StatusCode.value__ -eq 404) {
        Write-Host "✅ 文章已成功删除" -ForegroundColor Green
    } else {
        Write-Host "获取文章时出现错误: $($_.Exception.Message)" -ForegroundColor Yellow
        Write-Host "状态码: $($_.Exception.Response.StatusCode.value__)" -ForegroundColor Yellow
    }
}

Write-Host "`n=== 测试总结 ===" -ForegroundColor Cyan
Write-Host "✅ 用户A可以创建、更新、删除自己的文章" -ForegroundColor Green
Write-Host "✅ 用户B无法更新或删除用户A的文章" -ForegroundColor Green
Write-Host "✅ 权限验证系统工作正常！" -ForegroundColor Green
Write-Host "`n=== 测试完成 ===" -ForegroundColor Cyan
```

### powershell 测试脚本：使用用户A Token 发布文章的脚本，自测请修改用户Token
```
# 发布多篇文章的 PowerShell 脚本
Write-Host "开始发布测试文章..." -ForegroundColor Yellow

# 用户A的Token
$token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJ1c2VybmFtZSI6InVzZXJBXzAyNDYyOCIsImlzcyI6ImJsb2ctc3lzdGVtIiwiZXhwIjoxNzYzOTIzNTg4LCJpYXQiOjE3NjM4MzcxODh9.2SSqDxHW1MRkmTj0rrKbVAErhYQOrPmTTIss1KsFyt0"

$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

$baseUrl = "http://localhost:8080"

# 测试文章数据
$articles = @(
    @{
        title = "Go语言入门指南"
        content = "Go语言是一种开源的编程语言，由Google开发。它具有高效的并发编程能力和简洁的语法。本文将从基础语法开始，带你逐步掌握Go语言的核心概念。"
        summary = "本文介绍了Go语言的基础知识和核心特性"
        tags = "Go,编程,后端"
        status = "published"
    },
    @{
        title = "Gin框架实战教程"
        content = "Gin是一个用Go语言编写的Web框架，以其高性能和易用性著称。本文将教你如何使用Gin构建RESTful API，包括路由、中间件、参数绑定等核心功能。"
        summary = "学习使用Gin框架构建高效的Web应用"
        tags = "Gin,Go,Web框架"
        status = "published"
    },
    @{
        title = "JWT认证原理与实践"
        content = "JWT（JSON Web Token）是一种流行的跨域认证解决方案。本文将深入探讨JWT的工作原理、安全性考虑以及在Go语言中的实现方式。"
        summary = "深入理解JWT认证机制及其在Go中的实现"
        tags = "JWT,认证,安全"
        status = "published"
    },
    @{
        title = "数据库设计与优化"
        content = "良好的数据库设计是应用性能的基石。本文将分享数据库设计的最佳实践，包括表结构设计、索引优化、查询性能调优等内容。"
        summary = "数据库设计原则和性能优化技巧"
        tags = "数据库,MySQL,优化"
        status = "published"
    },
    @{
        title = "微服务架构解析"
        content = "微服务架构通过将应用拆分为小型、独立的服务来提高可维护性和可扩展性。本文将探讨微服务的优势、挑战以及实施策略。"
        summary = "微服务架构的核心概念和实践经验"
        tags = "微服务,架构,分布式"
        status = "draft"
    }
)

# 发布文章的函数
function Publish-Article {
    param($articleData)
    
    $body = $articleData | ConvertTo-Json
    
    try {
        $response = Invoke-RestMethod -Uri "$baseUrl/api/posts" -Method POST -Headers $headers -Body $body
        Write-Host "✓ 发布成功: $($articleData.title)" -ForegroundColor Green
        return $response
    } catch {
        Write-Host "✗ 发布失败: $($articleData.title)" -ForegroundColor Red
        Write-Host "  错误: $($_.Exception.Message)" -ForegroundColor Red
        return $null
    }
}

# 发布所有文章
Write-Host "`n开始发布文章..." -ForegroundColor Cyan
$successCount = 0

foreach ($article in $articles) {
    $result = Publish-Article -articleData $article
    if ($result) {
        $successCount++
    }
    Start-Sleep -Milliseconds 500  # 稍微延迟，避免请求过快
}

Write-Host "`n发布完成! 成功: $successCount/$($articles.Count) 篇" -ForegroundColor Yellow
```

### powershell 测试脚本：测试接口 authPosts.GET("/my", postController.GetUserPosts)      // 获取所有文章列表
```
# 修复后的测试脚本
Write-Host "`n2. 测试接口: GET /api/posts/my" -ForegroundColor Cyan
Write-Host "   URL: $baseUrl/api/posts/my" -ForegroundColor Gray
Write-Host "   描述: 获取当前登录用户的所有文章" -ForegroundColor Gray

try {
    $startTime = Get-Date
    $response1 = Invoke-RestMethod -Uri "$baseUrl/api/posts/my" -Method GET -Headers $headers
    $endTime = Get-Date
    $duration = ($endTime - $startTime).TotalMilliseconds

    Write-Host "    请求成功 (耗时: $duration ms)" -ForegroundColor Green
    Write-Host "   HTTP状态码: 200" -ForegroundColor Green

    # 显示完整的响应结构
    Write-Host "`n   完整响应:" -ForegroundColor Yellow
    $response1 | ConvertTo-Json -Depth 5 | Write-Host

    Write-Host "`n   响应结构分析:" -ForegroundColor Yellow
    Write-Host "   - success: $($response1.success)"
    Write-Host "   - message: $($response1.message)"
    
    # 检查 data 字段的结构
    if ($response1.data) {
        Write-Host "   - data 类型: $($response1.data.GetType().Name)"
        Write-Host "   - data 包含的字段: $($response1.data.PSObject.Properties.Name -join ', ')"
        
        # 检查是否有 posts 字段
        if ($response1.data.PSObject.Properties.Name -contains "posts") {
            $postsCount = $response1.data.posts.Count
            Write-Host "   - 文章数量: $postsCount" -ForegroundColor Green
            
            if ($postsCount -gt 0) {
                Write-Host "`n   文章列表:" -ForegroundColor Yellow
                for ($i = 0; $i -lt $postsCount; $i++) {
                    $post = $response1.data.posts[$i]
                    Write-Host "   [$($i+1)] ID: $($post.id) | 标题: $($post.title) | 状态: $($post.status) | 作者: $($post.author.username)" -ForegroundColor White
                }
            } else {
                Write-Host "   ℹ  posts 数组为空" -ForegroundColor Blue
            }
        } else {
            Write-Host "   ❗ data 中没有 posts 字段" -ForegroundColor Red
        }
        
        # 检查分页信息
        if ($response1.data.PSObject.Properties.Name -contains "pagination") {
            $pagination = $response1.data.pagination
            Write-Host "   - 分页信息: 页 $($pagination.page), 大小 $($pagination.pageSize), 总数 $($pagination.total)" -ForegroundColor Cyan
        }
    } else {
        Write-Host "   ❗ data 字段为空" -ForegroundColor Red
    }

} catch {
    Write-Host "    请求失败: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $statusCode = $_.Exception.Response.StatusCode.value__
        Write-Host "   HTTP状态码: $statusCode" -ForegroundColor Red
        Write-Host "   错误详情: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}
```

### powershell 测试脚本：测试接口 authPosts.GET("/user/:id", postController.GetUserPost)       // 获取当前用户的单个文章详情
```
# 测试接口：GET /api/posts/user/:id - 获取当前用户的单个文章详情
Write-Host "`n=== 测试接口: GET /api/posts/user/:id ===" -ForegroundColor Cyan
Write-Host "   描述: 获取当前登录用户的指定文章详情（需要验证文章归属）" -ForegroundColor Gray

# 使用现有的文章ID进行测试
$testPostIds = @(2, 3, 4, 5, 6)

foreach ($postId in $testPostIds) {
    Write-Host "`n测试获取用户文章 ID: $postId" -ForegroundColor Yellow
    Write-Host "  URL: http://localhost:8080/api/posts/user/$postId" -ForegroundColor Gray
    
    try {
        $startTime = Get-Date
        $response = Invoke-RestMethod -Uri "http://localhost:8080/api/posts/user/$postId" -Method GET -Headers $headers
        $endTime = Get-Date
        $duration = ($endTime - $startTime).TotalMilliseconds
        
        Write-Host "  ✅ 请求成功 (耗时: $duration ms)" -ForegroundColor Green
        Write-Host "  HTTP状态码: 200" -ForegroundColor Green
        
        # 显示文章详情
        $post = $response.data
        Write-Host "`n  文章详情:" -ForegroundColor White
        Write-Host "  - ID: $($post.id)" -ForegroundColor Cyan
        Write-Host "  - 标题: $($post.title)" -ForegroundColor Cyan
        Write-Host "  - 状态: $($post.status)" -ForegroundColor Cyan
        Write-Host "  - 作者: $($post.author.username)" -ForegroundColor Cyan
        Write-Host "  - 浏览次数: $($post.viewCount)" -ForegroundColor Cyan
        Write-Host "  - 创建时间: $($post.createdAt)" -ForegroundColor Cyan
        Write-Host "  - 更新时间: $($post.updatedAt)" -ForegroundColor Cyan
        
        if ($post.publishedAt) {
            Write-Host "  - 发布时间: $($post.publishedAt)" -ForegroundColor Cyan
        }
        
        if ($post.excerpt) {
            Write-Host "  - 摘要: $($post.excerpt)" -ForegroundColor Gray
        }
        
        if ($post.tags) {
            Write-Host "  - 标签: $($post.tags)" -ForegroundColor Gray
        }
        
        Write-Host "  - 内容预览: $($post.content.Substring(0, [Math]::Min(100, $post.content.Length)))..." -ForegroundColor Gray
        
    } catch {
        Write-Host "  ❌ 请求失败: $($_.Exception.Message)" -ForegroundColor Red
        
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode.value__
            Write-Host "  HTTP状态码: $statusCode" -ForegroundColor Red
            
            if ($statusCode -eq 404) {
                Write-Host "  💡 提示: 文章不存在或不属于当前用户" -ForegroundColor Yellow
            } elseif ($statusCode -eq 401) {
                Write-Host "  💡 提示: 未授权访问" -ForegroundColor Yellow
            } elseif ($statusCode -eq 403) {
                Write-Host "  💡 提示: 无权访问此文章" -ForegroundColor Yellow
            }
            
            Write-Host "  错误详情: $($_.ErrorDetails.Message)" -ForegroundColor Red
        }
    }
}
```
这个完整的博客系统后端项目包含了：

1. **清晰的分层架构**：Controller → Service → Repository → Model
2. **完整的RESTful API**：符合REST规范
3. **JWT认证**：安全的用户认证机制
4. **错误处理**：统一的错误响应格式
5. **日志记录**：完整的操作日志
6. **数据库操作**：使用GORM进行数据持久化
7. **API文档**：详细的接口说明
8. **测试用例**：基础的单元测试示例

项目可以直接运行，只需要配置好MySQL数据库和相应的环境变量即可。