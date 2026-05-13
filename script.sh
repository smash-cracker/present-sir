#!/usr/bin/env bash

EMAIL="$1"
PASSWORD="$2"

BASE_URL="https://archeios.cezentech.com"

if [ -z "$EMAIL" ] || [ -z "$PASSWORD" ]; then
  echo "Usage:"
  echo "./attendance.sh <email> <password>"
  exit 1
fi

LOGIN_RESPONSE=$(
  curl -s -X POST "$BASE_URL/api/login" \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"$EMAIL\",
      \"password\": \"$PASSWORD\"
    }"
)

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')

ATTENDANCE_RESPONSE=$(
  curl -s -X POST "$BASE_URL/api/attendance/mark" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{}'
)

echo "$ATTENDANCE_RESPONSE"