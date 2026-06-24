import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { teamAPI } from '@/api'
import type { TeamMe } from '@/types'

export const useTeamStore = defineStore('team', () => {
  const me = ref<TeamMe | null>(null)
  const loaded = ref(false)
  const loading = ref(false)

  const canViewUsage = computed(() => me.value?.can_view_usage === true)

  async function fetchMe(force = false): Promise<TeamMe | null> {
    if (loading.value) return me.value
    if (loaded.value && !force) return me.value
    loading.value = true
    try {
      me.value = await teamAPI.getMe()
      loaded.value = true
      return me.value
    } catch (error) {
      me.value = null
      loaded.value = true
      throw error
    } finally {
      loading.value = false
    }
  }

  function clear(): void {
    me.value = null
    loaded.value = false
    loading.value = false
  }

  return {
    me,
    loaded,
    loading,
    canViewUsage,
    fetchMe,
    clear
  }
})
