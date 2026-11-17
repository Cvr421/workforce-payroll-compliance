<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ApiError } from '../api/client'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()
const email = ref('')
const password = ref('')
const error = ref('')
const pending = ref(false)
const demos = [
  ['Admin (Tenant A)', 'admin@demo.test', 'Admin@12345'],
  ['Payroll Manager 1', 'payroll1@demo.test', 'Payroll@12345'],
  ['Payroll Manager 2', 'payroll2@demo.test', 'Payroll@12345'],
  ['Compliance Officer', 'compliance@demo.test', 'Compliance@12345'],
  ['Auditor', 'auditor@demo.test', 'Auditor@12345'],
  ['Admin (Tenant B)', 'other@demo.test', 'Other@12345'],
]
async function submit() {
  pending.value = true
  error.value = ''
  try {
    await auth.login(email.value, password.value)
    await router.push('/')
  } catch (reason) {
    error.value = reason instanceof ApiError ? reason.message : 'Unable to sign in'
  } finally {
    pending.value = false
  }
}
function fill(account: string[]) {
  email.value = account[1]
  password.value = account[2]
}
</script>

<template>
  <main class="grid min-h-screen lg:grid-cols-2">
    <section class="relative hidden overflow-hidden bg-[#a69cf4] p-12 text-white lg:flex lg:flex-col lg:justify-between">
      <div class="absolute inset-0 opacity-30 [background-image:linear-gradient(rgba(255,255,255,.25)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,.25)_1px,transparent_1px)] [background-size:42px_42px]" />
      <div class="relative flex items-center gap-3 text-xl font-semibold"><span class="grid h-9 w-9 place-items-center rounded-xl bg-white/20">C</span>comply<span class="text-[#eeeaff]">lite</span></div>
      <div class="relative max-w-lg"><p class="mb-5 text-sm uppercase tracking-[.25em] text-white/70">Compliance, made clear</p><h1 class="text-5xl font-semibold leading-tight tracking-[-.05em]">A calmer way to run your back office.</h1><p class="mt-6 max-w-md text-base leading-7 text-white/75">Keep payroll workflows, evidence, and audit readiness moving in one focused workspace.</p></div>
      <p class="relative text-xs text-white/60">Fictional demo workspace · Not legal advice</p>
    </section>
    <section class="flex items-center justify-center bg-[#f7f6f3] px-6 py-12 dark:bg-[#18171b]">
      <div class="w-full max-w-md">
        <div class="mb-10 lg:hidden"><span class="font-semibold">comply<span class="text-[#8f83ec]">lite</span></span></div>
        <p class="text-sm text-[#99949a]">Welcome back</p><h2 class="mt-2 text-3xl font-semibold tracking-tight">Sign in to your workspace</h2>
        <form class="mt-8 space-y-5" @submit.prevent="submit">
          <label class="block text-sm font-medium">Email<input v-model="email" type="email" autocomplete="email" class="mt-2 h-11 w-full rounded-xl border border-[#dfdcd7] bg-white px-3 outline-none focus:border-[#9b90ef]" required /></label>
          <label class="block text-sm font-medium">Password<input v-model="password" type="password" autocomplete="current-password" class="mt-2 h-11 w-full rounded-xl border border-[#dfdcd7] bg-white px-3 outline-none focus:border-[#9b90ef]" required /></label>
          <p v-if="error" class="rounded-lg bg-[#fff0ed] px-3 py-2 text-sm text-[#c76052]">{{ error }}</p>
          <button :disabled="pending" class="h-11 w-full rounded-xl bg-[#8f83ec] font-medium text-white transition hover:bg-[#796ddd] disabled:opacity-60">{{ pending ? 'Signing in…' : 'Sign in' }}</button>
        </form>
        <div class="mt-8 border-t border-[#e4e0db] pt-6"><p class="text-xs font-semibold uppercase tracking-[.16em] text-[#aaa5a0]">Demo accounts</p><div class="mt-3 flex flex-wrap gap-2"><button v-for="account in demos" :key="account[0]" @click="fill(account)" class="rounded-lg border border-[#dfdcd7] bg-white px-3 py-2 text-xs font-medium hover:border-[#9b90ef]">{{ account[0] }}</button></div></div>
      </div>
    </section>
  </main>
</template>
