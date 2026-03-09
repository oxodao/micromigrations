# MicroMigrations

## What's this

A simple bundle to make database migrations easy.

This library is tailored at my usage so it might not fit your needs but if it does, I'm happy that it solves this issue for you.

Note that there will not be much updates, not because its unmaintained, but because it is feature-complete (Except the downgrade one but I'm too lazy to implement it properly).

Note that this package currently only supports SQLite and mysql/mariadb but it would be quite easy to support more.

You can implement your own `DatabaseDependentQuery` struct, the only DBMS-dependent query should currently be the `FindMigrationTable` method. Do not hesitate to make a pull request if you implement one.

## Usage

```go
mmLogger := micromigrations.NewGenericLogger()
db := sql.Connect(....) // Or get the sql.DB instance from your database library
allowDowngrades := false // False, it's not implemented yet anyway
queriesAdapter := micromigrations.NewSqliteQueriesAdapter() // Use your DBMS adapter to have the correct queries

// Finally implements your migrations
migrations := []micromigrations.Migration{
    {
        Name: "The migration name",
        Up: `
            CREATE TABLE my_table(id INTEGER NOT NULL, name TEXT);
            CREATE TABLE my_table_2(id INTEGER NOT NULL, name TEXT);
        `,
        Down: `
            DROP TABLE my_table_2;
            DROP TABLE my_table;
        `,
        MigrationTS: 1725887251, // 2024-09-09 @ 15:07
    },
}

// And at the start of your software, right after connecting to the database, run the migrations
micromigrations.MustApplyMigrations(
    mmLogger,
    queriesAdapter,
    db,
    migrations,
    allowDowngrades,
)

```

## Logger

By default, you have two Logger adapter which are NoopLogger (Does nothing, useful for tests) and GenericLogger (Simply print stuff, uses the golang's default log package).

I'm using [Uber's zap](https://github.com/uber-go/zap) in my projects, you can find the adapter at [MicroMigrationsZapAdapter](https://github.com/oxodao/micromigrations-zap-adapter).

## Roadmap

At some point, I would like to do the following:
- [ ] Postgres support
- [ ] Downgrade migrations

## License

This project is licensed under the LGPLv3 license.

You can find the full text in the file [LICENSE.md](LICENSE.md)