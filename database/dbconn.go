package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDatabase() *sql.DB {

	godotenv.Load()
	
	DBHOST := os.Getenv("DBHOST")
	DBPORT := os.Getenv("DBPORT")
	DBNAME := os.Getenv("DBNAME")
	DBUSER := os.Getenv("DBUSER")
	DBPASSWORD := os.Getenv("DBPASSWORD")
	DBSSLMODE := os.Getenv("DBSSLMODE")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",DBHOST, DBPORT, DBUSER, DBPASSWORD, DBNAME, DBSSLMODE)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Info("Error Connecting to Database", "message", err.Error())
	}

	err = db.Ping()
	if err != nil {
		slog.Error("error pinging the database", "error", err.Error())
	}

	data,err  := os.ReadFile("database/init.sql")
	if err != nil {
		slog.Error("error reading initdb file content", "error", err.Error())
	}

	fmt.Println("sql data" , string(data))
	if _, err := db.Exec(string(data)); err != nil {
		slog.Error("error executing init sql", "error", err.Error())
	}
	return db
}
