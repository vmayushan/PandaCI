package utilsCSV

import (
	"bytes"
	"encoding/csv"

	"github.com/rs/zerolog/log"
)

func FormatCSVRow(row []string) (string, error) {
	var buf bytes.Buffer

	w := csv.NewWriter(&buf)
	if err := w.Write(row); err != nil {
		log.Error().Err(err).Msg("writing csv row")
		return "", err
	}
	w.Flush()

	if err := w.Error(); err != nil {
		log.Error().Err(err).Msg("writing csv row")
		return "", err
	}

	return buf.String(), nil
}
