import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export type FiltersState = {
  title: string
  active: 'all' | 'true' | 'false'
  appliedTitle: string
  appliedActive: 'all' | 'true' | 'false'
}

const initialState: FiltersState = {
  title: '',
  active: 'all',
  appliedTitle: '',
  appliedActive: 'all',
}

const filtersSlice = createSlice({
  name: 'filters',
  initialState,
  reducers: {
    setTitle(state, action: PayloadAction<string>) {
      state.title = action.payload
    },
    setActive(state, action: PayloadAction<'all' | 'true' | 'false'>) {
      state.active = action.payload
    },
    resetTitle(state) {
      state.title = ''
    },
    apply(state) {
      state.appliedTitle = state.title
      state.appliedActive = state.active
    },
    hydrate(state, action: PayloadAction<Partial<FiltersState>>) {
      return { ...state, ...action.payload }
    }
  }
})

export const { setTitle, setActive, resetTitle, apply, hydrate } = filtersSlice.actions
export default filtersSlice.reducer
