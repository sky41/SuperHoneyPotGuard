# SuperHoneyPotGuard 用户管理功能

## 提交目的
开发用户管理功能，包括用户注册、登录、注销、权限管理等核心功能。

## 提交内容
- 数据库表结构设计（用户表、角色表、权限表、关联表、操作日志表）
- 后端 API 服务（Go + Gin）
  - 认证接口（注册、登录、注销、获取当前用户）
  - 用户管理接口（列表、详情、创建、编辑、删除、状态更新、重置密码）
  - 角色管理接口（列表、详情、创建、编辑、删除）
  - 权限管理接口（树形结构、列表、详情、创建、编辑、删除）
- 前端 Web 应用（Vue3 + Ant Design Vue）
  - 登录页面
  - 主布局（侧边栏导航、顶部栏）
  - 首页仪表盘
  - 用户管理页面
- 中间件（认证、权限、日志、限流）
- 工具类（日志、JWT、密码加密、HTTP 请求）

## 提交时间
2026-01-17

## 项目结构

```
SuperHoneyPotGuard/
├── db/
│   ├── user.sql
│   └── user_management.sql
├── code/
│   ├── api/
│   │   ├── main.go
│   │   ├── go.mod
│   │   ├── .env.example
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── database/
│   │   │   └── database.go
│   │   ├── models/
│   │   │   └── models.go
│   │   ├── controllers/
│   │   │   ├── authController.go
│   │   │   ├── userController.go
│   │   │   ├── roleController.go
│   │   │   └── permissionController.go
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── permission.go
│   │   │   ├── log.go
│   │   │   └── rateLimit.go
│   │   ├── routes/
│   │   │   └── routes.go
│   │   └── utils/
│   │       ├── jwt.go
│   │       ├── password.go
│   │       └── response.go
│   ├── web/
│   │   ├── index.html
│   │   ├── package.json
│   │   ├── vite.config.js
│   │   └── src/
│   │       ├── api/
│   │       │   └── index.js
│   │       ├── layouts/
│   │       │   └── MainLayout.vue
│   │       ├── router/
│   │       │   └── index.js
│   │       ├── views/
│   │       │   ├── Login.vue
│   │       │   ├── Dashboard.vue
│   │       │   └── UserManage.vue
│   │       ├── utils/
│   │       │   └── request.js
│   │       ├── App.vue
│   │       ├── main.js
│   │       └── style.css
│   └── README.md
└── docs/
    └── 技术栈.md
```

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 16.x+
- MySQL 5.7+ / 8.0+

### 数据库配置

1. 创建数据库：
```sql
CREATE DATABASE superhoneypotguard DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 导入数据库表结构：
```bash
mysql -u root -p superhoneypotguard < db/user_management.sql
```

### 后端服务配置

1. 进入 API 目录：
```bash
cd code/api
```

2. 安装依赖：
```bash
go mod download
```

3. 配置环境变量：
```bash
cp .env.example .env
```

编辑 `.env` 文件，配置数据库连接信息：
```
PORT=3000
GIN_MODE=debug

DB_HOST=localhost
DB_PORT=3306
DB_NAME=superhoneypotguard
DB_USER=root
DB_PASSWORD=your_password

JWT_SECRET=your_jwt_secret_key_here
JWT_EXPIRES_IN=24h

BCRYPT_COST=10

RATE_LIMIT_WINDOW_MS=900000
RATE_LIMIT_MAX_REQUESTS=100

LOG_LEVEL=info
LOG_FILE_PATH=logs/

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

4. 启动后端服务：
```bash
go run main.go
```

或者编译后运行：
```bash
go build -o superhoneypotguard-api
./superhoneypotguard-api
```

### 前端应用配置

1. 进入 Web 目录：
```bash
cd code/web
```

2. 安装依赖：
```bash
npm install
```

3. 启动前端应用：
```bash
npm run dev
```

4. 访问应用：
打开浏览器访问 http://localhost:5173

### 默认账号

- 用户名：admin
- 密码：admin123

## API 文档

### 认证接口

- POST `/api/auth/register` - 用户注册
- POST `/api/auth/login` - 用户登录
- POST `/api/auth/logout` - 用户注销
- GET `/api/auth/current` - 获取当前用户信息

### 用户管理接口

- GET `/api/user/list` - 获取用户列表
- GET `/api/user/:id` - 获取用户详情
- POST `/api/user` - 创建用户
- PUT `/api/user/:id` - 更新用户
- DELETE `/api/user/:id` - 删除用户
- PATCH `/api/user/:id/status` - 更新用户状态
- POST `/api/user/:id/reset-password` - 重置用户密码

### 角色管理接口

- GET `/api/role/list` - 获取角色列表
- GET `/api/role/all` - 获取所有角色
- GET `/api/role/:id` - 获取角色详情
- POST `/api/role` - 创建角色
- PUT `/api/role/:id` - 更新角色
- DELETE `/api/role/:id` - 删除角色

### 权限管理接口

- GET `/api/permission/tree` - 获取权限树
- GET `/api/permission/all` - 获取所有权限
- GET `/api/permission/:id` - 获取权限详情
- POST `/api/permission` - 创建权限
- PUT `/api/permission/:id` - 更新权限
- DELETE `/api/permission/:id` - 删除权限

## 功能特性

### 用户管理
- 用户注册、登录、注销
- 用户列表查询（支持分页、搜索）
- 用户信息编辑
- 用户状态管理（启用/禁用）
- 用户密码重置
- 用户删除

### 角色管理
- 角色列表查询
- 角色创建、编辑、删除
- 角色权限分配

### 权限管理
- 权限树形结构展示
- 权限创建、编辑、删除
- 权限类型分类（菜单、按钮、接口）

### 安全特性
- JWT 身份认证
- BCrypt 密码加密
- 接口权限验证
- 请求频率限制
- 操作日志记录
- SQL 注入防护
- XSS 防护

## 技术栈

### 后端
- Go 1.21+
- Gin Web 框架
- GORM ORM
- JWT 认证
- BCrypt 密码加密
- MySQL 驱动
- 令牌桶限流

### 前端
- Vue 3.4+
- Ant Design Vue 4.0+
- Vue Router 4.2+
- Pinia 2.1+
- Axios 1.6+
- Vite 5.0+

### 数据库
- MySQL 5.7+ / 8.0+

详细技术栈说明请参考 [技术栈文档](../../docs/技术栈.md)

## 开发说明

### 提交代码注释规范
在每一次提交代码时，都需要添加注释，注释中需要包含以下内容：
- 提交的目的
- 提交的内容
- 提交的时间

### 笔记属性规范
在每一个创建笔记的文件中都要包含笔记属性，笔记属性包括：
- 笔记标题
- 笔记标签
- 笔记创建时间
- 笔记最后修改时间
