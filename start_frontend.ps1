# Quick frontend start script
Set-Location "$PSScriptRoot\frontend"

if (-not (Test-Path "node_modules")) {
  npm install
}

npm run dev
