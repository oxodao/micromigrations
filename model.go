package micromigrations

// @TODO: migration_ts should be PK

type Migration struct {
	Name        string `db:"name"`
	Up          string `db:"-"`
	Down        string `db:"down_script"`
	MigrationTS int64  `db:"migration_ts"`
	AppliedAtTS int64  `db:"applied_at_ts"`
}
