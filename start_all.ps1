$root = $PSScriptRoot

function Start-Terminal($title, $command) {
	Start-Process powershell -ArgumentList @(
		'-NoProfile',
		'-Command',
		$command
	) -WorkingDirectory $root | Out-Null
	Write-Host "Started: $title" -ForegroundColor Green
}

Write-Host "Запуск полного стека (ЛР6/7/8)..." -ForegroundColor Cyan

# 1) Infra (docker) — postgres/redis/minio
if (Get-Command docker -ErrorAction SilentlyContinue) {
	Start-Terminal "infra (docker)" "& '$root\start_infra.ps1'"
} else {
	Write-Host "Docker не найден. Пропускаю infra." -ForegroundColor Yellow
}

# 2) Backend
Start-Terminal "backend" "& '$root\start_backend.ps1'"

# 3) Async service (Rust)
Start-Terminal "asyncservice" "& '$root\start_async.ps1'"

# 4) Frontend (dev)
Start-Terminal "frontend dev" "& '$root\start_frontend.ps1'"

# 5) PWA preview (HTTPS)
$pwaCmd = "Set-Location '$root\frontend'; if (-not (Test-Path 'node_modules')) { npm install }; npm run build; npm run preview:https -- --host 127.0.0.1 --port 5197"
Start-Terminal "pwa preview" $pwaCmd

# 6) Tauri (release)
$tauriCmd = "Set-Location '$root\frontend'; if (-not (Test-Path 'node_modules')) { npm install }; if (-not (Test-Path '$root\frontend\src-tauri\target\release\app.exe')) { npm run tauri build }; Start-Process -FilePath '$root\frontend\src-tauri\target\release\app.exe'"
Start-Terminal "tauri" $tauriCmd

Write-Host 'Готово. Открой: http://localhost:5196 (frontend), https://127.0.0.1:5197 (PWA)' -ForegroundColor Cyan
