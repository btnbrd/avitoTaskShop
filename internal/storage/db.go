package storage

import (
	"database/sql"
	"fmt"
	"github.com/btnbrd/avitoshop/internal/storage/dbconfig"
	_ "github.com/lib/pq"
	"log"
)

//import "github.com/jmoiron/sqlx"
//
//type DB struct {
//	db *sqlx.DB
//}

func NewDBConnection() (*sql.DB, error) {
	conf, err := dbconfig.NewDBConfig()
	if err != nil {
		return nil, err
	}
	fmt.Printf("config %+v\n", conf)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Name)

	fmt.Printf("%+v\n", psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	//defer db.Close()
	//
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = InitDB(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitDB(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			coins INTEGER NOT NULL DEFAULT 1000,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS coin_transfers (
			id SERIAL PRIMARY KEY,
			from_user_id INTEGER NOT NULL REFERENCES users(id),
			to_user_id INTEGER NOT NULL REFERENCES users(id),
			amount INTEGER NOT NULL CHECK (amount > 0),
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS purchases (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id),
			item VARCHAR(255) NOT NULL,
			price INTEGER NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);`,
	}
	log.Println("tables created")

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	log.Println("tables created")
	return nil
}
