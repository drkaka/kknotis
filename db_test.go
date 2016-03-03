package kknotis

import (
	"testing"

	"github.com/jackc/pgx"
)

func testTableGeneration(t *testing.T) {
	var dbname pgx.NullString
	if err := dbPool.QueryRow("SELECT 'public.notification'::regclass;").Scan(&dbname); err != nil {
		t.Fatal(err)
	}

	if dbname.String != "notification" {
		t.Fatal("dbname is not correct.")
	}
}


