package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB holds databse connection pool
type DB struct {
	SQL *sql.DB
}

var dbConnection = &DB{}

const maxDBConnections = 10
const maxIdolConnections = 5
const maxDBLifeTime = 5 * time.Minute

// ConnectSQL create a database pool  for Postgres
func ConnectSQL(dns string) (*DB, error) {
	db, err := NewDB(dns)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(maxDBConnections)
	db.SetMaxIdleConns(maxIdolConnections)
	db.SetConnMaxLifetime(maxDBLifeTime)

	dbConnection.SQL = db

	err = testDB(db)
	if err != nil {
		return nil, err
	}

	return dbConnection, nil
}

// NewDB create a new database for application
func NewDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// testDB tries to ping the database
func testDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}
