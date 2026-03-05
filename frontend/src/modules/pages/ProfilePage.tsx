import { FormEvent, useEffect, useState } from 'react'
import { Alert, Button, Card, Form, Spinner } from 'react-bootstrap'
import { useAppDispatch, useAppSelector } from '../../store'
import { setUser, setError as setAuthError } from '../../store/authSlice'
import * as authApi from '../services/authApi'

export function ProfilePage() {
  const dispatch = useAppDispatch()
  const { user } = useAppSelector(s => s.auth)
  const [loginField, setLoginField] = useState(user?.login || '')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [ok, setOk] = useState<string | null>(null)

  useEffect(() => { setLoginField(user?.login || '') }, [user])

  const onSave = async (e: FormEvent) => {
    e.preventDefault()
    setLoading(true); setError(null); setOk(null)
    try {
      const user = await authApi.updateProfile({ login: loginField })
      dispatch(setUser(user))
      setOk('Профиль обновлён')
    } catch (e) {
      const msg = String(e)
      setError(msg)
      dispatch(setAuthError(msg))
    } finally { setLoading(false) }
  }

  return (
    <Card className="mx-auto" style={{ maxWidth: 520 }}>
      <Card.Body>
        <Card.Title>Личный кабинет</Card.Title>
        {error && <Alert variant="danger">{error}</Alert>}
        {ok && <Alert variant="success">{ok}</Alert>}
        <Form onSubmit={onSave}>
          <Form.Group className="mb-3">
            <Form.Label>Логин</Form.Label>
            <Form.Control value={loginField} onChange={e => setLoginField(e.target.value)} />
          </Form.Group>
          <div className="d-grid">
            <Button type="submit" disabled={loading}>{loading ? (<><Spinner size="sm"/> Сохранение…</>) : 'Сохранить'}</Button>
          </div>
        </Form>
        <hr/>
        <div className="text-muted small">Сброс пароля может быть реализован через отдельную форму (эндпоинт /api/users/password).</div>
      </Card.Body>
    </Card>
  )
}
