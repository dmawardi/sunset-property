package models

import (
	"github.com/dmawardi/Go-Template/internal/db"
)

type CreateMaintenanceRequest struct {
	WorkDefinition string  `json:"work_definition,omitempty" valid:"required,in(Repair|Replacement|Project|Investigation|Pest Control|Other)"`
	Type           string  `json:"type,omitempty" valid:"in(Electrical|Plumbing|Painting|HVAC|Civil|Other)"`
	Notes          string  `json:"notes,omitempty" valid:"length(5|500)"`
	Scale          string  `json:"scale,omitempty" valid:"required,in(Urgent|High|Medium|Low)"`
	TotalCost      float64 `json:"total_cost,omitempty" valid:"float"`
	Tax            float64 `json:"tax,omitempty" valid:"float"`

	// Relationships (Not editable through update)
	Property db.Property `json:"property,omitempty" valid:"required"`
	Task     db.Task     `json:"task,omitempty" valid:"required"`
}

type UpdateMaintenanceRequest struct {
	WorkDefinition string  `json:"work_definition,omitempty" valid:"in(Repair|Replacement|Project|Investigation|Pest Control|Other)"`
	Type           string  `json:"type,omitempty" valid:"in(Electrical|Plumbing|Painting|HVAC|Civil|Other)"`
	Notes          string  `json:"notes,omitempty" valid:"length(5|500)"`
	Scale          string  `json:"scale,omitempty" valid:"in(Urgent|High|Medium|Low)"`
	TotalCost      float64 `json:"total_cost,omitempty" valid:"float"`
	Tax            float64 `json:"tax,omitempty" valid:"float"`
	// Relationships (Not editable through update)
	Property db.Property `json:"property,omitempty" valid:""`
}
