import { beforeEach, describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'

import KeyUsageView from '../KeyUsageView.vue'

const { checkAuth, fetchPublicSettings } = vi.hoisted(() => ({
  checkAuth: vi.fn(),
  fetchPublicSettings: vi.fn(),
}))

const messages: Record<string, string> = {
  'home.viewDocs': 'Documentation',
  'keyUsage.accessTitle': 'Usage requires sign-in.',
  'keyUsage.accessDescription': 'Usage is available only in the authenticated console.',
  'keyUsage.accessLogin': 'Sign in',
  'keyUsage.accessDashboard': 'Open console',
  'keyUsage.securityTitle': 'Credential safety',
  'keyUsage.securityDescription': 'Do not enter credentials on a public page.',
  'keyUsage.securityReport': 'Report suspicious requests through the published operator contact.',
}

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({ t: (key: string) => messages[key] ?? key }),
  }
})

vi.mock('@/stores', () => ({
  useAppStore: () => ({
    cachedPublicSettings: null,
    siteName: 'Sub2API',
    siteLogo: '',
    docUrl: '',
    publicSettingsLoaded: true,
    fetchPublicSettings,
  }),
  useAuthStore: () => ({
    isAuthenticated: false,
    isAdmin: false,
    checkAuth,
  }),
}))

describe('KeyUsageView', () => {
  beforeEach(() => {
    checkAuth.mockReset()
    fetchPublicSettings.mockReset()
  })

  it('does not collect or submit API keys on the public route', () => {
    const wrapper = mount(KeyUsageView, {
      global: {
        stubs: {
          RouterLink: { template: '<a><slot /></a>' },
          LocaleSwitcher: true,
        },
      },
    })

    expect(wrapper.find('input').exists()).toBe(false)
    expect(wrapper.find('textarea').exists()).toBe(false)
    expect(wrapper.text()).toContain('Usage requires sign-in.')
    expect(wrapper.text()).toContain('Credential safety')
    expect(checkAuth).toHaveBeenCalledOnce()
  })
})
