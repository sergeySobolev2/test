import { useEffect, useState } from 'react'
import { Alert, Badge, Button, Col, Form, Row, Spinner, Table } from 'react-bootstrap'
import { useAppDispatch, useAppSelector } from '../../store'
import { fetchCart, listRequests } from '../../store/requestsSlice'
import { Link, useNavigate } from 'react-router-dom'
import { RequestsService } from '../../api/generated'

export function CalculationsListPage() {
  const dispatch = useAppDispatch()
  const { list, loading, error, draftId, draftCount } = useAppSelector(s => s.requests)
  const { user } = useAppSelector(s => s.auth)
  const nav = useNavigate()
  const [status, setStatus] = useState<string>('')
  const [from, setFrom] = useState<string>('')
  const [to, setTo] = useState<string>('')
  const [actionLoading, setActionLoading] = useState<number | null>(null)

  const isModerator = user?.is_moderator || false

  //Short Polling - обновление каждые 2 секунды
  useEffect(() => {
    const hasToken = !!user || !!localStorage.getItem('app:token')
    if (!hasToken) { nav('/login'); return }
    
    // Первоначальная загрузка
    dispatch(listRequests({ status: status || undefined, from: from || undefined, to: to || undefined }))
    dispatch(fetchCart())
    
    // Short polling - обновляем список каждые 2 секунды
    const intervalId = setInterval(() => {
      dispatch(listRequests({ status: status || undefined, from: from || undefined, to: to || undefined }))
    }, 2000)
    
    // Очистка при размонтировании компонента
    return () => clearInterval(intervalId)
  }, [dispatch, user, status, from, to, nav])

  // Модератор - смена статуса заявки
  const handleCompleteRequest = async (requestId: number, newStatus: 'завершен' | 'отклонен') => {
    if (!confirm(`Вы уверены, что хотите ${newStatus === 'завершен' ? 'завершить' : 'отклонить'} заявку #${requestId}?`)) {
      return
    }

    setActionLoading(requestId)
    try {
      if (newStatus === 'завершен') {
        const area = prompt('Введите площадь помещения (м²):')
        const noise = prompt('Введите требуемое снижение шума (дБ):')
        if (!area || !noise) {
          alert('Необходимо заполнить все поля')
          setActionLoading(null)
          return
        }
        await RequestsService.complete(requestId, {
          status: newStatus,
          room_area: parseFloat(area),
          noise_reduction_db: parseFloat(noise),
        })
      } else {
        const comment = prompt('Комментарий (необязательно):') || undefined
        await RequestsService.complete(requestId, {
          status: newStatus,
          expert_comment: comment,
        })
      }
      dispatch(listRequests({ status: status || undefined, from: from || undefined, to: to || undefined }))
    } catch (err: any) {
      alert(`Ошибка: ${err.message || 'Неизвестная ошибка'}`)
    } finally {
      setActionLoading(null)
    }
  }

  // Модератор - запустить асинхронный расчет
  const handleTriggerCalculation = async (requestId: number) => {
    const roomArea = prompt('Введите площадь помещения (м²):')
    if (!roomArea || isNaN(parseFloat(roomArea)) || parseFloat(roomArea) <= 0) {
      alert('Необходимо ввести корректную площадь помещения')
      return
    }

    const thickness = prompt('Введите толщину перегородки (см):')
    if (!thickness || isNaN(parseFloat(thickness)) || parseFloat(thickness) <= 0) {
      alert('Необходимо ввести корректную толщину перегородки')
      return
    }

    if (!confirm(`Запустить расчет для заявки #${requestId} с площадью ${roomArea} м² и толщиной ${thickness} см?`)) {
      return
    }

    setActionLoading(requestId)
    try {
      const response = await fetch(`/api/requests/${requestId}/calculate`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('app:token')}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          room_area: parseFloat(roomArea),
          thickness_cm: parseFloat(thickness)
        })
      })

      if (!response.ok) {
        throw new Error('Ошибка запуска расчета')
      }

      alert('Расчет запущен. Результат появится через ~7 секунд.')
      
      // Обновляем список
      dispatch(listRequests({ status: status || undefined, from: from || undefined, to: to || undefined }))
    } catch (err: any) {
      alert(`Ошибка: ${err.message || 'Неизвестная ошибка'}`)
    } finally {
      setActionLoading(null)
    }
  }

  return (
    <div>
      <h2>{isModerator ? 'Все заявки (Модератор)' : 'Мои заявки'}</h2>
      {draftId && !list.some((r: any) => r.id === draftId) && (
        <Alert variant="info" className="d-flex justify-content-between align-items-center">
          <div>
            Есть черновик заявки #{draftId}
            {typeof draftCount === 'number' ? ` (перегородок: ${draftCount})` : ''}
          </div>
          <Button as={Link as any} to={`/calculations/${draftId}`} size="sm" variant="outline-primary">
            Продолжить
          </Button>
        </Alert>
      )}
      <Row className="g-2 mb-3">
        <Col md={4} lg={3}>
          <Form.Label>Статус</Form.Label>
          <Form.Select value={status} onChange={e => setStatus(e.target.value)}>
            <option value="">Все статусы</option>
            <option value="сформирован">Сформирован</option>
            <option value="завершен">Завершен</option>
            <option value="отклонен">Отклонен</option>
          </Form.Select>
        </Col>
        <Col md={4} lg={3}>
          <Form.Label>Дата от</Form.Label>
          <Form.Control type="date" value={from} onChange={e => setFrom(e.target.value)} placeholder="дд.мм.гггг" />
        </Col>
        <Col md={4} lg={3}>
          <Form.Label>Дата до</Form.Label>
          <Form.Control type="date" value={to} onChange={e => setTo(e.target.value)} placeholder="дд.мм.гггг" />
        </Col>
      </Row>
      {false && loading && <div className="my-2"><Spinner size="sm" /> Загрузка…</div>}
      {error && <Alert variant="danger">{error}</Alert>}
      <Table striped hover responsive size="sm">
        <thead>
          <tr>
            <th>ID</th>
            {isModerator && <th>Пользователь</th>}
            <th>Статус</th>
            <th>Дата создания</th>
            <th>Дата отправки</th>
            <th>Площадь (м²)</th>
            <th>Снижение шума (дБ)</th>
            <th>Толщина (см)</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {list.filter((r: any) => r.status !== 'черновик').map((r: any) => (
            <tr key={r.id}>
              <td>{r.id}</td>
              {isModerator && <td>{r.user?.login || `User #${r.user_id}`}</td>}
              <td>
                <Badge bg={r.status === 'черновик' ? 'secondary' : r.status === 'сформирован' ? 'info' : r.status === 'завершен' ? 'success' : 'danger'}>
                  {r.status}
                </Badge>
              </td>
              <td>{new Date(r.created_at).toLocaleString('ru-RU', { 
                year: 'numeric', 
                month: '2-digit', 
                day: '2-digit', 
                hour: '2-digit', 
                minute: '2-digit' 
              })}</td>
              <td>{r.formed_at ? new Date(r.formed_at).toLocaleString('ru-RU', { 
                year: 'numeric', 
                month: '2-digit', 
                day: '2-digit', 
                hour: '2-digit', 
                minute: '2-digit' 
              }) : '—'}</td>
              <td>{typeof r.room_area === 'number' ? r.room_area : '—'}</td>
              <td>{typeof r.noise_reduction_db === 'number' ? r.noise_reduction_db : '—'}</td>
              <td>{typeof r.required_thickness === 'number' ? r.required_thickness : '—'}</td>
              <td className="text-end">
                <Button 
                  as={Link as any} 
                  to={`/calculations/${r.id}`} 
                  size="sm" 
                  variant="outline-primary"
                  className="me-1"
                >
                  Открыть
                </Button>
                {/* Кнопки модератора */}
                {isModerator && (r.status === 'сформирован' || r.status === 'черновик') && (
                  <>
                    <Button 
                      size="sm" 
                      variant="outline-success"
                      className="me-1"
                      onClick={() => handleCompleteRequest(r.id, 'завершен')}
                      disabled={actionLoading === r.id}
                      title="Принять"
                    >
                      {actionLoading === r.id ? <Spinner size="sm" /> : 'Принять'}
                    </Button>
                    <Button 
                      size="sm" 
                      variant="outline-danger"
                      className="me-1"
                      onClick={() => handleCompleteRequest(r.id, 'отклонен')}
                      disabled={actionLoading === r.id}
                      title="Отклонить"
                    >
                      {actionLoading === r.id ? <Spinner size="sm" /> : 'Отклонить'}
                    </Button>
                    <Button 
                      size="sm" 
                      variant="outline-info"
                      onClick={() => handleTriggerCalculation(r.id)}
                      disabled={actionLoading === r.id}
                      title="Запустить расчет снижения шума"
                    >
                      {actionLoading === r.id ? <Spinner size="sm" /> : '🧮'}
                    </Button>
                  </>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
      {isModerator && (
        <Alert variant="info" className="mt-3">
          <strong>Режим модератора:</strong> Список обновляется автоматически каждые 2 секунды (short polling).
          Вы можете завершать/отклонять заявки и запускать асинхронный расчет снижения шума (🧮) с указанием площади помещения.
        </Alert>
      )}
    </div>
  )
}
