package model

import "time"

type Order struct {
	ID          int64
	RecipientID int64
	Weight      int
	Price       int
	Cover       string
	ExpireDate  time.Time
}

type OrderInput struct {
	ID          int64
	RecipientID int64
	Weight      int
	Price       int
	Cover       string
	ExpireDate  string
}
