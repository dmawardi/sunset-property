package models

import "github.com/dmawardi/Go-Template/internal/db"

type CreateProperty struct {
	Postcode         int          `json:"postcode" valid:"length(5|6),numeric"`
	Property_Name    string       `json:"property_name" valid:"length(6|25),required"`
	Suburb           string       `json:"suburb" valid:"length(4|25)"`
	City             string       `json:"city" valid:"length(4|25)"`
	Street_Address_1 string       `json:"street_address_1" valid:"length(6|32),required"`
	Street_Address_2 string       `json:"street_address_2" valid:"length(6|32)"`
	Bedrooms         float32      `json:"bedrooms" valid:"float"`
	Bathrooms        float32      `json:"bathrooms" valid:"float"`
	Land_Area        float64      `json:"land_area" valid:"float"`
	Land_Metric      string       `json:"land_metric" valid:"length(2|32)"`
	Description      string       `json:"description" valid:"length(5|250)"`
	Notes            string       `json:"notes" valid:"length(5|250)"`
	Features         []db.Feature `json:"features" valid:""`
}

type UpdateProperty struct {
	Postcode         int          `json:"postcode,omitempty" valid:"length(5|6),number"`
	Property_Name    string       `json:"property_name,omitempty" valid:"length(6|25)"`
	Suburb           string       `json:"suburb,omitempty" valid:"length(4|25)"`
	City             string       `json:"city,omitempty" valid:"length(4|25)"`
	Street_Address_1 string       `json:"street_address_1,omitempty" valid:"length(6|32)"`
	Street_Address_2 string       `json:"street_address_2,omitempty" valid:"length(6|32)"`
	Bedrooms         float32      `json:"bedrooms,omitempty" valid:"number"`
	Bathrooms        float32      `json:"bathrooms,omitempty" valid:"number"`
	Land_Area        float64      `json:"land_area,omitempty" valid:"number"`
	Land_Metric      string       `json:"land_metric,omitempty" valid:"length(2|32),number"`
	Description      string       `json:"description,omitempty" valid:"length(5|250)"`
	Notes            string       `json:"notes,omitempty" valid:"length(5|250)"`
	Features         []db.Feature `json:"features,omitempty" valid:""`
}
