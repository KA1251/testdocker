package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	dbHost  = "cockroachdb" // имя сервиса CockroachDB в сети Docker
	dbPort  = 26257
	dbUser  = "root"
	dbName  = "test_db"
	sslMode = "disable"
)

func main() {
	// Формирование строки подключения к базе данных
	connStr := fmt.Sprintf("postgresql://%s@%s:%d/%s?sslmode=%s", dbUser, dbHost, dbPort, dbName, sslMode)

	// Подключение к базе данных
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ { // Пробуем подключиться 10 раз с интервалом в 1 секунду
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Error connecting to the database: %v. Retrying...", err)
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	if db == nil {
		log.Fatal("Failed to connect to the database after multiple attempts")
	}
	defer db.Close()

	// Создание таблицы
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name STRING)`); err != nil {
		log.Fatal("Error creating table: ", err)
	}

	// Вставка данных
	if _, err := db.Exec(`INSERT INTO users (name) VALUES ($1)`, "John Doe"); err != nil {
		log.Fatal("Error inserting data: ", err)
	}

	// Чтение данных
	rows, err := db.Query(`SELECT id, name FROM users`)
	if err != nil {
		log.Fatal("Error querying data: ", err)
	}
	defer rows.Close()

	// Вывод данных
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal("Error scanning row: ", err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating rows: ", err)
	}
}
