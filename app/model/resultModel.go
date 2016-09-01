package model

// ResultModel 返回结果
type ResultModel struct {
	ErrorCode int         `json:"error_code"`
	ErrorMsg  string      `json:"error_msg"`
	Extra     interface{} `json:"extra"`
}
