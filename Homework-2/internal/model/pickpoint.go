package model

type PickPoint struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	Address string `db:"address"`
	Contact string `db:"contacts"`
}
