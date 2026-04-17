package database

import "database/sql"

// HandleQueryError signals the health monitor to switch to fast polling
// when a real connectivity error occurs, as opposed to expected errors
// like sql.ErrNoRows which don't indicate DB unavailability.
func HandleQueryError(err error) {
	if err == nil || err == sql.ErrNoRows {
		return
	}
	MarkDegraded()
}
