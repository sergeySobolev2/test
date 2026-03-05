import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Alert, Badge, Button, Col, Container, Form, Row, Spinner, Table } from 'react-bootstrap'
import { RequestsService, RequestPartitionsService } from '../../api/generated'
import { useAppDispatch, useAppSelector } from '../../store'
import { formRequest, fetchCart } from '../../store/requestsSlice'

type PartitionRow = {
  id: number
  title: string
  image_url?: string
  public_image_url?: string
  noise_reduction?: string
  thickness?: string
  material?: string
  price_per_sqm?: string
  comment?: string
}

const normalizeTitle = (title: string, thickness?: string) => {
  if (!thickness) return title
  const cleaned = title.replace(thickness, '').replace(/\bмм\b/gi, '').replace(/\s{2,}/g, ' ').trim()
  return cleaned || title
}

export function CalculationDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const user = useAppSelector(s => s.auth.user)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [status, setStatus] = useState<string>('черновик')
  const [createdAt, setCreatedAt] = useState<string>('')
  const [ownerId, setOwnerId] = useState<number | null>(null)
  const [roomArea, setroomArea] = useState<string>('')
  const [expertComment, setexpertComment] = useState<string>('')
  const [noiseReductionDB, setNoiseReductionDB] = useState<number | null>(null)
  const [requiredThickness, setRequiredThickness] = useState<number | null>(null)
  const [rows, setRows] = useState<PartitionRow[]>([])

  const isDraft = status === 'черновик'
  const canEdit = isDraft && !!user && (user.is_moderator || ownerId === user.id)

  const applyRequestData = (data: any) => {
    const req = data.calculation || data.request || data
    const syms = data.partitions || []

    setStatus(req.status || 'черновик')
    setOwnerId(typeof req.user_id === 'number' ? req.user_id : (typeof req.user?.id === 'number' ? req.user.id : null))
    setCreatedAt(req.created_at || '')
    setroomArea(req.room_area ? String(req.room_area) : '')
    setexpertComment(typeof req.expert_comment === 'string' ? req.expert_comment : (req.expert_comment?.string || ''))
    setNoiseReductionDB(typeof req.noise_reduction_db === 'number' ? req.noise_reduction_db : null)
    setRequiredThickness(typeof req.required_thickness === 'number' ? req.required_thickness : null)

    const items: PartitionRow[] = syms.map((x: any) => ({
      id: x.id,
      title: x.title,
      image_url: x.image_url,
      public_image_url: x.public_image_url,
      noise_reduction: x.noise_reduction,
      thickness: x.thickness,
      material: x.material,
      price_per_sqm: x.price_per_sqm,
      comment: x.comment ?? '',
    }))
    setRows(items)
  }

  const reloadRequest = async (silent = false) => {
    const rid = Number(id)
    if (!rid) return
    if (!silent) {
      setLoading(true)
      setError(null)
    }
    try {
      const r = await RequestsService.get(rid)
      const data = r.data?.data || r.data
      applyRequestData(data)
    } catch (e: any) {
      if (!silent) setError(String(e))
    } finally {
      if (!silent) setLoading(false)
    }
  }

  useEffect(() => {
    const rid = Number(id)
    if (!rid) return
    let cancel = false
    setLoading(true)
    setError(null)
    RequestsService.get(rid)
      .then(r => {
        if (cancel) return
        const data = r.data?.data || r.data
        applyRequestData(data)
      })
      .catch(e => { if (!cancel) setError(String(e)) })
      .finally(() => { if (!cancel) setLoading(false) })

    const intervalId = setInterval(() => {
      if (cancel) return
      reloadRequest(true)
    }, 2000)

    return () => {
      cancel = true
      clearInterval(intervalId)
    }
  }, [id])

  const onRemove = async (sid: number) => {
    if (!id) return
    setLoading(true)
    try {
      await RequestPartitionsService.remove(Number(id), sid)
      setRows(prev => prev.filter(r => r.id !== sid))
      dispatch(fetchCart())
    } catch (e: any) {
      setError(String(e))
    } finally {
      setLoading(false)
    }
  }

  const onUpdate = (row: PartitionRow) => {
    setRows(prev => prev.map(r => r.id === row.id ? { ...r, comment: row.comment } : r))
  }

  const onSaveComments = async () => {
    if (!id || rows.length === 0) return
    setLoading(true)
    try {
      await Promise.all(rows.map(r => RequestPartitionsService.update(Number(id), r.id, { comment: r.comment || '' })))
    } catch (e: any) {
      setError(String(e))
    } finally {
      setLoading(false)
    }
  }

  const onSaveWeight = async () => {
    if (!id) return
    const hasWeight = roomArea.trim() !== ''
    const hasComment = expertComment.trim() !== ''
    if (!hasWeight && !hasComment) return
    setLoading(true)
    try {
      const payload: any = {}
      if (hasWeight) {
        payload.room_area = parseFloat(roomArea)
      }
      if (hasComment) {
        payload.expert_comment = expertComment
      }
      await RequestsService.update(Number(id), payload)
      await reloadRequest(true)
    } catch (e: any) {
      setError(String(e))
    } finally {
      setLoading(false)
    }
  }

  const onForm = async () => {
    if (!id) return
    setLoading(true)
    try {
      const res = await dispatch(formRequest(Number(id)))
      if ((res as any).meta?.requestStatus === 'fulfilled') {
        setStatus('сформирован')
        dispatch(fetchCart())
        alert('Заявка сформирована!')
        navigate('/calculations')
      }
    } catch (e: any) {
      setError(String(e))
    } finally {
      setLoading(false)
    }
  }

  return (
    <Container className="mt-4">
      <div className="d-flex justify-content-between align-items-center mb-3">
        <h2>
          Заявка #{id} 
          {status && (
            <Badge 
              bg={isDraft ? 'secondary' : status === 'сформирован' ? 'info' : status === 'завершен' ? 'success' : 'danger'} 
              className="ms-2"
            >
              {status}
            </Badge>
          )}
        </h2>
      </div>

      {error && <Alert variant="danger" dismissible onClose={() => setError(null)}>{error}</Alert>}
      
      <div className="text-muted mb-3">
        Создана: {createdAt ? new Date(createdAt).toLocaleString('ru-RU') : '—'}
      </div>

      {/* Поля веса пациента и комментария */}
      <Form className="mb-4">
        <Row className="g-3">
          <Col md={4}>
            <Form.Label>Площадь помещения (м²):</Form.Label>
            <Form.Control
              type="number"
              step="0.1"
              placeholder="м²"
              value={roomArea}
              disabled={!canEdit || loading}
              onChange={e => setroomArea(e.target.value)}
              onBlur={onSaveWeight}
            />
          </Col>
          {noiseReductionDB !== null && (
            <Col md={4}>
              <Form.Label>Требуемое снижение шума (дБ):</Form.Label>
              <div className="form-control-plaintext fw-bold">
                {noiseReductionDB.toFixed(1)}
              </div>
            </Col>
          )}
          {requiredThickness !== null && (
            <Col md={4}>
              <Form.Label>Рекомендуемая толщина (см):</Form.Label>
              <div className="form-control-plaintext fw-bold">
                {requiredThickness.toFixed(1)}
              </div>
            </Col>
          )}
        </Row>
        <Row className="mt-3">
          <Col>
            <Form.Label>Комментарий эксперта:</Form.Label>
            <Form.Control
              as="textarea"
              rows={3}
              placeholder="Введите комментарий к заявке..."
              value={expertComment}
              disabled={!canEdit || loading}
              onChange={e => setexpertComment(e.target.value)}
              onBlur={onSaveWeight}
            />
          </Col>
        </Row>
      </Form>

      {/* Таблица симптомов */}
      <div className="mb-4">
        <h5 className="mb-3">Перегородки в заявке</h5>
        
        {rows.length === 0 ? (
          <Alert variant="info">
            В заявке пока нет перегородок. Добавьте их на странице перегородок.
          </Alert>
        ) : (
          <Table bordered hover responsive className="align-middle">
            <thead className="table-light">
              <tr>
                <th style={{ width: '250px' }}>Перегородка</th>
                <th style={{ width: '120px' }}>Ширина</th>
                <th style={{ width: '220px' }}>Комментарий</th>
                {isDraft && <th style={{ width: '100px' }}>Действия</th>}
              </tr>
            </thead>
            <tbody>
              {rows.map((r) => (
                <tr key={r.id}>
                  <td>
                    <div className="d-flex align-items-center gap-2">
                      <img 
                        src={r.public_image_url || r.image_url || '/placeholder-partition.svg'} 
                        alt={r.title}
                        style={{ width: 100, height: 100, objectFit: 'cover', borderRadius: 4 }}
                      />
                      <div>
                        <strong>{normalizeTitle(r.title, r.thickness)}</strong>
                      </div>
                    </div>
                  </td>
                  <td className="text-center">
                    {r.thickness || '—'}
                  </td>
                  <td>
                    {canEdit ? (
                      <Form.Control
                        value={r.comment || ''}
                        onChange={e => onUpdate({ ...r, comment: e.target.value })}
                      />
                    ) : (r.comment || '—')}
                  </td>
                  {canEdit && (
                    <td className="text-center">
                      <Button 
                        variant="outline-danger" 
                        size="sm" 
                        disabled={loading} 
                        onClick={() => onRemove(r.id)}
                      >
                        Удалить
                      </Button>
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
          </Table>
        )}
      </div>

      {/* Кнопки действий */}
      <div className="d-flex gap-2 flex-wrap">
        {/* 1. Сохранить поля */}
        {canEdit && (
          <Button 
            variant="info"
            onClick={onSaveWeight}
            disabled={loading}
          >
            Сохранить черновик
          </Button>
        )}
        {canEdit && (
          <Button
            variant="outline-primary"
            onClick={onSaveComments}
            disabled={loading || rows.length === 0}
          >
            Сохранить комментарии
          </Button>
        )}
        {/* 2. Удалить заявку */}
        {canEdit && (
          <Button 
            variant="danger"
            onClick={async () => {
              if (!id) return
              if (!window.confirm('Удалить заявку?')) return
              setLoading(true)
              try {
                await RequestsService.remove(Number(id))
                dispatch(fetchCart())
                navigate('/calculations')
              } catch (e: any) {
                setError(String(e))
              } finally {
                setLoading(false)
              }
            }}
            disabled={loading}
          >
            Удалить заявку
          </Button>
        )}
        {/* 3. Оформить заявку */}
        {isDraft && (
          <Button 
            variant="primary" 
            onClick={onForm} 
            disabled={loading || rows.length === 0}
          >
            {loading ? (
              <>
                <Spinner size="sm" className="me-2" />
                Оформление…
              </>
            ) : (
              'Оформить заявку'
            )}
          </Button>
        )}
      </div>
    </Container>
  )
}
