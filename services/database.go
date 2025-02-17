package services

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Driver   string
	Host     string
	Username string
	Password string
	Port     string
	Database string
}

type Database struct {
	*gorm.DB
}


var DB *gorm.DB

func NewDB(config *DatabaseConfig) (*gorm.DB, error) {
	var err error

	switch strings.ToLower(config.Driver) {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=UTC",
			config.Username, config.Password, config.Host, config.Port, config.Database)
			DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		break
	case "postgresql", "postgres":
		dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s TimeZone=UTC", config.Username, config.Password, config.Database, config.Host, config.Port)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		break
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(os.Getenv("DB_DATABASE")+".db"), &gorm.Config{})
		break
	}

	return DB, err
}
