package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DBHost     = "USER_DB_HOST"
	DBUser     = "USER_DB_USER"
	DBPassword = "USER_DB_PASSWORD"
	DBName     = "USER_DB_NAME"
	DBPort     = "USER_DB_PORT"
)

func InitDB() *mongo.Database {
	host := os.Getenv(DBHost)
	user := os.Getenv(DBUser)
	password := os.Getenv(DBPassword)
	dbName := os.Getenv(DBName)
	port := os.Getenv(DBPort)

	if host == "" || dbName == "" || port == "" {
		log.Fatal("Missing required database environment variables (host, name, port)")
	}

	var uri string
	var clientOptions *options.ClientOptions

	if user != "" && password != "" {
		authSource := os.Getenv("USER_DB_AUTH_SOURCE")
		if authSource == "" {
			authSource = dbName
		}

		uri = fmt.Sprintf("mongodb:%s%s%s:%s@%s:%s/%s?authSource=%s", "/", "/", user, password, host, port, dbName, authSource)

		clientOptions = options.Client().ApplyURI(uri)
	} else {
		uri = fmt.Sprintf("mongodb:%s%s%s:%s/%s", "/", "/", host, port, dbName)
		clientOptions = options.Client().ApplyURI(uri)
		log.Println("Warning: Connecting to MongoDB without authentication")
	}

	clientOptions.SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v\nURI pattern: mongodb:%s%s<host>:<port>", err, "/", "/")
	}

	log.Printf("Successfully connected to MongoDB at %s:%s", host, port)

	return client.Database(dbName)
}
