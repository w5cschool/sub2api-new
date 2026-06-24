<template>
  <AppLayout>
    <div v-if="!canViewUsage" class="mx-auto max-w-3xl p-6">
      <div class="rounded-lg border border-gray-200 bg-white p-10 text-center dark:border-dark-700 dark:bg-dark-900">
        <Icon name="users" size="xl" class="mx-auto mb-4 text-gray-400" />
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('teamUsage.unavailableTitle') }}</h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">{{ t('teamUsage.unavailableDescription') }}</p>
      </div>
    </div>

    <TablePageLayout v-else>
      <template #actions>
        <div class="grid gap-4 lg:grid-cols-[1fr_320px]">
          <div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
            <div class="card p-4">
              <div class="flex items-center gap-3">
                <div class="rounded-lg bg-sky-100 p-2 dark:bg-sky-900/30">
                  <Icon name="document" size="md" class="text-sky-600 dark:text-sky-400" />
                </div>
                <div>
                  <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('usage.totalRequests') }}</p>
                  <p class="text-xl font-bold text-gray-900 dark:text-white">{{ formatNumber(stats?.total_requests || 0) }}</p>
                </div>
              </div>
            </div>
            <div class="card p-4">
              <div class="flex items-center gap-3">
                <div class="rounded-lg bg-amber-100 p-2 dark:bg-amber-900/30">
                  <Icon name="cube" size="md" class="text-amber-600 dark:text-amber-400" />
                </div>
                <div>
                  <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('usage.totalTokens') }}</p>
                  <p class="text-xl font-bold text-gray-900 dark:text-white">{{ formatNumber(stats?.total_tokens || 0) }}</p>
                </div>
              </div>
            </div>
            <div class="card p-4">
              <div class="flex items-center gap-3">
                <div class="rounded-lg bg-emerald-100 p-2 dark:bg-emerald-900/30">
                  <Icon name="dollar" size="md" class="text-emerald-600 dark:text-emerald-400" />
                </div>
                <div>
                  <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('usage.totalCost') }}</p>
                  <p class="text-xl font-bold text-emerald-600 dark:text-emerald-400">{{ formatCurrency(stats?.total_actual_cost || 0) }}</p>
                </div>
              </div>
            </div>
            <div class="card p-4">
              <div class="flex items-center gap-3">
                <div class="rounded-lg bg-violet-100 p-2 dark:bg-violet-900/30">
                  <Icon name="clock" size="md" class="text-violet-600 dark:text-violet-400" />
                </div>
                <div>
                  <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('usage.avgDuration') }}</p>
                  <p class="text-xl font-bold text-gray-900 dark:text-white">{{ formatDuration(stats?.average_duration_ms || 0) }}</p>
                </div>
              </div>
            </div>
          </div>

          <div class="card p-4">
            <div class="mb-3 flex items-center justify-between">
              <h3 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('teamUsage.membersSummary') }}</h3>
              <span class="text-xs text-gray-500 dark:text-gray-400">{{ membersSummary.length }}</span>
            </div>
            <div class="max-h-40 space-y-2 overflow-y-auto">
              <div v-for="item in membersSummary.slice(0, 5)" :key="item.user_id" class="flex items-center justify-between gap-3">
                <div class="min-w-0">
                  <div class="truncate text-sm font-medium text-gray-900 dark:text-white">{{ item.email || item.username || `#${item.user_id}` }}</div>
                  <div class="text-xs text-gray-500 dark:text-gray-400">{{ formatNumber(item.total_requests) }} {{ t('teamUsage.requests') }}</div>
                </div>
                <div class="text-right text-sm font-semibold text-emerald-600 dark:text-emerald-400">
                  {{ formatCurrency(item.total_actual_cost) }}
                </div>
              </div>
              <div v-if="membersSummary.length === 0" class="py-4 text-center text-sm text-gray-500 dark:text-gray-400">
                {{ t('teamUsage.noSummary') }}
              </div>
            </div>
          </div>
        </div>
      </template>

      <template #filters>
        <div class="card p-4">
          <div class="flex flex-wrap items-end gap-4">
            <div class="w-full sm:w-60">
              <label class="input-label">{{ t('teamUsage.member') }}</label>
              <Select v-model="filters.member_id" :options="memberOptions" @change="applyFilters" />
            </div>
            <div>
              <label class="input-label">{{ t('usage.timeRange') }}</label>
              <DateRangePicker v-model:start-date="startDate" v-model:end-date="endDate" @change="onDateRangeChange" />
            </div>
            <div class="w-full sm:w-52">
              <label class="input-label">{{ t('usage.model') }}</label>
              <input v-model="filters.model" class="input" :placeholder="t('teamUsage.modelPlaceholder')" @keyup.enter="applyFilters" />
            </div>
            <div class="ml-auto flex items-center gap-2">
              <button class="btn btn-secondary" :disabled="loading" @click="applyFilters">{{ t('common.refresh') }}</button>
              <button class="btn btn-secondary" @click="resetFilters">{{ t('common.reset') }}</button>
            </div>
          </div>
        </div>
      </template>

      <template #table>
        <DataTable
          :columns="columns"
          :data="logs"
          :loading="loading"
          :server-side-sort="true"
          default-sort-key="created_at"
          default-sort-order="desc"
          @sort="handleSort"
        >
          <template #cell-member="{ row }">
            <div class="min-w-[180px]">
              <div class="truncate font-medium text-gray-900 dark:text-white">{{ memberName(row.user_id, row.user) }}</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">#{{ row.user_id }}</div>
            </div>
          </template>
          <template #cell-model="{ value }">
            <span class="font-medium text-gray-900 dark:text-white">{{ value || '-' }}</span>
          </template>
          <template #cell-request_type="{ row }">
            <span class="inline-flex items-center rounded px-2 py-0.5 text-xs font-medium bg-gray-100 text-gray-700 dark:bg-dark-700 dark:text-gray-300">
              {{ requestTypeLabel(row) }}
            </span>
          </template>
          <template #cell-tokens="{ row }">
            <span class="font-medium text-gray-900 dark:text-white">{{ formatNumber(totalTokens(row)) }}</span>
          </template>
          <template #cell-cost="{ row }">
            <span class="font-medium text-emerald-600 dark:text-emerald-400">{{ formatCurrency(row.actual_cost || 0) }}</span>
          </template>
          <template #cell-duration_ms="{ value }">
            <span>{{ formatDuration(value) }}</span>
          </template>
          <template #cell-created_at="{ value }">
            <span class="text-gray-600 dark:text-gray-300">{{ formatDateTime(value) }}</span>
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination
          v-if="pagination.total > 0"
          :page="pagination.page"
          :total="pagination.total"
          :page-size="pagination.page_size"
          @update:page="handlePageChange"
          @update:pageSize="handlePageSizeChange"
        />
      </template>
    </TablePageLayout>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import Select from '@/components/common/Select.vue'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import Icon from '@/components/icons/Icon.vue'
import type { Column } from '@/components/common/types'
import { teamAPI } from '@/api/team'
import { useAppStore, useTeamStore } from '@/stores'
import type { TeamMember, TeamMemberUsageSummary, TeamUsageQueryParams, UsageLog, UsageStatsResponse, User } from '@/types'
import { formatCurrency, formatDateTime } from '@/utils/format'

const { t } = useI18n()
const appStore = useAppStore()
const teamStore = useTeamStore()

const members = ref<TeamMember[]>([])
const membersSummary = ref<TeamMemberUsageSummary[]>([])
const stats = ref<UsageStatsResponse | null>(null)
const logs = ref<UsageLog[]>([])
const loading = ref(false)

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0,
  pages: 0
})

const sortState = reactive({
  sort_by: 'created_at',
  sort_order: 'desc' as 'asc' | 'desc'
})

const formatLocalDate = (date: Date): string => {
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

const now = new Date()
const weekAgo = new Date(now)
weekAgo.setDate(weekAgo.getDate() - 6)
const startDate = ref(formatLocalDate(weekAgo))
const endDate = ref(formatLocalDate(now))

const filters = ref<TeamUsageQueryParams>({
  member_id: undefined,
  model: '',
  start_date: startDate.value,
  end_date: endDate.value
})

const canViewUsage = computed(() => teamStore.canViewUsage)

const columns = computed<Column[]>(() => [
  { key: 'member', label: t('teamUsage.member'), sortable: false },
  { key: 'model', label: t('usage.model'), sortable: true },
  { key: 'request_type', label: t('usage.type'), sortable: false },
  { key: 'tokens', label: t('usage.tokens'), sortable: false },
  { key: 'cost', label: t('usage.cost'), sortable: false },
  { key: 'duration_ms', label: t('usage.duration'), sortable: true },
  { key: 'created_at', label: t('usage.time'), sortable: true }
])

const memberOptions = computed(() => [
  { value: null, label: t('teamUsage.allMembers') },
  ...members.value.map((member) => ({
    value: member.user_id,
    label: member.user?.email || member.user?.username || `#${member.user_id}`
  }))
])

function normalizeError(error: unknown): string {
  if (typeof error === 'object' && error !== null && 'message' in error) {
    return String((error as { message?: unknown }).message || t('common.error'))
  }
  return t('common.error')
}

function currentParams(): TeamUsageQueryParams {
  return {
    member_id: filters.value.member_id || undefined,
    model: String(filters.value.model || '').trim() || undefined,
    start_date: startDate.value,
    end_date: endDate.value
  }
}

async function loadMembers(): Promise<void> {
  members.value = await teamAPI.listMembers()
}

async function loadSummary(): Promise<void> {
  membersSummary.value = await teamAPI.getMembersSummary()
}

async function loadStats(): Promise<void> {
  stats.value = await teamAPI.getUsageStats(currentParams())
}

async function loadUsage(): Promise<void> {
  const result = await teamAPI.queryUsage({
    ...currentParams(),
    page: pagination.page,
    page_size: pagination.page_size,
    sort_by: sortState.sort_by,
    sort_order: sortState.sort_order
  })
  logs.value = result.items
  pagination.total = result.total
  pagination.pages = result.pages
}

async function loadAll(): Promise<void> {
  if (!teamStore.canViewUsage) return
  loading.value = true
  try {
    await Promise.all([loadMembers(), loadSummary(), loadStats(), loadUsage()])
  } catch (error) {
    appStore.showError(normalizeError(error))
  } finally {
    loading.value = false
  }
}

function applyFilters(): void {
  pagination.page = 1
  loadAll()
}

function resetFilters(): void {
  filters.value.member_id = undefined
  filters.value.model = ''
  startDate.value = formatLocalDate(weekAgo)
  endDate.value = formatLocalDate(now)
  applyFilters()
}

function onDateRangeChange(range: { startDate: string; endDate: string }): void {
  startDate.value = range.startDate
  endDate.value = range.endDate
  applyFilters()
}

function handlePageChange(page: number): void {
  pagination.page = page
  loadAll()
}

function handlePageSizeChange(pageSize: number): void {
  pagination.page_size = pageSize
  pagination.page = 1
  loadAll()
}

function handleSort(key: string, order: 'asc' | 'desc'): void {
  sortState.sort_by = key
  sortState.sort_order = order
  loadAll()
}

function memberName(userID: number, user?: User): string {
  const member = members.value.find((item) => item.user_id === userID)
  return user?.email || member?.user?.email || user?.username || member?.user?.username || `#${userID}`
}

function totalTokens(row: UsageLog): number {
  return (row.input_tokens || 0) + (row.output_tokens || 0) + (row.cache_creation_tokens || 0) + (row.cache_read_tokens || 0) + (row.image_output_tokens || 0)
}

function requestTypeLabel(row: UsageLog): string {
  if (row.request_type) return row.request_type
  return row.stream ? 'stream' : 'sync'
}

function formatNumber(value: number): string {
  return Number(value || 0).toLocaleString()
}

function formatDuration(ms: number | null | undefined): string {
  if (!ms) return '-'
  if (ms < 1000) return `${Math.round(ms)}ms`
  return `${(ms / 1000).toFixed(2)}s`
}

onMounted(async () => {
  try {
    await teamStore.fetchMe(true)
    if (teamStore.canViewUsage) {
      await loadAll()
    }
  } catch (error) {
    appStore.showError(normalizeError(error))
  }
})
</script>
