package database

import (
    "database/sql"
    "log"
    _ "github.com/mattn/go-sqlite3"
)

func ConnectDatabase() *sql.DB {
    db, err := sql.Open("sqlite3", "database/blog.db")
    if err != nil {
        log.Println("error: ", err.Error())
    }
    return db
}
