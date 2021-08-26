package services

import (
	"github.com/merico-dev/lake/api/models"
	"github.com/merico-dev/lake/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var connectionString = config.V.GetString("DB_URL")
	var err error
	db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	migrateDB()
}

func migrateDB() {
	err := db.AutoMigrate(&models.Source{}, &models.Task{})
	if err != nil {
		panic(err)
	}
	// TODO: create customer migration here
}
