import { Container } from 'react-bootstrap'
import { Navigate, Route, Routes } from 'react-router-dom'
import { AppNavbar } from './components/AppNavbar'
import { AppBreadcrumbs } from './components/AppBreadcrumbs'
import { HomePage } from './pages/HomePage'
import { PartitionsListPage } from './pages/PartitionsListPage'
import { PartitionDetailPage } from './pages/PartitionDetailPage'
import { LoginPage } from './pages/LoginPage'
import { RegisterPage } from './pages/RegisterPage'
import { CalculationsListPage } from './pages/CalculationsListPage'
import { CalculationDetailPage } from './pages/CalculationDetailPage'
import { AccountPage } from './pages/AccountPage'
import { isGuestMode } from '../config/appConfig'

export default function App() {
  const guest = isGuestMode()
  return (
    <>
      <AppNavbar />
      <Container className="my-3">
        <AppBreadcrumbs />
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/partitions" element={<PartitionsListPage />} />
          <Route path="/partitions/:id" element={<PartitionDetailPage />} />
          {!guest && <Route path="/login" element={<LoginPage />} />}
          {!guest && <Route path="/register" element={<RegisterPage />} />}
          {!guest && <Route path="/calculations" element={<CalculationsListPage />} />}
          {!guest && <Route path="/calculations/:id" element={<CalculationDetailPage />} />}
          {!guest && <Route path="/account" element={<AccountPage />} />}
          <Route path="*" element={<Navigate to={guest ? '/partitions' : '/'} replace />} />
        </Routes>
      </Container>
    </>
  )
}
