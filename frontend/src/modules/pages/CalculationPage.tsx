import { useEffect, useState } from 'react'
import { Alert, Badge, Button, Card, Col, Form, Row, Spinner, Table } from 'react-bootstrap'
import { useNavigate, useParams } from 'react-router-dom'
import { RequestsService, RequestPartitionsService } from '../../api/generated'
import { useAppSelector } from '../../store'

type Item = {
  partition: { id: number; title: string; thickness?: string }
  comment?: string
}

const normalizeTitle = (title: string, thickness?: string) => {
  if (!thickness) return title
  const cleaned = title.replace(thickness, '').replace(/мм/gi, '').replace(/\s{2,}/g, ' ').trim()
  return cleaned || title
}

export function CalculationPage() {
  const { id } = useParams()
  const nav = useNavigate()
  const { user } = useAppSelector(s => s.auth)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [status, setStatus] = useState<string>('')
  const [items, setItems] = useState<Item[]>([])
  const [roomArea, setroomArea] = useState<string>('')
  const [expertComment, setexpertComment] = useState<string>('')

  // Восстановление черновика из localStorage
  useEffect(() => {
    if (!id) return;
    const draft = localStorage.getItem(`draft_request_${id}`);
    if (draft) {
      try {
        const data = JSON.parse(draft);
        setStatus(data.status || 'черновик');
        setroomArea(data.roomArea || '');
        setexpertComment(data.expertComment || '');
        setItems(data.items || []);
      } catch {}
    }
  }, [id]);

  // Сохранять черновик в localStorage при изменениях
  useEffect(() => {
    if (!id) return;
    localStorage.setItem(`draft_request_${id}`,
      JSON.stringify({ status, roomArea, expertComment, items })
    );
  }, [id, status, roomArea, expertComment, items]);

  useEffect(() => { if (!user) nav('/login') }, [user])
  useEffect(() => {
    if (!id) return
    setLoading(true); setError(null)
    RequestsService.get(Number(id))
      .then(res => {
        const data: any = res.data
        const rq = data?.calculation || data?.request || data?.Request || data
        setStatus(rq.status)
        setroomArea(rq.room_area ?? '')
        setexpertComment(rq.expert_comment?.string ?? '')
        const rs: any[] = (data?.partitions || data?.items || []) as any[]
        const normalized: Item[] = rs.map((x: any) => ({
          partition: x.partition || { id: x.partition_id, title: x.title, thickness: x.thickness },
          comment: x.comment,
        }))
        setItems(normalized)
      })
      .catch(e => setError(String(e)))
      .finally(() => setLoading(false))
  }, [id])

  const editable = status === 'черновик'

  const onRemove = async (sid: number) => {
    if (!editable || !id) return
    setLoading(true)
    try {
      await RequestPartitionsService.remove(Number(id), sid)
      setItems(prev => prev.filter(i => i.partition.id !== sid))
      // Сохраняем изменения в localStorage
      localStorage.setItem(`draft_request_${id}`,
        JSON.stringify({ status, roomArea, expertComment, items: items.filter(i => i.partition.id !== sid) })
      );
    } catch (e) {
      setError(String(e))
    } finally { setLoading(false) }
  }

  const updateLocalItemComment = (sid: number, value: string) => {
    if (!id) return
    setItems(prev => {
      const next = prev.map(i => i.partition.id === sid ? { ...i, comment: value } : i)
      localStorage.setItem(`draft_request_${id}`,
        JSON.stringify({ status, roomArea, expertComment, items: next })
      )
      return next
    })
  }

  const onSaveComment = async (sid: number, value: string) => {
    if (!editable || !id) return
    setLoading(true)
    try {
      await RequestPartitionsService.update(Number(id), sid, { comment: value })
    } catch (e) {
      setError(String(e))
    } finally { setLoading(false) }
  }

  const onForm = async () => {
    if (!id) return
    setLoading(true)
    try {
      // Сначала сохраняем площадь и комментарий
      await RequestsService.update(Number(id), {
        room_area: roomArea ? Number(roomArea) : undefined,
        expert_comment: expertComment || undefined,
      })
      // Затем формируем заявку
      await RequestsService.form(Number(id))
      setStatus('сформирован')
      // Не сбрасываем черновик, сохраняем его
      localStorage.setItem(`draft_request_${id}`,
        JSON.stringify({ status: 'сформирован', roomArea, expertComment, items })
      );
    } catch (e) { setError(String(e)) } finally { setLoading(false) }
  }

  const onSave = async () => {
    if (!id) return
    setLoading(true)
    try {
      await RequestsService.update(Number(id), {
        room_area: roomArea ? Number(roomArea) : undefined,
        expert_comment: expertComment || undefined,
      })
      // Сохраняем изменения в localStorage
      localStorage.setItem(`draft_request_${id}`,
        JSON.stringify({ status, roomArea, expertComment, items })
      );
    } catch (e) { setError(String(e)) } finally { setLoading(false) }
  }

  // Кнопка удаления заявки
  const onDelete = async () => {
    if (!id) return;
    if (!window.confirm('Удалить заявку?')) return;
    setLoading(true);
    try {
      await RequestsService.remove(Number(id));
      localStorage.removeItem(`draft_request_${id}`);
      nav('/calculations');
    } catch (e) { setError(String(e)) } finally { setLoading(false); }
  }

  if (error) return <Alert variant="danger">Ошибка: {error}</Alert>

  return (
    <div>
      <div className="d-flex justify-content-between align-items-center mb-3">
        <h2>Заявка #{id} {status && <Badge bg={editable ? 'warning' : 'secondary'}>{status}</Badge>}</h2>
        {loading && <span><Spinner size="sm"/> Обновление…</span>}
      </div>
      <Row className="g-3">
        <Col md={8}>
          <Table bordered responsive>
            <thead><tr><th>Перегородка</th><th>Ширина</th><th>Комментарий</th><th></th></tr></thead>
            <tbody>
              {items.map(it => (
                <tr key={it.partition.id}>
                  <td>{normalizeTitle(it.partition.title, it.partition.thickness)}</td>
                  <td style={{ width: 140 }}>{it.partition.thickness || '—'}</td>
                  <td>
                    {editable ? (
                      <Form.Control
                        value={it.comment || ''}
                        onChange={e => updateLocalItemComment(it.partition.id, e.target.value)}
                        onBlur={e => onSaveComment(it.partition.id, e.target.value)}
                      />
                    ) : (it.comment || '—')}
                  </td>
                  <td className="text-end" style={{ width: 120 }}>
                    {editable && <Button size="sm" variant="outline-danger" onClick={() => onRemove(it.partition.id)}>Удалить</Button>}
                  </td>
                </tr>
              ))}
              {items.length === 0 && (
                <tr><td colSpan={4} className="text-muted text-center">Пусто</td></tr>
              )}
            </tbody>
          </Table>
        </Col>
        <Col md={4}>
          <Card>
            <Card.Body>
              <Card.Title>Параметры</Card.Title>
              <Form.Group className="mb-2">
                <Form.Label>Площадь помещения (м²)</Form.Label>
                <Form.Control type="number" step="0.01" value={roomArea} disabled={!editable} onChange={e => setroomArea(e.target.value)} />
              </Form.Group>
              <Form.Group className="mb-3">
                <Form.Label>Комментарий эксперта</Form.Label>
                <Form.Control as="textarea" rows={3} value={expertComment} disabled={!editable} onChange={e => setexpertComment(e.target.value)} />
              </Form.Group>
              {editable ? (
                <div className="d-grid gap-2">
                  <Button onClick={onSave} disabled={loading} variant="secondary">Сохранить черновик</Button>
                  <Button onClick={onForm} disabled={loading || items.length === 0}>Подтвердить заявку</Button>
                  <Button onClick={onDelete} disabled={loading} variant="outline-danger">Удалить заявку</Button>
                </div>
              ) : (
                <div className="text-muted">Редактирование недоступно для статуса «{status}»</div>
              )}
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </div>
  )
}
