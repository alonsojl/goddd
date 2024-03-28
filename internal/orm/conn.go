package orm

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMySQL() (*gorm.DB, error) {
	var (
		dbUser     = os.Getenv("DB_USER")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbHost     = os.Getenv("DB_HOST")
		dbPort     = os.Getenv("DB_PORT")
		dbName     = os.Getenv("DB_NAME")
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}

	//Doc: https://go.dev/doc/database/manage-connections
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(1000)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
