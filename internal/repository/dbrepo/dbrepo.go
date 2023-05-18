package dbrepo

import (
	"database/sql"

	"github.com/shahinm95/bookings/internal/config"
	"github.com/shahinm95/bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB *sql.DB
}


func NewPostgresRepo (conn *sql.DB , app *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App : app, 
		DB : conn,
	}
}