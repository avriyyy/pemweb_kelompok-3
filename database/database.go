package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Ganti dengan kredensial MySQL Anda:
	// user:password@tcp(127.0.0.1:3306)/nama_database
	// dsn := "toktik:Toktik123!@tcp(43.157.228.4:3306)/toktik?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:#Bahtiar28@tcp(127.0.0.1:3306)/toktik_dev?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database! \n", err)
	}
	DB = database
	log.Println("Koneksi database berhasil!")
}
