import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { Alert, Badge, Card, Col, Placeholder, Row } from 'react-bootstrap'
import { useAppDispatch, useAppSelector } from '../../store'
import { getPartition } from '../services/partitionsApi'
import { start, setCurrent, setError } from '../../store/partitionsSlice'

const PLACEHOLDER = '/placeholder-partition.svg'

export function PartitionDetailPage() {
  const { id } = useParams()
  const dispatch = useAppDispatch()
  const { current: partition, loading, error } = useAppSelector(s => s.partitions)
  const [img, setImg] = useState<string>(PLACEHOLDER)

  useEffect(() => {
    if (!id) return
    const num = Number(id)
    if (!Number.isFinite(num)) return
    let cancelled = false
    dispatch(start())
    getPartition(num)
      .then(r => {
        if (cancelled) return
        const s = r.data.partition
        dispatch(setCurrent(s))
        const u = s.image_url || s.public_image_url || PLACEHOLDER
        setImg(u)
      })
      .catch(e => { if (!cancelled) dispatch(setError(String(e))) })
    return () => { cancelled = true }
  }, [id])

  if (error) return <Alert variant="danger">Не удалось загрузить: {error}</Alert>

  return (
    <Row>
      <Col md={6}>
        <Card>
          {partition ? (
            <Card.Img className="partition-detail-img" src={img} alt={partition.title} />
          ) : (
            <Placeholder as={Card.Img} animation="wave" />
          )}
        </Card>
        {partition && (
          <div className="text-muted mt-2" style={{ wordBreak: 'break-all' }}>
            Источник изображения: {img !== PLACEHOLDER ? img : 'по умолчанию (плейсхолдер)'}
          </div>
        )}
      </Col>
      <Col md={6}>
        <h2 className="mb-2">{partition ? partition.title.split(' ')[0] : <Placeholder xs={6} />}</h2>
        {partition && !partition.is_active && <Badge bg="secondary" className="mb-2">скрыт</Badge>}
        {partition?.thickness && (
          <div className="mb-3">
            <strong>Ширина:</strong> {partition.thickness}
          </div>
        )}
        {partition?.description && (
          <div>
            <h5>Описание:</h5>
            <p style={{ whiteSpace: 'pre-line' }}>{partition.description}</p>
          </div>
        )}
      </Col>
    </Row>
  )
}
