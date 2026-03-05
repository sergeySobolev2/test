import { api, ApiListResponse } from '../http'
import type { partition } from './types'

export const SymptomsService = {
  async list(params: { title?: string; active?: boolean } = {}): Promise<ApiListResponse<partition[]>> {
    const search = new URLSearchParams()
    if (params.title) search.set('title', params.title)
    if (typeof params.active === 'boolean') search.set('active', String(params.active))
    const { data } = await api.get(`/api/partitions?${search.toString()}`)
    return data
  },
  async get(id: number): Promise<ApiListResponse<{ partition: partition; public_image_url?: string }>> {
    const { data } = await api.get(`/api/partitions/${id}`)
    return data
  }
}
