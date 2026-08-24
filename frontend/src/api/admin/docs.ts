import { apiClient } from '../client'
import type { DocumentDetail, DocumentStatus, DocumentSummary } from '../docs'

export interface DocumentMutation {
  slug: string
  title: string
  status: DocumentStatus
  sort_order: number
  content: string
}

export interface DocumentImageUpload {
  filename: string
  markdown: string
}

export async function list(): Promise<DocumentSummary[]> {
  const { data } = await apiClient.get<DocumentSummary[] | null>('/admin/docs')
  return Array.isArray(data) ? data : []
}

export async function get(slug: string): Promise<DocumentDetail> {
  const { data } = await apiClient.get<DocumentDetail>(`/admin/docs/${encodeURIComponent(slug)}`)
  return data
}

export async function create(payload: DocumentMutation): Promise<DocumentDetail> {
  const { data } = await apiClient.post<DocumentDetail>('/admin/docs', payload)
  return data
}

export async function update(slug: string, payload: DocumentMutation): Promise<DocumentDetail> {
  const { data } = await apiClient.put<DocumentDetail>(`/admin/docs/${encodeURIComponent(slug)}`, payload)
  return data
}

export async function remove(slug: string): Promise<{ message: string }> {
  const { data } = await apiClient.delete<{ message: string }>(`/admin/docs/${encodeURIComponent(slug)}`)
  return data
}

export async function uploadImage(slug: string, file: File): Promise<DocumentImageUpload> {
  const form = new FormData()
  form.append('file', file)
  const { data } = await apiClient.post<DocumentImageUpload>(
    `/admin/docs/${encodeURIComponent(slug)}/images`,
    form,
    { headers: { 'Content-Type': 'multipart/form-data' } },
  )
  return data
}

export default { list, get, create, update, remove, uploadImage }
