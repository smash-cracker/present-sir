import requests
import sys

BASE_URL = "https://archeios.cezentech.com"

if len(sys.argv) < 3:
    print("Usage:")
    print("attendance <email> <password>")
    sys.exit(1)

email = sys.argv[1]
password = sys.argv[2]

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