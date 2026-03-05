import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { User } from '../modules/services/http'

export type AuthState = {
  user: User | null
  token: string | null
  sessionId: string | null
  loading: boolean
  error: string | null
}

export const AUTH_TOKEN_KEY = 'app:token'
export const AUTH_USER_KEY = 'app:user'

const initialState: AuthState = {
  user: null,
  token: null,
  sessionId: null,
  loading: false,
  error: null,
}

const slice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    hydrate(state, action) {
      Object.assign(state, action.payload || {})
    },
    start(state) {
      state.loading = true
      state.error = null
    },
    setAuth(state, action: PayloadAction<{ user: User | null; token?: string | null; sessionId?: string | null }>) {
      state.loading = false
      state.error = null
      state.user = action.payload.user
      state.token = action.payload.token ?? null
      state.sessionId = action.payload.sessionId ?? null
    },
    setUser(state, action: PayloadAction<User | null>) {
      state.loading = false
      state.error = null
      state.user = action.payload
    },
    setError(state, action: PayloadAction<string | null>) {
      state.loading = false
      state.error = action.payload
    },
    clearError(state) { state.error = null },
    reset(state) {
      state.user = null
      state.token = null
      state.sessionId = null
      state.loading = false
      state.error = null
    },
  },
})

export const { clearError, reset, start, setAuth, setUser, setError, hydrate } = slice.actions
export default slice.reducer
