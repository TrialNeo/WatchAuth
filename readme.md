## 项目介绍

WatchAuth 是一个功能完善的授权管理系统，专为软件授权和机器管理设计。该系统提供了用户管理、机器管理、授权控制等核心功能，旨在解决软件授权的安全和管理问题。

### 主要功能
- 用户管理：用户注册、登录、权限控制
- 机器管理：机器注册、状态监控、日志管理
- 授权管理：软件授权发放、验证、管理
- 日志系统：详细的操作日志和机器日志
- 数据统计：用户活跃度、机器使用情况等

## 快速开始

### 环境要求
- Go 1.20+
- Node.js 16+
- PostgreSQL 18+
- Redis 7+

### 后端运行

#### 开发环境
```shell
# 进入后端目录
cd backend/cmd

# 安装依赖
go mod tidy
# 运行开发服务器
go run main.go
```

#### 生产环境
```shell
# 构建
python build.py
```

### 前端运行

#### 开发环境
```shell
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 运行开发服务器
npm run dev
```

#### 生产环境
```shell
cd dist

# 进入前端目录
cd frontend

# 安装依赖
npm install

# 部署 dist 目录到服务器
```

## 项目结构

```
WatchAuth/
├── backend/            # 后端代码
│   ├── app/            # 应用逻辑
│   ├── config/         # 配置管理
│   ├── internal/       # 内部包
│   ├── middleware/     # 中间件
│   ├── model/          # 数据模型
│   ├── router/         # 路由
│   ├── utils/          # 工具函数
│   ├── main.go         # 主入口
│   └── go.mod          # Go 模块定义
├── frontend/           # 前端代码
│   ├── src/            # 源代码
│   │   ├── api/        # API 调用
│   │   ├── components/ # 组件
│   │   ├── config/     # 配置
│   │   ├── mocks/      # 模拟数据
│   │   ├── types/      # 类型定义
│   │   ├── utils/      # 工具函数
│   │   ├── views/      # 页面视图
│   │   └── main.ts     # 主入口
│   ├── index.html      # HTML 模板
│   ├── package.json    # 依赖管理
│   └── vite.config.ts  # Vite 配置
├── data/               # 数据文件
└── readme.md           # 项目说明
```

## 技术栈

### 后端
1. **Go 1.20+** - 主要开发语言
2. **Fiber** - 高性能 Web 框架
3. **GORM** - ORM 库，用于数据库操作
4. **JWT** - 用于身份认证
5. **Redis** - 用于缓存和会话管理
6. **PostgreSQL** - 主数据库
7. **Viper** - 配置管理

### 前端
1. **Vue 3** - 前端框架
2. **TypeScript** - 类型系统
3. **Element Plus** - UI 组件库
4. **Vite** - 构建工具
6. **IndexedDB** - 本地数据持久化

## 配置说明

### 后端配置
后端配置文件位于 `backend/config/config.yaml`，主要配置项包括：
- 数据库连接信息
- Redis 连接信息
- JWT 密钥
- 服务器端口

### 前端配置
前端配置文件位于 `frontend/src/config`，主要配置项包括：
- API 基础 URL
- 主题配置
- 路由配置

## 未来规划

1. **防破解对抗** - 增强授权系统的安全性，防止破解和盗版
2. **AI 风控** - 利用 AI 技术识别异常使用行为，提高系统安全性
4. **API 文档** - 完善 API 文档，便于集成
5. **监控系统** - 增加系统监控和告警功能

## 参考资料

1. [go-redis 指南](https://redis.ac.cn/docs/latest/develop/clients/go/)
2. [GORM 指南](https://gorm.io/zh_CN/docs/index.html#特性)
3. [IDE & AI 编程助手 | TRAE](https://www.trae.cn/)
4. [Trial Blog](http://blog.trialpro.top/)
5. [PostgreSQL: The world's most advanced open source database](https://www.postgresql.org/)
6. [Vue 3 官方文档](https://v3.vuejs.org/)
7. [Element Plus 文档](https://element-plus.org/zh-CN/)

## 贡献

欢迎提交 Issue 和 Pull Request 来帮助改进这个项目！

