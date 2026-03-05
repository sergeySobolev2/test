Set-Location "$PSScriptRoot\asyncservice_rust"

# Запуск Rust asyncservice (устойчиво для Insomnia):
# - собираем release
# - запускаем бинарник отдельным процессом, чтобы он не умирал вместе с cargo
$env:BACKEND_RESULT_URL = "http://127.0.0.1:8082/api/async/result"
$env:LISTEN_PORT = "8083"
$env:HOST = "127.0.0.1"

cargo build --release
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

$exe = Join-Path $PWD "target\release\asyncservice_rust.exe"
if (-not (Test-Path $exe)) {
	Write-Host "Binary not found: $exe" -ForegroundColor Red
	exit 1
}

$p = Start-Process -FilePath $exe -WorkingDirectory $PWD -PassThru
Write-Host "Rust asyncservice started. PID=$($p.Id) URL=http://127.0.0.1:8083" -ForegroundColor Green
