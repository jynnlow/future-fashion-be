package infra

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"future-fashion/models"
)

//connect db
func InitDB() (*gorm.DB, error) {
	//to connect mysql db - username:password@protocol(address)/dbname?param=value
	dsn := "root:04110203@tcp(127.0.0.1:3306)/future_fashion_app?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Credential{}, &models.Product{}, &models.Order{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
