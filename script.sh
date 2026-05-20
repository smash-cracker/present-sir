#!/bin/bash

EMAIL="$1"
PASSWORD="$2"

BASE_URL="https://archeios.cezentech.com"

if [ -z "$EMAIL" ] || [ -z "$PASSWORD" ]; then
  echo "Usage:"
  echo "./attendance.sh <email> <password>"
  exit 1
fi

# Check dependencies
command -v jq >/dev/null 2>&1 || {
  echo "jq is required but not installed."
  exit 1
}

command -v notify-send >/dev/null 2>&1 || {
  echo "notify-send is required but not installed."
  exit 1
}

# Current date in DD-MM-YY format
TODAY=$(date +"%d-%m-%y")

# Holidays to skip
HOLIDAYS=(
  "01-01-26" # NEW YEAR
  "14-01-26" # MAKARA SANKRANTI
  "26-01-26" # REPUBLIC DAY
  "19-03-26" # UGADI FESTIVAL
  "21-03-26" # RAMZAN
  "03-04-26" # GOOD FRIDAY
  "15-04-26" # VISHU
  "01-05-26" # MAY DAY
  "28-05-26" # BAKRID
  "15-08-26" # INDEPENDENCE DAY
  "26-08-26" # ONAM
  "14-09-26" # GANESH CHATHURTHI
  "02-10-26" # GANDHI JAYANTHI
  "31-10-26" # KANNADA RAJYOTSAV
  "10-11-26" # DIWALI
  "25-12-26" # CHRISTMAS DAY
)

# Skip Sundays
DAY_OF_WEEK=$(date +%u)
if [ "$DAY_OF_WEEK" -eq 7 ]; then
  echo "Today is Sunday. Skipping attendance."
  notify-send "Attendance" "Skipped: Sunday"
  exit 0
fi

# Skip holidays
for HOLIDAY in "${HOLIDAYS[@]}"; do
  if [ "$TODAY" = "$HOLIDAY" ]; then
    echo "Today ($TODAY) is a holiday. Skipping attendance."
    notify-send "Attendance" "Skipped: Holiday ($TODAY)"
    exit 0
  fi
done

LOGIN_RESPONSE=$(
  curl -s -X POST "$BASE_URL/api/login" \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"$EMAIL\",
      \"password\": \"$PASSWORD\"
    }"
)

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo "Login failed"
  echo "$LOGIN_RESPONSE"

  notify-send "Attendance Failed" "Login failed"
  exit 1
fi

ATTENDANCE_RESPONSE=$(
  curl -s -X POST "$BASE_URL/api/attendance/mark" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{}'
)

SUCCESS=$(echo "$ATTENDANCE_RESPONSE" | jq -r '.success // empty')
MESSAGE=$(echo "$ATTENDANCE_RESPONSE" | jq -r '.message // empty')

echo "$ATTENDANCE_RESPONSE"

if [ "$MESSAGE" = "Attendance already marked for today" ]; then
  exit 0
elif [ "$SUCCESS" = "true" ]; then
  notify-send "Attendance Success" "Attendance marked successfully"
else
  MESSAGE=$(echo "$ATTENDANCE_RESPONSE" | jq -r '.message // "Attendance marking failed"')
  notify-send "Attendance Failed" "$MESSAGE"
fi
