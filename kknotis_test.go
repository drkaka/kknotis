package kknotis

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx"
)

func TestMain(t *testing.T) {
	DBName := os.Getenv("dbname")
	DBHost := os.Getenv("dbhost")
	DBUser := os.Getenv("dbuser")
	DBPassword := os.Getenv("dbpassword")

	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     DBHost,
			User:     DBUser,
			Password: DBPassword,
			Database: DBName,
			Dial:     (&net.Dialer{KeepAlive: 5 * time.Minute, Timeout: 5 * time.Second}).Dial,
		},
		MaxConnections: 10,
	}

	var err error
	var pool *pgx.ConnPool
	if pool, err = pgx.NewConnPool(connPoolConfig); err != nil {
		t.Fatal(err)
	}

	if err = Use(pool); err != nil {
		t.Fatal(err)
	}
	testTableGeneration(t)

	testDBMethods(t)

	if _, err = dbPool.Exec("DROP TABLE notification"); err != nil {
		t.Error(err)
	}
}
