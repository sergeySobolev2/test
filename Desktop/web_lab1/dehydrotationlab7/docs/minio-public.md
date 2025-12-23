# MinIO: публичный доступ и кастомный публичный URL

Переменные окружения (Windows PowerShell):

```powershell
$env:MINIO_HOST = '127.0.0.1'
$env:MINIO_PORT = '9000'
$env:MINIO_ACCESS_KEY = 'minioadmin'
$env:MINIO_SECRET_KEY = 'minioadmin'
$env:MINIO_BUCKET = 'img'
# Включить HTTPS к MinIO, если нужен
$env:MINIO_USE_SSL = 'false'

# Кастомный публичный base URL (например, CDN/прокси/домен)
$env:MINIO_PUBLIC_BASE = 'https://cdn.example.com'

# Сделать бакет публично читаемым (анонимный GET на объекты)
$env:MINIO_PUBLIC_READ = 'true'
```

Что делает код:
- Если `MINIO_PUBLIC_READ=true`, на бакет устанавливается политика публичного чтения (List/Get и GetObject).
- Публичные ссылки строятся как `<MINIO_PUBLIC_BASE>/<bucket>/<key>`. Если `MINIO_PUBLIC_BASE` пустой — используется `http(s)://<MINIO_HOST>:<MINIO_PORT>`.

Где менять:
- Конфиг и чтение переменных: `internal/app/config/config.go`
- Инициализация MinIO: `cmd/main/main.go`
- Клиент MinIO и установка политики: `internal/app/pkg/storage/minio.go`
