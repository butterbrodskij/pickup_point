package model

type PickPoint struct {
	ID      int64  `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Address string `db:"address" json:"address"`
	Contact string `db:"contacts" json:"contacts"`
}

type PickPointAdd struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contacts"`
}
