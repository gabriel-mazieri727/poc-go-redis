package main

import (
	"database/sql"
	"fmt"
)

// PostgreSQL

func (p *PostgreSQL) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(`
		host=%s
		port=%s
		user=%s
		password=%s
		dbname=%s
		sslmode=disable
	`,
		p.Host,
		p.Port,
		p.User,
		p.Password,
		p.DB,
	))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (p *PostgreSQL) Close(db *sql.DB) error {
	return db.Close()
}

func (p *PostgreSQL) CreateUsersTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			email VARCHAR(255)
		);
	`)
	if err != nil {
		return err
	}
	return nil
}
