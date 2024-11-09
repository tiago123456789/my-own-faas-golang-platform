package models

import "time"

type Function struct {
	ID            int       `json:"id"`
	LambdaName    string    `json:"name"`
	Runtime       string    `json:"runtime"`
	LambdaPath    string    `json:"path"`
	Cpu           string    `json:"cpu"`
	Memory        string    `json:"memory"`
	BuildProgress string    `json:"buildProgress"`
	LastExecution time.Time `json:"last_execution"`
	Interval      int       `json: "interval"`
	Trigger       string    `json:"trigger"`
}
