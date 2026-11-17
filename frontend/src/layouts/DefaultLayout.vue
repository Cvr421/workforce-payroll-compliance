<script setup lang="ts">
import { computed, ref } from 'vue'
import { useUiStore } from '../stores/ui'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import { Bell, ChevronDown, CircleHelp, FileCheck2, FileSearch, LayoutDashboard, Moon, Search, Sun, Users, WalletCards } from 'lucide-vue-next'
import { useNotifications } from '../stores/notifications'
import { RouterView } from 'vue-router'

const ui = useUiStore()
const auth = useAuthStore()
const router = useRouter()
const notifications = useNotifications()
const notificationItems = notifications.items
const showNotifications = ref(false)
const greeting = computed(() => new Date().getHours() < 12 ? 'Good morning' : 'Good afternoon')
const navigation = [
  { label: 'Dashboard', path: '/dashboard', icon: LayoutDashboard, roles: ['admin', 'hr_admin', 'payroll_manager', 'compliance_officer', 'auditor'] },
  { label: 'Organizations', path: '/organizations', icon: Users, roles: ['admin', 'hr_admin', 'payroll_manager', 'compliance_officer', 'auditor'] },
  { label: 'Employees', path: '/employees', icon: Users, roles: ['admin', 'hr_admin', 'payroll_manager', 'compliance_officer', 'auditor'] },
  { label: 'Attendance', path: '/attendance', icon: FileCheck2, roles: ['admin', 'hr_admin', 'payroll_manager', 'compliance_officer', 'auditor'] },
  { label: 'Payroll', path: '/payroll', icon: WalletCards, roles: ['admin', 'payroll_manager', 'auditor'] },
  { label: 'Compliance', path: '/compliance', icon: FileCheck2, roles: ['admin', 'compliance_officer', 'auditor'] },
  { label: 'Compliance calendar', path: '/compliance/calendar', icon: FileCheck2, roles: ['admin', 'compliance_officer', 'auditor'] },
  { label: 'Obligations', path: '/compliance/obligations', icon: FileCheck2, roles: ['admin', 'compliance_officer', 'auditor'] },
  { label: 'Reports', path: '/reports', icon: FileSearch, roles: ['admin', 'payroll_manager', 'compliance_officer', 'auditor'] },
  { label: 'Notifications', path: '/notifications', icon: Bell, roles: ['admin', 'hr_admin', 'payroll_manager', 'compliance_officer', 'auditor'] },
  { label: 'Profile', path: '/profile', icon: Users, roles: ['admin', 'hr_admin', 'payroll_manager', 'compliance_officer', 'auditor'] },
  { label: 'User administration', path: '/admin/users', icon: Users, roles: ['admin'] },
  { label: 'Audits', path: '/audits', icon: FileSearch, roles: ['admin', 'compliance_officer', 'auditor'] },
  { label: 'Audit log', path: '/audit/log', icon: FileSearch, roles: ['admin', 'auditor'] },
]
</script>

<template>
  <div :class="{ dark: ui.dark }" class="min-h-screen bg-[#f7f6f3] text-[#25232a] transition-colors dark:bg-[#18171b] dark:text-[#f3f0eb]">
    <aside class="fixed inset-y-0 left-0 z-10 flex w-[248px] flex-col border-r border-[#e7e4df] bg-[#fbfaf8] px-5 py-7 dark:border-white/10 dark:bg-[#201f24]">
      <div class="mb-12 flex items-center gap-3 px-2"><div class="grid h-9 w-9 place-items-center rounded-xl bg-[#a69cf4] font-bold text-white">C</div><span class="text-xl font-semibold tracking-tight">comply<span class="text-[#8f83ec]">lite</span></span></div>
      <p class="mb-3 px-3 text-[10px] font-bold uppercase tracking-[.2em] text-[#aaa5a0]">Workspace</p>
      <nav class="space-y-1">
        <button
          v-for="item in navigation.filter((entry) => entry.roles.includes(auth.user?.role ?? ''))"
          :key="item.label"
          @click="ui.activePage = item.label; router.push(item.path)"
          :class="ui.activePage === item.label ? 'bg-[#eeebff] text-[#665ad3] dark:bg-[#393350] dark:text-[#c7c1ff]' : 'text-[#77737a] hover:bg-[#f1efec] dark:hover:bg-white/5'"
          class="flex w-full items-center gap-3 rounded-xl px-3 py-3 text-sm font-medium transition"
        >
          <component :is="item.icon" :size="18" />{{ item.label }}
        </button>
      </nav>
      <div class="mt-auto rounded-2xl bg-[#f0eefb] p-4 dark:bg-[#353142]">
        <p class="text-xs font-semibold">Phase 0 workspace</p>
        <p class="mt-1 text-[11px] leading-5 text-[#77737a] dark:text-[#beb9c9]">The application shell is ready for feature pages.</p>
        <button class="mt-3 flex items-center gap-2 text-xs font-semibold text-[#665ad3]">Open guide <CircleHelp :size="14" /></button>
      </div>
      <div class="mt-5 flex items-center gap-3 border-t border-[#e7e4df] px-2 pt-5 dark:border-white/10">
        <div class="grid h-9 w-9 place-items-center rounded-full bg-[#dec8b6] text-sm font-bold text-[#604a38]">{{ auth.user?.full_name.slice(0, 2).toUpperCase() }}</div>
        <div class="min-w-0">
          <p class="truncate text-xs font-semibold">{{ auth.user?.full_name }}</p>
          <p class="text-[11px] text-[#99949a]">{{ auth.user?.role }}</p>
        </div>
        <button class="ml-auto text-[#99949a]" aria-label="Log out" @click="auth.logout(); router.push('/login')"><ChevronDown :size="15" /></button>
      </div>
    </aside>
    <main class="ml-[248px] px-10 py-8">
      <header class="mb-12 flex items-center justify-between">
        <div>
          <p class="text-sm text-[#99949a]">{{ greeting }}</p>
          <h1 class="mt-1 text-[30px] font-semibold tracking-[-.04em]">{{ ui.activePage }}</h1>
        </div>
        <div class="flex items-center gap-3">
          <div class="flex items-center gap-2 rounded-xl border border-[#e2dfda] bg-white px-3 py-2 text-sm text-[#a19da0] dark:border-white/10 dark:bg-white/5">
            <Search :size="16" />
            <span>Search anything</span>
            <kbd class="ml-8 rounded bg-[#f2f0ed] px-1.5 py-0.5 text-[10px] dark:bg-white/10">⌘ K</kbd>
          </div>
          <div class="relative">
            <button aria-label="Notifications" class="relative rounded-xl border border-[#e2dfda] bg-white p-2.5 text-[#77737a] dark:border-white/10 dark:bg-white/5" @click="showNotifications = !showNotifications">
              <Bell :size="18" />
              <span v-if="notifications.unreadCount" class="absolute -right-1 -top-1 grid min-h-4 min-w-4 place-items-center rounded-full bg-[#665ad3] px-1 text-[9px] text-white">{{ notifications.unreadCount }}</span>
            </button>
            <div v-if="showNotifications" class="absolute right-0 top-12 z-30 w-80 rounded-2xl border border-[#e2dfda] bg-white p-3 shadow-xl dark:border-white/10 dark:bg-[#242328]">
              <div class="flex items-center justify-between px-2 py-1">
                <span class="text-sm font-semibold">Notifications</span>
                <button class="text-xs text-[#665ad3]" @click="notifications.markAllRead">Mark all read</button>
              </div>
              <p v-if="!notificationItems.length" class="p-5 text-xs text-[#99949a]">No notifications.</p>
              <button
                v-for="item in notificationItems.slice(0, 5)"
                :key="item.id"
                class="mt-2 block w-full rounded-xl p-2 text-left hover:bg-[#f5f3f0] dark:hover:bg-white/5"
                @click="notifications.markRead(item.id)"
              >
                <p class="text-xs font-semibold">{{ item.title }}</p>
                <p class="mt-1 line-clamp-2 text-[11px] text-[#99949a]">{{ item.body }}</p>
              </button>
            </div>
          </div>
          <button aria-label="Toggle dark mode" @click="ui.toggleDark" class="rounded-xl border border-[#e2dfda] bg-white p-2.5 text-[#77737a] dark:border-white/10 dark:bg-white/5">
            <component :is="ui.dark ? Sun : Moon" :size="18" />
          </button>
        </div>
      </header>
      <RouterView />
      <footer class="mt-10 pb-4 text-center text-xs text-[#aaa5a0]">Demo data and demo rules. Not legal advice.</footer>
    </main>
  </div>
</template>
