package models

import "github.com/dmawardi/Go-Template/internal/db"

// Struct required by Property Log service
type CreateTaskLog struct {
	User       db.User `json:"user" valid:"required"`
	Task       db.Task `json:"task" valid:"required"`
	LogMessage string  `json:"log_message" valid:"required,length(3|300)"`
	Type       string  `json:"type"`
}

// Struct received by controller/handler
type RecvTaskLog struct {
	LogMessage string  `json:"log_message" valid:"required,length(3|300)"`
	Task       db.Task `json:"task" valid:"required"`
}

// Struct received by service
type UpdateTaskLog struct {
	LogMessage string `json:"log_message" valid:"required,length(3|300)"`
}
