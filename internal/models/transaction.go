package models

import (
	"time"

	"github.com/dmawardi/Go-Template/internal/db"
)

// CreateTransaction is used to create a transaction
type CreateTransaction struct {
	// Required fields
	// Limit values to: buy, sell, rent, lease
	Type string `json:"type,omitempty" valid:"required, in(buy, sell, rent, lease)"`
	// Limit values to: own, other
	Agency string `json:"agency,omitempty" valid:"required, in(own, other), length(2|36)"`
	// Limit values to: created, open, pending, cancelled, processing, active, completed, archived
	Status string `json:"status,omitempty" valid:"required, length(2|36), in(created, open, pending, cancelled, processing, active, completed, archived)"`

	// Task fields
	Snoozed     bool      `json:"snoozed,omitempty" valid:"bool"`
	SnoozedTill time.Time `json:"snoozed_till,omitempty" valid:"time"`

	// Relationships
	// Many to one (requires uint for key and Property for object data)
	Property db.Property `json:"property,omitempty" valid:"required"`
}

// UpdateTransaction is used to update a transaction (note: Property cannot be updated)
type UpdateTransaction struct {
	// Required fields
	// Limit values to: buy, sell, rent, lease
	Type string `json:"type,omitempty" valid:"required, in(buy, sell, rent, lease)"`
	// Limit values to: own, other
	Agency string `json:"agency,omitempty" valid:"required, in(own, other), length(2|36)"`
	// Limit values to: created, open, pending, cancelled, processing, active, completed, archived
	Status string `json:"status,omitempty" valid:"length(2|36), in(created, open, pending, cancelled, processing, active, completed, archived)"`

	// Relationships
	// Many to many
	Contacts []db.Contact `json:"contacts,omitempty"`
}
