export type partition = {
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

type ApiListResponse<T> = {
  data: T
  total: number
  filters?: Record<string, unknown>
}

async function safeFetch<T>(input: RequestInfo, init?: RequestInit, fallback?: () => Promise<T>): Promise<T> {
  try {
    const res = await fetch(input, init)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    return await res.json()
  } catch (e) {
    if (fallback) return await fallback()
    throw e
  }
}

const API_BASE = (import.meta as any).env?.VITE_API_BASE || ''
const USE_MOCK = ((import.meta as any).env?.VITE_USE_MOCK || '').toString().toLowerCase() === 'true'

function buildUrl(path: string) {
  if (!API_BASE) return path // dev: используем proxy
  return API_BASE.replace(/\/$/, '') + path
}

function filterMock(list: partition[], params: { title?: string; active?: boolean }) {
  const q = (params.title || '').trim().toLowerCase()
  return list.filter(x => {
    const matchTitle = q ? ((x.title || '').toLowerCase().includes(q) || (x.description || '').toLowerCase().includes(q)) : true
    const matchActive = typeof params.active === 'boolean' ? x.is_active === params.active : true
    return matchTitle && matchActive
  })
}

export async function listPartitions(params: { title?: string; active?: boolean } = {}): Promise<ApiListResponse<partition[]>> {
  const search = new URLSearchParams()
  if (params.title) search.set('title', params.title)
  if (typeof params.active === 'boolean') search.set('active', String(params.active))
  const path = `/api/partitions?${search.toString()}`
  const url = buildUrl(path)

  if (USE_MOCK) {
    const data = filterMock(mockPartitions, params)
    return { data, total: data.length, filters: Object.fromEntries(search.entries()) }
  }

  return safeFetch<ApiListResponse<partition[]>>(
    url,
    { headers: { 'Accept': 'application/json' } },
    async () => {
      const data = filterMock(mockPartitions, params)
      return { data, total: data.length, filters: Object.fromEntries(search.entries()) }
    }
  )
}

export async function getPartition(id: number): Promise<ApiListResponse<{ partition: partition; public_image_url?: string }>> {
  const url = buildUrl(`/api/partitions/${id}`)
  if (USE_MOCK) {
    const s = mockPartitions.find(x => x.id === id)
    if (!s) throw new Error('not found')
    return { data: { partition: s, public_image_url: s.public_image_url }, total: 1, filters: { id } }
  }
  return safeFetch<ApiListResponse<{ partition: partition; public_image_url?: string }>>(
    url,
    { headers: { 'Accept': 'application/json' } },
    async () => {
      const s = mockPartitions.find(x => x.id === id)
      if (!s) throw new Error('not found')
      return { data: { partition: s, public_image_url: s.public_image_url }, total: 1, filters: { id } }
    }
  )
}

// Mock data fallback
export const mockPartitions: partition[] = [
  {
    id: 1,
    title: 'Гипсокартон ГКЛ 12.5мм + минвата 50мм',
    category: 'Легкие конструкции',
    description: 'Легкая каркасная перегородка с хорошей звукоизоляцией для офисов и квартир.',
    is_active: true,
    image_url: '/partitions/гипсокартонная.png'
  },
  {
    id: 2,
    title: 'Кирпичная кладка 120мм',
    category: 'Капитальные',
    description: 'Капитальная перегородка из керамического кирпича.',
    is_active: true,
    image_url: '/partitions/кирпичная.png'
  },
  {
    id: 3,
    title: 'Газобетонные блоки D500 100мм',
    category: 'Блочные',
    description: 'Современный материал с хорошей теплоизоляцией.',
    is_active: true,
    image_url: '/partitions/газобетонная.jpg'
  },
  {
    id: 4,
    title: 'Сэндвич-панели акустические',
    category: 'Специализированные',
    description: 'Многослойные панели с акустическим наполнителем.',
    is_active: true,
    image_url: '/partitions/сэндвич-панель.jpg'
  },
  {
    id: 5,
    title: 'Двойная каркасная перегородка с виброподвесами',
    category: 'Профессиональные',
    description: 'Максимальная звукоизоляция для домашних кинотеатров.',
    is_active: true,
    image_url: '/partitions/двойная.png'
  }
]
