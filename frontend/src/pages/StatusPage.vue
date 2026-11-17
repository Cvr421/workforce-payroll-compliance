<script setup lang="ts">
interface Props {
  status: number
  title: string
  message: string
}
const props = defineProps<Props>()
</script>

<template>
  <main class="flex min-h-screen items-center justify-center bg-[#f7f6f3] px-6 dark:bg-[#18171b]">
    <div class="text-center">
      <p class="text-6xl font-semibold text-[#8f83ec]">{{ status }}</p>
      <h1 class="mt-4 text-2xl font-semibold">{{ title }}</h1>
      <p class="mt-2 text-[#99949a]">{{ message }}</p>
      <a href="/" class="mt-6 inline-block text-[#665ad3] hover:underline">Go home</a>
    </div>
  </main>
</template>
