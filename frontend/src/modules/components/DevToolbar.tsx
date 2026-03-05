import { useEffect, useMemo } from 'react'
import { Button, ButtonGroup } from 'react-bootstrap'
import { useAppDispatch, useAppSelector } from '../../store'
import { addPartitionToDraft, fetchCart, formRequest } from '../../store/requestsSlice'
import { setAuth, setError, start } from '../../store/authSlice'
import * as authApi from '../services/authApi'
import { getApiBase } from '../../config/appConfig'

function useDevEnabled() {
  const enabled = useMemo(() => {
    if (!import.meta.env.DEV) return false
    const params = new URLSearchParams(window.location.search)
    if (params.get('dev') === '1') return true
    try { if (localStorage.getItem('dev') === '1') return true } catch {}
    return false
  }, [])
  useEffect(() => {
    // one-time marker for demo convenience
    if (enabled) try { localStorage.setItem('dev', '1') } catch {}
  }, [enabled])
  return enabled
}

export function DevToolbar() {
  const enabled = useDevEnabled()
  const dispatch = useAppDispatch()
  const { draftId, loading } = useAppSelector(s => s.requests)
  const apiBase = getApiBase() || '(proxy /api)'
  const host = typeof window !== 'undefined' ? window.location.host : ''
  if (!enabled) return null
  return (
    <div style={{ position: 'fixed', right: 16, bottom: 16, zIndex: 1030 }}>
      <div className="bg-dark text-white rounded px-3 py-2 shadow">
        <div className="small mb-2">Dev панель</div>
        <div className="small text-muted" style={{ maxWidth: 240 }}>
          API: {apiBase}
          <br />
          Host: {host}
        </div>
        <ButtonGroup size="sm">
          <Button
            variant="outline-info"
            onClick={async () => {
              dispatch(start())
              try {
                const res = await authApi.login({ login: 'demo', password: 'demo123' })
                localStorage.setItem('app:token', res.token || '')
                localStorage.setItem('app:user', JSON.stringify(res.user))
                dispatch(setAuth({ user: res.user, token: res.token, sessionId: res.sessionId }))
              } catch (e: any) {
                dispatch(setError(e?.response?.data?.error || String(e)))
              }
            }}
          >Login demo</Button>
          <Button variant="outline-light" disabled={loading} onClick={async () => { await dispatch(addPartitionToDraft(1)); dispatch(fetchCart()) }}>+ перегородка #1</Button>
          <Button variant="outline-success" disabled={!draftId || loading} onClick={async () => { await dispatch(formRequest(draftId!)) }}>Сформировать</Button>
        </ButtonGroup>
      </div>
    </div>
  )
}
