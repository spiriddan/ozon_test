package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type Config struct {
	Host, Port, User, Password, Dbname string
}

func GetPostgresConnection(config Config) (*sql.DB, error) {
	psqInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)

	db, err := sql.Open("postgres", psqInfo)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3) // это можно было бы в конфиг запихнуть
	db.SetMaxOpenConns(90)
	db.SetMaxIdleConns(90)

	return db, nil
}
