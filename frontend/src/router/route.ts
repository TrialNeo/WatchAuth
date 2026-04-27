/**
 * name: 路由名称, 也用于keepAlive缓存，需要与组件名称一致
 * meta.keepAlive: 是否需要缓存
 *
 */
export const staticRoutes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/login/index.vue'),
    meta: { keepAlive: false },
  },
  // 重定向路由(暂时注释掉，因为redirect路由会导致加载缓慢)
  // {
  //   path: '/redirect/:path(.*)',
  //   name: 'redirect',
  //   component: () => import('@/views/redirect/index.vue'),
  //   meta: { hidden: true },
  // },
  {
    path: '/',
    name: 'layout',
    component: () => import('@/layouts/index.vue'),
    children: [
      {
        path: '/profile',
        name: 'ProfileView',
        component: () => import('@/views/profile/index.vue'),
        meta: { title: '个人中心', icon: 'HOutline:UserCircleIcon', keepAlive: true },
      },
      {
        path: '/statistics',
        name: 'StatisticsView',
        component: () => import('@/views/dashboard/statistics/index.vue'),
        meta: { title: '数据统计', icon: 'HOutline:ChartBarSquareIcon', keepAlive: true },
      },
      {
        path: '/license',
        name: 'LicenseView',
        component: () => import('@/views/license/index.vue'),
        meta: { title: '授权管理', icon: 'HOutline:KeyIcon', keepAlive: true },
      },
      {
        path: '/agent',
        name: 'AgentView',
        component: () => import('@/views/agent/index.vue'),
        meta: { title: '代理管理', icon: 'HOutline:UserGroupIcon', keepAlive: true },
      },
      {
        path: '/system-config',
        name: 'SystemConfigView',
        component: () => import('@/views/systemConfig/index.vue'),
        meta: { title: '系统配置', icon: 'HOutline:Cog6ToothIcon', keepAlive: true },
      },
      {
        path: '/operation-log',
        name: 'OperationLogView',
        component: () => import('@/views/operationLog/index.vue'),
        meta: { title: '操作日志', icon: 'HOutline:ClipboardDocumentListIcon', keepAlive: true },
      },
      {
        path: '/announcement',
        name: 'AnnouncementView',
        component: () => import('@/views/announcement/index.vue'),
        meta: { title: '公告管理', icon: 'HOutline:MegaphoneIcon', keepAlive: true },
      },
      {
        path: '/exception/403',
        name: '403',
        component: () => import('@/views/exception/403/index.vue'),
        meta: { title: '403', icon: 'HOutline:NoSymbolIcon', keepAlive: true },
      },
      {
        path: '/:pathMatch(.*)*',
        name: '404',
        component: () => import('@/views/exception/404/index.vue'),
        meta: { title: '404', icon: 'HOutline:QuestionMarkCircleIcon', keepAlive: true },
      },
    ],
  },
]
