import { http, ApiListResponse } from './http'
import { isGuestMode } from '../../config/appConfig'

export type Partition = {
  id: number
  title: string
  category?: string
  description?: string
  noise_reduction?: string
  thickness?: string
  material?: string
  price_per_sqm?: string
  image_url?: string
  public_image_url?: string
  is_active: boolean
}

const MOCK_PARTITIONS: Partition[] = [
  {
    id: 1,
    title: 'Кирпичная перегородка',
    category: 'Строительная',
    description: 'Прочная кирпичная перегородка с отличной звукоизоляцией',
    noise_reduction: '55-60 дБ',
    thickness: '120 мм',
    material: 'Кирпич',
    price_per_sqm: '2500 ₽',
    public_image_url: '/partition_front/partitions/кирпичная.png',
    is_active: true
  },
  {
    id: 2,
    title: 'Гипсокартонная перегородка',
    category: 'Легкая конструкция',
    description: 'Легкая и быстромонтируемая гипсокартонная перегородка',
    noise_reduction: '45-50 дБ',
    thickness: '100 мм',
    material: 'Гипсокартон',
    price_per_sqm: '1800 ₽',
    public_image_url: '/partition_front/partitions/гипсокартонная.png',
    is_active: true
  },
  {
    id: 3,
    title: 'Газобетонная перегородка',
    category: 'Строительная',
    description: 'Современная газобетонная перегородка с хорошей теплоизоляцией',
    noise_reduction: '50-55 дБ',
    thickness: '100 мм',
    material: 'Газобетон',
    price_per_sqm: '2000 ₽',
    public_image_url: '/partition_front/partitions/газобетонная.jpg',
    is_active: true
  }
]

export async function listPartitions(params: { title?: string; active?: boolean } = {}): Promise<ApiListResponse<Partition[]>> {
  if (isGuestMode()) {
    // Mock mode для GitHub Pages и Tauri
    let filtered = MOCK_PARTITIONS
    if (params.title) {
      const searchLower = params.title.toLowerCase()
      filtered = MOCK_PARTITIONS.filter(p => 
        p.title.toLowerCase().includes(searchLower) ||
        p.description?.toLowerCase().includes(searchLower)
      )
    }
    return Promise.resolve({
      data: filtered,
      total: filtered.length,
      message: 'Mock data'
    })
  }
  
  const res = await http.get('/api/partitions', { params })
  return res.data as ApiListResponse<Partition[]>
}

export async function getPartition(id: number): Promise<ApiListResponse<{ partition: Partition; public_image_url?: string }>> {
  if (isGuestMode()) {
    const partition = MOCK_PARTITIONS.find(p => p.id === id)
    if (!partition) {
      return Promise.reject(new Error('Partition not found'))
    }
    return Promise.resolve({
      data: { partition, public_image_url: partition.public_image_url },
      total: 1,
      message: 'Mock data'
    })
  }
  
  const res = await http.get(`/api/partitions/${id}`)
  return res.data as ApiListResponse<{ partition: Partition; public_image_url?: string }>
}

