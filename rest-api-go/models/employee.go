package models

type Employee struct {
	EmployeeId string `json:"employee_id,omitempty" bson:"employee_id,omitempty"`
	Name       string `json:"name,omitempty" bson:"name,omitempty"`
	Department string `json:"department,omitempty" bson:"department,omitempty"`
}
