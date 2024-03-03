package model

import "time"

type Order struct {
	ID         int
	Recipient  int
	ExpireDate time.Time
	IsReturned bool
	IsGiven    bool
}
