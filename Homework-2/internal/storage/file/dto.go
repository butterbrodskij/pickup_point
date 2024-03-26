package storage

import "time"

type OrderDTO struct {
	ID          int64     `json:"id"`
	RecipientID int64     `json:"recipient"`
	Weight      int       `json:"weight"`
	Price       int       `json:"price"`
	Cover       string    `json:"cover"`
	ExpireDate  time.Time `json:"expires"`
	IsReturned  bool      `json:"is_returned"`
	IsGiven     bool      `json:"is_given"`
	GivenTime   time.Time `json:"given"`
}
