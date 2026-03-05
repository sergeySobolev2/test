$env:DB_HOST = 'localhost'
$env:DB_PORT = '5432'
$env:DB_USER = 'postgres'
$env:DB_PASS = '123'
$env:DB_NAME = 'dehydration'

Push-Location "$PSScriptRoot"
try {
  go run ./cmd/migrate
} finally {
  Pop-Location
}
