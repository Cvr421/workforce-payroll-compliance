<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { apiFetch } from '../api/client'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useNotifications } from '../stores/notifications'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const notifications = useNotifications()
const loading = ref(true)
const error = ref('')
const data = ref<Record<string, unknown> | null>(null)
const rows = ref<Record<string, unknown>[]>([])

const isDashboard = computed(() => route.path === '/' || route.path === '/dashboard')
const isOrganizations = computed(() => route.path === '/organizations')
const isEmployees = computed(() => route.path === '/employees')
const isAttendance = computed(() => route.path === '/attendance')
const isPayroll = computed(() => route.path === '/payroll')
const isCompliance = computed(() => route.path === '/compliance')
const isComplianceCalendar = computed(() => route.path === '/compliance/calendar')
const isComplianceObligations = computed(() => route.path === '/compliance/obligations')
const isReports = computed(() => route.path === '/reports')
const isNotifications = computed(() => route.path === '/notifications')
const isProfile = computed(() => route.path === '/profile')
const isUsers = computed(() => route.path === '/admin/users')
const isAudits = computed(() => route.path === '/audits' || route.path === '/audit/findings')
const isOrganizationDetail = computed(() => route.path.startsWith('/organizations/'))
const isEmployeeDetail = computed(() => route.path.startsWith('/employees/'))
const isPayrollDetail = computed(() => route.path.startsWith('/payroll/runs/'))
const isComplianceDetail = computed(() => route.path.startsWith('/compliance/tasks/'))
const isAuditDetail = computed(() => route.path.startsWith('/audits/'))
const isFindingDetail = computed(() => route.path.startsWith('/findings/'))

const title = computed(() => {
  const map: Record<string, string> = {
    '/dashboard': 'Dashboard',
    '/organizations': 'Organizations',
    '/employees': 'Employees',
    '/attendance': 'Attendance',
    '/payroll': 'Payroll',
    '/compliance': 'Compliance',
    '/reports': 'Reports',
    '/notifications': 'Notifications',
    '/profile': 'Profile',
    '/admin/users': 'User Administration',
    '/audits': 'Audits',
  }
  for (const [key, val] of Object.entries(map)) {
    if (route.path.startsWith(key)) return val
  }
  return 'Page'
})

const subtitle = computed(() => {
  const map: Record<string, string> = {
    '/dashboard': 'Overview of your compliance workspace',
    '/organizations': 'Manage organizations and their settings',
    '/employees': 'Employee directory and records',
    '/attendance': 'Time tracking and leave management',
    '/payroll': 'Payroll runs and calculations',
    '/compliance': 'Compliance tasks and obligations',
    '/reports': 'Generated reports and exports',
    '/notifications': 'System and compliance notifications',
    '/profile': 'Your account and preferences',
    '/admin/users': 'Manage workspace users and roles',
    '/audits': 'Audit planning and evidence',
  }
  for (const [key, val] of Object.entries(map)) {
    if (route.path.startsWith(key)) return val
  }
  return 'Feature page under construction'
})

const endpoints: Record<string, string> = {
  '/organizations': '/organizations',
  '/': '/dashboard/summary',
  '/dashboard': '/dashboard/summary',
  '/employees': '/employees',
  '/attendance': '/attendance?year=2026&month=10',
  '/payroll': '/payroll/runs',
  '/compliance': '/compliance/tasks',
  '/compliance/obligations': '/compliance/obligations',
  '/reports': '/reports/headcount',
  '/audits': '/audits',
  '/audit/findings': '/audits',
  '/audit/log': '/admin/audit-log',
  '/admin/audit-log': '/admin/audit-log',
  '/admin/users': '/users',
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    let path = ''
    if (isDashboard.value) path = '/dashboard/summary'
    else if (isOrganizationDetail.value) path = `/organizations/${route.params.id}`
    else if (isEmployeeDetail.value) path = `/employees/${route.params.id}`
    else if (isPayrollDetail.value) path = `/payroll/runs/${route.params.id}`
    else if (isComplianceDetail.value) path = `/compliance/tasks/${route.params.id}`
    else if (isAuditDetail.value) path = `/audits/${route.params.id}`
    else if (isFindingDetail.value) path = `/findings/${route.params.id}`
    else path = endpoints[route.path] || ''

    if (path) {
      const res = await apiFetch<{ data: unknown }>(path)
      data.value = (isDashboard.value || isOrganizationDetail.value || isEmployeeDetail.value || isPayrollDetail.value || isComplianceDetail.value || isAuditDetail.value || isFindingDetail.value) ? res.data as Record<string, unknown> : null
      rows.value = Array.isArray(res.data) ? res.data as Record<string, unknown>[] :
        Array.isArray((res.data as { employees?: unknown[] })?.employees) ? (res.data as { employees: Record<string, unknown>[] }).employees : []
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Unable to load this page'
  } finally { loading.value = false }
}

onMounted(() => { void load() })
watch(() => route.path, () => { void load() })
</script>

<template>
  <div>
    <div class="mb-8">
      <p class="text-sm text-[#99949a]" v-if="subtitle">{{ subtitle }}</p>
      <h1 class="mt-1 text-[30px] font-semibold tracking-[-.04em]">{{ title }}</h1>
    </div>

    <div v-if="loading" class="flex h-64 items-center justify-center text-[#99949a]">Loading…</div>

    <div v-else-if="error" class="rounded-xl border border-[#fecaca] bg-[#fff0ed] p-6 text-red-700 dark:border-red-900 dark:bg-red-900/20">{{ error }}</div>

    <div v-else-if="isDashboard && data" class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
      <div class="rounded-2xl border border-[#e2dfda] bg-white p-6 dark:border-white/10 dark:bg-white/5">
        <p class="text-sm text-[#99949a]">Organizations</p>
        <p class="mt-2 text-3xl font-semibold">{{ (data as any).organizations_count ?? 0 }}</p>
      </div>
      <div class="rounded-2xl border border-[#e2dfda] bg-white p-6 dark:border-white/10 dark:bg-white/5">
        <p class="text-sm text-[#99949a]">Employees</p>
        <p class="mt-2 text-3xl font-semibold">{{ (data as any).employees_count ?? 0 }}</p>
      </div>
      <div class="rounded-2xl border border-[#e2dfda] bg-white p-6 dark:border-white/10 dark:bg-white/5">
        <p class="text-sm text-[#99949a]">Open tasks</p>
        <p class="mt-2 text-3xl font-semibold">{{ (data as any).open_tasks_count ?? 0 }}</p>
      </div>
      <div class="rounded-2xl border border-[#e2dfda] bg-white p-6 dark:border-white/10 dark:bg-white/5">
        <p class="text-sm text-[#99949a]">Compliance score</p>
        <p class="mt-2 text-3xl font-semibold">{{ (data as any).compliance_score ?? '--' }}%</p>
      </div>
    </div>

    <div v-else-if="rows.length" class="space-y-4">
      <div v-for="row in rows" :key="row.id || row.employee_id || row.organization_id" class="rounded-xl border border-[#e2dfda] bg-white p-4 dark:border-white/10 dark:bg-white/5">
        <div class="flex items-center justify-between">
          <div>
            <p class="font-medium">{{ row.name || row.full_name || row.organization_name || row.title || row.code || 'Item' }}</p>
            <p class="text-xs text-[#99949a]">{{ row.email || row.category || row.status || '' }}</p>
          </div>
          <span class="rounded-full px-2 py-1 text-xs" v-if="row.status" :class="row.status === 'OPEN' ? 'bg-[#e6f4ea] text-[#1e7e34]' : row.status === 'IN_PROGRESS' ? 'bg-[#e8f0fe] text-[#1a73e8]' : row.status === 'IN_REVIEW' ? 'bg-[#fef7e0] text-[#b06000]' : 'bg-[#fce8e6] text-[#c76052]'">{{ row.status }}</span>
        </div>
      </div>
    </div>

    <div v-else class="rounded-2xl border border-[#e2dfda] bg-white p-8 dark:border-white/10 dark:bg-white/5">
      <div class="flex items-center justify-center min-h-[300px] text-[#99949a]">
        <div class="text-center">
          <svg class="mx-auto mb-4 text-[#aaa5a0]" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" /><path d="M9 12l2 2 4-4" /></svg>
          <p class="text-lg font-medium">Page under construction</p>
          <p class="mt-1 text-sm">{{ subtitle }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
