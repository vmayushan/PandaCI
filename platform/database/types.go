package database

import (
	"github.com/jmoiron/sqlx"
)

func NewPostgresDatabase() (*sqlx.DB, error) {
	return PostgreSQLConnection()
}
