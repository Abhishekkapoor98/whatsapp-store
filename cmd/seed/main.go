package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	db, err := sql.Open("postgres", "postgres://admin:password@localhost:6432/store_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	// Purane saare admins ko delete kar do taaki naya wala set ho sake
	_, err = db.Exec("DELETE FROM admins")
	if err != nil {
		log.Fatal("Could not clear admins: ", err)
	}

	// Naye password ko bcrypt se hash karo
	password := "ShimlaAdmin@1234"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Naya admin insert karo
	_, err = db.Exec("INSERT INTO admins (username, password_hash) VALUES ($1, $2)", "Admin", string(hash))
	if err != nil {
		log.Fatal("Could not create admin: ", err)
	}

	log.Println("✅ Admin credentials updated successfully!")
	log.Println("👤 Username: Admin")
	log.Println("🔑 Password: ShimlaAdmin@1234")
}
