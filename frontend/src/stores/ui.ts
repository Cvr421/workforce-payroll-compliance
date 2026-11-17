import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUiStore = defineStore('ui', () => {
  const dark = ref(false)
  const activePage = ref('Dashboard')
  const toggleDark = () => { dark.value = !dark.value }
  return { dark, activePage, toggleDark }
})
