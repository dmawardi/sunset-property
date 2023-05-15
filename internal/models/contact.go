package models

type CreateContact struct {
	FirstName    string `json:"first_name,omitempty" valid:"length(2|36),required"`
	LastName     string `json:"last_name,omitempty" valid:"length(2|36)"`
	ContactType  string `json:"contact_type,omitempty" valid:"required"`
	Email        string `json:"email,omitempty" valid:"email;length(6|36)"`
	Phone        string `json:"phone,omitempty" valid:"length(7|16),numeric"`
	Mobile       string `json:"mobile,omitempty" valid:"length(7|16),numeric"`
	ContactNotes string `json:"notes,omitempty" valid:"length(5|320)"`
}

type UpdateContact struct {
	FirstName    string `json:"first_name,omitempty" valid:"length(2|36)"`
	LastName     string `json:"last_name,omitempty" valid:"length(2|36)"`
	ContactType  string `json:"contact_type,omitempty" valid:""`
	Email        string `json:"email,omitempty" valid:"email;length(6|36)"`
	Phone        string `json:"phone,omitempty" valid:"length(7|16),numeric"`
	Mobile       string `json:"mobile,omitempty" valid:"length(7|16),numeric"`
	ContactNotes string `json:"notes,omitempty" valid:"length(5|320)"`
}
