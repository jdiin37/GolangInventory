package model

type Result struct {
	Code    int         `json:"code" example:"000"`
	Message string      `json:"message" example:"回傳訊息"`
	Data    interface{} `json:"data" `
}
