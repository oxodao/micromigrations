package micromigrations

import "time"

// @TODO: migration_ts should be PK

type Migration struct {
	Name        string    `db:"name"`
	Up          string    `db:"-"`
	Down        string    `db:"down_script"`
	MigrationTS int       `db:"migration_ts"`
	AppliedAt   time.Time `db:"applied_at"`
}
