package persistance

import (
	"database/sql"
	"fmt"
	"task-manager/src/internal/config"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {

	config, err := config.LoadConfig()

	if err != nil {
		return nil, err
	}

	dataBaseUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME, config.DB_SSLMODE)

	fmt.Println("DATABSE URL:", dataBaseUrl)
	db, err := sql.Open("postgres", dataBaseUrl)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
