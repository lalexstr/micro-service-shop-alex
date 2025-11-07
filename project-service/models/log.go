package models

import "time"

type Log struct {
	ID        uint `gorm:"primaryKey"`
	UserID    *uint
	Action    string
	CreatedAt time.Time
}
