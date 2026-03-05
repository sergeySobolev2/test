import axios from 'axios'

const API_BASE = (import.meta as any).env?.VITE_API_BASE || ''

export const http = axios.create({
  baseURL: API_BASE || '/',
  withCredentials: true,
})

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('app:token')
  if (token) {
    config.headers = config.headers || {}
    ;(config.headers as any)['Authorization'] = `Bearer ${token}`
  }
  return config
})

export type ApiListResponse<T> = {
  data: T
  total: number
  filters?: Record<string, unknown>
}

export type User = {
  id: number
  login: string
  is_moderator: boolean
}

export type Calculation = {
  id: number
  user_id: number
  status: string
  created_at: string
  formed_at?: string | null
  completed_at?: string | null
  moderator_id?: number | null
  room_area?: number | null
  Soundproofing_percent?: number | null
  required_thickness?: number | null
  expert_comment?: { string: string; valid: boolean }
}
