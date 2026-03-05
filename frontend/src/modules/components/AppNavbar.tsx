import { Container, Nav, Navbar, Button } from 'react-bootstrap'
import { Link, NavLink, useNavigate } from 'react-router-dom'
import { useAppDispatch, useAppSelector } from '../../store'
import { reset as resetAuth } from '../../store/authSlice'
import * as authApi from '../services/authApi'
import { isGuestMode } from '../../config/appConfig'

export function AppNavbar() {
  const dispatch = useAppDispatch()
  const nav = useNavigate()
  const guest = isGuestMode()
  const user = useAppSelector(s => s.auth.user)
  let displayName = user?.login
  let isAuthed = !!user
  return (
    <Navbar bg="dark" data-bs-theme="dark" expand="md">
      <Container>
        <Navbar.Brand as={Link} to="/">Partition Soundproofing</Navbar.Brand>
        <Navbar.Toggle aria-controls="main-navbar" />
        <Navbar.Collapse id="main-navbar">
          <Nav className="me-auto">
            <Nav.Link as={NavLink} to="/">Главная</Nav.Link>
            <Nav.Link as={NavLink} to="/partitions">Перегородки</Nav.Link>
            {!guest && <Nav.Link as={NavLink} to="/calculations">Расчеты</Nav.Link>}
          </Nav>
          {!guest && (
            <Nav className="ms-auto align-items-center" variant="pills">
            {isAuthed ? (
              <>
                <Nav.Link as={NavLink} to="/account">👤 {displayName || 'Аккаунт'}</Nav.Link>
                <Button
                  size="sm"
                  variant="outline-light"
                  className="ms-2"
                  onClick={async () => {
                    try { await authApi.logout() } catch {}
                    localStorage.removeItem('app:token')
                    localStorage.removeItem('app:user')
                    dispatch(resetAuth())
                    nav('/login')
                  }}
                >Выход</Button>
              </>
            ) : (
              <>
                <Nav.Link as={NavLink} to="/login">Войти</Nav.Link>
                <Nav.Link as={NavLink} to="/register">Регистрация</Nav.Link>
              </>
            )}
            </Nav>
          )}
        </Navbar.Collapse>
      </Container>
    </Navbar>
  )
}
