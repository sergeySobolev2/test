import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit'
import { RequestsService, RequestPartitionsService } from '../api/generated'

export type RequestsState = {
  draftId: number | null
  draftCount: number
  loading: boolean
  error: string | null
  list: Array<any>
}

const initialState: RequestsState = {
  draftId: null,
  draftCount: 0,
  loading: false,
  error: null,
  list: [],
}

export const fetchCart = createAsyncThunk('requests/cart', async () => {
  return await RequestsService.cart()
})

export const addPartitionToDraft = createAsyncThunk('requests/addPartition', async (partitionId: number, { getState }) => {
  const state = getState() as any
  const calculationId = state.requests?.draftId as number | undefined
  const res = await RequestPartitionsService.add(partitionId, calculationId)
  const cart = await RequestsService.cart()
  return { res, cart }
})

export const removePartitionFromDraft = createAsyncThunk('requests/removePartition', async (payload: { calculationId: number, partitionId: number }) => {
  const res = await RequestPartitionsService.remove(payload.calculationId, payload.partitionId)
  const cart = await RequestsService.cart()
  return { res, cart }
})

export const listRequests = createAsyncThunk('requests/list', async (filters: { status?: string; from?: string; to?: string } | undefined) => {
  const res = await RequestsService.list(filters || {})
  return res
})

export const getRequest = createAsyncThunk('requests/get', async (id: number) => {
  const res = await RequestsService.get(id)
  return res
})

export const formRequest = createAsyncThunk('requests/form', async (id: number) => {
  const res = await RequestsService.form(id)
  return res
})

const slice = createSlice({
  name: 'requests',
  initialState,
  reducers: {
    setLoading(state, action: PayloadAction<boolean>) {
      state.loading = action.payload
    },
    setRequestsError(state, action: PayloadAction<string | null>) {
      state.error = action.payload
    },
    setCartState(state, action: PayloadAction<{ draftId: number | null; draftCount: number }>) {
      state.draftId = action.payload.draftId
      state.draftCount = action.payload.draftCount
    },
    resetState() { return initialState },
  },
  extraReducers: builder => {
    builder
      .addCase(fetchCart.pending, (s) => { s.loading = true; s.error = null })
      .addCase(fetchCart.fulfilled, (s, a) => {
        s.loading = false
        const d = (a.payload as any)?.data || a.payload
        s.draftId = (d?.draft_id as number | null) ?? null
        s.draftCount = (d?.count as number) ?? 0
      })
      .addCase(fetchCart.rejected, (s, a) => { s.loading = false; s.error = a.error.message || 'Ошибка корзины' })

      .addCase(addPartitionToDraft.pending, (s) => { s.loading = true; s.error = null })
      .addCase(addPartitionToDraft.fulfilled, (s, a) => {
        s.loading = false
        const d = (a.payload as any)?.cart?.data || (a.payload as any)?.cart
        s.draftId = (d?.draft_id as number | null) ?? s.draftId
        s.draftCount = (d?.count as number) ?? s.draftCount
      })
      .addCase(addPartitionToDraft.rejected, (s, a) => { s.loading = false; s.error = a.error.message || 'Не удалось добавить' })

      .addCase(listRequests.pending, (s) => { s.loading = true; s.error = null })
      .addCase(listRequests.fulfilled, (s, a) => { s.loading = false; s.list = a.payload.data as any[] })
      .addCase(listRequests.rejected, (s, a) => { s.loading = false; s.error = a.error.message || 'Не удалось загрузить заявки' })

      .addCase(getRequest.pending, (s) => { s.loading = true; s.error = null })
      .addCase(getRequest.fulfilled, (s) => { s.loading = false })
      .addCase(getRequest.rejected, (s, a) => { s.loading = false; s.error = a.error.message || 'Не удалось загрузить заявку' })

      .addCase(formRequest.pending, (s) => { s.loading = true; s.error = null })
      .addCase(formRequest.fulfilled, (s) => { s.loading = false; s.draftId = null; s.draftCount = 0 })
      .addCase(removePartitionFromDraft.pending, (s) => { s.loading = true; s.error = null })
      .addCase(removePartitionFromDraft.fulfilled, (s, a) => {
        s.loading = false
        const d = (a.payload as any)?.cart?.data || (a.payload as any)?.cart
        s.draftId = (d?.draft_id as number | null) ?? s.draftId
        s.draftCount = (d?.count as number) ?? s.draftCount
      })
      .addCase(formRequest.rejected, (s, a) => { s.loading = false; s.error = a.error.message || 'Не удалось сформировать' })
  }
})

export const { resetState: resetRequests, setLoading, setRequestsError, setCartState } = slice.actions
export default slice.reducer
