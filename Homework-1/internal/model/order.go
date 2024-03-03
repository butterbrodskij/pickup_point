package model

import "time"

type OrderInput struct {
	ID         int
	Recipient  int
	ExpireDate time.Time
}
