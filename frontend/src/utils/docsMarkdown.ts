import DOMPurify from 'dompurify'
import { marked } from 'marked'
import { buildApiUrl } from '@/api/client'

export interface DocumentHeading {
  id: string
  text: string
  level: number
}

export interface RenderedDocument {
  html: string
  headings: DocumentHeading[]
}

function isRelativeAsset(source: string): boolean {
  const value = source.trim()
  if (!value || /^[a-z][a-z0-9+.-]*:/i.test(value) || value.startsWith('//') || value.startsWith('/')) {
    return false
  }
  return value.split('/').every((part) => part !== '..' && !part.includes('\\'))
}

function documentImageUrl(slug: string, source: string): string {
  const [pathPart, suffix = ''] = source.trim().split(/([?#].*)/, 2)
  const encoded = pathPart
    .split('/')
    .filter((part) => part && part !== '.')
    .map(encodeURIComponent)
    .join('/')
  return buildApiUrl(`/docs/${encodeURIComponent(slug)}/images/${encoded}${suffix}`)
}

function headingID(text: string, index: number): string {
  const normalized = text
    .toLowerCase()
    .replace(/<[^>]+>/g, '')
    .replace(/[^\w\u4e00-\u9fff]+/g, '-')
    .replace(/^-+|-+$/g, '')
  return normalized ? `${normalized}-${index}` : `heading-${index}`
}

export function renderDocumentMarkdown(content: string, slug: string): RenderedDocument {
  const rewritten = content.replace(
    /!\[([^\]]*)\]\(([^)]+)\)/g,
    (match, alt: string, source: string) => isRelativeAsset(source)
      ? `![${alt}](${documentImageUrl(slug, source)})`
      : match,
  )
  const sanitized = DOMPurify.sanitize(marked.parse(rewritten) as string)
  const headings: DocumentHeading[] = []
  let index = 0
  const html = sanitized.replace(
    /<(h[1-4])[^>]*>(.*?)<\/h[1-4]>/gi,
    (_match, tag: string, inner: string) => {
      const text = inner.replace(/<[^>]+>/g, '').trim()
      const id = headingID(text, index++)
      headings.push({ id, text, level: Number(tag[1]) })
      return `<${tag} id="${id}">${inner}</${tag}>`
    },
  )
  return { html, headings }
}
