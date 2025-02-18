package database

import "time"

// Task model
type Task struct {
	ID              int       `gorm:"primaryKey" json:"id"`
	TaskName        string    `gorm:"not null" json:"task_name"`
	TaskDescription string    `json:"task_description"`
	Importance      string    `gorm:"default:critical" json:"importance"`
	IsDone          bool      `gorm:"default:false" json:"is_done"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
