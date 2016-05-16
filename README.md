# kknotis
[![Build Status](https://travis-ci.org/drkaka/kknotis.svg)](https://travis-ci.org/drkaka/kknotis)
[![Coverage Status](https://codecov.io/github/drkaka/kknotis/coverage.svg?branch=master)](https://codecov.io/github/drkaka/kknotis?branch=master) 

The notification database module for golang.

## Database
It is using PostgreSQL as the database and will create a table:

```sql  
CREATE TABLE IF NOT EXISTS notification (
	id uuid primary key,
	userid integer,
    type smallint,
    read boolean DEFAULT false,
    at integer,
    value text
);
CREATE INDEX IF NOT EXISTS index_notification_userid ON notification (userid);
CREATE INDEX IF NOT EXISTS index_notification_at ON notification (at);
```

## Dependence

```Go
go get github.com/jackc/pgx
go get github.com/satori/go.uuid
```

## Usage 

####First need to use the module with the pgx pool passed in:
```Go
err := kknotis.Use(pool)
```

####Get notifications:
```Go
result, err := kknotis.GetNotifications(3, 0);
```

####Get unread notifications count:
```Go
count, err := GetUnreadCount(3);
```

####Read one notification:
```Go
err := kknotis.ReadNotification(notisid);
```

####Read all notifications:
```Go
err := kknotis.ReadAllNotifications(3);
```

####Read notifications of a type
```Go
err := kknotis.ReadNotificationsByType(3, 0);
```

####Delete one notification:
```Go
err := kknotis.DeleteNotification(notisid);
```

####Delete all notifications:
```Go
err := kknotis.DeleteAllNotifications(3);
```

####Delete notifications of a type
```Go
err := kknotis.DeleteNotificaitonByType(3, 0);
```
