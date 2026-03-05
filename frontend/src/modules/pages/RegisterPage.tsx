import { useState } from 'react'
import { Alert, Button, Card, Form, Spinner } from 'react-bootstrap'
import { useNavigate } from 'react-router-dom'
import { useAppDispatch, useAppSelector } from '../../store'
import { clearError, start, setError } from '../../store/authSlice'
import * as authApi from '../services/authApi'

export function RegisterPage() {
  const dispatch = useAppDispatch()
  const nav = useNavigate()
  const { loading, error } = useAppSelector(s => s.auth)
  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')

  const handleRegister = async () => {
    dispatch(start())
    try {
      await authApi.register({ login, password })
      nav('/login')
    } catch (e: any) {
      dispatch(setError(e?.response?.data?.error || String(e)))
    }
  }

  return (
    <Card className="mx-auto" style={{ maxWidth: 480 }}>
      <Card.Body>
        <Card.Title>Регистрация</Card.Title>
        {error && <Alert variant="danger" onClose={() => dispatch(clearError())} dismissible>{error}</Alert>}
        <Form onSubmit={(e) => { e.preventDefault(); handleRegister() }}>
          <Form.Group className="mb-3">
            <Form.Label>Логин</Form.Label>
            <Form.Control value={login} onChange={e => setLogin(e.target.value)} required />
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Пароль</Form.Label>
            <Form.Control type="password" value={password} onChange={e => setPassword(e.target.value)} required />
          </Form.Group>
          <Button type="button" disabled={loading} onClick={handleRegister}>
            {loading ? (<><Spinner size="sm" /> Регистрируем…</>) : 'Зарегистрироваться'}
          </Button>
        </Form>
      </Card.Body>
    </Card>
  )
}
