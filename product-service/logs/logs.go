package logs

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type LogEvent struct {
	UserID uint   `json:"user_id"`
	Action string `json:"action"`
}

func SendLog(userID uint, action string) {
	event := LogEvent{UserID: userID, Action: action}
	data, _ := json.Marshal(event)

	http.Post("http://localhost:8081/api/logs", "application/json", bytes.NewBuffer(data))
}
