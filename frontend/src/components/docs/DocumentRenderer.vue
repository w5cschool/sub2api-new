<template>
  <article ref="container" class="document-markdown" v-html="rendered.html"></article>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { renderDocumentMarkdown, type DocumentHeading } from '@/utils/docsMarkdown'

const props = defineProps<{
  content: string
  slug: string
}>()

const emit = defineEmits<{
  headingsChange: [headings: DocumentHeading[]]
}>()

const container = ref<HTMLElement | null>(null)
const rendered = computed(() => renderDocumentMarkdown(props.content || '', props.slug || 'preview'))

watch(
  () => rendered.value.headings,
  (headings) => emit('headingsChange', headings),
  { immediate: true },
)

function scrollToHeading(id: string) {
  container.value?.querySelector(`#${CSS.escape(id)}`)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

defineExpose({ scrollToHeading })
</script>

<style>
.document-markdown { color: inherit; font-size: 0.98rem; line-height: 1.8; }
.document-markdown > :first-child { margin-top: 0; }
.document-markdown h1 { margin: 0 0 1.5rem; padding-bottom: 0.75rem; border-bottom: 1px solid rgb(229 231 235); font-size: 2.1rem; font-weight: 800; line-height: 1.2; letter-spacing: -0.025em; }
.document-markdown h2 { margin: 2.25rem 0 1rem; font-size: 1.55rem; font-weight: 750; line-height: 1.3; }
.document-markdown h3 { margin: 1.8rem 0 0.75rem; font-size: 1.25rem; font-weight: 700; }
.document-markdown h4 { margin: 1.5rem 0 0.6rem; font-size: 1.08rem; font-weight: 700; }
.document-markdown p { margin: 0.85rem 0; }
.document-markdown ul { margin: 0.9rem 0; padding-left: 1.6rem; list-style: disc; }
.document-markdown ol { margin: 0.9rem 0; padding-left: 1.6rem; list-style: decimal; }
.document-markdown li { margin: 0.35rem 0; }
.document-markdown a { color: rgb(16 185 129); text-decoration: underline; text-underline-offset: 3px; }
.document-markdown blockquote { margin: 1.2rem 0; border-left: 4px solid rgb(52 211 153); padding: 0.5rem 0 0.5rem 1rem; color: rgb(75 85 99); background: rgb(236 253 245 / 0.55); border-radius: 0 0.65rem 0.65rem 0; }
.document-markdown img { display: block; max-width: 100%; height: auto; margin: 1.5rem auto; border-radius: 0.85rem; box-shadow: 0 8px 30px rgb(15 23 42 / 0.1); }
.document-markdown table { display: block; width: 100%; overflow-x: auto; margin: 1.2rem 0; border-collapse: collapse; }
.document-markdown th, .document-markdown td { border: 1px solid rgb(209 213 219); padding: 0.6rem 0.8rem; text-align: left; }
.document-markdown th { background: rgb(249 250 251); font-weight: 700; }
.document-markdown code { border-radius: 0.35rem; background: rgb(243 244 246); padding: 0.15rem 0.38rem; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 0.88em; }
.document-markdown pre { overflow-x: auto; margin: 1.2rem 0; border-radius: 0.85rem; background: rgb(15 23 42); padding: 1rem 1.15rem; color: rgb(226 232 240); }
.document-markdown pre code { background: transparent; padding: 0; color: inherit; }
.document-markdown hr { margin: 2rem 0; border: 0; border-top: 1px solid rgb(229 231 235); }
.dark .document-markdown h1 { border-color: rgb(55 65 81); }
.dark .document-markdown blockquote { color: rgb(209 213 219); background: rgb(6 78 59 / 0.18); }
.dark .document-markdown th, .dark .document-markdown td { border-color: rgb(75 85 99); }
.dark .document-markdown th, .dark .document-markdown code { background: rgb(31 41 55); }
.dark .document-markdown hr { border-color: rgb(55 65 81); }
</style>
