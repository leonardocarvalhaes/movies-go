package database

import "database/sql"

type Connector interface {
	openDB(dsn string) (*sql.DB, error)
	connectToDB() (*sql.DB, error)
}
