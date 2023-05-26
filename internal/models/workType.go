package models

type CreateWorkType struct {
	Name string `json:"work_type_name" valid:"required,length(2|30)"`
}

type UpdateWorkType struct {
	Name string `json:"work_type_name" valid:"required,length(2|30)"`
}
