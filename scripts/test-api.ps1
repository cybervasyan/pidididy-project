param(
    [string]$GrpcurlPath = ".\bin\grpcurl.exe"
)

$ErrorActionPreference = 'Stop'

function Grpc {
    param([string]$Method, [string]$Data)
    $result = cmd /c "echo $Data | `"$GrpcurlPath`" -plaintext -d @ localhost:50051 $Method" 2>&1
    return $result | Out-String
}

Write-Host 'Test 1: ListParts'
$partsResponse = Grpc 'inventory.v1.InventoryService/ListParts' '{"filter":{}}'
if ([string]::IsNullOrWhiteSpace($partsResponse) -or $partsResponse -match 'Error invoking') {
    Write-Host "FAIL: $partsResponse"; exit 1
}
$partUuid = ([regex]'"uuid":\s*"([^"]+)"').Match($partsResponse).Groups[1].Value
if ([string]::IsNullOrEmpty($partUuid)) {
    Write-Host "FAIL: no part UUID. Response: $partsResponse"; exit 1
}
Write-Host "OK: part UUID=$partUuid"

Write-Host 'Test 2: GetPart'
$getPartData = '{"uuid":"' + $partUuid + '"}'
$partResponse = Grpc 'inventory.v1.InventoryService/GetPart' $getPartData
if ([string]::IsNullOrWhiteSpace($partResponse) -or $partResponse -match 'Error invoking') {
    Write-Host "FAIL: $partResponse"; exit 1
}
$partName = ([regex]'"name":\s*"([^"]+)"').Match($partResponse).Groups[1].Value
if ([string]::IsNullOrEmpty($partName)) {
    Write-Host "FAIL: no part name. Response: $partResponse"; exit 1
}
Write-Host "OK: part name=$partName"

Write-Host 'Test 3: Generate user UUID'
$userUuid = [System.Guid]::NewGuid().ToString().ToLower()
Write-Host "OK: user UUID=$userUuid"

Write-Host 'Test 4: CreateOrder'
$orderBody = '{"user_uuid":"' + $userUuid + '","part_uuids":["' + $partUuid + '"]}'
try {
    $orderObj = Invoke-RestMethod -Method Post -Uri 'http://localhost:8080/api/v1/orders' -ContentType 'application/json' -Body $orderBody
} catch {
    Write-Host "FAIL: $_"; exit 1
}
$orderUuid = $orderObj.order_uuid
if ([string]::IsNullOrEmpty($orderUuid)) { $orderUuid = $orderObj.uuid }
if ([string]::IsNullOrEmpty($orderUuid)) {
    Write-Host "FAIL: no order UUID. Response: $($orderObj | ConvertTo-Json)"; exit 1
}
Write-Host "OK: order UUID=$orderUuid"

Write-Host 'Test 5: Check PENDING_PAYMENT'
$orderInfo = Invoke-RestMethod -Method Get -Uri "http://localhost:8080/api/v1/orders/$orderUuid"
if ($orderInfo.status -notmatch 'PENDING_PAYMENT') {
    Write-Host "FAIL: expected PENDING_PAYMENT, got $($orderInfo.status)"; exit 1
}
Write-Host "OK: status=$($orderInfo.status)"

Write-Host 'Test 6: PayOrder'
try {
    Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/v1/orders/$orderUuid/pay" `
        -ContentType 'application/json' -Body '{"payment_method":"CARD"}' | Out-Null
} catch {
    Write-Host "FAIL: $_"; exit 1
}
Write-Host 'OK: paid'

Write-Host 'Test 7: Check PAID'
$orderInfo = Invoke-RestMethod -Method Get -Uri "http://localhost:8080/api/v1/orders/$orderUuid"
if ($orderInfo.status -notmatch 'PAID' -and $orderInfo.status -notmatch 'ASSEMBLED') {
    Write-Host "FAIL: expected PAID/ASSEMBLED, got $($orderInfo.status)"; exit 1
}
Write-Host "OK: status=$($orderInfo.status)"

Write-Host 'Test 8: Create second order'
try {
    $order2Obj = Invoke-RestMethod -Method Post -Uri 'http://localhost:8080/api/v1/orders' -ContentType 'application/json' -Body $orderBody
} catch {
    Write-Host "FAIL: $_"; exit 1
}
$order2Uuid = $order2Obj.order_uuid
if ([string]::IsNullOrEmpty($order2Uuid)) { $order2Uuid = $order2Obj.uuid }
if ([string]::IsNullOrEmpty($order2Uuid)) {
    Write-Host "FAIL: no order2 UUID"; exit 1
}
Write-Host "OK: order2 UUID=$order2Uuid"

$order2Info = Invoke-RestMethod -Method Get -Uri "http://localhost:8080/api/v1/orders/$order2Uuid"
if ($order2Info.status -notmatch 'PENDING_PAYMENT') {
    Write-Host "FAIL: expected PENDING_PAYMENT, got $($order2Info.status)"; exit 1
}
Write-Host "OK: status=$($order2Info.status)"

Write-Host 'Test 9: CancelOrder'
Start-Sleep -Seconds 2
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/v1/orders/$order2Uuid/cancel" -ErrorAction SilentlyContinue | Out-Null

$order2Info = Invoke-RestMethod -Method Get -Uri "http://localhost:8080/api/v1/orders/$order2Uuid"
if ($order2Info.status -notmatch 'CANCELLED') {
    Write-Host "FAIL: expected CANCELLED, got $($order2Info.status)"; exit 1
}
Write-Host "OK: status=$($order2Info.status)"

Write-Host ''
Write-Host 'All tests passed!'