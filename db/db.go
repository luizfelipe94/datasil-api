package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/luizfelipe94/datasil/configs"
)

func NewDB(cfg *configs.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(0)
	db.SetConnMaxLifetime(time.Nanosecond)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}

func CountRows(rows *sql.Rows) (count int) {
	for rows.Next() {
		rows.Scan(&count)
	}
	return count
}
