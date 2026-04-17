//go:build !production

package database

// SetHealthForTest directly sets the DB health flag.
// Only compiled in non-production builds - use in tests only.
func SetHealthForTest(healthy bool) {
	if healthy {
		dbHealthy = 1
	} else {
		dbHealthy = 0
	}
}
