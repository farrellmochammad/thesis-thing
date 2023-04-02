package models

type FraudTransaction struct {
	TransactionID string  `json:"transactionid"`
	IndexFraud    float64 `json:"indexfraud"`
	FraudCategory bool    `json:"fraudcategory"`
}
