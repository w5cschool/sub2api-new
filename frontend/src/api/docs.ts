import { apiClient } from './client'

export type DocumentStatus = 'draft' | 'published'

export interface DocumentSummary {
  slug: string
  title: string
  status: DocumentStatus
  sort_order: number
  created_at: string
  updated_at: string
}

export interface DocumentDetail extends DocumentSummary {
  content: string
}

export async function listDocuments(): Promise<DocumentSummary[]> {
  const { data } = await apiClient.get<DocumentSummary[]>('/docs')
  return data
}

export async function getDocument(slug: string): Promise<DocumentDetail> {
  const { data } = await apiClient.get<DocumentDetail>(`/docs/${encodeURIComponent(slug)}`)
  return data
}

export default { list: listDocuments, get: getDocument }
