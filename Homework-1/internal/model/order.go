package model

import "time"

type Order struct {
	ID         int
	Recipient  int
	ExpireDate time.Time
}

type OrderInput struct {
	ID         int
	Recipient  int
	ExpireDate string
}
