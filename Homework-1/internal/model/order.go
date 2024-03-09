package model

import "time"

type Order struct {
	ID          int64
	RecipientID int64
	ExpireDate  time.Time
}

type OrderInput struct {
	ID          int64
	RecipientID int64
	ExpireDate  string
}
