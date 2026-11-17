# GoSpike 管理后台（前端）

GoSpike 秒杀系统的运营管理前端，基于 **Vue 3 + TypeScript + Vite**。

## 技术栈

- Vue 3（`<script setup>`）+ TypeScript
- Vite（开发服务器端口 `3000`，`/api` 代理到 `http://localhost:8888`）
- Vue Router（路由 + 登录守卫）
- Pinia 状态管理
- Arco Design Vue 组件库
- Tailwind CSS
- Axios（统一拦截器：注入 JWT、`401` 跳登录、`429` 提示限流）

## 开发

```bash
npm install
npm run dev      # http://localhost:3000
```

启动前请确保后端网关已在 `localhost:8888` 运行，参见仓库根目录 README。

## 构建

```bash
npm run build    # 类型检查 + 产出到 dist/
npm run preview  # 本地预览构建结果
```

## 目录结构

```
src
├── api/         # 按业务拆分的接口封装（user / product / order / flashsale / health）
├── components/  # 通用组件
├── constants/   # Token key、订单状态等常量
├── hooks/       # 组合式函数
├── layouts/     # AdminLayout 后台布局
├── pages/       # 页面：登录、仪表盘、商品、订单、秒杀、个人中心
├── router/      # 路由表与登录守卫
├── services/    # 业务服务层
├── store/       # Pinia store
├── types/       # 接口与数据模型类型定义
└── utils/       # axios 实例、格式化等工具
```

## 页面

| 路径 | 页面 |
| --- | --- |
| `/login`、`/register` | 登录 / 注册 |
| `/dashboard` | 仪表盘 |
| `/product`、`/product/create`、`/product/:id/edit` | 商品列表 / 创建 / 编辑 |
| `/order`、`/order/:orderId` | 订单列表 / 详情 |
| `/flashsale` | 秒杀活动管理 |
| `/user/profile` | 个人中心 |

除登录注册外，其余路由均需登录后访问。
