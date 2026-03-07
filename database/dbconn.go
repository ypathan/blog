package database

import (
    "database/sql"
    "log/slog"
    _ "github.com/mattn/go-sqlite3"
)

func ConnectDatabase() *sql.DB {
    db, err := sql.Open("sqlite3", "database/blog.db")
    if err != nil {
        slog.Info("Error Connecting to Database","message", err.Error())
    }
    return db
}
