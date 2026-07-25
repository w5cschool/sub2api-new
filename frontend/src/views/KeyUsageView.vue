<template>
  <div class="min-h-screen bg-[#f7f8f3] text-slate-900 dark:bg-[#10140f] dark:text-stone-100">
    <header class="border-b border-stone-200/80 bg-[#f7f8f3]/95 dark:border-stone-800 dark:bg-[#10140f]/95">
      <nav class="mx-auto flex min-h-16 max-w-5xl items-center justify-between gap-4 px-6 py-3">
        <RouterLink to="/home" class="flex min-w-0 items-center gap-3">
          <span class="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-xl border border-stone-200 bg-white p-1.5 dark:border-stone-700 dark:bg-stone-900">
            <img :src="siteLogo || '/logo.svg'" :alt="`${siteName} 标识`" class="h-full w-full object-contain" />
          </span>
          <span class="truncate text-base font-bold tracking-[-0.03em]">{{ siteName }}</span>
        </RouterLink>
        <div class="flex items-center gap-3">
          <LocaleSwitcher />
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="hidden text-sm font-medium text-stone-600 transition hover:text-stone-950 sm:inline dark:text-stone-300 dark:hover:text-stone-100"
          >
            {{ t('home.viewDocs') }}
          </a>
        </div>
      </nav>
    </header>

    <main class="mx-auto flex min-h-[calc(100vh-4rem)] max-w-5xl items-center px-6 py-16">
      <section class="grid w-full gap-12 md:grid-cols-[1.15fr_0.85fr] md:items-end">
        <div class="border-l-2 border-[#7d9b72] pl-6 md:pl-8">
          <p class="text-sm font-semibold tracking-[0.14em] text-[#52694b] dark:text-[#b8d5ab]">ACCOUNT-ONLY ACCESS</p>
          <h1 class="mt-5 text-4xl font-bold leading-[1.08] tracking-[-0.05em] text-stone-950 dark:text-stone-50 md:text-5xl">
            {{ t('keyUsage.accessTitle') }}
          </h1>
          <p class="mt-6 max-w-xl text-base leading-8 text-stone-600 dark:text-stone-300">
            {{ t('keyUsage.accessDescription') }}
          </p>
          <RouterLink
            :to="primaryTarget"
            class="mt-8 inline-flex items-center gap-2 rounded-lg bg-[#42633b] px-5 py-3 text-sm font-semibold text-white transition hover:bg-[#355330] dark:bg-[#c5dfbb] dark:text-[#173019] dark:hover:bg-[#d9edcf]"
          >
            {{ isAuthenticated ? t('keyUsage.accessDashboard') : t('keyUsage.accessLogin') }}
            <svg class="h-4 w-4" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.8">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4 10h12m-4-4 4 4-4 4" />
            </svg>
          </RouterLink>
        </div>

        <aside class="border-y border-stone-200 py-7 dark:border-stone-800">
          <h2 class="text-base font-bold tracking-[-0.02em]">{{ t('keyUsage.securityTitle') }}</h2>
          <p class="mt-3 text-sm leading-7 text-stone-600 dark:text-stone-300">{{ t('keyUsage.securityDescription') }}</p>
          <p class="mt-5 text-sm leading-7 text-stone-600 dark:text-stone-300">{{ t('keyUsage.securityReport') }}</p>
        </aside>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import { useAppStore, useAuthStore } from '@/stores'
import { sanitizeUrl } from '@/utils/url'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const isAuthenticated = computed(() => authStore.isAuthenticated)
const primaryTarget = computed(() => {
  if (!isAuthenticated.value) return '/login'
  return authStore.isAdmin ? '/admin/dashboard' : '/dashboard'
})

onMounted(() => {
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>
