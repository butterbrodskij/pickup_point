package storage

import (
	"time"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type OrderDTO struct {
	ID           int64       `json:"id"`
	RecipientID  int64       `json:"recipient"`
	WeightGrams  int64       `json:"weight"`
	PriceKopecks int64       `json:"price"`
	Cover        model.Cover `json:"cover"`
	ExpireDate   time.Time   `json:"expires"`
	IsReturned   bool        `json:"is_returned"`
	IsGiven      bool        `json:"is_given"`
	GivenTime    time.Time   `json:"given"`
}
