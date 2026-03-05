# Quick backend start script
Set-Location "$PSScriptRoot\backend"

$env:DB_HOST = "localhost"
$env:DB_PORT = "5433"
$env:DB_USER = "partition_user"
$env:DB_PASS = "123"
$env:DB_NAME = "partition"
$env:ASYNC_SERVICE_URL = "http://127.0.0.1:8083"
$env:CORS_ALLOW_ALL = "true"

# HTTPS (mkcert) для локального фронта/PWA/Tauri
$certPath = "$PSScriptRoot\frontend\certs\localhost.pem"
$keyPath = "$PSScriptRoot\frontend\certs\localhost-key.pem"
if ((Test-Path $certPath) -and (Test-Path $keyPath)) {
	$env:TLS_CERT_FILE = $certPath
	$env:TLS_KEY_FILE = $keyPath
} else {
	Remove-Item Env:TLS_CERT_FILE -ErrorAction SilentlyContinue
	Remove-Item Env:TLS_KEY_FILE -ErrorAction SilentlyContinue
}

go mod download

Set-Location cmd\migrate
go run main.go
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Set-Location ..\seed
go run main.go
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Set-Location ..\..
go run cmd\main\main.go
