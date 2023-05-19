package db

import (
	"time"

	"gorm.io/gorm"
)

// Schemas
// Users
type User struct {
	// gorm.Model `json:"-"`
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index,omitempty"`
	Name      string         `json:"name,omitempty"`
	Username  string         `json:"username,omitempty"`
	Email     string         `json:"email,omitempty" gorm:"uniqueIndex"`
	Password  string         `json:"-"`
	Role      string         `json:"role,omitempty" gorm:"default:user"`
	// Foreign keys
	PropertyLogs []PropertyLog `json:"property_logs"`
	TaskLogs     []TaskLog     `json:"task_logs"`
	Tasks        []Task        `json:"tasks" gorm:"many2many:user_tasks"`
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
	// One to many
	PropertyLogs []PropertyLog `json:"property_logs" gorm:"foreignKey:PropertyID"`
	// Many to many
	Features []Feature `json:"features" gorm:"many2many:prop_features"`
	Contacts []Contact `json:"contacts" gorm:"many2many:contact_properties"`
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
	Type       string `json:"type" gorm:"not null;enum:INPUT,GEN"`
	// Use PropertyID as foreign key and Property as object for relationship data
	PropertyID uint     `json:"property_id" gorm:"not null"`
	Property   Property `json:"property" gorm:"not null;foreignKey:PropertyID"`
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
	// Relationships
	Properties []Property `json:"properties" gorm:"many2many:contact_properties"`
	// Transactions []Transaction `json:"transactions" gorm:"many2many:contact_transactions"`
	// Vendors []Vendor `json:"vendors" gorm:"many2many:contact_vendors"`
}

// Tasks
type Task struct {
	ID        uint           `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index,omitempty"`
	// Required fields
	TaskName string `json:"task_name,omitempty" gorm:"not null"`
	Type     string `json:"type,omitempty" gorm:"not null, enum:Maintenance,Inspection,Transaction,Other"`
	// Default fields
	Status    string `json:"status,omitempty" gorm:"default:created;enum:Created,Open,Pending,Cancelled,Processing,Active,Completed,Archived"`
	Notes     string `json:"notes,omitempty" gorm:"default:null"`
	Snoozed   bool   `json:"snoozed,omitempty" gorm:"default:false"`
	Completed bool   `json:"completed,omitempty" gorm:"default:false"`
	// Optional fields
	SnoozedTill time.Time `json:"snoozed_till,omitempty"`
	// Relationship fields
	// Many to many
	Assignment []User `json:"assignment,omitempty" gorm:"many2many:user_tasks"`
	// One to many
	Log []TaskLog `json:"log,omitempty" gorm:"foreignKey:TaskID"`
	// Note: Relationship between tasks and properties are handled within transactions
	// One to one
	// TransactionID uint     `json:"property_id,omitempty"`
	// Transaction   Transaction `json:"property,omitempty" gorm:"foreignKey:TransactionID"`
}

type TaskLog struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	LogMessage string         `json:"log_message" gorm:"not null"`
	Type       string         `json:"type" gorm:"not null;enum:INPUT,GEN"`
	// Use UserID as foreign key and User as object for relationship data
	UserID uint `json:"user_id"`
	User   User `json:"user" gorm:"not null;foreignKey:UserID"`
	// Use TaskID as foreign key and Task as object for relationship data
	TaskID uint `json:"task_id" gorm:"not null"`
	Task   Task `json:"task" gorm:"not null;foreignKey:TaskID"`
}

// Types of tasks: Transactions, Maintenance Requests, Inspections, Appraisals, Other
// Transactions
// type Transaction struct {
// 	ID        uint           `json:"id,omitempty" gorm:"primaryKey"`
// 	CreatedAt time.Time      `json:"created_at,omitempty"`
// 	UpdatedAt time.Time      `json:"updated_at,omitempty"`
// 	DeletedAt gorm.DeletedAt `gorm:"index,omitempty"`
// 	Type      string         `json:"type,omitempty" gorm:"not null"`
// 	Agency    string         `json:"agency,omitempty" gorm:"not null"`
// 	Status    string         `json:"status,omitempty" gorm:"not null"`
// 	Notes     string         `json:"notes,omitempty"`
// 	// Many to one (requires uint for key and Property for object data)
// 	PropertyID uint     `json:"property_id,omitempty" gorm:"not null"`
// 	Property   Property `json:"property,omitempty" gorm:"foreignKey:PropertyID"`
// 	// Many to many
// 	Contacts []Contact `json:"contacts,omitempty" gorm:"many2many:contact_transactions"`
// }
