package entities

import "time"

type LogEntry struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	IPAddress string    `json:"ip_address"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Headers   string    `json:"headers"`
	Body      string    `json:"body"`
	Response  string    `json:"response"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (LogEntry) TableName() string {
	return "log"
}
