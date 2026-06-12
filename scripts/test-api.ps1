param(
    [string]$GrpcurlPath = ".\bin\grpcurl.exe"
)

$ErrorActionPreference = 'Stop'

Write-Host 'Тестирование API микросервисов через gRPC и REST'

# Тест 1
Write-Host ''
Write-Host 'Тест 1: Получение списка деталей из Inventory'
$partsResponse = & $GrpcurlPath -plaintext -d '{"filter":{}}' localhost:50051 inventory.v1.InventoryService/ListParts 2>&1 | Out-String
if ([string]::IsNullOrWhiteSpace($partsResponse) -or $partsResponse -match '"error"') {
    Write-Host 'Не удалось получить список деталей.'
    Write-Host "Ответ сервера: $partsResponse"
    exit 1
}
$partUuid = ([regex]'"uuid":\s*"([^"]+)"').Match($partsResponse).Groups[1].Value
if ([string]::IsNullOrEmpty($partUuid)) {
    Write-Host 'Не удалось найти UUID детали.'
    Write-Host "Ответ: $partsResponse"
    exit 1
}
Write-Host "Список деталей получен. Первая UUID: $partUuid"

# Тест 2
Write-Host ''
Write-Host 'Тест 2: Получение информации о детали по UUID'
$partResponse = & $GrpcurlPath -plaintext -d "{`"uuid`":`"$partUuid`"}" localhost:50051 inventory.v1.InventoryService/GetPart 2>&1 | Out-String
if ([string]::IsNullOrWhiteSpace($partResponse) -or $partResponse -match '"error"') {
    Write-Host 'Не удалось получить информацию о детали.'
    Write-Host "Ответ: $partResponse"
    exit 1
}
$partName = ([regex]'"name":\s*"([^"]+)"').Match($partResponse).Groups[1].Value
if ([string]::IsNullOrEmpty($partName)) {
    Write-Host 'Не удалось извлечь имя детали.'
    exit 1
}
Write-Host "Деталь получена: $partName"

# Тест 3
Write-Host ''
Write-Host 'Тест 3: Генерация UUID пользователя'
$userUuid = [System.Guid]::NewGuid().ToString().ToLower()
Write-Host "UUID пользователя: $userUuid"

# Тест 4
Write-Host ''
Write-Host 'Тест 4: Создание заказа (REST API)'
$orderBody = "{`"user_uuid`":`"$userUuid`",`"part_uuids`":[`"$partUuid`"]}"
try {
    $orderObj = Invoke-RestMethod -Method Post -Uri 'http://localhost:8080/api/v1/orders' -ContentType 'application/json' -Body $orderBody
} catch {
    Write-Host "Не удалось создать заказ: $_"; exit 1
}
$orderUuid = $orderObj.uuid
if ([string]::IsNullOrEmpty($orderUuid)) {
    Write-Host 'Не удалось извлечь UUID заказа.'; exit 1
}
Write-Host "Заказ создан. UUID: $orderUuid"

# Тест 5
Write-Host ''
Write-Host 'Тест 5: Проверка начального статуса (ожидается PENDING_PAYMENT)'
$orderInfo = Invoke-RestMethod -Method Get -Uri "http://localhost:8080/api/v1/orders/$orderUuid"
if ($orderInfo.status -notmatch 'PENDING_PAYMENT') {
    Write-Host "Ожидался PENDING_PAYMENT, получен: $($orderInfo.status)"; exit 1
}
Write-Host "Статус корректный: $($orderInfo.status)"

# Тест 6
Write-Host ''
Write-Host 'Тест 6: Оплата заказа'
try {
    Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/v1/orders/$orderUuid/pay" `
        -ContentType 'application/json' -Body '{"payment_method":"PAYMENT_METHOD_CARD"}' | Out-Null
} catch {
    Write-Host "Ошибка при оплате: $_"; exit 1
}
Write-Host 'Заказ оплачен'

# Тест 7
Write-Host ''
Write-Host 'Тест 7: Проверка статуса после оплаты (ожидается PAID или ASSEMBLED)'
$orderInfo = Invoke-RestMethod -Method Get -Uri "http://localhost:8080/api/v1/orders/$orderUuid"
if ($orderInfo.status -notmatch 'PAID' -and $orderInfo.status -notmatch 'ASSEMBLED') {
    Write-Host "Ожидался PAID/ASSEMBLED, получен: $($orderInfo.status)"; exit 1
}
Write-Host "Статус после оплаты: $($orderInfo.status)"

# Тест 8
Write-Host ''
Write-Host 'Тест 8: Создание второго заказа для отмены'
try {
    $order2Obj = Invoke-RestMethod -Method Post -Uri 'http://localhost:8080/api/v1/orders' -ContentType 'application/json' -Body $orderBody
} catch {
    Write-Host "Не удалось создать второй заказ: $_"; exit 1
}
$order2Uuid = $order2Obj.uuid
if ([string]::IsNullOrEmpty($order2Uuid)) {
    Write-Host 'Не удалось извлечь UUID второго заказа.'; exit 1
}
Write-Host "Второй заказ создан. UUID: $order2Uuid"

$order2Info = Invoke-RestMethod -Method Get -Uri "http://localhost:8080/api/v1/orders/$order2Uuid"
if ($order2Info.status -notmatch 'PENDING_PAYMENT') {
    Write-Host "Ожидался PENDING_PAYMENT, получен: $($order2Info.status)"; exit 1
}
Write-Host "Начальный статус второго заказа: $($order2Info.status)"

# Тест 9
Write-Host ''
Write-Host 'Тест 9: Отмена второго заказа'
Write-Host 'Ожидаем 2 секунды...'
Start-Sleep -Seconds 2
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/v1/orders/$order2Uuid/cancel" -ErrorAction SilentlyContinue | Out-Null

$order2Info = Invoke-RestMethod -Method Get -Uri "http://localhost:8080/api/v1/orders/$order2Uuid"
if ($order2Info.status -notmatch 'CANCELLED') {
    Write-Host "Ожидался CANCELLED, получен: $($order2Info.status)"
    Write-Host "Детали: $($order2Info | ConvertTo-Json)"
    exit 1
}
Write-Host "Статус после отмены: $($order2Info.status)"

Write-Host ''
Write-Host 'Все тесты API успешно выполнены!'
