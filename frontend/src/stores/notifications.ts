import { computed, onMounted, onUnmounted, ref } from 'vue'
import { apiFetch, ApiError } from '../api/client'

export interface NotificationItem {
  id: string
  kind: string
  title: string
  body: string
  link: string
  read: boolean
  created_at: string
}

export function useNotifications() {
  const items = ref<NotificationItem[]>([])
  const loading = ref(false)
  let timer: number | undefined

  const unreadCount = computed(() => items.value.filter((item) => !item.read).length)
  async function refresh() {
    if (document.visibilityState !== 'visible') return
    if (!localStorage.getItem('complylite_token')) return
    loading.value = true
    try {
      const response = await apiFetch<{ data: NotificationItem[] }>('/notifications?unread=false')
      items.value = response.data
    } catch (err) {
      if (err instanceof ApiError && err.status === 401) {
        items.value = []
      }
    } finally { loading.value = false }
  }
  async function markRead(id: string) {
    await apiFetch(`/notifications/${id}/read`, { method: 'POST' })
    const item = items.value.find((entry) => entry.id === id)
    if (item) item.read = true
  }
  async function markAllRead() {
    await apiFetch('/notifications/read-all', { method: 'POST' })
    items.value.forEach((item) => { item.read = true })
  }
  onMounted(() => { void refresh(); timer = window.setInterval(() => void refresh(), 60_000) })
  onUnmounted(() => { if (timer) window.clearInterval(timer) })
  return { items, loading, unreadCount, refresh, markRead, markAllRead }
}
