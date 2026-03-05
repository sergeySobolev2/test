import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import type { Partition } from '../modules/services/partitionsApi'

export type PartitionsState = {
  list: Partition[]
  current: Partition | null
  total: number
  loading: boolean
  error: string | null
}

const initialState: PartitionsState = {
  list: [],
  current: null,
  total: 0,
  loading: false,
  error: null,
}

const slice = createSlice({
  name: 'partitions',
  initialState,
  reducers: {
    start(state) {
      state.loading = true
      state.error = null
    },
    setList(state, action: PayloadAction<{ items: Partition[]; total?: number }>) {
      state.loading = false
      state.list = action.payload.items
      state.total = action.payload.total ?? action.payload.items.length
    },
    setCurrent(state, action: PayloadAction<Partition | null>) {
      state.loading = false
      state.current = action.payload
    },
    setError(state, action: PayloadAction<string | null>) {
      state.loading = false
      state.error = action.payload
    },
    clearError(state) {
      state.error = null
    },
  },
})

export const { start, setList, setCurrent, setError, clearError } = slice.actions
export default slice.reducer
