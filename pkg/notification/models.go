package notification

import "time"

type Notification struct {
	ID       uint      `gorm:"primary_key" json:"id,omitempty"`
	UserID   string    `json:"user_id"`
	UserType string    `json:"user_type"` //staff or doctor or organization
	Title    string    `json:"title"`
	Message  string    `json:"message"`
	Priority string    `json:"priority"` // high or regular
	CreateAt time.Time `json:"created_at"`
}
