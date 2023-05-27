package models

import "github.com/dmawardi/Go-Template/internal/db"

type CreateVendor struct {
	// Required fields
	CompanyName string `json:"company_name,omitempty" valid:"required"`
	NPWP        string `json:"npwp,omitempty" valid:"required;length(12|20)"`
	NIB         string `json:"nib,omitempty" valid:"length(12|20)"`
	Email       string `json:"email,omitempty" valid:"email"`
	Phone       string `json:"phone,omitempty" valid:"length(8|20)"`
	// Address fields
	Street_Address_1 string `json:"street_address_1,omitempty" valid:"length(2|80)"`
	Street_Address_2 string `json:"street_address_2,omitempty" valid:"length(2|80)"`
	City             string `json:"city,omitempty" valid:"length(2|80)"`
	Postal_Code      string `json:"postal_code,omitempty" valid:"numeric"`
	Suburb           string `json:"suburb,omitempty" valid:""`
	Province         string `json:"province,omitempty" valid:"in(Aceh|Bali|Banten|Bengkulu|Central Java|Central Kalimantan|Central Sulawesi|East Java|East Kalimantan|East Nusa Tenggara|Gorontalo|Jakarta Special Capital Region|Jambi|Lampung|Maluku|North Kalimantan|North Maluku|North Sulawesi|North Sumatra|Papua|Riau|Riau Islands|South Kalimantan|South Sulawesi|South Sumatra|Southeast Sulawesi|West Java|West Kalimantan|West Nusa Tenggara|West Papua|West Sulawesi|West Sumatra|Yogyakarta Special Region)"`
	// Relationships
	WorkTypes []db.WorkType `json:"work_types,omitempty" valid:""`
}

type UpdateVendor struct {
	// Required fields
	CompanyName string `json:"company_name,omitempty" valid:""`
	NPWP        string `json:"npwp,omitempty" valid:"length(12|20)"`
	NIB         string `json:"nib,omitempty" valid:"length(12|20)"`
	Email       string `json:"email,omitempty" valid:"email"`
	Phone       string `json:"phone,omitempty" valid:"length(8|20)"`
	// Address fields
	Street_Address_1 string `json:"street_address_1,omitempty" valid:"length(2|80)"`
	Street_Address_2 string `json:"street_address_2,omitempty" valid:"length(2|80)"`
	City             string `json:"city,omitempty" valid:"length(2|80)"`
	Postal_Code      string `json:"postal_code,omitempty" valid:"numeric"`
	Suburb           string `json:"suburb,omitempty" valid:""`
	Province         string `json:"province,omitempty" valid:"in(Aceh|Bali|Banten|Bengkulu|Central Java|Central Kalimantan|Central Sulawesi|East Java|East Kalimantan|East Nusa Tenggara|Gorontalo|Jakarta Special Capital Region|Jambi|Lampung|Maluku|North Kalimantan|North Maluku|North Sulawesi|North Sumatra|Papua|Riau|Riau Islands|South Kalimantan|South Sulawesi|South Sumatra|Southeast Sulawesi|West Java|West Kalimantan|West Nusa Tenggara|West Papua|West Sulawesi|West Sumatra|Yogyakarta Special Region)"`
	// Relationships
	WorkTypes []db.WorkType `json:"work_types,omitempty" valid:""`
}
