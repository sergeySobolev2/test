import { api, ApiListResponse } from '../http'
import type { User } from './types'

export const AuthService = {
  async register(payload: { login: string; password: string }): Promise<{ user: User }> {
    const { data } = await api.post('/api/users/register', payload)
    return data
  },
  async login(payload: { login: string; password: string }): Promise<{ user: User; token: string; session_id: string }> {
    const { data } = await api.post('/api/users/login', payload)
    return data
  },
  async logout(): Promise<{ message: string }> {
    const { data } = await api.post('/api/users/logout')
    return data
  },
  async profile(): Promise<{ user: User }> {
    const { data } = await api.get('/api/users/profile')
    return data
  },
  async updateProfile(payload: { login?: string }): Promise<{ user: User }> {
    const { data } = await api.put('/api/users/profile', payload)
    return data
  },
}
