package models

import "time"

type ContactStatus string

const (
	ContactNew      ContactStatus = "new"
	ContactAnswered ContactStatus = "answered"
	ContactDeleted  ContactStatus = "deleted"
)

type ContactRequest struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	Contact   string        `json:"contact"`
	Status    ContactStatus `gorm:"type:text;default:new" json:"status"`
	AdminID   *uint         `json:"admin_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
