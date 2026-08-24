<template>
  <component
    :is="authStore.isAuthenticated ? AppLayout : 'div'"
    :class="authStore.isAuthenticated ? undefined : 'min-h-screen bg-gray-50 px-4 py-6 dark:bg-dark-950 sm:px-6'"
  >
    <header
      v-if="!authStore.isAuthenticated"
      class="mx-auto mb-6 flex max-w-7xl items-center justify-between rounded-2xl border border-gray-200 bg-white px-5 py-4 shadow-sm dark:border-dark-700 dark:bg-dark-900"
    >
      <router-link to="/home" class="flex items-center gap-3 font-semibold text-gray-900 dark:text-white">
        <Icon name="book" size="md" class="text-primary-500" />
        <span>Sub2API · {{ t('docs.title') }}</span>
      </router-link>
      <router-link to="/login" class="btn btn-primary btn-sm">{{ t('auth.signIn') }}</router-link>
    </header>

    <div class="mx-auto flex min-h-[calc(100vh-8rem)] max-w-7xl overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900">
      <aside class="hidden w-72 flex-shrink-0 border-r border-gray-200 bg-gray-50/70 p-4 dark:border-dark-700 dark:bg-dark-900/60 md:block">
        <div class="mb-4 flex items-center gap-2 px-2">
          <Icon name="book" size="md" class="text-primary-500" />
          <h2 class="font-semibold text-gray-900 dark:text-white">{{ t('docs.directory') }}</h2>
        </div>
        <div v-if="loadingList" class="space-y-2 px-2">
          <div v-for="i in 4" :key="i" class="h-9 animate-pulse rounded-lg bg-gray-200 dark:bg-dark-700"></div>
        </div>
        <nav v-else class="space-y-1">
          <button
            v-for="item in documents"
            :key="item.slug"
            type="button"
            class="w-full rounded-xl px-3 py-2.5 text-left text-sm transition"
            :class="item.slug === activeSlug
              ? 'bg-primary-50 font-semibold text-primary-700 dark:bg-primary-900/20 dark:text-primary-300'
              : 'text-gray-600 hover:bg-white hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white'"
            @click="selectDocument(item.slug)"
          >
            {{ item.title }}
          </button>

          <div v-if="headings.length" class="ml-3 mt-2 border-l border-gray-200 py-1 pl-3 dark:border-dark-600">
            <button
              v-for="heading in headings"
              :key="heading.id"
              type="button"
              class="block w-full truncate py-1 text-left text-xs text-gray-500 hover:text-primary-600 dark:text-dark-400 dark:hover:text-primary-400"
              :style="{ paddingLeft: `${Math.max(0, heading.level - 1) * 8}px` }"
              @click="scrollToHeading(heading.id)"
            >
              {{ heading.text }}
            </button>
          </div>
        </nav>
      </aside>

      <main class="min-w-0 flex-1">
        <div class="border-b border-gray-100 p-4 dark:border-dark-700 md:hidden">
          <select v-model="activeSlug" class="input" @change="selectDocument(activeSlug)">
            <option v-for="item in documents" :key="item.slug" :value="item.slug">{{ item.title }}</option>
          </select>
        </div>

        <div v-if="loadingDocument" class="flex min-h-[30rem] items-center justify-center">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
        </div>
        <div v-else-if="document" class="mx-auto max-w-4xl px-6 py-8 md:px-10 md:py-12">
          <div class="mb-8 flex flex-wrap items-center justify-between gap-3 border-b border-gray-100 pb-5 text-sm text-gray-500 dark:border-dark-700 dark:text-dark-400">
            <span>{{ t('docs.lastUpdated') }} {{ formatDate(document.updated_at) }}</span>
            <span class="rounded-full bg-emerald-50 px-3 py-1 text-xs font-medium text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300">Markdown</span>
          </div>
          <DocumentRenderer
            ref="renderer"
            :content="document.content"
            :slug="document.slug"
            @headings-change="headings = $event"
          />
        </div>
        <div v-else class="flex min-h-[32rem] items-center justify-center px-6 text-center">
          <div class="max-w-md">
            <Icon name="book" size="xl" class="mx-auto mb-4 text-gray-300 dark:text-dark-600" />
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('docs.emptyTitle') }}</h2>
            <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">{{ t('docs.emptyDescription') }}</p>
          </div>
        </div>
      </main>
    </div>
  </component>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { docsAPI, type DocumentDetail, type DocumentSummary } from '@/api'
import type { DocumentHeading } from '@/utils/docsMarkdown'
import AppLayout from '@/components/layout/AppLayout.vue'
import DocumentRenderer from '@/components/docs/DocumentRenderer.vue'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { extractApiErrorMessage } from '@/utils/apiError'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()

const documents = ref<DocumentSummary[]>([])
const document = ref<DocumentDetail | null>(null)
const activeSlug = ref('')
const headings = ref<DocumentHeading[]>([])
const renderer = ref<InstanceType<typeof DocumentRenderer> | null>(null)
const loadingList = ref(true)
const loadingDocument = ref(false)

function formatDate(value: string) {
  return new Intl.DateTimeFormat(undefined, { year: 'numeric', month: 'short', day: 'numeric' }).format(new Date(value))
}

async function loadDocument(slug: string) {
  if (!slug) {
    document.value = null
    headings.value = []
    return
  }
  loadingDocument.value = true
  try {
    document.value = await docsAPI.get(slug)
    activeSlug.value = slug
  } catch (error) {
    document.value = null
    headings.value = []
    appStore.showError(extractApiErrorMessage(error, t('docs.loadFailed')))
  } finally {
    loadingDocument.value = false
  }
}

async function selectDocument(slug: string) {
  if (!slug) return
  if (route.params.slug !== slug) {
    await router.push(`/docs/${encodeURIComponent(slug)}`)
  } else if (document.value?.slug !== slug) {
    await loadDocument(slug)
  }
}

function scrollToHeading(id: string) {
  renderer.value?.scrollToHeading(id)
}

onMounted(async () => {
  try {
    documents.value = await docsAPI.list()
    const routeSlug = typeof route.params.slug === 'string' ? route.params.slug : ''
    const initial = documents.value.some((item) => item.slug === routeSlug) ? routeSlug : documents.value[0]?.slug || ''
    if (initial) {
      await selectDocument(initial)
    }
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('docs.loadFailed')))
  } finally {
    loadingList.value = false
  }
})

watch(
  () => route.params.slug,
  (slug) => {
    if (typeof slug === 'string' && slug && slug !== document.value?.slug) {
      void loadDocument(slug)
    }
  },
)
</script>
