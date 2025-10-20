package Config

import (
	"eventManager/Model"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Bucharest",
		getEnv("EVENT_MANAGER_DB_HOST", "localhost"),
		getEnv("EVENT_MANAGER_DB_USER", "postgres"),
		getEnv("EVENT_MANAGER_DB_PASSWORD", "postgres"),
		getEnv("EVENT_MANAGER_DB_NAME", "eventdb"),
		getEnv("EVENT_MANAGER_DB_PORT", "5432"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
}

func Migrate() {
	DB.AutoMigrate(
		&Model.Event{},
		&Model.EventPacket{},
		&Model.PacketEventRelation{},
		&Model.PacketEventSeats{},
	)
	log.Println("✅ Database migration completed")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
