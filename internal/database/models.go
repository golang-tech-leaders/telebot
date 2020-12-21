package database

import (
	"database/sql"
)

var (
	RU = 1
	EN = 2
)

type TelebotMessage struct {
	ID      sql.NullInt32  `json:"id"`
	Message sql.NullString `json:"message"`
}
