package models

import (
	"time"

	"github.com/dmawardi/Go-Template/internal/db"
)

type CreateTransaction struct {
	// Required fields
	Type   string `json:"type,omitempty" valid:"required,in(Sale|Lease|Management|Other)"`
	Agency string `json:"agency,omitempty" valid:"required,in(Own|Other)"`
	// Optional fields
	IsLease          bool    `json:"is_lease,omitempty" valid:"bool"`
	TenancyType      string  `json:"tenancy_type,omitempty" valid:"in(Monthly|LongTerm|ShortTerm|Commercial|NA)"`
	AgencyName       string  `json:"agency_name,omitempty" valid:"length(2|36)"`
	TransactionNotes string  `json:"transaction_notes,omitempty" valid:"length(5|320)"`
	TransactionValue float64 `json:"transaction_value,omitempty" valid:"float"`
	Fee              float32 `json:"fee,omitempty" valid:"float"`

	// Relationships (Not editable through update)
	Property db.Property `json:"property,omitempty" valid:"required"`
	Task     db.Task     `json:"task,omitempty" valid:"required"`
}
type UpdateTransaction struct {
	Type   string `json:"type,omitempty" valid:"in(Sale|Lease|Management|Other)"`
	Agency string `json:"agency,omitempty" valid:"in(Own|Other)"`
	// Optional fields
	IsLease          bool    `json:"is_lease,omitempty" valid:"bool"`
	TenancyType      string  `json:"tenancy_type,omitempty" valid:"in(Monthly|LongTerm|ShortTerm|Commercial|NA)"`
	AgencyName       string  `json:"agency_name,omitempty" valid:"length(2|36)"`
	TransactionNotes string  `json:"transaction_notes,omitempty" valid:"length(5|320)"`
	TransactionValue float64 `json:"transaction_value,omitempty" valid:"float"`
	Fee              float32 `json:"fee,omitempty" valid:"float"`
	// Editable only through update
	TransactionCompletion time.Time `json:"transaction_completion,omitempty" valid:"time"`

	// Relationships (Can only update contacts)
	Contacts []db.Contact `json:"contacts,omitempty"`
}
