package config

import (
	"fmt"
	"log"
	"web_sarpras/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB 

func ConnectDatabase() {
	dsn := "root:kyano211209@tcp(127.0.0.1:3306)/web_sarpras?charset=utf8mb4&parseTime=True&loc=Local"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database: ", err)
	}

	
	err = database.AutoMigrate(&models.User{}, &models.Aset{})
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database: ", err)
	}

	DB = database
	fmt.Println("Migrasi database berhasil!")
}