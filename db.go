package kknotis

import "github.com/jackc/pgx"

// dbPool the pgx database pool.
var dbPool *pgx.ConnPool

// prepareDB to prepare the database.
func prepareDB() error {
	s := `CREATE TABLE IF NOT EXISTS notification (
	id uuid primary key,
	userid integer,
    type smallint,
    read boolean DEFAULT false,
    at integer,
    value JSONB);
    CREATE INDEX IF NOT EXISTS index_notification_userid ON notification (userid);
    CREATE INDEX IF NOT EXISTS index_notification_at ON notification (at);`

	_, err := dbPool.Exec(s)
	return err
}
