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
	return `INSERT INTO migrations (name, down_script, migration_ts, applied_at_ts) VALUES ($1, $2, $3, $4);`
}

func (q SqliteQueries) DeleteMigrationByTimestamp() string {
	return `DELETE FROM migrations WHERE migration_ts = $1;`
}

func (q SqliteQueries) CreateMigrationTable() string {
	return `CREATE TABLE migrations(name VARCHAR(255), down_script TEXT, migration_ts BIGINT, applied_at_ts BIGINT);`
}

func (q SqliteQueries) DropMigrationTable() string {
	return `DROP TABLE migrations;`
}

func (q SqliteQueries) FindMigrationTable() string {
	return `SELECT name FROM sqlite_master WHERE type='table' AND tbl_name='migrations';`
}

func (q SqliteQueries) ListMigrations() string {
	return `
		SELECT name, down_script, migration_ts, applied_at_ts
		FROM migrations
		ORDER BY migration_ts ASC;
	`
}

func NewSqliteQueriesAdapter() DatabaseDependentQuery {
	return SqliteQueries{}
}

type MysqlQueries struct{}

func (q MysqlQueries) InsertMigration() string {
	return `INSERT INTO migrations (name, down_script, migration_ts, applied_at_ts) VALUES (?, ?, ?, ?);`
}

func (q MysqlQueries) DeleteMigrationByTimestamp() string {
	return `DELETE FROM migrations WHERE migration_ts = ?;`
}

func (q MysqlQueries) CreateMigrationTable() string {
	return `CREATE TABLE migrations(name VARCHAR(255), down_script TEXT, migration_ts BIGINT, applied_at_ts BIGINT);`
}

func (q MysqlQueries) DropMigrationTable() string {
	return `DROP TABLE migrations;`
}

func (q MysqlQueries) FindMigrationTable() string {
	return `
		SELECT TABLE_NAME
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'migrations';
	`
}

func (q MysqlQueries) ListMigrations() string {
	return `
		SELECT name, down_script, migration_ts, applied_at_ts
		FROM migrations
		ORDER BY migration_ts ASC;
	`
}

func NewMysqlQueriesAdapter() DatabaseDependentQuery {
	return MysqlQueries{}
}
