package entity

import (
	"mayfly-go/base/model"
)

type Doc struct {
	model.Model

	Title    string `json:"title"`
	Category string `json:"category"`
	Content  string `json:"content"`
}

func (a *Doc) TableName() string {
	return "t_doc"
}
