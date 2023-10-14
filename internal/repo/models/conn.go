package models

import "time"

type ConnLog struct {
	UserID    int64     `db:"user_id"`
	Addr      string    `db:"addr"`
	CreatedAt time.Time `db:"created_at"`
}
