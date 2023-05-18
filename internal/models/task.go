package models

import (
	"time"

	"github.com/dmawardi/Go-Template/internal/db"
)

type CreateTask struct {
	// Required fields
	TaskName string `json:"task_name,omitempty" valid:"required, length(3|36)"`
	Status   string `json:"status,omitempty" valid:"length(3|36), in(Created|Open|Pending|Cancelled|Processing|Active|Completed|Archived)"`
	Type     string `json:"type,omitempty" valid:"required,in(Maintenance|Inspection|Transaction|Other)"`
	// Optional fields
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
	Status   string `json:"status,omitempty" valid:"length(2|36), in(Created|Open|Pending|Cancelled|Processing|Active|Completed|Archived)"`
	Type     string `json:"type,omitempty" valid:"in(Maintenance|Inspection|Transaction|Other)"`
	// Optional fields
	Notes       string    `json:"notes,omitempty" valid:"length(5|320)"`
	SnoozedTill time.Time `json:"snoozed_till,omitempty" valid:"time"`
	Snoozed     bool      `json:"snoozed,omitempty" valid:"bool"`
	Completed   bool      `json:"completed,omitempty" valid:"bool"`
	// Relationships
	Assignment []db.User `json:"assignment,omitempty"`
}
