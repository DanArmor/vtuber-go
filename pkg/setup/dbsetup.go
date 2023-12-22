package setup

import (
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/DanArmor/vtuber-go/ent"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func MustDatabaseSetup(driverName string, sqlUrl string) *ent.Client {
	db, err := sql.Open(driverName, sqlUrl)
	if err != nil {
		panic(err)
	}
	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))
	return client
}
