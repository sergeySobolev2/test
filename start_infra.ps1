Set-Location "$PSScriptRoot"

# Запускаем инфраструктуру и инициализацию MinIO (публичный bucket image)
docker compose up -d postgres redis minio minio-init
