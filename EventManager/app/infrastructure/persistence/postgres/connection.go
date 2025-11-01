package postgres

import (
	gormmodel "eventManager/infrastructure/persistence/postgres/gormModel"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DBHost     = "EVENT_MANAGER_DB_HOST"
	DBUser     = "EVENT_MANAGER_DB_USER"
	DBPassword = "EVENT_MANAGER_DB_PASSWORD"
	DBName     = "EVENT_MANAGER_DB_NAME"
	DBPort     = "EVENT_MANAGER_DB_PORT"
)

func buildDSN() string {
	host := os.Getenv(DBHost)
	user := os.Getenv(DBUser)
	password := os.Getenv(DBPassword)
	dbname := os.Getenv(DBName)
	port := os.Getenv(DBPort)

	if host == "" || user == "" || password == "" || dbname == "" || port == "" {

		log.Fatal("FATAL: One or more critical database environment variables (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT) are missing.")
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port,
	)
}

func InitDB() *gorm.DB {
	dsn := buildDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection successfully established.")

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	err = db.AutoMigrate(&gormmodel.GormEvent{})
	if err != nil {
		log.Fatalf("FATAL: Failed to run migrations: %v", err)
	}

	fmt.Println("Database schema migrated successfully.")
	return db
}
