package dbrepo

import (
	"database/sql"
	"github.com/ismail118/bookings-app/internal/config"
	"github.com/ismail118/bookings-app/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(db *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  db,
	}
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
