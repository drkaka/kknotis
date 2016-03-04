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
    value JSONB
);
CREATE INDEX IF NOT EXISTS index_notification_userid ON notification (userid);
CREATE INDEX IF NOT EXISTS index_notification_at ON notification (at);
```

## Dependence

```Go
go get github.com/jackc/pgx
go get github.com/satori/go.uuid
```