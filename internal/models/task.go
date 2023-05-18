package models

import (
	"time"

	"github.com/dmawardi/Go-Template/internal/db"
)

type CreateTask struct {
	// Required fields
	TaskName string `json:"task_name,omitempty" valid:"required, length(2|36)"`
	Type     string `json:"type,omitempty" valid:"required,in(maintenance|inspection|transaction|other)"`
	// Optional fields
	Status      string    `json:"status,omitempty" valid:"length(2|36), in(created|open|pending|cancelled|processing|active|completed|archived)"`
	Notes       string    `json:"notes,omitempty" valid:"length(5|320)"`
	SnoozedTill time.Time `json:"snoozed_till,omitempty" valid:"time"`
	Snoozed     bool      `json:"snoozed,omitempty" valid:"bool"`
	Completed   bool      `json:"completed,omitempty" valid:"bool"`
	// Relationships
	Assignment []db.User `json:"assignment,omitempty"`
}

type UpdateTask struct {
	// Required fields
	TaskName string `json:"task_name,omitempty" valid:"length(2|36)"`
	Status   string `json:"status,omitempty" valid:"length(2|36), in(created|open|pending|cancelled|processing|active|completed|archived)"`
	Type     string `json:"type,omitempty" valid:"in(maintenance|inspection|transaction|other)"`
	// Optional fields
	Notes       string    `json:"notes,omitempty" valid:"length(5|320)"`
	SnoozedTill time.Time `json:"snoozed_till,omitempty" valid:"time"`
	Snoozed     bool      `json:"snoozed,omitempty" valid:"bool"`
	Completed   bool      `json:"completed,omitempty" valid:"bool"`
	// Relationships
	Assignment []db.User `json:"assignment,omitempty"`
}
