package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/pandaci-com/pandaci/pkg/utils/env"

	_ "github.com/lib/pq"
)

// PostgreSQLConnection func for connection to PostgreSQL database.
func PostgreSQLConnection() (*sqlx.DB, error) {
	// Define database connection settings.
	// maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	// maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	// maxLifetimeConn, _ := time.ParseDuration(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	dsn, err := env.GetPostgresDSN()
	if err != nil {
		return nil, err
	}

	// Define database connection for PostgreSQL.
	db, err := sqlx.Connect("postgres", *dsn)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, connection_url: %s, %w", *dsn, err)
	}

	// Set database connection settings:
	// 	- SetMaxOpenConns: the default is 0 (unlimited)
	// 	- SetMaxIdleConns: defaultMaxIdleConns = 2
	// 	- SetConnMaxLifetime: 0, connections are reused forever
	// db.SetMaxOpenConns(maxConn)
	// db.SetMaxIdleConns(maxIdleConn)
	// db.SetConnMaxLifetime(maxLifetimeConn)

	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer db.Close() // close database connection
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	log.Info().Msg("Connected to postgres database")

	return db.Unsafe(), nil
}
