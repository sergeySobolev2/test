import { useEffect, useState } from 'react'
import { Alert, Button, Card, Form, Spinner } from 'react-bootstrap'
import { useAppDispatch, useAppSelector } from '../../store'
import { start, setUser, setError, reset as resetAuth } from '../../store/authSlice'
import * as authApi from '../services/authApi'

export function AccountPage() {
  const dispatch = useAppDispatch()
  const { user, loading, error } = useAppSelector(s => s.auth)
  const [login, setLogin] = useState(user?.login || '')

  useEffect(() => {
    dispatch(start())
    authApi.profile()
      .then(u => dispatch(setUser(u)))
      .catch(e => dispatch(setError(String(e))))
  }, [dispatch])
  useEffect(() => { setLogin(user?.login || '') }, [user])

  const onSave = async (e: React.FormEvent) => {
    e.preventDefault()
    dispatch(start())
    try {
      const u = await authApi.updateProfile({ login })
      dispatch(setUser(u))
    } catch (e: any) {
      dispatch(setError(e?.response?.data?.error || String(e)))
    }
  }

  const onLogout = async () => {
    try { await authApi.logout() } catch {}
    localStorage.removeItem('app:token')
    localStorage.removeItem('app:user')
    dispatch(resetAuth())
  }

  return (
    <Card className="mx-auto" style={{ maxWidth: 560 }}>
      <Card.Body>
        <Card.Title>Личный кабинет</Card.Title>
        {error && <Alert variant="danger">{error}</Alert>}
        {!user && <Alert variant="warning">Вы не авторизованы</Alert>}
        {user && (
          <Form onSubmit={onSave}>
            <Form.Group className="mb-3">
              <Form.Label>Логин</Form.Label>
              <Form.Control value={login} onChange={e => setLogin(e.target.value)} />
            </Form.Group>
            <div className="d-flex gap-2">
              <Button type="submit" disabled={loading}>{loading ? (<><Spinner size="sm" /> Сохраняю…</>) : 'Сохранить'}</Button>
              <Button variant="outline-danger" onClick={onLogout}>Выход</Button>
            </div>
          </Form>
        )}
      </Card.Body>
    </Card>
  )
}
