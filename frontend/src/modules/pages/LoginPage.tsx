import { useState } from 'react'
import { Alert, Button, Card, Form, Spinner } from 'react-bootstrap'
import { useNavigate } from 'react-router-dom'
import { useAppDispatch, useAppSelector } from '../../store'
import { clearError, start, setAuth, setError, setUser } from '../../store/authSlice'
import { fetchCart } from '../../store/requestsSlice'
import * as authApi from '../services/authApi'

export function LoginPage() {
  const dispatch = useAppDispatch()
  const nav = useNavigate()
  const { loading, error, user } = useAppSelector(s => s.auth)
  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')

  const handleLogin = async () => {
    dispatch(start())
    try {
      const res = await authApi.login({ login, password })
      localStorage.setItem('app:token', res.token || '')
      localStorage.setItem('app:user', JSON.stringify(res.user))
      dispatch(setAuth({ user: res.user, token: res.token, sessionId: res.sessionId }))
      const profile = await authApi.profile().catch(() => res.user)
      dispatch(setUser(profile))
      await dispatch(fetchCart())
      nav('/')
    } catch (e: any) {
      dispatch(setError(e?.response?.data?.error || String(e)))
    }
  }

  return (
    <Card className="mx-auto" style={{ maxWidth: 480 }}>
      <Card.Body>
        <Card.Title>Вход</Card.Title>
        {error && <Alert variant="danger" onClose={() => dispatch(clearError())} dismissible>{error}</Alert>}
        {user && <Alert variant="success">Вы вошли как {user.login}</Alert>}
        <Form onSubmit={(e) => { e.preventDefault(); handleLogin() }}>
          <Form.Group className="mb-3">
            <Form.Label>Логин</Form.Label>
            <Form.Control value={login} onChange={e => setLogin(e.target.value)} required />
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Пароль</Form.Label>
            <Form.Control type="password" value={password} onChange={e => setPassword(e.target.value)} required />
          </Form.Group>
          <Button type="button" disabled={loading} onClick={handleLogin}>
            {loading ? (<><Spinner size="sm" /> Входим…</>) : 'Войти'}
          </Button>
        </Form>
      </Card.Body>
    </Card>
  )
}
