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
	testGetNoNotifications(t)
	insertNotifications(t)
	testGetNotificationsAndRead(t)
	testReadEmpty(t)
	testDeleteNotifications(t)
	testDeleteEmpty(t)

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

	var two Notification
	two.NotificationID = uuid.NewV1().String()
	two.Type = 0
	two.Userid = 2
	two.At = int32(time.Now().Unix())
	two.Value, _ = json.Marshal(testMsg)
	insertNotification(&two)

	var three Notification
	three.NotificationID = uuid.NewV1().String()
	three.Type = 0
	three.Userid = 2
	three.At = int32(time.Now().Unix())
	three.Value, _ = json.Marshal(testMsg)
	insertNotification(&three)
}

func testGetNoNotifications(t *testing.T) {
	result, err := getNotifications(2, 0)
	if err != nil {
		t.Error(err)
	}

	if len(result) != 0 {
		t.Error("result is wrong")
	}
}

func testGetNotificationsAndRead(t *testing.T) {
	result, err := getNotifications(2, 0)
	if err != nil {
		t.Error(err)
	}

	if len(result) != 3 {
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

	// read one notification.
	notisid := result[0].NotificationID
	if err := readNotification(notisid); err != nil {
		t.Error(err)
	}

	if oneReadResult, err := getNotifications(2, 0); err != nil {
		t.Error(err)
	} else {
		for _, one := range oneReadResult {
			if one.NotificationID == notisid {
				if !one.Read {
					t.Error("result is wrong.")
				}
			} else {
				if one.Read {
					t.Error("result is wrong.")
				}
			}
		}
	}

	// read all notifications.
	if err := readAllNotification(2); err != nil {
		t.Error(err)
	}

	if oneReadResult, err := getNotifications(2, 0); err != nil {
		t.Error(err)
	} else {
		for _, one := range oneReadResult {
			if !one.Read {
				t.Error("result is wrong.")
			}
		}
	}
}

// testReadEmpty to read none existed notifications.
func testReadEmpty(t *testing.T) {
	// read one notification.
	if err := readNotification("notisid"); err != nil {
		t.Error(err)
	}

	// read all notifications.
	if err := readAllNotification(3); err != nil {
		t.Error(err)
	}
}

func testDeleteNotifications(t *testing.T) {
	result, err := getNotifications(2, 0)
	if err != nil {
		t.Error(err)
	}

	if len(result) != 3 {
		t.Error("result is wrong")
	}

	// delete one notification.
	notisid := result[0].NotificationID
	if err := deleteNotification(notisid); err != nil {
		t.Error(err)
	}

	if oneReadResult, err := getNotifications(2, 0); err != nil {
		t.Error(err)
	} else {
		for _, one := range oneReadResult {
			if one.NotificationID == notisid {
				t.Error("notification should be deleted.")
			}
		}
	}

	// delete all notifications.
	if err := deleteAllNotification(2); err != nil {
		t.Error(err)
	}

	if oneReadResult, err := getNotifications(2, 0); err != nil {
		t.Error(err)
	} else {
		if len(oneReadResult) > 0 {
			t.Error("There should be no result.")
		}
	}
}

// testDeleteEmpty to delete none existed notifications.
func testDeleteEmpty(t *testing.T) {
	// delete one notification.
	if err := deleteNotification("notisid"); err != nil {
		t.Error(err)
	}

	// delete all notifications.
	if err := deleteAllNotification(3); err != nil {
		t.Error(err)
	}
}

func truncate(t *testing.T) {
	if _, err := dbPool.Exec("TRUNCATE notification"); err != nil {
		t.Error(err)
	}
}
