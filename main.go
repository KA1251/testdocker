package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=disable", dbUser, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	// Создаем базу данных, если ее еще не существует
	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)); err != nil {
		log.Fatal("Error creating database:", err)
	}
	fmt.Println("created")
	// Проверяем существует ли база данных
	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = $1)`, dbName).Scan(&exists)
	if err != nil {
		log.Fatal("Error checking database existence:", err)
	}
	if !exists {
		log.Fatalf("Database %s does not exist11", dbName)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name STRING)`); err != nil {
		log.Fatal("Error creating table:", err)
	}

	if _, err := db.Exec(`INSERT INTO users (name) VALUES ($1)`, "John Doe"); err != nil {
		log.Fatal("Error inserting data:", err)
	}

	rows, err := db.Query(`SELECT id, name FROM users`)
	if err != nil {
		log.Fatal("Error querying data:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal("Error scanning row:", err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating rows:", err)
	}
}
