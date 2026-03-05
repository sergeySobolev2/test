$ErrorActionPreference = 'Stop'

$root = Split-Path -Parent $PSScriptRoot
$toolsDir = Join-Path $root 'tools'
$certDir = Join-Path $root 'certs'

New-Item -ItemType Directory -Force -Path $toolsDir | Out-Null
New-Item -ItemType Directory -Force -Path $certDir | Out-Null

$mkcert = Join-Path $toolsDir 'mkcert.exe'

if (-not (Test-Path $mkcert)) {
  # Официальный релиз mkcert для Windows (amd64)
  $url = 'https://github.com/FiloSottile/mkcert/releases/download/v1.4.4/mkcert-v1.4.4-windows-amd64.exe'
  Write-Host "Downloading mkcert from $url"
  Invoke-WebRequest -UseBasicParsing -Uri $url -OutFile $mkcert
}

Write-Host 'Installing local CA (may require permission prompt)...'
& $mkcert -install

$certFile = Join-Path $certDir 'localhost.pem'
$keyFile  = Join-Path $certDir 'localhost-key.pem'

Write-Host 'Generating certificate for localhost/127.0.0.1/::1'
& $mkcert -cert-file $certFile -key-file $keyFile localhost 127.0.0.1 ::1

Write-Host "Done. Generated: $certFile and $keyFile"
