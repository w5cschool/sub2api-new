import { apiClient } from './client'
import type {
  PaginatedResponse,
  TeamMe,
  TeamMember,
  TeamMemberUsageSummary,
  TeamUsageQueryParams,
  UsageLog,
  UsageStatsResponse
} from '@/types'

export async function getMe(): Promise<TeamMe> {
  const { data } = await apiClient.get<TeamMe>('/team/me')
  return data
}

export async function listMembers(): Promise<TeamMember[]> {
  const { data } = await apiClient.get<TeamMember[]>('/team/members')
  return data
}

export async function queryUsage(
  params: TeamUsageQueryParams & { sort_by?: string; sort_order?: 'asc' | 'desc' },
  options: { signal?: AbortSignal } = {}
): Promise<PaginatedResponse<UsageLog>> {
  const { data } = await apiClient.get<PaginatedResponse<UsageLog>>('/team/usage', {
    params,
    signal: options.signal
  })
  return data
}

export async function getUsageStats(params: TeamUsageQueryParams = {}): Promise<UsageStatsResponse> {
  const { data } = await apiClient.get<UsageStatsResponse>('/team/usage/stats', { params })
  return data
}

export async function getMembersSummary(): Promise<TeamMemberUsageSummary[]> {
  const { data } = await apiClient.get<TeamMemberUsageSummary[]>('/team/usage/members-summary')
  return data
}

export const teamAPI = {
  getMe,
  listMembers,
  queryUsage,
  getUsageStats,
  getMembersSummary
}

export default teamAPI
