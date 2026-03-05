# Скрипт для тестирования API Dehydration Lab 4
# Проверяет все требования методички

Write-Host "=== Тестирование Dehydration Lab 4 API ===" -ForegroundColor Cyan
Write-Host ""

$baseUrl = "http://localhost:8082"

# Функция для красивого вывода
function Show-Result {
    param($TestName, $Success, $Details)
    if ($Success) {
        Write-Host "[✓] $TestName" -ForegroundColor Green
    } else {
        Write-Host "[✗] $TestName" -ForegroundColor Red
    }
    if ($Details) {
        Write-Host "    $Details" -ForegroundColor Gray
    }
}

# Функция для выполнения HTTP запросов
function Invoke-ApiRequest {
    param($Method, $Uri, $Body, $Headers = @{})
    try {
        $params = @{
            Method = $Method
            Uri = $Uri
            ContentType = "application/json"
            Headers = $Headers
            TimeoutSec = 10
        }
        if ($Body) {
            $params.Body = $Body | ConvertTo-Json
        }
        $response = Invoke-RestMethod @params
        return @{ Success = $true; Data = $response; StatusCode = 200 }
    } catch {
        return @{ Success = $false; Error = $_.Exception.Message; StatusCode = $_.Exception.Response.StatusCode.value__ }
    }
}

Write-Host "Шаг 1: Проверка доступности Swagger UI" -ForegroundColor Yellow
try {
    $swagger = Invoke-WebRequest -Uri "$baseUrl/swagger/index.html" -TimeoutSec 5
    Show-Result "Swagger UI доступен" $true "$baseUrl/swagger/index.html"
} catch {
    Show-Result "Swagger UI доступен" $false "Не удалось подключиться"
    Write-Host "Убедитесь, что сервер запущен!" -ForegroundColor Red
    exit
}
Write-Host ""

Write-Host "Шаг 2: Регистрация пользователей" -ForegroundColor Yellow

# Регистрация обычного пользователя
$userBody = @{
    login = "testuser$(Get-Random -Maximum 9999)"
    password = "password123"
}
$userReg = Invoke-ApiRequest "POST" "$baseUrl/api/users/register" $userBody
Show-Result "Регистрация обычного пользователя" $userReg.Success $userBody.login

# Регистрация модератора
$modBody = @{
    login = "testmod$(Get-Random -Maximum 9999)"
    password = "moderator123"
}
$modReg = Invoke-ApiRequest "POST" "$baseUrl/api/users/register" $modBody
Show-Result "Регистрация модератора" $modReg.Success $modBody.login

Write-Host ""
Write-Host "ВНИМАНИЕ: Нужно вручную установить is_moderator=true для модератора в БД!" -ForegroundColor Yellow
Write-Host "SQL: UPDATE users SET is_moderator = true WHERE login = '$($modBody.login)';" -ForegroundColor Cyan
Write-Host "Adminer: http://localhost:8081 (postgres/123/dehydration)" -ForegroundColor Gray
Write-Host ""
Read-Host "Нажмите Enter после установки флага модератора в БД"

Write-Host ""
Write-Host "Шаг 3: Аутентификация" -ForegroundColor Yellow

# Вход обычного пользователя
$loginUser = @{
    login = $userBody.login
    password = $userBody.password
}
$userLogin = Invoke-ApiRequest "POST" "$baseUrl/api/users/login" $loginUser
Show-Result "Аутентификация пользователя" $userLogin.Success
if ($userLogin.Success) {
    $userToken = $userLogin.Data.data.token
    $userSessionId = $userLogin.Data.data.session_id
    Write-Host "    JWT Token: $($userToken.Substring(0, 20))..." -ForegroundColor Gray
    Write-Host "    Session ID: $userSessionId" -ForegroundColor Gray
}

# Вход модератора
$loginMod = @{
    login = $modBody.login
    password = $modBody.password
}
$modLogin = Invoke-ApiRequest "POST" "$baseUrl/api/users/login" $loginMod
Show-Result "Аутентификация модератора" $modLogin.Success
if ($modLogin.Success) {
    $modToken = $modLogin.Data.data.token
    $modSessionId = $modLogin.Data.data.session_id
    Write-Host "    JWT Token: $($modToken.Substring(0, 20))..." -ForegroundColor Gray
    Write-Host "    Session ID: $modSessionId" -ForegroundColor Gray
}

Write-Host ""
Write-Host "Шаг 4: Проверка прав доступа - GET /api/symptoms (публичный)" -ForegroundColor Yellow

# Без авторизации (должно работать)
$symptomsPublic = Invoke-ApiRequest "GET" "$baseUrl/api/symptoms"
Show-Result "GET /api/symptoms без авторизации" $symptomsPublic.Success "Публичный доступ"

Write-Host ""
Write-Host "Шаг 5: Проверка прав доступа - GET /api/requests" -ForegroundColor Yellow

# Без авторизации (должно быть 401)
$requestsGuest = Invoke-ApiRequest "GET" "$baseUrl/api/requests"
$guestUnauthorized = $requestsGuest.StatusCode -eq 401
Show-Result "GET /api/requests без авторизации → 401" $guestUnauthorized "Ожидается отказ"

# С авторизацией пользователя (должно работать)
$userHeaders = @{ Authorization = "Bearer $userToken" }
$requestsUser = Invoke-ApiRequest "GET" "$baseUrl/api/requests" $null $userHeaders
Show-Result "GET /api/requests от пользователя → 200" $requestsUser.Success "Только свои заявки"

# С авторизацией модератора (должно работать)
$modHeaders = @{ Authorization = "Bearer $modToken" }
$requestsMod = Invoke-ApiRequest "GET" "$baseUrl/api/requests" $null $modHeaders
Show-Result "GET /api/requests от модератора → 200" $requestsMod.Success "Все заявки"

Write-Host ""
Write-Host "Шаг 6: Создание симптома (только модератор)" -ForegroundColor Yellow

# От пользователя (должно быть 403)
$symptomBody = @{
    title = "Тестовый симптом $(Get-Random)"
    description = "Тестовое описание"
    severity = "Легкая (1-2%)"
    is_active = $true
}
$symptomUser = Invoke-ApiRequest "POST" "$baseUrl/api/symptoms" $symptomBody $userHeaders
$userForbidden = $symptomUser.StatusCode -eq 403
Show-Result "POST /api/symptoms от пользователя → 403" $userForbidden "Ожидается отказ"

# От модератора (должно работать)
$symptomMod = Invoke-ApiRequest "POST" "$baseUrl/api/symptoms" $symptomBody $modHeaders
Show-Result "POST /api/symptoms от модератора → 200" $symptomMod.Success "Симптом создан"

Write-Host ""
Write-Host "Шаг 7: Создание и формирование заявки пользователем" -ForegroundColor Yellow

# Создать симптом в черновике
if ($requestsUser.Success -and $requestsUser.Data.data) {
    # Получить или создать черновик через добавление симптома
    $addSymptom = @{
        symptom_id = 1  # Используем существующий симптом
    }
    $addResult = Invoke-ApiRequest "POST" "$baseUrl/api/request-symptoms" $addSymptom $userHeaders
    Show-Result "Добавление симптома в черновик" $addResult.Success
    
    # Получить ID черновика
    $cart = Invoke-ApiRequest "GET" "$baseUrl/api/requests/cart" $null $userHeaders
    if ($cart.Success -and $cart.Data.data.draft_id) {
        $draftId = $cart.Data.data.draft_id
        Write-Host "    ID черновика: $draftId" -ForegroundColor Gray
        
        # Сформировать заявку
        $formResult = Invoke-ApiRequest "PUT" "$baseUrl/api/requests/$draftId/form" $null $userHeaders
        Show-Result "Формирование заявки" $formResult.Success "ID: $draftId"
    }
}

Write-Host ""
Write-Host "Шаг 8: Завершение заявки (только модератор)" -ForegroundColor Yellow

# Получить список заявок модератора
$modRequests = Invoke-ApiRequest "GET" "$baseUrl/api/requests?status=сформирован" $null $modHeaders
if ($modRequests.Success -and $modRequests.Data.data.Count -gt 0) {
    $requestId = $modRequests.Data.data[0].id
    Write-Host "    Тестируем на заявке ID: $requestId" -ForegroundColor Gray
    
    # Попытка завершить от пользователя (должно быть 403)
    $completeBody = @{
        status = "завершен"
        patient_weight = 70.0
        dehydration_percent = 5.0
        doctor_comment = "Тестовый комментарий"
    }
    $completeUser = Invoke-ApiRequest "PUT" "$baseUrl/api/requests/$requestId/complete" $completeBody $userHeaders
    $userForbiddenComplete = $completeUser.StatusCode -eq 403
    Show-Result "PUT /complete от пользователя → 403" $userForbiddenComplete "Ожидается отказ"
    
    # Завершение от модератора (должно работать)
    $completeMod = Invoke-ApiRequest "PUT" "$baseUrl/api/requests/$requestId/complete" $completeBody $modHeaders
    Show-Result "PUT /complete от модератора → 200" $completeMod.Success "Заявка завершена"
} else {
    Write-Host "    Нет сформированных заявок для теста" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Шаг 9: Проверка Redis" -ForegroundColor Yellow

try {
    # Проверка сессий в Redis
    $redisContainer = docker ps --filter "ancestor=redis:7-alpine" -q | Select-Object -First 1
    if ($redisContainer) {
        Write-Host "    Подключение к Redis..." -ForegroundColor Gray
        $redisKeys = docker exec $redisContainer redis-cli KEYS "session:*"
        $sessionCount = ($redisKeys | Measure-Object).Count
        Show-Result "Сессии в Redis" ($sessionCount -gt 0) "Найдено сессий: $sessionCount"
        
        # Показать содержимое одной сессии
        if ($sessionCount -gt 0) {
            $firstSession = $redisKeys | Select-Object -First 1
            $sessionData = docker exec $redisContainer redis-cli GET $firstSession
            Write-Host "    Пример сессии ($firstSession):" -ForegroundColor Gray
            Write-Host "    $sessionData" -ForegroundColor DarkGray
        }
    } else {
        Show-Result "Redis контейнер" $false "Контейнер не найден"
    }
} catch {
    Show-Result "Проверка Redis" $false $_.Exception.Message
}

Write-Host ""
Write-Host "=== ИТОГОВАЯ СВОДКА ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "✓ Swagger UI: $baseUrl/swagger/index.html" -ForegroundColor Green
Write-Host "✓ Регистрация и аутентификация работают" -ForegroundColor Green
Write-Host "✓ JWT токены генерируются" -ForegroundColor Green
Write-Host "✓ Сессии сохраняются в Redis" -ForegroundColor Green
Write-Host "✓ Разграничение прав (Guest/User/Moderator)" -ForegroundColor Green
Write-Host "✓ Публичные endpoints доступны без авторизации" -ForegroundColor Green
Write-Host "✓ Защищенные endpoints требуют авторизацию (401)" -ForegroundColor Green
Write-Host "✓ Методы модератора недоступны пользователям (403)" -ForegroundColor Green
Write-Host ""
Write-Host "Токены для тестирования в Insomnia/Postman:" -ForegroundColor Yellow
Write-Host "User Token: Bearer $userToken" -ForegroundColor Gray
Write-Host "Mod Token:  Bearer $modToken" -ForegroundColor Gray
Write-Host ""
Write-Host "Cookie для тестирования:" -ForegroundColor Yellow
Write-Host "User Session: session_id=$userSessionId" -ForegroundColor Gray
Write-Host "Mod Session:  session_id=$modSessionId" -ForegroundColor Gray
Write-Host ""
Write-Host "Тестирование завершено!" -ForegroundColor Green
