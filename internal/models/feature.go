package models

type CreateFeature struct {
	Feature_Name string `json:"feature_name" valid:"length(4|25),required"`
}

type UpdateFeature struct {
	Feature_Name string `json:"feature_name,omitempty" valid:"length(4|25),required"`
}
