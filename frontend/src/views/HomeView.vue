<template>
  <div class="min-h-screen bg-[#f7f8f3] text-slate-900 dark:bg-[#10140f] dark:text-stone-100">
    <header class="border-b border-stone-200/80 bg-[#f7f8f3]/95 dark:border-stone-800 dark:bg-[#10140f]/95">
      <nav class="mx-auto flex min-h-16 max-w-6xl items-center justify-between gap-4 px-6 py-3">
        <div class="flex min-w-0 items-center gap-3">
          <div
            v-if="siteLogo"
            class="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-xl border border-stone-200 bg-white p-1.5 dark:border-stone-700 dark:bg-stone-900"
          >
            <img :src="siteLogo" :alt="`${siteName} 标识`" class="h-full w-full object-contain" />
          </div>
          <div class="min-w-0">
            <div class="truncate text-base font-bold tracking-[-0.03em]">{{ siteName }}</div>
            <div class="mt-0.5 text-xs text-stone-500 dark:text-stone-400">API 管理与接入服务</div>
          </div>
        </div>

        <div class="flex shrink-0 items-center gap-2 sm:gap-3">
          <LocaleSwitcher />
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="hidden rounded-lg px-3 py-2 text-sm font-medium text-stone-600 transition hover:bg-stone-200/70 hover:text-stone-900 sm:inline-flex dark:text-stone-300 dark:hover:bg-stone-800 dark:hover:text-stone-100"
          >
            使用文档
          </a>
          <router-link
            :to="primaryTarget"
            class="inline-flex rounded-lg bg-[#263b2d] px-4 py-2 text-sm font-semibold text-white transition hover:bg-[#1c3023] dark:bg-[#dcebd6] dark:text-[#193220] dark:hover:bg-white"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main>
      <section class="mx-auto max-w-6xl px-6 pb-16 pt-16 md:pb-24 md:pt-24">
        <div class="max-w-3xl border-l-2 border-[#7d9b72] pl-6 md:pl-8">
          <p class="text-sm font-semibold tracking-[0.14em] text-[#52694b] dark:text-[#b8d5ab]">INDEPENDENT SERVICE</p>
          <h1 class="mt-5 text-4xl font-bold leading-[1.06] tracking-[-0.055em] text-stone-950 dark:text-stone-50 md:text-6xl">
            清晰管理 API 接入，
            <span class="block text-[#42633b] dark:text-[#c5dfbb]">让每一次连接有据可查。</span>
          </h1>
          <p class="mt-7 max-w-2xl text-base leading-8 text-stone-600 dark:text-stone-300 md:text-lg">
            {{ siteName }} 是由本站运营方独立提供的 API 管理与接入服务，用于帮助已获授权的用户配置、管理和审计自己的 API 调用。
          </p>
          <div class="mt-8 flex flex-wrap gap-3">
            <router-link
              :to="primaryTarget"
              class="inline-flex items-center gap-2 rounded-lg bg-[#42633b] px-5 py-3 text-sm font-semibold text-white transition hover:bg-[#355330] dark:bg-[#c5dfbb] dark:text-[#173019] dark:hover:bg-[#d9edcf]"
            >
              {{ isAuthenticated ? '进入控制台' : '前往登录' }}
              <svg class="h-4 w-4" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4 10h12m-4-4 4 4-4 4" />
              </svg>
            </router-link>
            <a
              v-if="docUrl"
              :href="docUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center rounded-lg border border-stone-300 px-5 py-3 text-sm font-semibold text-stone-700 transition hover:border-stone-400 hover:bg-white dark:border-stone-700 dark:text-stone-200 dark:hover:bg-stone-900"
            >
              阅读使用文档
            </a>
          </div>
        </div>
      </section>

      <section class="border-y border-stone-200 bg-white/60 dark:border-stone-800 dark:bg-stone-900/30">
        <div class="mx-auto grid max-w-6xl gap-x-12 gap-y-10 px-6 py-14 md:grid-cols-3">
          <article>
            <div class="text-sm font-semibold text-[#42633b] dark:text-[#c5dfbb]">01 / 服务范围</div>
            <h2 class="mt-3 text-xl font-bold tracking-[-0.035em]">面向已获授权的 API 使用</h2>
            <p class="mt-3 text-sm leading-7 text-stone-600 dark:text-stone-300">
              提供密钥、用量与接入配置的管理能力。请仅连接您有权使用的服务与账户，并遵守适用的法律、合同和服务条款。
            </p>
          </article>
          <article>
            <div class="text-sm font-semibold text-[#42633b] dark:text-[#c5dfbb]">02 / 品牌边界</div>
            <h2 class="mt-3 text-xl font-bold tracking-[-0.035em]">独立运营，不冒充第三方</h2>
            <p class="mt-3 text-sm leading-7 text-stone-600 dark:text-stone-300">
              本站与任何第三方模型、支付、身份或账户服务品牌不存在隶属或授权关系，除非双方以正式书面协议另有明确约定。
            </p>
          </article>
          <article>
            <div class="text-sm font-semibold text-[#42633b] dark:text-[#c5dfbb]">03 / 安全边界</div>
            <h2 class="mt-3 text-xl font-bold tracking-[-0.035em]">不会通过首页索取敏感信息</h2>
            <p class="mt-3 text-sm leading-7 text-stone-600 dark:text-stone-300">
              本站不会在首页要求您提供第三方账号密码、验证码、银行卡信息、私钥或其他平台的 API 密钥；遇到此类请求请停止操作并联系本站支持。
            </p>
          </article>
        </div>
      </section>

      <section class="mx-auto grid max-w-6xl gap-10 px-6 py-16 md:grid-cols-[1.1fr_0.9fr] md:py-24">
        <div>
          <p class="text-sm font-semibold tracking-[0.14em] text-[#52694b] dark:text-[#b8d5ab]">HOW TO USE</p>
          <h2 class="mt-4 text-3xl font-bold tracking-[-0.045em] md:text-4xl">从文档开始，而不是从承诺开始。</h2>
          <p class="mt-5 max-w-xl text-base leading-8 text-stone-600 dark:text-stone-300">
            使用前，请先确认服务范围、访问权限与数据处理要求。配置步骤和支持渠道以本站公开文档与控制台内的说明为准。
          </p>
        </div>
        <aside class="border-t-2 border-[#7d9b72] pt-6 dark:border-[#9dbb91]">
          <div class="space-y-5 text-sm leading-7 text-stone-600 dark:text-stone-300">
            <p><span class="font-semibold text-stone-900 dark:text-stone-100">需要帮助？</span> 请通过下方公开联系方式联系本站运营方。</p>
            <p><span class="font-semibold text-stone-900 dark:text-stone-100">怀疑异常？</span> 请不要提交任何敏感信息，并立即停止相关操作。</p>
            <p v-if="docUrl"><a :href="docUrl" target="_blank" rel="noopener noreferrer" class="font-semibold text-[#42633b] underline decoration-[#a9c89b] underline-offset-4 dark:text-[#c5dfbb]">查看使用与配置文档</a></p>
          </div>
        </aside>
      </section>

      <section class="border-t border-stone-200 bg-[#e8eee4] dark:border-stone-800 dark:bg-[#172017]">
        <div class="mx-auto flex max-w-6xl flex-col gap-5 px-6 py-12 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <p class="text-sm font-semibold tracking-[0.14em] text-[#52694b] dark:text-[#b8d5ab]">OPERATOR CONTACT</p>
            <h2 class="mt-3 text-2xl font-bold tracking-[-0.04em]">联系本站运营方</h2>
            <p class="mt-3 text-sm leading-7 text-stone-600 dark:text-stone-300">服务问题、安全报告或账号协助，请使用以下公开联系方式。</p>
          </div>
          <div class="rounded-lg border border-stone-300 bg-white/75 px-4 py-3 font-medium text-stone-800 dark:border-stone-700 dark:bg-stone-900/70 dark:text-stone-100">
            {{ contactValue }}
          </div>
        </div>
      </section>
    </main>

    <footer class="border-t border-stone-200 px-6 py-7 text-center text-sm text-stone-500 dark:border-stone-800 dark:text-stone-400">
      © {{ currentYear }} {{ siteName }} · 独立运营的 API 管理服务
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore, useAuthStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import { sanitizeUrl } from '@/utils/url'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const contactValue = computed(() => appStore.cachedPublicSettings?.contact_info?.trim() || '请联系本站管理员')
const isAuthenticated = computed(() => authStore.isAuthenticated)
const primaryTarget = computed(() => {
  if (!isAuthenticated.value) return '/login'
  return authStore.isAdmin ? '/admin/dashboard' : '/dashboard'
})
const currentYear = computed(() => new Date().getFullYear())

onMounted(() => {
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>
