package model

// PagesModel 返回结果
type PagesModel struct {
	Total     int         `json:"total"`
	PageIndex int         `json:"page_index"`
	PageSize  int         `json:"page_size"`
	Data      interface{} `json:"data"`
}
