<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <div class="relative w-full md:w-72">
            <Icon name="search" size="md" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input
              v-model="searchQuery"
              type="text"
              class="input pl-10"
              :placeholder="t('admin.teams.searchPlaceholder')"
              @input="scheduleLoadTeams"
            />
          </div>
          <div class="w-full sm:w-36">
            <Select v-model="statusFilter" :options="statusOptions" @change="loadTeams" />
          </div>
          <div class="ml-auto flex items-center gap-2">
            <button class="btn btn-secondary px-3" :disabled="loading" :title="t('common.refresh')" @click="loadTeams">
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
            </button>
            <button class="btn btn-primary" @click="openCreateDialog">
              <Icon name="plus" size="md" class="mr-2" />
              {{ t('admin.teams.create') }}
            </button>
          </div>
        </div>
      </template>

      <template #table>
        <DataTable :columns="columns" :data="teams" :loading="loading">
          <template #cell-name="{ row }">
            <div class="min-w-[220px]">
              <div class="font-medium text-gray-900 dark:text-white">{{ row.name }}</div>
              <div class="mt-1 max-w-md whitespace-normal text-xs text-gray-500 dark:text-gray-400">
                {{ row.description || t('admin.teams.noDescription') }}
              </div>
            </div>
          </template>

          <template #cell-status="{ row }">
            <span
              class="inline-flex items-center rounded px-2 py-0.5 text-xs font-medium"
              :class="row.status === 'active' ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300' : 'bg-gray-100 text-gray-700 dark:bg-dark-700 dark:text-gray-300'"
            >
              {{ row.status === 'active' ? t('common.active') : t('common.inactive') }}
            </span>
          </template>

          <template #cell-members="{ row }">
            <div class="text-sm text-gray-700 dark:text-gray-300">
              <span class="font-medium">{{ row.members?.length || 0 }}</span>
              <span class="mx-1 text-gray-400">/</span>
              <span>{{ leaderCount(row) }} {{ t('admin.teams.leaders') }}</span>
            </div>
          </template>

          <template #cell-created_at="{ value }">
            <span class="text-sm text-gray-600 dark:text-gray-300">{{ formatDateTime(value) }}</span>
          </template>

          <template #cell-actions="{ row }">
            <div class="flex items-center justify-end gap-2">
              <button class="btn btn-ghost btn-sm" :title="t('admin.teams.members')" @click="openMembersDialog(row)">
                <Icon name="users" size="sm" />
              </button>
              <button class="btn btn-ghost btn-sm" :title="t('common.edit')" @click="openEditDialog(row)">
                <Icon name="edit" size="sm" />
              </button>
              <button class="btn btn-ghost btn-sm text-red-600 hover:text-red-700" :title="t('common.delete')" @click="deleteTeam(row)">
                <Icon name="trash" size="sm" />
              </button>
            </div>
          </template>

          <template #empty>
            <div class="flex flex-col items-center py-8">
              <Icon name="users" size="xl" class="mb-3 text-gray-400" />
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('admin.teams.empty') }}</p>
            </div>
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

    <BaseDialog :show="showTeamDialog" :title="teamDialogTitle" @close="closeTeamDialog">
      <div class="space-y-4">
        <div>
          <label class="input-label">{{ t('admin.teams.name') }}</label>
          <input v-model="teamForm.name" class="input" maxlength="120" />
        </div>
        <div>
          <label class="input-label">{{ t('admin.teams.descriptionLabel') }}</label>
          <textarea v-model="teamForm.description" class="input min-h-[96px] resize-y" maxlength="500"></textarea>
        </div>
        <div>
          <label class="input-label">{{ t('common.status') }}</label>
          <Select v-model="teamForm.status" :options="statusEditOptions" />
        </div>
      </div>
      <template #footer>
        <button class="btn btn-secondary" @click="closeTeamDialog">{{ t('common.cancel') }}</button>
        <button class="btn btn-primary" :disabled="savingTeam" @click="saveTeam">{{ t('common.save') }}</button>
      </template>
    </BaseDialog>

    <BaseDialog :show="showMembersDialog" :title="membersDialogTitle" width="wide" @close="closeMembersDialog">
      <div class="space-y-5">
        <div class="grid gap-3 md:grid-cols-[1fr_auto]">
          <div class="relative">
            <Icon name="search" size="md" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input
              v-model="userSearch"
              class="input pl-10"
              :placeholder="t('admin.teams.searchUsers')"
              @input="scheduleUserSearch"
            />
          </div>
          <button class="btn btn-secondary" :disabled="searchingUsers" @click="searchUsers">
            {{ t('common.search') }}
          </button>
        </div>

        <div v-if="userSearchResults.length > 0" class="rounded-lg border border-gray-200 dark:border-dark-700">
          <button
            v-for="user in userSearchResults"
            :key="user.id"
            class="flex w-full items-center justify-between border-b border-gray-100 px-4 py-3 text-left last:border-b-0 hover:bg-gray-50 dark:border-dark-700 dark:hover:bg-dark-800"
            @click="addMember(user)"
          >
            <div class="min-w-0">
              <div class="truncate text-sm font-medium text-gray-900 dark:text-white">{{ user.email }}</div>
              <div class="truncate text-xs text-gray-500 dark:text-gray-400">{{ user.username || `#${user.id}` }}</div>
            </div>
            <Icon name="plus" size="sm" class="text-primary-500" />
          </button>
        </div>

        <div>
          <div class="mb-2 flex items-center justify-between">
            <h4 class="text-sm font-medium text-gray-900 dark:text-white">{{ t('admin.teams.selectedMembers') }}</h4>
            <span class="text-xs text-gray-500 dark:text-gray-400">{{ memberDrafts.length }} {{ t('admin.teams.members') }}</span>
          </div>
          <div class="overflow-hidden rounded-lg border border-gray-200 dark:border-dark-700">
            <div v-if="memberDrafts.length === 0" class="px-4 py-8 text-center text-sm text-gray-500 dark:text-gray-400">
              {{ t('admin.teams.noMembers') }}
            </div>
            <div
              v-for="member in memberDrafts"
              :key="member.user_id"
              class="grid grid-cols-[1fr_140px_40px] items-center gap-3 border-b border-gray-100 px-4 py-3 last:border-b-0 dark:border-dark-700"
            >
              <div class="min-w-0">
                <div class="truncate text-sm font-medium text-gray-900 dark:text-white">{{ memberLabel(member) }}</div>
                <div class="text-xs text-gray-500 dark:text-gray-400">#{{ member.user_id }}</div>
              </div>
              <Select v-model="member.role" :options="roleOptions" />
              <button class="btn btn-ghost btn-sm text-red-600" :title="t('common.delete')" @click="removeMember(member.user_id)">
                <Icon name="trash" size="sm" />
              </button>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <button class="btn btn-secondary" @click="closeMembersDialog">{{ t('common.cancel') }}</button>
        <button class="btn btn-primary" :disabled="savingMembers" @click="saveMembers">{{ t('common.save') }}</button>
      </template>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import type { Column } from '@/components/common/types'
import { teamsAPI } from '@/api/admin/teams'
import usersAPI from '@/api/admin/users'
import { useAppStore } from '@/stores'
import type { AdminUser, Team, TeamMemberInput, TeamRole, User } from '@/types'
import { formatDateTime } from '@/utils/format'

interface MemberDraft extends TeamMemberInput {
  user?: User | AdminUser
}

const { t } = useI18n()
const appStore = useAppStore()

const teams = ref<Team[]>([])
const loading = ref(false)
const savingTeam = ref(false)
const savingMembers = ref(false)
const searchingUsers = ref(false)
const searchQuery = ref('')
const statusFilter = ref<string | null>('')
const showTeamDialog = ref(false)
const showMembersDialog = ref(false)
const editingTeam = ref<Team | null>(null)
const membersTeam = ref<Team | null>(null)
const userSearch = ref('')
const userSearchResults = ref<AdminUser[]>([])
const memberDrafts = ref<MemberDraft[]>([])
let loadTimer: ReturnType<typeof setTimeout> | null = null
let userSearchTimer: ReturnType<typeof setTimeout> | null = null

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0,
  pages: 0
})

const teamForm = reactive({
  name: '',
  description: '',
  status: 'active' as 'active' | 'inactive'
})

const columns = computed<Column[]>(() => [
  { key: 'name', label: t('admin.teams.name'), sortable: false },
  { key: 'status', label: t('common.status'), sortable: false },
  { key: 'members', label: t('admin.teams.members'), sortable: false },
  { key: 'created_at', label: t('common.createdAt'), sortable: false },
  { key: 'actions', label: t('common.actions'), sortable: false, class: 'text-right' }
])

const statusOptions = computed(() => [
  { value: '', label: t('admin.teams.allStatus') },
  { value: 'active', label: t('common.active') },
  { value: 'inactive', label: t('common.inactive') }
])

const statusEditOptions = computed(() => [
  { value: 'active', label: t('common.active') },
  { value: 'inactive', label: t('common.inactive') }
])

const roleOptions = computed(() => [
  { value: 'member', label: t('admin.teams.member') },
  { value: 'leader', label: t('admin.teams.leader') }
])

const teamDialogTitle = computed(() => (editingTeam.value ? t('admin.teams.edit') : t('admin.teams.create')))
const membersDialogTitle = computed(() => (membersTeam.value ? `${t('admin.teams.members')} · ${membersTeam.value.name}` : t('admin.teams.members')))

function normalizeError(error: unknown): string {
  if (typeof error === 'object' && error !== null && 'message' in error) {
    return String((error as { message?: unknown }).message || t('common.error'))
  }
  return t('common.error')
}

async function loadTeams(): Promise<void> {
  loading.value = true
  try {
    const result = await teamsAPI.list(pagination.page, pagination.page_size, {
      search: searchQuery.value.trim() || undefined,
      status: statusFilter.value || undefined,
      sort_by: 'created_at',
      sort_order: 'desc'
    })
    teams.value = result.items
    pagination.total = result.total
    pagination.pages = result.pages
  } catch (error) {
    appStore.showError(normalizeError(error))
  } finally {
    loading.value = false
  }
}

function scheduleLoadTeams(): void {
  if (loadTimer) clearTimeout(loadTimer)
  loadTimer = setTimeout(() => {
    pagination.page = 1
    loadTeams()
  }, 250)
}

function handlePageChange(page: number): void {
  pagination.page = page
  loadTeams()
}

function handlePageSizeChange(pageSize: number): void {
  pagination.page_size = pageSize
  pagination.page = 1
  loadTeams()
}

function openCreateDialog(): void {
  editingTeam.value = null
  teamForm.name = ''
  teamForm.description = ''
  teamForm.status = 'active'
  showTeamDialog.value = true
}

function openEditDialog(team: Team): void {
  editingTeam.value = team
  teamForm.name = team.name
  teamForm.description = team.description || ''
  teamForm.status = team.status === 'inactive' ? 'inactive' : 'active'
  showTeamDialog.value = true
}

function closeTeamDialog(): void {
  showTeamDialog.value = false
}

async function saveTeam(): Promise<void> {
  const name = teamForm.name.trim()
  if (!name) {
    appStore.showError(t('admin.teams.nameRequired'))
    return
  }
  savingTeam.value = true
  try {
    const payload = {
      name,
      description: teamForm.description.trim() || null,
      status: teamForm.status
    }
    if (editingTeam.value) {
      await teamsAPI.update(editingTeam.value.id, payload)
      appStore.showSuccess(t('admin.teams.updated'))
    } else {
      await teamsAPI.create(payload)
      appStore.showSuccess(t('admin.teams.created'))
    }
    closeTeamDialog()
    loadTeams()
  } catch (error) {
    appStore.showError(normalizeError(error))
  } finally {
    savingTeam.value = false
  }
}

function leaderCount(team: Team): number {
  return team.members?.filter((member) => member.role === 'leader').length || 0
}

function openMembersDialog(team: Team): void {
  membersTeam.value = team
  memberDrafts.value = (team.members || []).map((member) => ({
    user_id: member.user_id,
    role: member.role,
    user: member.user
  }))
  userSearch.value = ''
  userSearchResults.value = []
  showMembersDialog.value = true
}

function closeMembersDialog(): void {
  showMembersDialog.value = false
}

function scheduleUserSearch(): void {
  if (userSearchTimer) clearTimeout(userSearchTimer)
  userSearchTimer = setTimeout(searchUsers, 250)
}

async function searchUsers(): Promise<void> {
  const query = userSearch.value.trim()
  if (!query) {
    userSearchResults.value = []
    return
  }
  searchingUsers.value = true
  try {
    const result = await usersAPI.list(1, 20, {
      search: query,
      include_subscriptions: false
    })
    userSearchResults.value = result.items
  } catch (error) {
    appStore.showError(normalizeError(error))
  } finally {
    searchingUsers.value = false
  }
}

function addMember(user: AdminUser): void {
  if (memberDrafts.value.some((member) => member.user_id === user.id)) return
  memberDrafts.value.push({ user_id: user.id, role: 'member', user })
}

function removeMember(userID: number): void {
  memberDrafts.value = memberDrafts.value.filter((member) => member.user_id !== userID)
}

function memberLabel(member: MemberDraft): string {
  return member.user?.email || member.user?.username || `#${member.user_id}`
}

async function saveMembers(): Promise<void> {
  if (!membersTeam.value) return
  savingMembers.value = true
  try {
    await teamsAPI.replaceMembers(membersTeam.value.id, {
      members: memberDrafts.value.map((member) => ({
        user_id: member.user_id,
        role: member.role as TeamRole
      }))
    })
    appStore.showSuccess(t('admin.teams.membersUpdated'))
    closeMembersDialog()
    loadTeams()
  } catch (error) {
    appStore.showError(normalizeError(error))
  } finally {
    savingMembers.value = false
  }
}

async function deleteTeam(team: Team): Promise<void> {
  if (!window.confirm(t('admin.teams.deleteConfirm', { name: team.name }))) return
  try {
    await teamsAPI.delete(team.id)
    appStore.showSuccess(t('admin.teams.deleted'))
    loadTeams()
  } catch (error) {
    appStore.showError(normalizeError(error))
  }
}

onMounted(loadTeams)
</script>
