# Deployment Diagram (Mermaid)

```mermaid
graph LR
  subgraph Client
    A[Browser/PWA]
    T[Tauri App]
  end

  subgraph GitHub
    GHP[GitHub Pages (Static)]
  end

  subgraph Backend
    N[Nginx / Go API (8082)]
    DB[(Postgres)]
    M[(MinIO)]
    R[(Redis)]
  end

  A -- HTTPS GET --> GHP
  A -- HTTPS fetch /api --> N
  T -- HTTP(S) fetch /api --> N
  N -- SQL --> DB
  N -- S3 API --> M
  N -- Sessions --> R
```

- PWA: сервис-воркер кэширует статику и `GET /api/symptoms` (NetworkFirst).
- Pages: билд фронта публикуется в `gh-pages` через GitHub Actions.
- Tauri: десктоп-клиент использует тот же API по LAN/IP или домену.
