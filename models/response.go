package models

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Count   int64       `json:"count"`
	Data    interface{} `json:"data"`
}
