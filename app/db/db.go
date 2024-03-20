package db

import (
	"fmt"

	"github.com/Tanakaryuki/brachio-backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() *gorm.DB {
	dbUser := config.DB_USER
	dbPassword := config.DB_PASSWORD
	mysqlHost := config.MYSQL_HOST
	dbName := config.DB_NAME
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		dbUser,
		dbPassword,
		mysqlHost,
		dbName,
	)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Could not connect to database.")
	}
	return DB
}
