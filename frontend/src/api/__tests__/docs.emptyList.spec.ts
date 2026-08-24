import { beforeEach, describe, expect, it, vi } from 'vitest'

const { get } = vi.hoisted(() => ({
  get: vi.fn(),
}))

vi.mock('../client', () => ({
  apiClient: { get },
}))

import { list as listAdminDocuments } from '@/api/admin/docs'
import { listDocuments } from '@/api/docs'

describe('docs list APIs', () => {
  beforeEach(() => {
    get.mockReset()
  })

  it('normalizes a null admin document list to an empty array', async () => {
    get.mockResolvedValue({ data: null })

    await expect(listAdminDocuments()).resolves.toEqual([])
    expect(get).toHaveBeenCalledWith('/admin/docs')
  })

  it('normalizes a null public document list to an empty array', async () => {
    get.mockResolvedValue({ data: null })

    await expect(listDocuments()).resolves.toEqual([])
    expect(get).toHaveBeenCalledWith('/docs')
  })
})
