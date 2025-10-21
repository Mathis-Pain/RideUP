package models

import (
	"database/sql"
	"time"
)

type Event struct {
	ID            int            `db:"id"`
	Title         string         `db:"title"`
	Description   sql.NullString `db:"description"`
	CreatedBy     int            `db:"created_by"`
	CreatorName   string
	CreatedAt     time.Time     `db:"created_at"`
	Latitude      float64       `db:"latitude"`
	Longitude     float64       `db:"longitude"`
	Address       string        `db:"address"`
	StartDatetime time.Time     `db:"start_datetime"`
	EndDatetime   sql.NullTime  `db:"end_datetime"` // 👈 si peut être NULL
	Participants  sql.NullInt64 `db:"participants"` // 👈 si peut être NULL
	Location      *SimpleAddress
	UserJoined    bool
}

func (e Event) FormattedStart() string {
	return e.StartDatetime.Format(" le 02/01/2006 à 15:04")
}
