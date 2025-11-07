package logs

import (
	"fmt"
	"ooolalex/project-service/db"
	"ooolalex/project-service/models"
)

func SendLog(userID uint, action string) {
	go db.DB.Create(&models.Log{
		UserID: &userID,
		Action: action,
	})
	fmt.Println("Log:", action)
}
