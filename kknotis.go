package kknotis

import "github.com/jackc/pgx"

// Use the pool to do further operations.
func Use(pool *pgx.ConnPool) error {
	dbPool = pool
	return prepareDB()
}
