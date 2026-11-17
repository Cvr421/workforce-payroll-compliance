import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: () => import('../pages/LoginPage.vue'), meta: { public: true } },
    { path: '/403', component: () => import('../pages/StatusPage.vue'), props: { status: 403, title: 'Access denied', message: 'Your role cannot access this page.' }, meta: { public: true } },
    {
      path: '/',
      component: () => import('../layouts/DefaultLayout.vue'),
      meta: { roles: ['admin', 'hr_admin', 'payroll_manager', 'compliance_officer', 'auditor'] },
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'hr_admin', 'payroll_manager', 'compliance_officer', 'auditor'] } },
        { path: 'organizations', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'hr_admin', 'payroll_manager', 'compliance_officer', 'auditor'] } },
        { path: 'organizations/:id', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'hr_admin', 'payroll_manager', 'compliance_officer', 'auditor'] } },
        { path: 'employees', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'hr_admin', 'payroll_manager', 'compliance_officer', 'auditor'] } },
        { path: 'employees/:id', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'hr_admin', 'payroll_manager', 'compliance_officer', 'auditor'] } },
        { path: 'attendance', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'hr_admin', 'payroll_manager', 'compliance_officer', 'auditor'] } },
        { path: 'payroll', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'payroll_manager', 'auditor'] } },
        { path: 'payroll/runs/:id', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'payroll_manager', 'auditor'] } },
        { path: 'compliance', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'compliance_officer', 'auditor'] } },
        { path: 'compliance/tasks/:id', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'compliance_officer', 'auditor'] } },
        { path: 'compliance/calendar', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'compliance_officer', 'auditor'] } },
        { path: 'compliance/obligations', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'compliance_officer', 'auditor'] } },
        { path: 'reports', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'payroll_manager', 'compliance_officer', 'auditor'] } },
        { path: 'notifications', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'hr_admin', 'payroll_manager', 'compliance_officer', 'auditor'] } },
        { path: 'profile', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'hr_admin', 'payroll_manager', 'compliance_officer', 'auditor'] } },
        { path: 'admin/users', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin'] } },
        { path: 'audits', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'compliance_officer', 'auditor'] } },
        { path: 'audits/:id', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'compliance_officer', 'auditor'] } },
        { path: 'findings/:id', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'compliance_officer', 'auditor'] } },
        { path: 'audit/findings', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'compliance_officer', 'auditor'] } },
        { path: 'audit/log', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'auditor'] } },
        { path: 'admin/audit-log', component: () => import('../components/FeaturePage.vue'), meta: { roles: ['admin', 'auditor'] } },
      ],
    },
    { path: '/:pathMatch(.*)*', component: () => import('../pages/StatusPage.vue'), props: { status: 404, title: 'Page not found', message: 'The page you requested does not exist.' }, meta: { public: true } },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (to.meta.public) {
    if (to.path === '/login' && auth.isAuthenticated) return '/'
    return true
  }
  if (!auth.user && auth.token) {
    try { await auth.loadUser() } catch { auth.logout() }
  }
  if (!auth.isAuthenticated) return { path: '/login', query: { redirect: to.fullPath } }
  const roles = to.meta.roles as string[] | undefined
  if (roles && auth.user && !roles.includes(auth.user.role)) return '/403'
  return true
})
