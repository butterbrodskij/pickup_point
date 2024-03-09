package storage

import "time"

type OrderDTO struct {
	ID          int64
	RecipientID int64
	ExpireDate  time.Time
	IsReturned  bool
	IsGiven     bool
	GivenTime   time.Time
}
