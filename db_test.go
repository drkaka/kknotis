package kknotis

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/jackc/pgx"
	"github.com/satori/go.uuid"
)

type TestFormat struct {
	Message string `json:"message"`
}

func testTableGeneration(t *testing.T) {
	var dbname pgx.NullString
	if err := dbPool.QueryRow("SELECT 'public.notification'::regclass;").Scan(&dbname); err != nil {
		t.Fatal(err)
	}

	if dbname.String != "notification" {
		t.Fatal("dbname is not correct.")
	}
}

// testDBMethods to test the database operation methods.
func testDBMethods(t *testing.T) {
	insertNotifications(t)
	testGetNotifications(t)

	truncate(t)
}

func insertNotifications(t *testing.T) {
	var testMsg TestFormat
	testMsg.Message = "你好"

	var one Notification
	one.NotificationID = uuid.NewV1().String()
	one.Type = 0
	one.Userid = 2
	one.At = int32(time.Now().Unix())
	one.Value, _ = json.Marshal(testMsg)
	insertNotification(&one)
}

func testGetNotifications(t *testing.T) {
	result, err := getNotifications(2, 0)
	if err != nil {
		t.Error(err)
	}

	if len(result) != 1 {
		t.Error("result is wrong")
	} else {
		var testMsg TestFormat
		if err := json.Unmarshal(result[0].Value, &testMsg); err != nil {
			t.Error(err)
		}
		if testMsg.Message != "你好" {
			t.Error("result is wrong.")
		}
	}
}

func truncate(t *testing.T) {
	if _, err := dbPool.Exec("TRUNCATE notification"); err != nil {
		t.Error(err)
	}
}
