package main

import (
	"database/sql"
	"fmt"
	"log"

	"askeladden/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("DROP TABLE IF EXISTS banned_bokmal_words, banned_bokmal_words_testing;")
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	fmt.Println("Tables dropped successfully.")
}
