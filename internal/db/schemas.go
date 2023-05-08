package db

import (
	"time"

	"gorm.io/gorm"
)

// Schemas
type User struct {
	// gorm.Model `json:"-"`
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `json:"name"`
	Username  string         `json:"username"`
	Email     string         `json:"email" gorm:"uniqueIndex"`
	Password  string         `json:"-"`
	Role      string         `json:"role" gorm:"default:user"`
}

type Property struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	Postcode         int            `json:"postcode"`
	Property_Name    string         `json:"property_name" gorm:"not null;uniqueIndex"`
	Suburb           string         `json:"suburb"`
	City             string         `json:"city"`
	Street_Address_1 string         `json:"street_address_1"`
	Street_Address_2 string         `json:"street_address_2"`
	Bedrooms         float32        `json:"bedrooms"`
	Bathrooms        float32        `json:"bathrooms"`
	Land_Area        float64        `json:"land_area"`
	Land_Metric      string         `json:"land_metric"`
	Description      string         `json:"description"`
	Notes            string         `json:"notes"`
	Features         []Feature      `json:"features" gorm:"many2many:prop_features"`
}

type Feature struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Feature_Name string         `json:"feature_name" gorm:"not null;uniqueIndex"`
	Properties   []Property     `json:"properties" gorm:"many2many:prop_features"`
}
