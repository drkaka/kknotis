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
	defer pool.Close()

	if err = Use(pool); err != nil {
		t.Fatal(err)
	}
	testTableGeneration(t)

	testDBMethods(t)

	testAddNotifications(t)
	testGetAndReadNotification(t)
	testReadAndDeleteNotifications(t)

	if _, err = dbPool.Exec("DROP TABLE notification"); err != nil {
		t.Error(err)
	}
}

func testAddNotifications(t *testing.T) {
	var err error
	if err = AddNotification(3, 0, nil); err != nil {
		t.Error(err)
	}

	if err = AddNotification(3, 1, nil); err != nil {
		t.Error(err)
	}

	if err = AddNotification(3, 2, nil); err != nil {
		t.Error(err)
	}
}

func testGetAndReadNotification(t *testing.T) {
	var err error
	var result []Notification
	if result, err = GetNotifications(3, 0); err != nil {
		t.Error(err)
	}

	if len(result) != 3 {
		t.Error("result length is wrong.")
	}

	if err = ReadNotification(result[0].NotificationID); err != nil {
		t.Error(err)
	}

	if oneResult, err := GetNotifications(3, 0); err != nil {
		t.Error(err)
	} else {
		count := 0
		for _, one := range oneResult {
			if !one.Read {
				count++
			}
		}
		if count != 2 {
			t.Error("result is wrong.")
		}
	}

	if err = ReadAllNotifications(3); err != nil {
		t.Error(err)
	}

	if oneResult, err := GetNotifications(3, 0); err != nil {
		t.Error(err)
	} else {
		count := 0
		for _, one := range oneResult {
			if !one.Read {
				count++
			}
		}
		if count != 0 {
			t.Error("result is wrong.")
		}
	}
}

func testReadAndDeleteNotifications(t *testing.T) {
	var err error
	var result []Notification
	if result, err = GetNotifications(3, 0); err != nil {
		t.Error(err)
	}

	if err = DeleteNotification(result[0].NotificationID); err != nil {
		t.Error(err)
	}

	if oneResult, err := GetNotifications(3, 0); err != nil {
		t.Error(err)
	} else if len(oneResult) != 2 {
		t.Error("result is wrong.")
	}

	if err = DeleteAllNotifications(3); err != nil {
		t.Error(err)
	}

	if oneResult, err := GetNotifications(3, 0); err != nil {
		t.Error(err)
	} else if len(oneResult) != 0 {
		t.Error("result is wrong.")
	}
}
