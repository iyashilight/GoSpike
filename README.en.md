# GoSpike · Flash-Sale E-Commerce System

> [中文](README.md) | **English**

A flash-sale (seckill) e-commerce demo built on [go-zero](https://github.com/zeromicro/go-zero).

The backend consists of **1 REST gateway + 4 gRPC microservices**; the frontend is an admin console written in **Vue 3 + TypeScript**.

The flash-sale path is built for peak shaving and oversell prevention: **atomic stock decrement in Redis via Lua → asynchronous order creation through RabbitMQ → conditional decrement in MySQL as the last line of defense**, combined with per-user rate limiting, JWT authentication and unified error codes.

---

## Features

- User registration / login / token refresh / profile update (JWT + password hashing)
- Product CRUD and paginated product listing
- Order creation, lookup, listing and cancellation (with status transitions)
- Flash-sale activity creation, stock pre-warming, flash-sale purchase
- Oversell prevention: atomic Lua decrement + MySQL `stock >= ?` conditional update
- Flash-sale rate limiting: go-zero `PeriodLimit`, 100 requests per second by default
- Unified response envelope and business error codes, automatic panic recovery
- One-command startup of middleware and all services via Docker Compose

---

## Architecture

```
                            ┌──────────────────────────┐
   Vue 3 admin console       │  http://localhost:3000   │
   (Vite proxy /api → 8888)  └────────────┬─────────────┘
                                         │ HTTP / JWT
                            ┌────────────▼─────────────┐
                            │  api gateway (go-zero)   │  8888
                            │  auth / rate limit / DTO │
                            └───┬────┬────┬────┬────────┘
                     gRPC       │    │    │    │
              ┌─────────────────┘    │    │    └─────────────────┐
              │                      │    │                      │
      ┌───────▼──────┐   ┌───────────▼──┐ │   ┌──────────────────▼─────┐
      │ user   8081  │   │ product 8082 │ │   │ flashsale 8084          │
      │ user service │   │ product svc  │ │   │ flash sale (Lua stock)  │
      └───────┬──────┘   └───────────▲──┘ │   └──────────┬──────────────┘
              │                      │    │              │ publish message
              │              ┌───────┴────▼────┐   ┌─────▼──────┐
              │              │ order     8083  │◀──│ RabbitMQ    │
              │              │ order service   │   │ order_flash │
              │              │ (MQ consumer)   │   │ _purchase   │
              │              └───────┬─────────┘   └────────────┘
              │                      │
        ┌─────▼──────────────────────▼─────┐      ┌──────────────┐
        │            MySQL 8.0             │      │   Redis 7    │
        │   user / product / orders / ...  │      │ stock, limits│
        └──────────────────────────────────┘      └──────────────┘
```

| Service | Type | Port | Description |
| --- | --- | --- | --- |
| `api` | REST gateway | 8888 | Single public entry point: JWT auth, rate limiting, validation, error mapping |
| `user` | gRPC | 8081 | Registration, login, user profile |
| `product` | gRPC | 8082 | Product CRUD, stock decrement |
| `order` | gRPC | 8083 | Order create / query / cancel, flash-sale MQ consumer |
| `flashsale` | gRPC | 8084 | Activities, stock pre-warming, atomic Lua decrement, message publishing |
| `frontend` | Vue 3 SPA | 3000 | Admin console (dev server) |

Infrastructure (Docker):

| Component | Host port | Description |
| --- | --- | --- |
| MySQL 8.0 | 3307 | Database `gospike`; `backend/init/sql` runs on first boot |
| Redis 7 | 6380 | Flash-sale stock, activity data, rate-limit counters |
| RabbitMQ 4 (management) | 5673 / 15673 | Queue `order_flash_purchase`, admin UI on `/15673` |

---

## Flash-sale flow

1. An operator creates a product and calls `POST /api/flashsale/create-activity` to define the activity window.
2. `POST /api/flashsale/cache-stock` **pre-warms** the stock into Redis (`stock:{product_id}` + `activity:{product_id}`).
3. A user calls `POST /api/flashsale/purchase`; the gateway applies per-user rate limiting first.
4. `flashsale` runs a Lua script that validates the activity window against **Redis server time** and atomically checks + decrements stock (`DECRBY`), making overselling impossible.
5. On success it publishes a message to the RabbitMQ queue `order_flash_purchase`.
6. The `order` consumer creates the order (`pending`) → calls the product service to decrement MySQL stock (`UPDATE ... WHERE stock >= ?`) → marks the order `paid`; if any step fails the order is marked `failed` and the message is requeued.

Order states: `pending` (created) → `paid` (stock decremented) / `failed` (e.g. insufficient stock).

---

## Tech stack

**Backend**

- Go 1.25 · go-zero 1.10 (`rest` / `zrpc` / `logx` / `sqlx` / `redis`)
- gRPC + Protocol Buffers
- MySQL 8.0 (go-zero `sqlx`, transactions, conditional decrement)
- Redis 7 (Lua scripts, `PeriodLimit` rate limiting)
- RabbitMQ (`amqp091-go`, manual Ack / Nack retry)
- JWT authentication

**Frontend**

- Vue 3 (`<script setup>`) + TypeScript + Vite
- Pinia state management, Vue Router
- Arco Design Vue + Tailwind CSS
- Axios with interceptors (token injection, 401 → login, error toasts)

---

## Getting started

### 1. Clone

```bash
git clone git@github.com:iyashilight/GoSpike.git
cd GoSpike
```

### 2. Prepare environment variables

Secrets (database password, JWT key, MQ credentials) are never committed and are read from `backend/.env`:

```bash
cd backend
cp .env.example .env
# then edit .env and set your own MYSQL_ROOT_PASSWORD / JWT_SECRET
```

| Variable | Purpose |
| --- | --- |
| `MYSQL_ROOT_PASSWORD` | MySQL root password, also used to build `MYSQL_DSN` |
| `MYSQL_DSN` | Database DSN (assembled automatically by Compose) |
| `JWT_SECRET` | JWT signing key, shared by `api` and `user` |
| `RABBITMQ_USER` / `RABBITMQ_PASS` | RabbitMQ credentials, defaults to `guest/guest` |

> Each service's `etc/*.yaml` uses placeholders such as `${MYSQL_DSN}` and `${JWT_SECRET}`; go-zero expands them from the environment when loading config.

### 3. Run

**Option A: everything in Docker**

```bash
cd backend
docker compose up -d --build
```

Gateway health check: `http://localhost:8888/api/health`

**Option B: middleware in Docker, services locally (better for debugging)**

```bash
cd backend
docker compose up -d mysql redis rabbitmq

set -a && source .env && set +a   # export env vars

go run ./user      -f user/etc/user.yaml
go run ./product   -f product/etc/product.yaml
go run ./order     -f order/etc/order.yaml
go run ./flashsale -f flashsale/etc/flashsale.yaml
go run ./api       -f api/etc/api-api.yaml
```

### 4. Run the frontend

```bash
cd frontend
npm install
npm run dev
```

Open `http://localhost:3000`; the dev server proxies `/api` to `http://localhost:8888`.

---

## API reference

All routes are prefixed with `/api`. Protected routes require `Authorization: Bearer <token>`.

| Method | Path | Auth | Description |
| --- | --- | :---: | --- |
| GET | `/api/health` | | Health check |
| POST | `/api/user/register` | | Register |
| POST | `/api/user/login` | | Login, returns JWT |
| POST | `/api/user/refresh` | | Refresh token |
| GET | `/api/product/list` | | Product list |
| GET | `/api/user/:id` | ✅ | User detail |
| PUT | `/api/user/:id` | ✅ | Update user |
| POST | `/api/product/create` | ✅ | Create product |
| GET | `/api/product/:id` | ✅ | Product detail |
| PUT | `/api/product/:id` | ✅ | Update product |
| DELETE | `/api/product/:id` | ✅ | Delete product |
| POST | `/api/order/create` | ✅ | Create order |
| GET | `/api/order/list` | ✅ | Order list |
| GET | `/api/order/:order_id` | ✅ | Order detail |
| POST | `/api/order/:order_id/cancel` | ✅ | Cancel order |
| POST | `/api/flashsale/create-activity` | ✅ | Create flash-sale activity |
| GET | `/api/flashsale/activity/:product_id` | ✅ | Get flash-sale activity |
| POST | `/api/flashsale/cache-stock` | ✅ | Pre-warm stock into Redis |
| POST | `/api/flashsale/purchase` | ✅ | Flash-sale purchase (rate limited) |

The API definition lives in `backend/api/api.api`.

### Response envelope

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

`code` `0` means success; anything else is a business error code (see `backend/api/errorx/errorx.go`):

| Code | Meaning |
| --- | --- |
| 1001 - 1003 | User not found / invalid credentials / user already exists |
| 2001 - 2006 | Product not found / insufficient stock / product CRUD or stock-cache failures |
| 3001 - 3002 | Order not found / failed to create order |
| 9999 | Unknown error |

Rate limiting returns HTTP `429` with a `Retry-After` header.

---

## Project layout

```
GoSpike
├── backend
│   ├── api/            # REST gateway: handler / logic / middleware / response / errorx
│   ├── user/           # user gRPC service
│   ├── product/        # product gRPC service
│   ├── order/          # order gRPC service + RabbitMQ consumer
│   ├── flashsale/      # flash-sale gRPC service (Lua stock decrement)
│   ├── init/sql/       # schema, executed automatically on first MySQL boot
│   ├── docker-compose.yml
│   ├── .env.example    # environment template
│   └── go.work         # Go workspace aggregating the 5 modules
└── frontend
    └── src
        ├── api/        # request wrappers per domain
        ├── pages/      # login, dashboard, product, order, flash sale, profile
        ├── layouts/    # admin layout
        ├── router/     # routes and auth guard
        └── utils/      # axios instance and interceptors
```

Every Go service follows go-zero layering: `internal/{config,svc,logic,handler,server,model,mqs}`, with protobuf definitions in a same-named subdirectory (e.g. `user/user.proto`).

---

## Notes

- **Never** commit `.env`; only `.env.example` is tracked.
- Before going live, rotate `JWT_SECRET` and the database password, and switch `Mode` from `dev` to production.
- `backend/bin/`, `backend/logs/`, `frontend/node_modules/` and `frontend/dist/` are already ignored by `.gitignore`.
- RabbitMQ management UI: `http://localhost:15673` (`guest/guest`).

---

## License

Released under [GPL-3.0](LICENSE). You are free to study, modify and redistribute it; derivative works must also be open-sourced under GPL-3.0.
