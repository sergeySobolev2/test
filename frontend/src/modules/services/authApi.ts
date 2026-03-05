import { http, User } from './http'

export async function register(payload: { login: string; password: string }): Promise<User> {
  const res = await http.post('/api/users/register', payload)
  return (res.data?.data?.user || res.data?.user) as User
}

export async function login(payload: { login: string; password: string }): Promise<{ user: User; token: string; sessionId: string }> {
  const res = await http.post('/api/users/login', payload)
  const d = res.data?.data || res.data
  return {
    user: d.user as User,
    token: d.token as string,
    sessionId: d.session_id as string,
  }
}

export async function logout(): Promise<void> {
  await http.post('/api/users/logout')
}

export async function profile(): Promise<User> {
  const res = await http.get('/api/users/profile')
  return (res.data?.data?.user || res.data?.user) as User
}

export async function updateProfile(payload: { login?: string }): Promise<User> {
  const res = await http.put('/api/users/profile', payload)
  return (res.data?.data?.user || res.data?.user) as User
}
