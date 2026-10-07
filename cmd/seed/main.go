package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://admin:password@localhost:6432/store_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database ping failed: ", err)
	}
	log.Println("✅ Connected to Database successfully!")

	// Safe table creation query so it never fails if tables already exist
	query := `
	CREATE TABLE IF NOT EXISTS admins (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		price NUMERIC(10, 2) NOT NULL,
		stock_count INT NOT NULL DEFAULT 0,
		image_url TEXT NOT NULL
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create tables: ", err)
	}
	log.Println("✅ Tables verified/created successfully!")

	// Purane admins clear karo
	_, err = db.Exec("DELETE FROM admins")
	if err != nil {
		log.Fatal("Could not clear admins: ", err)
	}

	// Naya Admin insert karo
	password := "ShimlaAdmin@1234"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	_, err = db.Exec("INSERT INTO admins (username, password_hash) VALUES ($1, $2)", "Admin", string(hash))
	if err != nil {
		log.Fatal("Could not create admin: ", err)
	}

	log.Println("✅ Admin user created successfully!")
	log.Println("👤 Username: Admin")
	log.Println("🔑 Password: ShimlaAdmin@1234")
}
