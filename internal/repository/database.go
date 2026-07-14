package repository

import (
	"log"
	"tg/bot/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) (*gorm.DB, error) {
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных.", err)
	}

	err = DB.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatal("Не удалось выполнить миграции.", err)
	}

	return DB, nil
}

func UserExists(db *gorm.DB, telegramID int64) (bool, error) {

	var count int64

	err := db.Model(&domain.User{}).Where("id = ?", telegramID).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
