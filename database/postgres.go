package database

import (
	"fmt"
	"github.com/guil95/ports-service/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
)

func NewPostgresDB() *sqlx.DB {
	dbSource := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBUser,
		config.AppConfig.DBName,
		config.AppConfig.DBPassword,
		config.AppConfig.DBSSLMode,
	)

	db, err := sqlx.Connect("postgres", dbSource)
	if err != nil {
		slog.Error("error to connect on db", "err", err)
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return db
}
