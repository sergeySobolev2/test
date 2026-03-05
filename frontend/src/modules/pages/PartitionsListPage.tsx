import { useEffect } from 'react'
import { Badge, Col, Form, Row } from 'react-bootstrap'
import { PartitionCard } from './components/PartitionCard'
import { useAppDispatch, useAppSelector } from '../../store'
import { listPartitions } from '../services/partitionsApi'
import { start, setList, setError } from '../../store/partitionsSlice'
import { fetchCart } from '../../store/requestsSlice'
import { Link } from 'react-router-dom'
import { apply, setTitle } from '../../store/filtersSlice'
import { isGuestMode } from '../../config/appConfig'

export function PartitionsListPage() {
  const dispatch = useAppDispatch()
  const { list, loading, error } = useAppSelector(s => s.partitions)
  const { draftId, draftCount } = useAppSelector(s => s.requests)
  const filters = useAppSelector(s => s.filters)
  const guest = isGuestMode()

  useEffect(() => {
    const title = filters.appliedTitle.trim() || undefined
    let cancelled = false
    dispatch(start())
    listPartitions({ title })
      .then(r => { if (!cancelled) dispatch(setList({ items: r.data, total: r.total })) })
      .catch(e => { if (!cancelled) dispatch(setError(String(e))) })
    return () => { cancelled = true }
  }, [dispatch, filters.appliedTitle])

  useEffect(() => {
    if (!guest) {
      dispatch(fetchCart())
    }
  }, [dispatch, guest])

  useEffect(() => {
    const t = setTimeout(() => {
      dispatch(apply())
    }, 300)
    return () => clearTimeout(t)
  }, [dispatch, filters.title])

  return (
    <>
      <h2 className="mb-3">Типы перегородок</h2>
      <Row className="g-2 mb-3">
        <Col md={8} lg={6}>
          <Form.Control
            placeholder="Поиск по названию..."
            value={filters.title}
            onChange={e => dispatch(setTitle(e.target.value))}
          />
        </Col>
      </Row>

      {(filters.appliedTitle) && (
        <div className="text-muted small mb-3">
          Применено: «{filters.appliedTitle || '—'}»
        </div>
      )}

      {error && <div className="alert alert-danger">Ошибка загрузки: {error}</div>}

      {/* Адаптивная сетка карточек услуг (Bootstrap):
          xs (<576px): 1 колонка
          sm (>=576px): 2 колонки
          md (>=768px): 2 колонки
          lg (>=992px): 3 колонки
          xl (>=1200px): 4 колонки */}
      <Row xs={1} sm={2} md={2} lg={3} xl={4} className="g-3">
        {list.map(s => (
          <Col key={s.id}><PartitionCard s={s} /></Col>
        ))}
      </Row>
      {!loading && list.length === 0 && !error && (
        <div className="text-muted mt-3">Ничего не найдено</div>
      )}

      {/* Корзина в левом нижнем углу */}
      {!guest && draftCount > 0 && (
        <Link
          to={`/calculations/${draftId}`}
          className="floating-cart position-fixed bottom-0 start-0"
          title={`Корзина: ${draftCount} перегородок`}
        >
          <div className="position-relative">
            <img 
              src="/partition-cart-icon.svg" 
              alt="Корзина" 
              width={64}
              height={64}
            />
            <Badge 
              bg="danger" 
              pill 
              className="position-absolute top-0 start-100 translate-middle"
              style={{ fontSize: '1rem', minWidth: '2rem' }}
            >
              {draftCount}
            </Badge>
          </div>
        </Link>
      )}
    </>
  )
}
