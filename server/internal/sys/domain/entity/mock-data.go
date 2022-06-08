package entity

import "mayfly-go/pkg/model"

type MockData struct {
	model.Model
	Service     string `json:"service"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Data        string `json:"data"`
	RealUrl     string `json:"realUrl"`
}

func (a *MockData) TableName() string {
	return "t_mock_data"
}
