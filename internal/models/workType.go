package models

type CreateWorkType struct {
	Name string `json:"work_type_name" valid:"required,length(2|255)"`
}

type UpdateWorkType struct {
	Name string `json:"work_type_name" valid:"required,length(2|255)"`
}
