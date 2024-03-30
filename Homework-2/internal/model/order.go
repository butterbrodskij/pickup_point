package model

import "time"

type Order struct {
	ID           int64
	RecipientID  int64
	WeightGrams  int64
	PriceKopecks int64
	Cover        Cover
	ExpireDate   time.Time
}

type OrderInput struct {
	ID           int64
	RecipientID  int64
	WeightGrams  int64
	PriceKopecks int64
	Cover        Cover
	ExpireDate   string
}
