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

	// test read and delete with single type.
	insertNotifications(t)
	testReadAndDeleteNotificationsWithType(t)

	truncate(t)
}

func insertNotifications(t *testing.T) {
	var testMsg TestFormat
	testMsg.Message = "你好"

	var v []byte
	var err error

	var one Notification
	one.NotificationID = uuid.NewV1().String()
	one.Type = 0
	one.Userid = 2
	one.At = int32(time.Now().Unix())
	if v, err = json.Marshal(testMsg); err != nil {
		t.Error(err)
	} else {
		one.Value = v
	}
	insertNotification(&one)

	var two Notification
	two.NotificationID = uuid.NewV1().String()
	two.Type = 0
	two.Userid = 2
	two.At = int32(time.Now().Unix())
	if v, err = json.Marshal(testMsg); err != nil {
		t.Error(err)
	} else {
		two.Value = v
	}
	insertNotification(&one)
	insertNotification(&two)

	var three Notification
	three.NotificationID = uuid.NewV1().String()
	three.Type = 1
	three.Userid = 2
	three.At = int32(time.Now().Unix())
	if v, err = json.Marshal(testMsg); err != nil {
		t.Error(err)
	} else {
		three.Value = v
	}
	insertNotification(&one)
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

func testReadAndDeleteNotificationsWithType(t *testing.T) {
	if count, err := getUnreadCount(2); err != nil {
		t.Error(err)
	} else if count != 3 {
		t.Error("Unread count is wrong.")
	}

	if err := readNotificationsByType(2, 1); err != nil {
		t.Error(err)
	}

	if count, err := getUnreadCount(2); err != nil {
		t.Error(err)
	} else if count != 2 {
		t.Error("Unread count is wrong.")
	}

	if result, err := getNotifications(2, 0); err != nil {
		count := 0
		for _, one := range result {
			if one.Read {
				count++
			}
		}
		if count != 2 {
			t.Error("Read is wrong.")
		}
	}

	if err := deleteNotificationsByType(2, 0); err != nil {
		t.Error(err)
	}

	if result, err := getNotifications(2, 0); err != nil {
		if len(result) != 1 {
			t.Error("Delete is wrong.")
		}
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
		if err := json.Unmarshal([]byte(result[0].Value), &testMsg); err != nil {
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
	if err := readAllNotifications(2); err != nil {
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
	// read one notification with invalid id.
	if err := readNotification("notisid"); err == nil {
		t.Error("Should have error with invalid uuid.")
	}

	// read one empty notification.
	if err := readNotification(uuid.NewV1().String()); err != nil {
		t.Error("Should have error with invalid uuid.")
	}

	// read all notifications.
	if err := readAllNotifications(3); err != nil {
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
	if err := deleteAllNotifications(2); err != nil {
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
	// delete one notification with invalid id.
	if err := deleteNotification("notisid"); err == nil {
		t.Error("Should have error with invalid uuid.")
	}

	// delete one empty notification.
	if err := deleteNotification(uuid.NewV1().String()); err != nil {
		t.Error(err)
	}

	// delete all notifications.
	if err := deleteAllNotifications(3); err != nil {
		t.Error(err)
	}
}

func truncate(t *testing.T) {
	if _, err := dbPool.Exec("TRUNCATE notification"); err != nil {
		t.Error(err)
	}
}
