import { Breadcrumb } from 'react-bootstrap'
import { Link, useLocation } from 'react-router-dom'

function toTitle(segment: string) {
  if (!segment) return ''
  if (Number.isInteger(Number(segment))) return `#${segment}`
  return segment.charAt(0).toUpperCase() + segment.slice(1)
}

export function AppBreadcrumbs() {
  const location = useLocation()
  const parts = location.pathname.split('/').filter(Boolean)
  const items = parts.map((p, idx) => {
    const url = '/' + parts.slice(0, idx + 1).join('/')
    const isLast = idx === parts.length - 1
    return (
      <Breadcrumb.Item key={url} active={isLast} linkAs={isLast ? undefined : Link} linkProps={isLast ? undefined : { to: url }}>
        {toTitle(p)}
      </Breadcrumb.Item>
    )
  })
  return (
    <Breadcrumb className="mb-3">
      <Breadcrumb.Item linkAs={Link} linkProps={{ to: '/' }}>Главная</Breadcrumb.Item>
      {items}
    </Breadcrumb>
  )
}
