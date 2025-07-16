package setup

import (
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	embed_migrations "github.com/DanArmor/vtuber-go"
	"github.com/DanArmor/vtuber-go/ent"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib" // Blank import for a driver
	"go.uber.org/zap"
)

// DatabaseSetup подключает к БД. !Также выполняет миграции.
func DatabaseSetup(driverName string, sqlURL string) (*ent.Client, error) {
	db, err := sql.Open(driverName, sqlURL)
	if err != nil {
		return nil, err
	}

	// For migrate
	zap.L().Info("Starting migrations")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	// Prepare migrations
	source, err := iofs.New(embed_migrations.Migrations, "ent/migrate/migrations")
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		panic(err)
	}

	// Apply migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	zap.L().Info("Migrations applied")

	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))
	return client, nil
}

// MustDatabaseSetup - это DatabaseSetup, но вызывает панику при ошибке.
func MustDatabaseSetup(driverName string, sqlURL string) *ent.Client {
	client, err := DatabaseSetup(driverName, sqlURL)
	if err != nil {
		panic(err)
	}
	return client
}
