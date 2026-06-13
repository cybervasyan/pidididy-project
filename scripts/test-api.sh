#!/bin/bash
set -e
GRPCURL="${1:-./bin/grpcurl}"

echo "Testing API..."

echo "Test 1: ListParts"
PARTS_RESPONSE=$("$GRPCURL" -plaintext -d '{"filter":{}}' localhost:50051 inventory.v1.InventoryService/ListParts)
PART_UUID=$(echo "$PARTS_RESPONSE" | grep -o '"uuid": "[^"]*' | head -1 | cut -d'"' -f4)
[ -z "$PART_UUID" ] && echo "FAIL: no part UUID" && exit 1
echo "OK: part UUID=$PART_UUID"

echo "Test 2: GetPart"
PART_RESPONSE=$("$GRPCURL" -plaintext -d "{\"uuid\":\"$PART_UUID\"}" localhost:50051 inventory.v1.InventoryService/GetPart)
PART_NAME=$(echo "$PART_RESPONSE" | grep -o '"name": "[^"]*' | cut -d'"' -f4)
[ -z "$PART_NAME" ] && echo "FAIL: no part name" && exit 1
echo "OK: part name=$PART_NAME"

USER_UUID=$(uuidgen | tr '[:upper:]' '[:lower:]')
echo "Test 3: user UUID=$USER_UUID"

echo "Test 4: CreateOrder"
ORDER_RESPONSE=$(curl -sf -X POST "http://localhost:8080/api/v1/orders" \
  -H "Content-Type: application/json" \
  -d "{\"user_uuid\":\"$USER_UUID\",\"part_uuids\":[\"$PART_UUID\"]}")
ORDER_UUID=$(echo "$ORDER_RESPONSE" | grep -o '"uuid":"[^"]*\|"uuid": "[^"]*' | head -1 | cut -d'"' -f4)
[ -z "$ORDER_UUID" ] && echo "FAIL: no order UUID" && exit 1
echo "OK: order UUID=$ORDER_UUID"

echo "Test 5: Check PENDING_PAYMENT"
STATUS=$(curl -sf "http://localhost:8080/api/v1/orders/$ORDER_UUID" | grep -o '"status":"[^"]*\|"status": "[^"]*' | cut -d'"' -f4)
[[ "$STATUS" != *"PENDING_PAYMENT"* ]] && echo "FAIL: expected PENDING_PAYMENT, got $STATUS" && exit 1
echo "OK: status=$STATUS"

echo "Test 6: PayOrder"
curl -sf -X POST "http://localhost:8080/api/v1/orders/$ORDER_UUID/pay" \
  -H "Content-Type: application/json" \
  -d '{"payment_method":"PAYMENT_METHOD_CARD"}' > /dev/null
echo "OK: paid"

echo "Test 7: Check PAID"
STATUS=$(curl -sf "http://localhost:8080/api/v1/orders/$ORDER_UUID" | grep -o '"status":"[^"]*\|"status": "[^"]*' | cut -d'"' -f4)
[[ "$STATUS" != *"PAID"* && "$STATUS" != *"ASSEMBLED"* ]] && echo "FAIL: expected PAID, got $STATUS" && exit 1
echo "OK: status=$STATUS"

echo "Test 8: Create second order"
ORDER2_RESPONSE=$(curl -sf -X POST "http://localhost:8080/api/v1/orders" \
  -H "Content-Type: application/json" \
  -d "{\"user_uuid\":\"$USER_UUID\",\"part_uuids\":[\"$PART_UUID\"]}")
ORDER2_UUID=$(echo "$ORDER2_RESPONSE" | grep -o '"uuid":"[^"]*\|"uuid": "[^"]*' | head -1 | cut -d'"' -f4)
[ -z "$ORDER2_UUID" ] && echo "FAIL: no order2 UUID" && exit 1
echo "OK: order2 UUID=$ORDER2_UUID"

echo "Test 9: CancelOrder"
sleep 2
curl -sf -X POST "http://localhost:8080/api/v1/orders/$ORDER2_UUID/cancel" > /dev/null || true
STATUS=$(curl -sf "http://localhost:8080/api/v1/orders/$ORDER2_UUID" | grep -o '"status":"[^"]*\|"status": "[^"]*' | cut -d'"' -f4)
[[ "$STATUS" != *"CANCELLED"* ]] && echo "FAIL: expected CANCELLED, got $STATUS" && exit 1
echo "OK: status=$STATUS"

echo "All tests passed!"
