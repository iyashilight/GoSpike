# GoSpike · 秒杀电商系统

基于 [go-zero](https://github.com/zeromicro/go-zero) 的秒杀（Flash Sale）电商示例项目。

后端由 **1 个 REST 网关 + 4 个 gRPC 微服务**组成，前端是 **Vue 3 + TypeScript** 写成的运营管理后台。

秒杀链路做了完整的削峰与防超卖处理：**Redis Lua 原子扣库存 → RabbitMQ 异步下单 → MySQL 条件扣减兜底**，并配合用户级限流、JWT 鉴权和统一错误码。

---

## 功能一览

- 用户注册 / 登录 / Token 刷新 / 资料修改（JWT + 密码哈希）
- 商品增删改查、商品列表分页
- 订单创建、查询、列表、取消（支持订单状态流转）
- 秒杀活动创建、库存预热、秒杀下单
- 秒杀防超卖：Lua 脚本原子扣减 + MySQL `stock >= ?` 条件扣减
- 秒杀限流：基于 go-zero `PeriodLimit`，默认 1 秒 100 次
- 统一响应体与业务错误码，panic 自动恢复
- Docker Compose 一键启动中间件与全部服务

---

## 系统架构

```
                            ┌──────────────────────────┐
   Vue 3 管理后台            │  http://localhost:3000   │
   (Vite 代理 /api → 8888)   └────────────┬─────────────┘
                                         │ HTTP / JWT
                            ┌────────────▼─────────────┐
                            │  api  网关 (go-zero REST) │  8888
                            │  鉴权 / 限流 / 统一响应    │
                            └───┬────┬────┬────┬────────┘
                     gRPC       │    │    │    │
              ┌─────────────────┘    │    │    └─────────────────┐
              │                      │    │                      │
      ┌───────▼──────┐   ┌───────────▼──┐ │   ┌──────────────────▼─────┐
      │ user   8081  │   │ product 8082 │ │   │ flashsale 8084          │
      │ 用户服务      │   │ 商品服务      │ │   │ 秒杀服务（Lua 扣库存）   │
      └───────┬──────┘   └───────────▲──┘ │   └──────────┬──────────────┘
              │                      │    │              │ 投递秒杀消息
              │              ┌───────┴────▼────┐   ┌─────▼──────┐
              │              │ order     8083  │◀──│ RabbitMQ    │
              │              │ 订单服务         │   │ order_flash │
              │              │ （MQ 消费者）    │   │ _purchase   │
              │              └───────┬─────────┘   └────────────┘
              │                      │
        ┌─────▼──────────────────────▼─────┐      ┌──────────────┐
        │            MySQL 8.0             │      │   Redis 7    │
        │   user / product / orders / ...  │      │ 库存、限流缓存 │
        └──────────────────────────────────┘      └──────────────┘
```

| 服务 | 类型 | 端口 | 说明 |
| --- | --- | --- | --- |
| `api` | REST 网关 | 8888 | 对外统一入口，负责 JWT 鉴权、限流、参数校验、错误码转换 |
| `user` | gRPC | 8081 | 注册、登录、用户信息 |
| `product` | gRPC | 8082 | 商品 CRUD、库存扣减 |
| `order` | gRPC | 8083 | 订单创建 / 查询 / 取消，消费秒杀 MQ 消息 |
| `flashsale` | gRPC | 8084 | 秒杀活动、库存预热、Lua 原子扣库存并投递消息 |
| `frontend` | Vue 3 SPA | 3000 | 运营管理后台（开发模式） |

基础设施（Docker）：

| 组件 | 宿主机端口 | 说明 |
| --- | --- | --- |
| MySQL 8.0 | 3307 | 库名 `gospike`，首次启动自动执行 `backend/init/sql` |
| Redis 7 | 6380 | 秒杀库存、活动信息、限流计数 |
| RabbitMQ 4 (management) | 5673 / 15673 | 队列 `order_flash_purchase`，管理台 `/15673` |

---

## 秒杀流程

1. 运营创建商品，并调用 `POST /api/flashsale/create-activity` 创建秒杀活动（起止时间）。
2. 调用 `POST /api/flashsale/cache-stock` 把库存**预热**到 Redis（`stock:{product_id}` + `activity:{product_id}`）。
3. 用户调用 `POST /api/flashsale/purchase`，网关先做用户级限流。
4. `flashsale` 执行 Lua 脚本，在 Redis 中以**服务器时间**校验活动窗口，并原子地判断 + 扣减库存（`DECRBY`），彻底避免超卖。
5. 扣减成功后投递消息到 RabbitMQ `order_flash_purchase`。
6. `order` 消费者创建订单（`pending`）→ 调用商品服务扣 MySQL 库存（`UPDATE ... WHERE stock >= ?` 兜底）→ 订单置为 `paid`；任一步失败则置为 `failed` 并重新入队。

订单状态：`pending`（已创建）→ `paid`（库存扣减成功）/ `failed`（库存不足等异常）。

---

## 技术栈

**后端**

- Go 1.25 · go-zero 1.10（`rest` / `zrpc` / `logx` / `sqlx` / `redis`）
- gRPC + Protocol Buffers
- MySQL 8.0（go-zero `sqlx`，含事务与 `WITH` 条件扣减）
- Redis 7（Lua 脚本、`PeriodLimit` 限流）
- RabbitMQ（`amqp091-go`，手动 Ack / Nack 重试）
- JWT 鉴权

**前端**

- Vue 3（`<script setup>`）+ TypeScript + Vite
- Pinia 状态管理、Vue Router
- Arco Design Vue 组件库 + Tailwind CSS
- Axios（统一拦截器：注入 Token、401 跳登录、错误提示）

---

## 快速开始

### 1. 克隆

```bash
git clone git@github.com:iyashilight/GoSpike.git
cd GoSpike
```

### 2. 准备环境变量

敏感配置（数据库密码、JWT 密钥、MQ 账号）不写入仓库，统一从 `backend/.env` 读取：

```bash
cd backend
cp .env.example .env
# 然后编辑 .env，把 MYSQL_ROOT_PASSWORD / JWT_SECRET 换成自己的值
```

| 变量 | 用途 |
| --- | --- |
| `MYSQL_ROOT_PASSWORD` | MySQL root 密码，同时用于拼装 `MYSQL_DSN` |
| `MYSQL_DSN` | 数据库连接串（compose 中由上一个变量自动拼装） |
| `JWT_SECRET` | JWT 签名密钥，`api` 与 `user` 共用 |
| `RABBITMQ_USER` / `RABBITMQ_PASS` | RabbitMQ 账号，缺省 `guest/guest` |

> 各服务的 `etc/*.yaml` 使用 `${MYSQL_DSN}`、`${JWT_SECRET}` 等占位符，go-zero 在加载配置时会用环境变量展开。

### 3. 启动

**方式一：全部跑在 Docker 里**

```bash
cd backend
docker compose up -d --build
```

网关地址：`http://localhost:8888/api/health`

**方式二：中间件用 Docker，服务本地跑（便于调试）**

```bash
cd backend
docker compose up -d mysql redis rabbitmq

set -a && source .env && set +a   # 导出环境变量

go run ./user      -f user/etc/user.yaml
go run ./product   -f product/etc/product.yaml
go run ./order     -f order/etc/order.yaml
go run ./flashsale -f flashsale/etc/flashsale.yaml
go run ./api       -f api/etc/api-api.yaml
```

### 4. 启动前端

```bash
cd frontend
npm install
npm run dev
```

打开 `http://localhost:3000`，开发服务器会把 `/api` 代理到 `http://localhost:8888`。

---

## API 一览

统一前缀 `/api`，需要鉴权的接口使用请求头 `Authorization: Bearer <token>`。

| 方法 | 路径 | 鉴权 | 说明 |
| --- | --- | :---: | --- |
| GET | `/api/health` | | 健康检查 |
| POST | `/api/user/register` | | 注册 |
| POST | `/api/user/login` | | 登录，返回 JWT |
| POST | `/api/user/refresh` | | 刷新 Token |
| GET | `/api/product/list` | | 商品列表 |
| GET | `/api/user/:id` | ✅ | 用户详情 |
| PUT | `/api/user/:id` | ✅ | 更新用户 |
| POST | `/api/product/create` | ✅ | 创建商品 |
| GET | `/api/product/:id` | ✅ | 商品详情 |
| PUT | `/api/product/:id` | ✅ | 更新商品 |
| DELETE | `/api/product/:id` | ✅ | 删除商品 |
| POST | `/api/order/create` | ✅ | 创建订单 |
| GET | `/api/order/list` | ✅ | 订单列表 |
| GET | `/api/order/:order_id` | ✅ | 订单详情 |
| POST | `/api/order/:order_id/cancel` | ✅ | 取消订单 |
| POST | `/api/flashsale/create-activity` | ✅ | 创建秒杀活动 |
| GET | `/api/flashsale/activity/:product_id` | ✅ | 查询秒杀活动 |
| POST | `/api/flashsale/cache-stock` | ✅ | 预热库存到 Redis |
| POST | `/api/flashsale/purchase` | ✅ | 秒杀下单（带限流） |

接口定义源码见 `backend/api/api.api`。

### 统一响应格式

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

`code` 为 `0` 表示成功，其余为业务错误码（见 `backend/api/errorx/errorx.go`）：

| 错误码 | 含义 |
| --- | --- |
| 1001 - 1003 | 用户不存在 / 凭证错误 / 用户已存在 |
| 2001 - 2006 | 商品不存在 / 库存不足 / 商品增删改、库存缓存失败 |
| 3001 - 3002 | 订单不存在 / 创建订单失败 |
| 9999 | 未知错误 |

限流触发时返回 HTTP `429`，并带 `Retry-After` 头。

---

## 目录结构

```
GoSpike
├── backend
│   ├── api/            # REST 网关：handler / logic / middleware / response / errorx
│   ├── user/           # 用户 gRPC 服务
│   ├── product/        # 商品 gRPC 服务
│   ├── order/          # 订单 gRPC 服务 + RabbitMQ 消费者
│   ├── flashsale/      # 秒杀 gRPC 服务（Lua 扣库存）
│   ├── init/sql/       # 建表脚本，MySQL 容器首次启动自动执行
│   ├── docker-compose.yml
│   ├── .env.example    # 环境变量模板
│   └── go.work         # Go workspace，聚合 5 个模块
└── frontend
    └── src
        ├── api/        # 按业务拆分的请求封装
        ├── pages/      # 登录、仪表盘、商品、订单、秒杀、个人中心
        ├── layouts/    # 后台布局
        ├── router/     # 路由与登录守卫
        └── utils/      # axios 实例与拦截器
```

每个 Go 服务内部采用 go-zero 的分层：`internal/{config,svc,logic,handler,server,model,mqs}`，protobuf 定义在同名子目录（如 `user/user.proto`）。

---

## 注意事项

- **不要**把 `.env` 提交到仓库；仓库里只保留 `.env.example`。
- 上线前请更换 `JWT_SECRET` 与数据库密码，并把 `Mode` 从 `dev` 改为生产模式。
- `backend/bin/`、`backend/logs/`、`frontend/node_modules/`、`frontend/dist/` 均已在 `.gitignore` 中忽略。
- RabbitMQ 管理台默认地址：`http://localhost:15673`（`guest/guest`）。

---

## License

本项目仅用于学习与演示。
