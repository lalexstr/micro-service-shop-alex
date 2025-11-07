package models

import "time"

type Log struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    *uint     `json:"user_id"`   // <-- добавили, чтобы можно было указывать администратора
	Action    string    `json:"action"`    // короткое описание действия
	Message   string    `json:"message"`   // подробности
	Timestamp time.Time `json:"timestamp"` // время записи
}
