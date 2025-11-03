# OJ Backend 启动教程

这是一个基于 Go 语言的 Online Judge(在线判题系统)后端服务,支持代码提交、判题和题目管理功能。

## 项目概述

- **技术栈**: Go + Fiber + PostgreSQL + GORM
- **功能**: 代码判题、题目管理、测试用例验证
- **支持语言**: C、C++、Python 等多种编程语言
- **运行端口**: 3000

## 快速开始

### 方式一:Docker Compose 一键部署(推荐新手)

1. **克隆项目**
```bash
git clone https://github.com/sanbei101/cau-oj-backend.git
cd cau-oj-backend
```

2. **启动 PostgreSQL 和后端服务**
```bash
docker-compose up -d
```

这将自动启动：
- PostgreSQL 数据库
- OJ 后端服务

3. **查看服务状态**
```bash
docker-compose ps
```

4. **查看运行日志**
```bash
docker-compose logs -f postgres
docker-compose logs -f oj-backend
```

5. **验证服务**
```bash
curl http://localhost:3000
# 应该返回: Hello, World!
```

6. **停止服务**
```bash
docker-compose down
```

### 方式三:本地开发环境

#### 1. 环境要求
- Go 1.25+
- PostgreSQL 数据库
- gcc、g++、python3

#### 2. 使用 Docker 启动 PostgreSQL（推荐）
```bash
docker run -d \
  --name oj-postgres \
  -e POSTGRES_DB=mydatabase \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=secretpassword \
  -p 5432:5432 \
  postgres:18
```

#### 3. 安装 Go 依赖
```bash
go mod download
```

#### 4. 配置数据库连接 DSN

编辑 `app/db/postgresql.go` 文件中的 DSN 连接字符串：


```go
const dsn string = "host=<ip> user=postgres password=secretpassword dbname=mydatabase sslmode=disable"
```
> `<ip>` 是 PostgreSQL 实例的 IP 地址,如果是本地开发,则使用 `localhost`


**DSN 参数详解：**
- `host`: 数据库地址
- `user`: 数据库用户名
- `password`: 数据库密码
- `dbname`: 数据库名称
- `port`: 数据库端口,默认为 `5432`
- `sslmode`: SSL 模式,默认为 `disable`

#### 5. 运行项目
```bash
go run cmd/main.go
```

服务将在 `http://localhost:3000` 启动。

## API 接口

### 基础接口
- `GET /` - 健康检查,返回 "Hello, World!"

### 判题接口
- `POST /judge/submit-code` - 提交代码进行判题

### 题目管理接口
- `GET /problem/get-all-problem` - 获取所有题目
- `GET /problem/get-problem-by-id` - 根据ID获取题目详情
- `GET /problem/get-problem-test-case` - 获取题目测试用例

## 项目结构

```
cau-oj-backend/
├── cmd/                    # 应用入口
│   └── main.go            # 主程序入口
├── app/                   # 应用层
│   ├── controller/        # 控制器层
│   ├── model/            # 数据模型
│   ├── service/          # 业务逻辑层
│   └── db/               # 数据库配置
├── pkg/                  # 包层
│   ├── judge/            # 判题逻辑
│   ├── parsers/          # 题目解析器
│   └── utils/            # 工具函数
├── data/                 # 题目数据
├── Dockerfile           # Docker 构建文件
├── go.mod               # Go 模块文件
└── README.md           # 项目说明
```

## 题目数据格式

题目数据存储在 `data/` 目录下的 TOML 文件中,格式如下:

```toml
[[problem]]
name = "题目名称"
description = """
题目描述,支持 Markdown 格式
"""
tags = ["标签1", "标签2"]

[[problem.test]]
input = "测试输入"
output = "期望输出"
```

## 开发说明

### 添加新题目
1. 在 `data/` 目录下创建或编辑 TOML 文件
2. 按照 TOML 格式添加题目信息和测试用例
3. 重启服务以加载新题目

### 判题流程
1. 用户通过 API 提交代码
2. 系统根据题目语言选择相应的判题器
3. 编译（如需要）并运行代码
4. 比较输出结果与期望结果
5. 返回判题结果

### 支持的编程语言
- C/C++(需要 gcc/g++)
- Python(需要 python3)
- 可通过扩展 `pkg/judge/` 目录添加更多语言支持

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查 PostgreSQL 是否运行
   - 验证数据库连接配置

2. **Docker 构建失败**
   - 确保网络连接正常
   - 检查 Go 模块依赖是否能正确下载

3. **判题失败**
   - 检查编译工具是否正确安装(gcc, g++, python3)
   - 确认测试用例格式正确

### 日志查看
Docker 运行时查看日志:
```bash
docker logs oj-backend
```