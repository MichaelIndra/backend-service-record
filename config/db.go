package config

import (
	"fmt"
	"log"
	"os"

	"service-record/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// 1. Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Ambil variabel dari .env
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// 3. Susun Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// 4. Buka koneksi ke MySQL
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal koneksi ke database!")
	}

	// 5. Auto Migration (Membuat tabel berdasarkan model secara otomatis)
	err = database.AutoMigrate(
		&models.Kendaraan{},
		&models.Bengkel{},
	)
	if err != nil {
		log.Fatal("Gagal migrasi tabel master:", err)
	}

	err = database.AutoMigrate(
		&models.Transaction{},
		&models.DetailService{},
		&models.DetailSparepart{},
	)
	if err != nil {
		log.Fatal("Gagal migrasi tabel transaksi:", err)
	}

	// 6. Simpan koneksi ke variabel global DB agar bisa dipakai di controller
	DB = database
	fmt.Println("Database terkoneksi dan migrasi berhasil!")
}
