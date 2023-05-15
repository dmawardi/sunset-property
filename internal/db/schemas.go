package db

import (
	"time"

	"gorm.io/gorm"
)

// Schemas
// Users
type User struct {
	// gorm.Model `json:"-"`
	ID           uint           `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Name         string         `json:"name"`
	Username     string         `json:"username"`
	Email        string         `json:"email" gorm:"uniqueIndex"`
	Password     string         `json:"-"`
	Role         string         `json:"role" gorm:"default:user"`
	PropertyLogs []PropertyLog  `json:"property_logs"`
}

// Properties
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
	PropertyLogs     []PropertyLog  `json:"property_logs" gorm:"foreignKey:PropertyID"`
}

type Feature struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Feature_Name string         `json:"feature_name" gorm:"not null;uniqueIndex"`
	Properties   []Property     `json:"properties" gorm:"many2many:prop_features"`
}

type PropertyLog struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// Use UserID as foreign key and User as object for relationship data
	UserID     uint   `json:"user_id"`
	User       User   `json:"user" gorm:"not null;foreignKey:UserID"`
	LogMessage string `json:"log_message" gorm:"not null"`
	Type       string `json:"type" gorm:"not null"`
	// Use PropertyID as foreign key and Property as object for relationship data
	PropertyID uint     `json:"property_id" gorm:"not null"`
	Property   Property `json:"property"`
}

// Contacts
type Contact struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	FirstName    string         `json:"first_name" gorm:"not null"`
	LastName     string         `json:"last_name"`
	ContactType  string         `json:"contact_type" gorm:"not null"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone"`
	Mobile       string         `json:"mobile"`
	ContactNotes string         `json:"notes"`
	Properties   []Property     `json:"properties" gorm:"many2many:contact_properties"`
	// Vendors []Vendor `json:"vendors" gorm:"many2many:contact_vendors"`
	// Transactions []Transaction `json:"transactions" gorm:"many2many:contact_transactions"`
}
