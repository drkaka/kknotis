package kknotis

import (
	"time"

	"github.com/jackc/pgx"
	"github.com/satori/go.uuid"
)

// Notification for private message.
type Notification struct {
	NotificationID string `json:"notificationid"`
	Userid         int32  `json:"userid,omitempty"`
	Type           int16  `json:"type"`
	Read           bool   `json:"read"`
	At             int32  `json:"at"`
	// Value can be JSON format
	Value string `json:"value"`
}

// Use the pool to do further operations.
func Use(pool *pgx.ConnPool) error {
	dbPool = pool
	return prepareDB()
}

// GetNotifications to get the notifications.
// utime the unixtime, the notifications will be got after that time.
func GetNotifications(userid, utime int32) ([]Notification, error) {
	return getNotifications(userid, utime)
}

// AddNotification to add a notification in database.
func AddNotification(userid int32, tp int16, value string) error {
	var notis Notification
	notis.NotificationID = uuid.NewV1().String()
	notis.Userid = userid
	notis.Type = tp
	notis.At = int32(time.Now().Unix())
	notis.Value = value
	return insertNotification(&notis)
}

// ReadNotification to read a notification.
func ReadNotification(notisid string) error {
	return readNotification(notisid)
}

// ReadAllNotification to read all notifications.
func ReadAllNotifications(userid int32) error {
	return readAllNotifications(userid)
}

// DeleteNotification to delete a notification.
func DeleteNotification(notisid string) error {
	return deleteNotification(notisid)
}

// DeleteAllNotification to delete all notifications.
func DeleteAllNotifications(userid int32) error {
	return deleteAllNotifications(userid)
}
