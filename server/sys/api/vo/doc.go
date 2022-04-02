package vo

import (
	"mayfly-go/base/model"
)

type DocListVO struct {
	model.Model
	Title    string `json:"title"`
	Category string `json:"category"`
	Content  string `json:"content"`
}

// 文档详细信息
type DocVO struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	Content  string `json:"content"`
}
