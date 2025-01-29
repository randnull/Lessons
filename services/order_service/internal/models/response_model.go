package models

type NewResponseModel struct {
	OrderId string `json:"order_id"`
}

type ResponseToBrokerModel struct {
	OrderId string `json:"order_id"`
}
