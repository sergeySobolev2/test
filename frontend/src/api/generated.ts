import { http, ApiListResponse, Calculation } from '../modules/services/http'

export const RequestsService = {
  async cart(): Promise<{ draft_id: number | null; count: number }> {
    const res = await http.get('/api/requests/cart')
    return res.data as { data: { draft_id: number | null; count: number }, total: number, filters: Record<string, unknown> } as any
  },
  async list(filters: { status?: string; from?: string; to?: string }): Promise<ApiListResponse<Calculation[]>> {
    const res = await http.get('/api/calculations', { params: filters })
    return res.data as ApiListResponse<Calculation[]>
  },
  async get(id: number): Promise<any> {
    const res = await http.get(`/api/requests/${id}`)
    return res.data
  },
  async update(id: number, payload: { room_area?: number; expert_comment?: string }): Promise<any> {
    const res = await http.put(`/api/requests/${id}`, payload)
    return res.data
  },
  async complete(id: number, payload: { status: 'завершен' | 'отклонен'; room_area?: number; noise_reduction_db?: number; expert_comment?: string }): Promise<any> {
    const res = await http.put(`/api/requests/${id}/complete`, payload)
    return res.data
  },
  async form(id: number): Promise<any> {
    const res = await http.put(`/api/requests/${id}/form`)
    return res.data
  },
  async remove(id: number): Promise<any> {
    const res = await http.delete(`/api/requests/${id}`)
    return res.data
  },
}

export const RequestPartitionsService = {
  async add(partitionId: number, calculationId?: number | null): Promise<{ added: number; calculation_id: number }> {
    const res = await http.post('/api/request-partitions', { partition_id: partitionId, calculation_id: calculationId || undefined })
    return res.data as { data: { added: number }, total: number, filters: { calculation_id: number } } as any
  },
  async remove(calculationId: number, partitionId: number): Promise<{ deleted: { calculation_id: number; partition_id: number } }> {
    const res = await http.delete('/api/request-partitions', { data: { calculation_id: calculationId, partition_id: partitionId } })
    return res.data as { data: { deleted: { calculation_id: number; partition_id: number } } } as any
  },
  async update(calculationId: number, partitionId: number, payload: { quantity?: number; comment?: string; is_main?: boolean }): Promise<{ updated: number }> {
    const res = await http.put('/api/request-partitions', { calculation_id: calculationId, partition_id: partitionId, ...payload })
    return res.data as { data: { updated: number } } as any
  },
}
