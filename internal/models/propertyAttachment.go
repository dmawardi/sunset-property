package models

import "github.com/dmawardi/Go-Template/internal/db"

// Create is only available to admin users (others must upload an attachment to create a record)
type CreatePropertyAttachment struct {
	Label     string      `json:"label,omitempty" valid:"required,length(4|32)"`
	FileName  string      `json:"file_name,omitempty" valid:"required,length(2|64)"`
	FileSize  int64       `json:"file_size,omitempty" valid:"required"`
	FileType  string      `json:"file_type,omitempty" valid:"required,length(2|15)"`
	ETag      string      `json:"etag,omitempty" valid:"required,length(6|64)"`
	ObjectKey string      `json:"object_key,omitempty" valid:"required,length(6|120)"`
	Property  db.Property `json:"property,omitempty" valid:"required"`
}

// Available to admin users only
type UpdatePropertyAttachment struct {
	Label     string      `json:"label,omitempty" valid:"length(4|32)"`
	FileName  string      `json:"file_name,omitempty" valid:"length(2|64)"`
	FileSize  int64       `json:"file_size,omitempty" valid:""`
	FileType  string      `json:"file_type,omitempty" valid:"length(2|15)"`
	ETag      string      `json:"etag,omitempty" valid:"length(6|64)"`
	ObjectKey string      `json:"object_key,omitempty" valid:"length(6|120)"`
	Property  db.Property `json:"property,omitempty" valid:""`
}

// Available to all users
type UpdatePropertyAttachmentLabel struct {
	Label string `json:"label" valid:"required"`
}
