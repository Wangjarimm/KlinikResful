package database

import (
	"fmt"
	"log"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sistem-informasi-klinik/config"
)

// Mendeklarasikan variabel untuk database
var DB *gorm.DB

// Fungsi ConnectDB() untuk menghubungkan ke database
func ConnectDB() {
    var err error
    p := config.Config("DB_PORT")
    port, err := strconv.ParseUint(p, 10, 32)
    if err != nil {
        log.Fatal("Error parsing DB_PORT")
    }

    // Cetak nilai-nilai yang digunakan untuk koneksi
    fmt.Printf("DB_HOST: %s\n", config.Config("DB_HOST"))
    fmt.Printf("DB_USER: %s\n", config.Config("DB_USER"))
    fmt.Printf("DB_PASSWORD: %s\n", config.Config("DB_PASSWORD"))
    fmt.Printf("DB_NAME: %s\n", config.Config("DB_NAME"))
    fmt.Printf("DB_PORT: %s\n", config.Config("DB_PORT"))

    // URL koneksi untuk menghubungkan ke database PostgreSQL (Supabase)
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=require TimeZone=Asia/Shanghai",
        config.Config("DB_HOST"), config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"), port)
    
    // Konfigurasi untuk mematikan caching prepared statements
    dbConfig := postgres.Config{
        DSN:                  dsn,
        PreferSimpleProtocol: true, // Disable prepared statement caching
    }

    // Menghubungkan ke DB dan menginisialisasikan variabel DB
    DB, err = gorm.Open(postgres.New(dbConfig), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            SingularTable: true,
        },
    })
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    fmt.Println("Connection Opened to Supabase Database")
}
