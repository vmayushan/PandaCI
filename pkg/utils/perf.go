package utils

import (
	"time"

	"github.com/rs/zerolog/log"
)

// Place at the beginning of a function and defer the returned function at the end of the function
func MeasureTime(start time.Time, name string) {
	elapsed := time.Since(start)

	log.Debug().Int("elapsed", int(elapsed.Milliseconds())).Str("name", name).Msg("Function execution time")
}
