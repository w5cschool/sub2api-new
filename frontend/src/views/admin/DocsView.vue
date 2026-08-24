<template>
  <AppLayout>
    <div class="mx-auto flex min-h-[calc(100vh-8rem)] max-w-[1600px] overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900">
      <aside class="w-72 flex-shrink-0 border-r border-gray-200 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/60">
        <div class="mb-4 flex items-center justify-between gap-2">
          <div class="flex items-center gap-2">
            <Icon name="book" size="md" class="text-primary-500" />
            <h2 class="font-semibold text-gray-900 dark:text-white">{{ t('admin.docs.title') }}</h2>
          </div>
          <button type="button" class="rounded-lg p-2 text-primary-600 hover:bg-primary-50 dark:text-primary-400 dark:hover:bg-primary-900/20" :title="t('admin.docs.newDocument')" @click="startCreate">
            <Icon name="plus" size="sm" />
          </button>
        </div>

        <div v-if="loadingList" class="space-y-2">
          <div v-for="i in 4" :key="i" class="h-14 animate-pulse rounded-xl bg-gray-200 dark:bg-dark-700"></div>
        </div>
        <div v-else class="space-y-1.5">
          <button
            v-for="item in documents"
            :key="item.slug"
            type="button"
            class="w-full rounded-xl px-3 py-2.5 text-left transition"
            :class="item.slug === selectedSlug && !creating
              ? 'bg-white shadow-sm ring-1 ring-gray-200 dark:bg-dark-800 dark:ring-dark-600'
              : 'hover:bg-white/80 dark:hover:bg-dark-800/70'"
            @click="selectDocument(item.slug)"
          >
            <div class="truncate text-sm font-medium text-gray-900 dark:text-white">{{ item.title }}</div>
            <div class="mt-1 flex items-center justify-between gap-2 text-xs">
              <span class="truncate text-gray-400">/{{ item.slug }}</span>
              <span :class="item.status === 'published' ? 'text-emerald-600 dark:text-emerald-400' : 'text-amber-600 dark:text-amber-400'">
                {{ item.status === 'published' ? t('admin.docs.published') : t('admin.docs.draft') }}
              </span>
            </div>
          </button>
          <button v-if="documents.length === 0" type="button" class="w-full rounded-xl border border-dashed border-gray-300 px-4 py-8 text-sm text-gray-500 hover:border-primary-400 hover:text-primary-600 dark:border-dark-600 dark:text-dark-400" @click="startCreate">
            {{ t('admin.docs.createFirst') }}
          </button>
        </div>
      </aside>

      <main class="min-w-0 flex-1">
        <div class="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 px-5 py-4 dark:border-dark-700">
          <div>
            <h1 class="text-lg font-semibold text-gray-900 dark:text-white">{{ creating ? t('admin.docs.newDocument') : form.title || t('admin.docs.editDocument') }}</h1>
            <p class="mt-0.5 text-xs text-gray-500 dark:text-dark-400">{{ t('admin.docs.description') }}</p>
          </div>
          <div class="flex items-center gap-2">
            <a v-if="!creating && form.status === 'published'" :href="`/docs/${encodeURIComponent(form.slug)}`" target="_blank" class="btn btn-secondary btn-sm">
              <Icon name="eye" size="sm" class="mr-1" />{{ t('admin.docs.view') }}
            </a>
            <button v-if="!creating" type="button" class="btn btn-danger btn-sm" :disabled="saving" @click="deleteCurrent">
              <Icon name="trash" size="sm" class="mr-1" />{{ t('common.delete') }}
            </button>
            <button type="button" class="btn btn-primary btn-sm" :disabled="saving || !form.title || !form.slug" @click="saveDocument">
              <Icon name="check" size="sm" class="mr-1" />{{ saving ? t('common.saving') : t('common.save') }}
            </button>
          </div>
        </div>

        <div class="grid gap-4 border-b border-gray-100 bg-gray-50/50 px-5 py-4 dark:border-dark-700 dark:bg-dark-900/40 md:grid-cols-[minmax(0,1fr)_minmax(0,1fr)_140px_160px]">
          <div>
            <label class="input-label">{{ t('admin.docs.documentTitle') }}</label>
            <input v-model="form.title" class="input" maxlength="200" @input="handleTitleInput" />
          </div>
          <div>
            <label class="input-label">Slug</label>
            <input v-model="form.slug" class="input font-mono" maxlength="64" :disabled="!creating" placeholder="getting-started" />
          </div>
          <div>
            <label class="input-label">{{ t('admin.docs.order') }}</label>
            <input v-model.number="form.sort_order" class="input" type="number" />
          </div>
          <div>
            <label class="input-label">{{ t('admin.docs.status') }}</label>
            <select v-model="form.status" class="input">
              <option value="draft">{{ t('admin.docs.draft') }}</option>
              <option value="published">{{ t('admin.docs.published') }}</option>
            </select>
          </div>
        </div>

        <div class="grid min-h-[calc(100vh-19rem)] grid-cols-1 divide-y divide-gray-200 dark:divide-dark-700 xl:grid-cols-2 xl:divide-x xl:divide-y-0">
          <section class="flex min-h-[34rem] min-w-0 flex-col">
            <div class="flex items-center justify-between border-b border-gray-100 px-4 py-2.5 dark:border-dark-700">
              <span class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-dark-400">Markdown</span>
              <div>
                <input ref="imageInput" type="file" accept="image/png,image/jpeg,image/webp,image/gif" class="hidden" @change="uploadImage" />
                <button type="button" class="btn btn-secondary btn-sm" :disabled="creating || uploading" @click="imageInput?.click()">
                  <Icon name="upload" size="sm" class="mr-1" />{{ uploading ? t('admin.docs.uploading') : t('admin.docs.uploadImage') }}
                </button>
              </div>
            </div>
            <textarea
              ref="editor"
              v-model="form.content"
              class="min-h-[34rem] flex-1 resize-none border-0 bg-white p-5 font-mono text-sm leading-7 text-gray-800 outline-none dark:bg-dark-900 dark:text-dark-100"
              :placeholder="t('admin.docs.editorPlaceholder')"
              spellcheck="false"
            ></textarea>
          </section>

          <section class="min-h-[34rem] min-w-0 overflow-auto bg-white dark:bg-dark-900">
            <div class="sticky top-0 z-10 border-b border-gray-100 bg-white/90 px-4 py-3 text-xs font-semibold uppercase tracking-wider text-gray-500 backdrop-blur dark:border-dark-700 dark:bg-dark-900/90 dark:text-dark-400">
              {{ t('admin.docs.livePreview') }}
            </div>
            <div class="p-6 lg:p-8">
              <DocumentRenderer :content="form.content" :slug="form.slug || 'preview'" />
            </div>
          </section>
        </div>
      </main>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { nextTick, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '@/api/admin'
import type { DocumentMutation } from '@/api/admin/docs'
import type { DocumentStatus, DocumentSummary } from '@/api/docs'
import { useAppStore } from '@/stores/app'
import AppLayout from '@/components/layout/AppLayout.vue'
import DocumentRenderer from '@/components/docs/DocumentRenderer.vue'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()
const appStore = useAppStore()
const documents = ref<DocumentSummary[]>([])
const selectedSlug = ref('')
const creating = ref(true)
const loadingList = ref(true)
const saving = ref(false)
const uploading = ref(false)
const editor = ref<HTMLTextAreaElement | null>(null)
const imageInput = ref<HTMLInputElement | null>(null)

const form = reactive<DocumentMutation>({
  slug: '',
  title: '',
  status: 'draft' as DocumentStatus,
  sort_order: 0,
  content: '# 新文档\n\n在这里开始编写内容。',
})

function errorMessage(error: any, fallback: string) {
  return error?.message || error?.response?.data?.message || fallback
}

async function loadList() {
  loadingList.value = true
  try {
    documents.value = await adminAPI.docs.list()
  } catch (error) {
    appStore.showError(errorMessage(error, t('admin.docs.loadFailed')))
  } finally {
    loadingList.value = false
  }
}

function resetForm() {
  Object.assign(form, {
    slug: '', title: '', status: 'draft', sort_order: documents.value.length,
    content: '# 新文档\n\n在这里开始编写内容。',
  })
}

function startCreate() {
  creating.value = true
  selectedSlug.value = ''
  resetForm()
}

function slugify(value: string) {
  const slug = value.toLowerCase().trim().replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '')
  return slug || `doc-${Date.now()}`
}

function handleTitleInput() {
  if (creating.value && (!form.slug || form.slug.startsWith('doc-'))) {
    form.slug = slugify(form.title)
  }
}

async function selectDocument(slug: string) {
  try {
    const detail = await adminAPI.docs.get(slug)
    Object.assign(form, detail)
    selectedSlug.value = slug
    creating.value = false
  } catch (error) {
    appStore.showError(errorMessage(error, t('admin.docs.loadFailed')))
  }
}

async function saveDocument() {
  saving.value = true
  try {
    const payload = { ...form, slug: form.slug.trim(), title: form.title.trim() }
    const saved = creating.value
      ? await adminAPI.docs.create(payload)
      : await adminAPI.docs.update(selectedSlug.value, payload)
    Object.assign(form, saved)
    selectedSlug.value = saved.slug
    creating.value = false
    await loadList()
    appStore.showSuccess(t('admin.docs.saved'))
  } catch (error) {
    appStore.showError(errorMessage(error, t('admin.docs.saveFailed')))
  } finally {
    saving.value = false
  }
}

async function deleteCurrent() {
  if (!selectedSlug.value || !window.confirm(t('admin.docs.deleteConfirm'))) return
  saving.value = true
  try {
    await adminAPI.docs.remove(selectedSlug.value)
    await loadList()
    appStore.showSuccess(t('admin.docs.deleted'))
    if (documents.value[0]) {
      await selectDocument(documents.value[0].slug)
    } else {
      startCreate()
    }
  } catch (error) {
    appStore.showError(errorMessage(error, t('admin.docs.deleteFailed')))
  } finally {
    saving.value = false
  }
}

async function uploadImage(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file || !selectedSlug.value) return
  uploading.value = true
  try {
    const result = await adminAPI.docs.uploadImage(selectedSlug.value, file)
    const textarea = editor.value
    const insertion = `\n${result.markdown}\n`
    if (textarea) {
      const start = textarea.selectionStart
      const end = textarea.selectionEnd
      form.content = form.content.slice(0, start) + insertion + form.content.slice(end)
      await nextTick()
      textarea.focus()
      textarea.selectionStart = textarea.selectionEnd = start + insertion.length
    } else {
      form.content += insertion
    }
    appStore.showSuccess(t('admin.docs.imageUploaded'))
  } catch (error) {
    appStore.showError(errorMessage(error, t('admin.docs.imageUploadFailed')))
  } finally {
    uploading.value = false
    input.value = ''
  }
}

onMounted(async () => {
  await loadList()
  if (documents.value[0]) {
    await selectDocument(documents.value[0].slug)
  } else {
    startCreate()
  }
})
</script>
