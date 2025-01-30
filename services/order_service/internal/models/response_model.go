package models

type NewResponseModel struct {
	OrderId string `json:"order_id"`
}

type ResponseToBrokerModel struct {
	UserId  int64  `json:"user_id"`
	OrderId string `json:"order_id"`
	ChatId  int64  `json:"chat_id"`
}
