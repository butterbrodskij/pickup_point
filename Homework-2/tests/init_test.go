package tests

import (
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	postgres_test "gitlab.ozon.dev/mer_marat/homework/tests/postgres"
)

var (
	db *postgres_test.TDB
)

func init() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	db = postgres_test.NewTDB(cfg)
}
