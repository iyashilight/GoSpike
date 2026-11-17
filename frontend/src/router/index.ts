import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { TOKEN_KEY } from '@/constants'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/auth/Login.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/pages/auth/Register.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/dashboard',
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/pages/dashboard/Index.vue'),
        meta: { title: '仪表盘', icon: 'dashboard' },
      },
      {
        path: 'product',
        name: 'Product',
        component: () => import('@/pages/product/Index.vue'),
        meta: { title: '商品管理', icon: 'shopping-cart' },
      },
      {
        path: 'product/create',
        name: 'ProductCreate',
        component: () => import('@/pages/product/Create.vue'),
        meta: { title: '创建商品', icon: 'plus' },
      },
      {
        path: 'product/:id/edit',
        name: 'ProductEdit',
        component: () => import('@/pages/product/Edit.vue'),
        meta: { title: '编辑商品', icon: 'edit' },
      },
      {
        path: 'order',
        name: 'Order',
        component: () => import('@/pages/order/Index.vue'),
        meta: { title: '订单管理', icon: 'file' },
      },
      {
        path: 'order/:orderId',
        name: 'OrderDetail',
        component: () => import('@/pages/order/Detail.vue'),
        meta: { title: '订单详情', icon: 'file' },
      },
      {
        path: 'flashsale',
        name: 'FlashSale',
        component: () => import('@/pages/flashsale/Index.vue'),
        meta: { title: '秒杀活动', icon: 'fire' },
      },
      {
        path: 'user/profile',
        name: 'UserProfile',
        component: () => import('@/pages/user/Profile.vue'),
        meta: { title: '个人中心', icon: 'user' },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/pages/NotFound.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem(TOKEN_KEY)
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth !== false)

  if (requiresAuth && !token) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if ((to.name === 'Login' || to.name === 'Register') && token) {
    next({ name: 'Dashboard' })
  } else {
    next()
  }
})

export default router
