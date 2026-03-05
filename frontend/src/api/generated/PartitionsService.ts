import type { ApiListResponse } from '../http'
import type { partition } from './types'
import { SymptomsService } from './SymptomsService'

export const PartitionsService = {
  async list(params: { title?: string; active?: boolean } = {}): Promise<ApiListResponse<partition[]>> {
    return SymptomsService.list(params)
  },
  async get(id: number): Promise<ApiListResponse<{ partition: partition; public_image_url?: string }>> {
    return SymptomsService.get(id)
  }
}
