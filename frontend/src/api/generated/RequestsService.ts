import { api, ApiListResponse } from '../http'
import type { Calculation, RequestWithCommentsCount, UpdateCalculationPayload } from './types'

export const RequestsService = {
  async cart(): Promise<{ draft_id: number | null, count: number }> {
    const { data } = await api.get('/api/requests/cart')
    return data.data
  },
  async list(filters: { status?: string; from?: string; to?: string } = {}): Promise<ApiListResponse<RequestWithCommentsCount[]>> {
    const search = new URLSearchParams()
    if (filters.status) search.set('status', filters.status)
    if (filters.from) search.set('from', filters.from)
    if (filters.to) search.set('to', filters.to)
    const { data } = await api.get(`/api/calculations?${search.toString()}`)
    return data
  },
  async get(id: number): Promise<ApiListResponse<unknown>> {
    const { data } = await api.get(`/api/requests/${id}`)
    return data
  },
  async update(
    id: number,
    payload: UpdateCalculationPayload
  ): Promise<ApiListResponse<Calculation>> {
    const { data } = await api.put(`/api/requests/${id}`, payload)
    return data
  },
  async form(id: number): Promise<ApiListResponse<Calculation>> {
    const { data } = await api.put(`/api/requests/${id}/form`)
    return data
  },
  async remove(id: number): Promise<ApiListResponse<{ deleted: number }>> {
    const { data } = await api.delete(`/api/requests/${id}`)
    return data
  }
}
