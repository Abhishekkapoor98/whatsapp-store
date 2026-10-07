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
