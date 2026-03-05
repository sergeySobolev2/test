# Тест добавления перегородки в корзину
$headers = @{
    "Content-Type" = "application/json"
    "Authorization" = "Bearer YOUR_TOKEN_HERE"
}

$body = @{
    partition_id = 1
    calculation_id = $null
} | ConvertTo-Json

Write-Host "Отправка запроса..."
Write-Host "Body: $body"

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8082/api/request-partitions" `
        -Method POST `
        -Headers $headers `
        -Body $body `
        -UseBasicParsing
    
    Write-Host "SUCCESS! Status: $($response.StatusCode)"
    Write-Host "Response: $($response.Content)"
} catch {
    Write-Host "ERROR! $($_.Exception.Message)"
    Write-Host "Response: $($_.Exception.Response)"
}
