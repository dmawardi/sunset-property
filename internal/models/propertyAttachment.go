package models

import "github.com/dmawardi/Go-Template/internal/db"

// Create is only available to admin users (others must upload an attachment to create a record)
type CreatePropertyAttachment struct {
	Label     string      `json:"label,omitempty" validate:"required,length(4|32)"`
	FileName  string      `json:"file_name,omitempty" validate:"required,length(2|32)"`
	FileSize  int64       `json:"file_size,omitempty" validate:"required"`
	FileType  string      `json:"file_type,omitempty" validate:"requiredlength(2|15)"`
	ETag      string      `json:"etag,omitempty" validate:"required,length(6|32)"`
	ObjectKey string      `json:"object_key,omitempty" validate:"required,length(6|32)"`
	Property  db.Property `json:"property,omitempty" validate:"required"`
}

// Available to admin users only
type UpdatePropertyAttachment struct {
	Label     string      `json:"label,omitempty" validate:"length(4|32)"`
	FileName  string      `json:"file_name,omitempty" validate:"length(2|32)"`
	FileSize  int64       `json:"file_size,omitempty" validate:""`
	FileType  string      `json:"file_type,omitempty" validate:"length(2|15)"`
	ETag      string      `json:"etag,omitempty" validate:"length(6|32)"`
	ObjectKey string      `json:"object_key,omitempty" validate:"length(6|32)"`
	Property  db.Property `json:"property,omitempty" validate:""`
}

// Available to all users
type UpdatePropertyAttachmentLabel struct {
	Label string `json:"label" validate:"required"`
}
