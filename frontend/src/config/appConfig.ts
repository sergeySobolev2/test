export type AppMode = 'full' | 'guest'

export function getAppMode(): AppMode {
  const raw = (import.meta as any).env?.VITE_APP_MODE
  return raw === 'guest' ? 'guest' : 'full'
}

export function isGuestMode(): boolean {
  return getAppMode() === 'guest' || isTauriRuntime() || isStandalonePwa()
}

export function getApiBase(): string {
  return String((import.meta as any).env?.VITE_API_BASE || '')
}

export function isTauriRuntime(): boolean {
  const w = window as any
  return Boolean(w.__TAURI__ || w.__TAURI_INTERNALS__)
}

export function isStandalonePwa(): boolean {
  if (typeof window === 'undefined') return false
  const w = window as any
  const isStandalone = w.matchMedia && w.matchMedia('(display-mode: standalone)').matches
  const isIosStandalone = typeof w.navigator !== 'undefined' && (w.navigator as any).standalone === true
  return Boolean(isStandalone || isIosStandalone)
}
