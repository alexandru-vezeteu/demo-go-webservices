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

	// If username and password are provided, use authenticated connection
	if user != "" && password != "" {
		// Construct MongoDB connection string with authentication
		// Using the database name as authSource by default (where the user was created)
		authSource := os.Getenv("USER_DB_AUTH_SOURCE")
		if authSource == "" {
			authSource = dbName // Use the target database as auth source
		}

		uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s", user, password, host, port, dbName, dbName)

		clientOptions = options.Client().ApplyURI(uri)
	} else {
		// No authentication (for local development)
		uri = fmt.Sprintf("mongodb://%s:%s", host, port)
		clientOptions = options.Client().ApplyURI(uri)
		log.Println("Warning: Connecting to MongoDB without authentication")
	}

	// Set additional client options
	clientOptions.SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(5 * time.Minute)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify connection
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v\nURI pattern: mongodb://%s:%s (check credentials and auth source)", err, host, port)
	}

	log.Printf("Successfully connected to MongoDB at %s:%s", host, port)

	return client.Database(dbName)
}
