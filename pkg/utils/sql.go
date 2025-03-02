package utils

import (
	"fmt"

	"github.com/lib/pq"
	"modernc.org/sqlite"
)

func CheckConstraintError(err error, constraint string) bool {
	if err == nil {
		return false
	}

	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "23505" {
			return pqErr.Constraint == constraint
		}
	}

	if sqliteErr, ok := err.(*sqlite.Error); ok {
		if sqliteErr.Code() == 2067 {
			// TODO - we should probably do this in a more robust way
			return sqliteErr.Error() == fmt.Sprintf("constraint failed: UNIQUE constraint failed: %s (2067)", constraint)
		}
	}

	return false
}
