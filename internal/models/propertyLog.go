package models

import "github.com/dmawardi/Go-Template/internal/db"

// Struct required by Property Log service
type CreatePropertyLog struct {
	User       db.User     `json:"user" valid:"required"`
	Property   db.Property `json:"property" valid:"required"`
	LogMessage string      `json:"log_message" valid:"required,length(3|300)"`
	Type       string      `json:"type"`
}

// Struct received by controller/handler
type RecvPropertyLog struct {
	LogMessage string      `json:"log_message" valid:"required,length(3|300)"`
	Property   db.Property `json:"property" valid:"required"`
}

type UpdatePropertyLog struct {
	LogMessage string `json:"log_message" valid:"required,length(3|300)"`
	Type       string `json:"type"`
}
