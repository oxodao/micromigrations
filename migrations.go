package micromigrations

import (
	"database/sql"
	"errors"
	"time"
)

func doApplyMigrations(
	logger Logger,
	dbmqQueries DatabaseDependentQuery,
	db *sql.DB,
	migrations []Migration,
	up bool,
) error {
	for _, migration := range migrations {
		way := "up"
		sql := migration.Up
		if !up {
			sql = migration.Down
			way = "down"
		}

		logger.Info(`- Migrating %v (%v)...`, migration.Name, way)
		_, err := db.Exec(sql)
		if err != nil {
			return err
		}

		if up {
			_, err = db.Exec(
				dbmqQueries.InsertMigration(),
				migration.Name,
				migration.Down,
				migration.MigrationTS,
				time.Now().Unix(),
			)

			if err != nil {
				return err
			}
		} else {
			_, err = db.Exec(
				dbmqQueries.DeleteMigrationByTimestamp(),
				migration.MigrationTS,
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func MustApplyMigrations(
	logger Logger,
	dbmqQueries DatabaseDependentQuery,
	db *sql.DB,
	migrations []Migration,
	allowDowngrades bool,
) {
	err := ApplyMigrations(logger, dbmqQueries, db, migrations, allowDowngrades)

	if err != nil {
		panic(err)
	}
}

func ApplyMigrations(
	logger Logger,
	dbmqQueries DatabaseDependentQuery,
	db *sql.DB,
	migrations []Migration,
	allowDowngrades bool,
) error {
	migrations = append(
		[]Migration{
			{
				Name:        "Migration table",
				Up:          dbmqQueries.CreateMigrationTable(),
				Down:        dbmqQueries.DropMigrationTable(),
				MigrationTS: 1725887251, // 2024-09-09 @ 15:07
			},
		},
		migrations...,
	)

	logger.Info("Checking for available migrations...")

	// Checking if the table "migrations" exists
	row := db.QueryRow(dbmqQueries.FindMigrationTable())

	if row.Err() != nil {
		return row.Err()
	}

	var tableName string
	err := row.Scan(&tableName)
	if err != nil {
		// If the table does not exists at all (brand new DB)
		// We pass all the migrations
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		logger.Warn("This database do not seems to have been initialized. Applying all migrations!")

		return doApplyMigrations(logger, dbmqQueries, db, migrations, true)
	}

	//#region Fetching currently applied migrations
	rows, err := db.Query(dbmqQueries.ListMigrations())
	if err != nil {
		return err
	}

	appliedMigrations := []Migration{}
	for rows.Next() {
		m := Migration{}

		err := rows.Scan(
			&m.Name,
			&m.Down,
			&m.MigrationTS,
			&m.AppliedAtTS,
		)

		if err != nil {
			return err
		}

		appliedMigrations = append(appliedMigrations, m)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	rows.Close()
	//#endregion

	amtAppliedMigrations := len(appliedMigrations)
	latestAvailableMigration := migrations[len(migrations)-1].MigrationTS
	latestAppliedMigration := int64(0)
	if amtAppliedMigrations > 0 {
		latestAppliedMigration = appliedMigrations[amtAppliedMigrations-1].MigrationTS
	}

	if amtAppliedMigrations > 0 && latestAppliedMigration > latestAvailableMigration {
		logger.Warn("This database was made with a newer version of the software")

		if allowDowngrades {
			tmpMigrations := []Migration{}

			for _, m := range appliedMigrations {
				if m.MigrationTS > latestAvailableMigration {
					tmpMigrations = append(tmpMigrations, m)
				}
			}

			reverse(tmpMigrations)

			// return doApplyMigrations(logger, db, false)
			return errors.New("downgrade db is not implemented yet! (Well it is but it doesn't work @TODO)")
		} else {
			return errors.New("won't downgrade the database automatically as it could result in data loss.")
		}
	}

	// Filter already applied migrations and put the other in the migrations array
	tmpMigrations := []Migration{}
	latestMigration := migrations[0].MigrationTS
	if amtAppliedMigrations > 0 {
		latestMigration = latestAppliedMigration
	}

	for _, m := range migrations {
		if m.MigrationTS > latestMigration {
			tmpMigrations = append(tmpMigrations, m)
		}
	}

	if len(tmpMigrations) == 0 {
		logger.Info("The database is up to date!")
	}

	return doApplyMigrations(logger, dbmqQueries, db, tmpMigrations, true)
}
