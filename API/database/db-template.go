package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

// Var global untuk koneksi
var DB *sql.DB

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func ConnectDB() error {

	var config = Config{
		Host:     "localhost",
		Port:     "5433",
		User:     "postgres",
		Password: "Password", // Ganti dengan password Anda
		DBName:   "POS_Restaurant",
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("gagal membuka koneksi: %v", err)
	}

	// Tambahkan validasi koneksi
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test Koneksi
	if err := DB.PingContext(ctx); err != nil {
		return fmt.Errorf("gagal ping database: %v", err)
	}

	log.Println("Koneksi database berhasil")
	return nil
}
