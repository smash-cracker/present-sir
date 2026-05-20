package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const baseURL = "https://archeios.cezentech.com"

var holidays = map[string]struct{}{
	"01-01-26": {}, // NEW YEAR
	"14-01-26": {}, // MAKARA SANKRANTI
	"26-01-26": {}, // REPUBLIC DAY
	"19-03-26": {}, // UGADI FESTIVAL
	"21-03-26": {}, // RAMZAN
	"03-04-26": {}, // GOOD FRIDAY
	"15-04-26": {}, // VISHU
	"01-05-26": {}, // MAY DAY
	"28-05-26": {}, // BAKRID
	"15-08-26": {}, // INDEPENDENCE DAY
	"26-08-26": {}, // ONAM
	"14-09-26": {}, // GANESH CHATHURTHI
	"02-10-26": {}, // GANDHI JAYANTHI
	"31-10-26": {}, // KANNADA RAJYOTSAV
	"10-11-26": {}, // DIWALI
	"25-12-26": {}, // CHRISTMAS DAY
}

type loginResponse struct {
	Token string `json:"token"`
}

type attendanceResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("attendance <email> <password>")
		os.Exit(1)
	}

	now := time.Now()
	today := now.Format("02-01-06")

	if now.Weekday() == time.Sunday {
		fmt.Println("Today is Sunday. Skipping attendance.")
		notify("Attendance", "Skipped: Sunday")
		return
	}

	if _, ok := holidays[today]; ok {
		fmt.Printf("Today (%s) is a holiday. Skipping attendance.\n", today)
		notify("Attendance", fmt.Sprintf("Skipped: Holiday (%s)", today))
		return
	}

	client := &http.Client{Timeout: 30 * time.Second}

	token, loginBody, err := login(client, os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println("Login failed")
		if len(loginBody) > 0 {
			fmt.Println(string(loginBody))
		}
		notify("Attendance Failed", "Login failed")
		os.Exit(1)
	}

	attendanceBody, attendance, err := markAttendance(client, token)
	if len(attendanceBody) > 0 {
		fmt.Println(string(attendanceBody))
	}
	if attendance.Message == "Attendance already marked for today" {
		return
	}
	if err != nil {
		message := "Attendance marking failed"
		if attendance.Message != "" {
			message = attendance.Message
		}
		notify("Attendance Failed", message)
		os.Exit(1)
	}

	switch {
	case attendance.Success:
		notify("Attendance Success", "Attendance marked successfully")
	default:
		message := attendance.Message
		if message == "" {
			message = "Attendance marking failed"
		}
		notify("Attendance Failed", message)
	}
}

func login(client *http.Client, email, password string) (string, []byte, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	body, err := postJSON(client, baseURL+"/api/login", "", payload)
	if err != nil {
		return "", body, err
	}

	var response loginResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", body, err
	}
	if response.Token == "" {
		return "", body, errors.New("missing token")
	}

	return response.Token, body, nil
}

func markAttendance(client *http.Client, token string) ([]byte, attendanceResponse, error) {
	body, err := postJSON(client, baseURL+"/api/attendance/mark", token, map[string]string{})
	var response attendanceResponse
	if len(body) > 0 {
		if unmarshalErr := json.Unmarshal(body, &response); unmarshalErr != nil {
			return body, response, unmarshalErr
		}
	}

	return body, response, err
}

func postJSON(client *http.Client, url, token string, payload any) ([]byte, error) {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return body, readErr
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return body, fmt.Errorf("request failed with status %s", response.Status)
	}

	return body, nil
}

func notify(title, message string) {
	if _, err := exec.LookPath("notify-send"); err != nil {
		return
	}

	_ = exec.Command("notify-send", title, message).Run()
}
