package model

type PickPoint struct {
	ID      int64  `db:"id" json:"id,omitempty"`
	Name    string `db:"name" json:"name"`
	Address string `db:"address" json:"address"`
	Contact string `db:"contacts" json:"contacts"`
}
