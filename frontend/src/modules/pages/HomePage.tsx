import { Card, Col, Row } from 'react-bootstrap'
import { getApiBase, isGuestMode } from '../../config/appConfig'

export function HomePage() {
  const apiBase = getApiBase()
  const guest = isGuestMode()
  return (
    <>
      <Row className="g-3">
        <Col md={8}>
          <h1 className="display-6">Partition Soundproofing</h1>
          <p className="text-muted">
            Этот сервис помогает рассчитать уровень звукоизоляции межкомнатных перегородок,
            подобрать материал и быстро оценить итоговые параметры по выбранным характеристикам.
          </p>
        </Col>
        <Col md={4}>
          <Card className="h-100">
            <Card.Body>
              <Card.Title>О проекте</Card.Title>
              <div className="text-muted small">
                Расчет звукоизоляции, сравнение материалов и формирование заявок на перегородки.
              </div>
            </Card.Body>
          </Card>
        </Col>
      </Row>

      <Row className="mt-4">
        <Col md={12}>
          <Card>
            <Card.Body>
              <Card.Title>Навигация</Card.Title>
              <ul className="mb-2">
                <li>Перегородки — список карточек с фильтрами</li>
                {!guest && <li>Авторизация — вход/регистрация пользователя</li>}
                {!guest && <li>Заявки — список и статус заявок пользователя</li>}
                {!guest && <li>Черновик — добавление перегородок и формирование заявки</li>}
              </ul>
              <div className="small text-muted">
                API base: <span className="fw-semibold">{apiBase || '(proxy / пусто)'}</span>
              </div>
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </>
  )
}
