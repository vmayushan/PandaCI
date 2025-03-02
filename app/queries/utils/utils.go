package queries_utils

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func Rollback(completed *bool, tx *sqlx.Tx) {
	if !*completed {
		if err := tx.Rollback(); err != nil {
			log.Error().Err(err).Msg("Rollback failed")
		}
	}
}
