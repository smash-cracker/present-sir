import requests
import sys
from datetime import datetime

BASE_URL = "https://archeios.cezentech.com"

if len(sys.argv) < 3:
    print("Usage:")
    print("attendance <email> <password>")
    sys.exit(1)

email = sys.argv[1]
password = sys.argv[2]

# Current date in DD-MM-YY format
today = datetime.now().strftime("%d-%m-%y")

# Holidays to skip
HOLIDAYS = {
    "01-01-26",  # NEW YEAR
    "14-01-26",  # MAKARA SANKRANTI
    "26-01-26",  # REPUBLIC DAY
    "19-03-26",  # UGADI FESTIVAL
    "21-03-26",  # RAMZAN
    "03-04-26",  # GOOD FRIDAY
    "15-04-26",  # VISHU
    "01-05-26",  # MAY DAY
    "28-05-26",  # BAKRID
    "15-08-26",  # INDEPENDENCE DAY
    "26-08-26",  # ONAM
    "14-09-26",  # GANESH CHATHURTHI
    "02-10-26",  # GANDHI JAYANTHI
    "31-10-26",  # KANNADA RAJYOTSAV
    "10-11-26",  # DIWALI
    "25-12-26",  # CHRISTMAS DAY
}

# Skip Sundays
if datetime.now().weekday() == 6:
    print("Today is Sunday. Skipping attendance.")
    sys.exit(0)

# Skip holidays
if today in HOLIDAYS:
    print(f"Today ({today}) is a holiday. Skipping attendance.")
    sys.exit(0)

# LOGIN
login_response = requests.post(
    f"{BASE_URL}/api/login",
    json={
        "email": email,
        "password": password
    }
)

login_json = login_response.json()

token = login_json.get("token")

if not token:
    print("Login failed")
    print(login_json)
    sys.exit(1)

# MARK ATTENDANCE
attendance_response = requests.post(
    f"{BASE_URL}/api/attendance/mark",
    headers={
        "Authorization": f"Bearer {token}"
    },
    json={}
)

print(attendance_response.text)