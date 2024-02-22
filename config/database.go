package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func ConnectDB() (*sql.DB, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	dbUser := os.Getenv("DBUser")
	dbPass := os.Getenv("DBPassword")
	dbName := os.Getenv("DBName")

	driverSourceName := fmt.Sprintf("%s:%s@/%s?parseTime=true", dbUser, dbPass, dbName)

	db, err := sql.Open("mysql", driverSourceName)
	if err != nil {
		return nil, err
	}

	return db, nil
}
