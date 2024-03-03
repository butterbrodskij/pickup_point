package storage

import "time"

type OrderDTO struct {
	ID         int
	Recipient  int
	ExpireDate time.Time
	IsReturned bool
	IsGiven    bool
	GivenTime  bool
}
