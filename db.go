package kknotis

import "github.com/jackc/pgx"

const (
	insert = "INSERT INTO notification(id,userid,type,at,value) VALUES($1,$2,$3,$4,$5)"
)

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

// insertNotification to insert a notification to database.
func insertNotification(notis *Notification) error {
	_, err := dbPool.Exec(insert, notis.NotificationID, notis.Userid, notis.Type, notis.At, notis.Value)
	return err
}

// getNotifications to get the notifications.
// utime the unixtime, the notifications will be got after that time.
func getNotifications(userid, utime int32) ([]Notification, error) {
	s := "SELECT id,type,read,at,value FROM notification WHERE userid=$1 AND at>$2"
	rows, _ := dbPool.Query(s, userid, utime)

	var result []Notification
	for rows.Next() {
		var one Notification
		err := rows.Scan(&(one.NotificationID), &(one.Type), &(one.Read), &(one.At), &(one.Value))
		if err != nil {
			return result, err
		}
		one.Userid = userid
		result = append(result, one)
	}

	return result, rows.Err()
}

// getUnreadCount to get the count of unread notifications.
func getUnreadCount(toid int32) (int32, error) {
	s := "SELECT COUNT(1) FROM notification WHERE userid=$1 AND read=false"
	var count int64
	if err := dbPool.QueryRow(s, toid).Scan(&count); err != nil {
		return 0, err
	}
	return int32(count), nil
}

// readNotification to mark a notification as read.
func readNotification(notisid string) error {
	s := "UPDATE notification SET read=true WHERE id=$1 AND read=false"
	_, err := dbPool.Exec(s, notisid)
	return err
}

// readNotificationsByType to delete notifications with a single type.
func readNotificationsByType(userid int32, tp int16) error {
	s := "UPDATE notification SET read=true WHERE userid=$1 AND type=$2"
	_, err := dbPool.Exec(s, userid, tp)
	return err
}

// readAllNotification to mark all notifications as read.
func readAllNotifications(userid int32) error {
	s := "UPDATE notification SET read=true WHERE userid=$1"
	_, err := dbPool.Exec(s, userid)
	return err
}

// deleteNotificationsByType to delete notifications with a single type.
func deleteNotificationsByType(userid int32, tp int16) error {
	s := "DELETE FROM notification WHERE userid=$1 AND type=$2"
	_, err := dbPool.Exec(s, userid, tp)
	return err
}

// deleteNotification to delete a notification.
func deleteNotification(notisid string) error {
	s := "DELETE FROM notification WHERE id=$1"
	_, err := dbPool.Exec(s, notisid)
	return err
}

// deleteAllNotifications to delete a notification.
func deleteAllNotifications(userid int32) error {
	s := "DELETE FROM notification WHERE userid=$1"
	_, err := dbPool.Exec(s, userid)
	return err
}
