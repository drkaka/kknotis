package kknotis

import "github.com/jackc/pgx"

// Notification for private message.
type Notification struct {
	NotificationID string `json:"notificationid"`
	Userid         int32  `json:"userid,omitempty"`
	Type           int16  `json:"type"`
	Read           bool   `json:"read"`
	At             int32  `json:"at"`
	// Value in JSON format
	Value []byte `json:"value"`
}

// Use the pool to do further operations.
func Use(pool *pgx.ConnPool) error {
	dbPool = pool
	return prepareDB()
}
