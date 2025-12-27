package model

type PaymentRequest struct {
	Amount     int64  `json:"amount"`
	ExternalId string `json:"external_id"`
}
type PaymentResponse struct {
	Data struct {
		Id         string `json:"id"`
		ExternalId string `json:"external_id"`
		Status     string `json:"status"`
	}
	Message string `json:"message"`
}
