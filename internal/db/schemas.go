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
	Managed          bool           `json:"managed" gorm:"default:false"`
	// One to many
	PropertyLogs []PropertyLog `json:"property_logs" gorm:"foreignKey:PropertyID"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:PropertyID"`

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
	Properties   []Property    `json:"properties" gorm:"many2many:contact_properties"`
	Transactions []Transaction `json:"transactions" gorm:"many2many:contact_transactions"`
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
	// One to one
	// Note: Relationship between tasks and properties are handled within transactions
	// TransactionID uint `json:"transaction_id,omitempty" gorm:"unique;default:null"`
	Transaction Transaction `json:"transaction,omitempty" gorm:"unique;foreignKey:TaskID"`
	// MaintenanceRequestID uint               `json:"maintenance_request_id,omitempty" gorm:"unique;default:null"`
	MaintenanceRequest MaintenanceRequest `json:"maintenance_request,omitempty" gorm:"unique;foreignKey:TaskID"`
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
type Transaction struct {
	ID        uint           `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index,omitempty"`
	// Required fields
	Type    string `json:"type,omitempty" gorm:"not null;enum:Sale,Lease,Management,Other"`
	Agency  string `json:"agency,omitempty" gorm:"not null;enum:Own,Other"`
	IsLease bool   `json:"is_lease,omitempty" gorm:"default:false"`
	// Optional fields
	TenancyType           string    `json:"tenancy_type,omitempty" gorm:"enum:Monthly,LongTerm,ShortTerm,Commercial,NA"`
	AgencyName            string    `json:"agency_name,omitempty" gorm:""`
	TransactionNotes      string    `json:"transaction_notes,omitempty" gorm:"default:null"`
	TransactionValue      float64   `json:"transaction_value,omitempty" gorm:"default:null"`
	TransactionCompletion time.Time `json:"transaction_completion,omitempty" gorm:"default:null"`
	Fee                   float32   `json:"fee,omitempty" gorm:"default:null"`

	// Many to one (requires uint for key and Property for object data)
	PropertyID uint     `json:"property_id,omitempty" gorm:"not null"`
	Property   Property `json:"property,omitempty" gorm:"foreignKey:PropertyID"`
	// Many to many
	Contacts []Contact `json:"contacts,omitempty" gorm:"many2many:contact_transactions"`
	// One to one
	TaskID uint `json:"task_id,omitempty" gorm:"unique"`
	// Task   Task `json:"task,omitempty" gorm:"foreignKey:TaskID"`
}

type MaintenanceRequest struct {
	ID        uint           `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index,omitempty"`
	// Required fields
	WorkDefinition string  `json:"work_definition,omitempty" gorm:"not null;enum:Repair,Replacement,Project,Investigation,Pest Control,Other"`
	Type           string  `json:"type,omitempty" gorm:"not null;enum:Electrical,Plumbing,Painting,HVAC,Civil,Other"`
	Notes          string  `json:"notes,omitempty" gorm:"default:null"`
	Scale          string  `json:"scale,omitempty" gorm:"not null;enum:Urgent,High,Medium,Low"`
	TotalCost      float64 `json:"cost,omitempty" gorm:"default:null"`
	Tax            float64 `json:"tax,omitempty" gorm:"default:null"`

	// Relationships
	// Many to one (requires uint for key and Property for object data)
	PropertyID uint     `json:"property_id,omitempty" gorm:"not null"`
	Property   Property `json:"property,omitempty" gorm:"foreignKey:PropertyID"`
	// One to one
	WorkTypeID uint     `json:"work_type_id,omitempty" gorm:""`
	WorkType   WorkType `json:"work_type,omitempty" gorm:"foreignKey:WorkTypeID"`
	TaskID     uint     `json:"task,omitempty" gorm:"unique"`
	// One to many
	// VendorID uint `json:"vendor_id,omitempty" gorm:""`
	// Vendor Vendor `json:"vendor_id,omitempty" gorm:""`
}

type WorkType struct {
	ID        uint           `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index,omitempty"`
	// Required fields
	Name string `json:"name,omitempty" gorm:"unique,not null"`
	// Relationships
	// One to many
	MaintenanceRequests []MaintenanceRequest `json:"maintenance_requests,omitempty" gorm:"foreignKey:WorkTypeID"`
	// Many to many
	Vendors []Vendor `json:"vendors,omitempty" gorm:"many2many:vendor_work_types"`
}

type Vendor struct {
	ID        uint           `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index,omitempty"`
	// Required fields
	CompanyName string `json:"company_name,omitempty" gorm:"unique,not null"`
	NPWP        string `json:"npwp,omitempty" gorm:"not null"`
	NIB         string `json:"nib,omitempty" gorm:""`
	Email       string `json:"email,omitempty" gorm:""`
	Phone       string `json:"phone,omitempty" gorm:""`
	// Address fields
	Street_Address_1 string `json:"street_address_1,omitempty" gorm:""`
	Street_Address_2 string `json:"street_address_2,omitempty" gorm:""`
	City             string `json:"city,omitempty" gorm:""`
	Province         string `json:"province,omitempty" gorm:"not null,default:Bali,enum:Aceh,Bali,Banten,Bengkulu,Central Java,Central Kalimantan,Central Sulawesi,East Java,East Kalimantan,East Nusa Tenggara,Gorontalo,Jakarta Special Capital Region,Jambi,Lampung,Maluku,North Kalimantan,North Maluku,North Sulawesi,North Sumatra,Papua,Riau,Riau Islands,South Kalimantan,South Sulawesi,South Sumatra,Southeast Sulawesi,West Java,West Kalimantan,West Nusa Tenggara,West Papua,West Sulawesi,West Sumatra,Yogyakarta Special Region"`
	Postal_Code      string `json:"postal_code,omitempty" gorm:""`
	Suburb           string `json:"suburb,omitempty" gorm:"not null,default:Badung"`
	// Relationships
	// One to many
	// MaintenanceRequests []MaintenanceRequest `json:"maintenance_requests,omitempty" gorm:"foreignKey:VendorID"`
	// Many to many
	WorkTypes []WorkType `json:"work_types,omitempty" gorm:"many2many:vendor_work_types"`
}
