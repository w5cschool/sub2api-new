import { apiClient } from '../client'
import type {
  CreateTeamRequest,
  PaginatedResponse,
  ReplaceTeamMembersRequest,
  Team,
  UpdateTeamRequest
} from '@/types'

export async function list(
  page = 1,
  pageSize = 20,
  filters?: {
    search?: string
    status?: string
    sort_by?: string
    sort_order?: 'asc' | 'desc'
  },
  options?: { signal?: AbortSignal }
): Promise<PaginatedResponse<Team>> {
  const { data } = await apiClient.get<PaginatedResponse<Team>>('/admin/teams', {
    params: {
      page,
      page_size: pageSize,
      ...filters
    },
    signal: options?.signal
  })
  return data
}

export async function getById(id: number): Promise<Team> {
  const { data } = await apiClient.get<Team>(`/admin/teams/${id}`)
  return data
}

export async function create(payload: CreateTeamRequest): Promise<Team> {
  const { data } = await apiClient.post<Team>('/admin/teams', payload)
  return data
}

export async function update(id: number, payload: UpdateTeamRequest): Promise<Team> {
  const { data } = await apiClient.put<Team>(`/admin/teams/${id}`, payload)
  return data
}

export async function deleteTeam(id: number): Promise<{ message: string }> {
  const { data } = await apiClient.delete<{ message: string }>(`/admin/teams/${id}`)
  return data
}

export async function replaceMembers(id: number, payload: ReplaceTeamMembersRequest): Promise<Team> {
  const { data } = await apiClient.put<Team>(`/admin/teams/${id}/members`, payload)
  return data
}

export const teamsAPI = {
  list,
  getById,
  create,
  update,
  delete: deleteTeam,
  replaceMembers
}

export default teamsAPI
