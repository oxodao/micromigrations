package micromigrations

type DatabaseDependentQuery interface {
	InsertMigration() string
	DeleteMigrationByTimestamp() string
	CreateMigrationTable() string
	DropMigrationTable() string
	FindMigrationTable() string
	ListMigrations() string
}

type SqliteQueries struct{}

func (q SqliteQueries) InsertMigration() string {
	return `INSERT INTO migrations (name, down_script, migration_ts, applied_at) VALUES ($1, $2, $3, $4);`
}

func (q SqliteQueries) DeleteMigrationByTimestamp() string {
	return `DELETE FROM migrations WHERE migration_ts = $1;`
}

func (q SqliteQueries) CreateMigrationTable() string {
	return `CREATE TABLE migrations(name VARCHAR(255), down_script TEXT, migration_ts BIGINT, applied_at TIMESTAMP);`
}

func (q SqliteQueries) DropMigrationTable() string {
	return `DROP TABLE migrations;`
}

func (q SqliteQueries) FindMigrationTable() string {
	return `SELECT name FROM sqlite_master WHERE type='table' AND tbl_name='migrations';`
}

func (q SqliteQueries) ListMigrations() string {
	return `
		SELECT name, down_script, migration_ts, applied_at
		FROM migrations
		ORDER BY migration_ts ASC;
	`
}

func NewSqliteQueriesAdapter() DatabaseDependentQuery {
	return SqliteQueries{}
}
