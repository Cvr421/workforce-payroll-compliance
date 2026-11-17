<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from './stores/auth'
import { router } from './router'

const auth = useAuthStore()

onMounted(async () => {
  if (auth.token && !auth.user) {
    try {
      await auth.loadUser()
    } catch {
      auth.logout()
      router.push('/login')
    }
  }
})
</script>

<template>
  <RouterView v-slot="{ Component }" >
    <component :is="Component" />
  </RouterView>
</template>
