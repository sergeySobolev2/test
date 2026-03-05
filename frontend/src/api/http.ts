import axios from 'axios'
import { store } from '../store'

const baseURL = (import.meta as any).env?.VITE_API_BASE || ''

export const api = axios.create({
  baseURL: baseURL || undefined,
  withCredentials: true,
})

api.interceptors.request.use((config) => {
  const state = store.getState() as any
  const token: string | undefined = state?.auth?.token || (typeof localStorage !== 'undefined' ? localStorage.getItem('app:token') || undefined : undefined)
  if (token) {
    config.headers = config.headers || {}
    ;(config.headers as any)['Authorization'] = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (r) => r,
  (error) => {
    // Если 401 — сбрасываем аутентификацию (мягко)
    if (error?.response?.status === 401) {
      try { store.dispatch({ type: 'auth/forceLoggedOut' }) } catch {}
    }
    return Promise.reject(error)
  }
)

export type ApiListResponse<T> = {
  data: T
  total: number
  filters?: Record<string, unknown>
}
