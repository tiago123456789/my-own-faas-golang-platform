package types

type NewFunction struct {
	ID         int    `json:"id"`
	Runtime    string `json:"runtime" validate:"required,min=3"`
	Name       string `json:"name" validate:"required,min=3"`
	ModuleName string `json:"moduleName" validate:"required"`
	LambdaPath string `json:"path"`
	Cpu        string `json:"cpu"`
	Memory     string `json:"memory"`
}
