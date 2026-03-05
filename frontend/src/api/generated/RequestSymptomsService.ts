import { api, ApiListResponse } from '../http'

export const RequestPartitionsService = {
  async add(partition_id: number, calculation_id?: number): Promise<{ added: number, calculation_id: number }> {
    // Формируем payload явно, убирая undefined
    const payload: any = { partition_id }
    if (calculation_id !== undefined && calculation_id !== null) {
      payload.calculation_id = calculation_id
    }
    
    console.log('RequestPartitionsService.add called with:', { partition_id, calculation_id })
    console.log('Sending payload:', payload)
    
    const { data } = await api.post('/api/request-partitions', payload)
    
    console.log('Response data:', data)
    
    // filters contains calculation_id
    return { added: data?.data?.added ?? partition_id, calculation_id: data?.filters?.calculation_id ?? calculation_id ?? 0 }
  },
  async remove(calculation_id: number, partition_id: number): Promise<{ deleted: { calculation_id: number; partition_id: number } }> {
    const { data } = await api.delete('/api/request-partitions', { data: { calculation_id, partition_id } })
    return { deleted: { calculation_id, partition_id } }
  },
  async update(payload: { calculation_id: number; partition_id: number; intensity?: number; comment?: string; is_main?: boolean }): Promise<{ updated: number }> {
    const { data } = await api.put('/api/request-partitions', payload)
    return data?.data ?? { updated: payload.partition_id }
  }
}
