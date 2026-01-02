package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"idmService/application/domain"
	"idmService/application/usecase"
	"idmService/infrastructure/blacklist"
	"idmService/infrastructure/persistence"
	"idmService/infrastructure/service"
	pb "idmService/proto"
	"idmService/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := ":50051"

	fmt.Println("Connecting to PostgreSQL database...")
	dbConfig := persistence.GetDatabaseConfigFromEnv()
	db, err := persistence.ConnectDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connected successfully")

	userRepo := persistence.NewPostgresUserRepository(db)

	fmt.Println("Running database migrations...")
	if err := userRepo.MigrateSchema(); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
	fmt.Println("Database migrations completed")

	serviceEmail := "clients_service@system.local"
	servicePassword := "service_secret_password"
	if err := userRepo.SeedServiceAccount(serviceEmail, servicePassword); err != nil {
		log.Printf("Warning: Failed to seed service account: %v", err)
	} else {
		fmt.Println("Service account seeded successfully")
	}

	var tokenBlacklist domain.TokenBlacklist
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr != "" {
		fmt.Println("Connecting to Redis for token blacklist...")
		redisBlacklist, err := blacklist.NewRedisBlacklist(&blacklist.RedisConfig{
			Addr:     redisAddr,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
			TTL:      24 * time.Hour,
		})
		if err != nil {
			log.Printf("Warning: Failed to connect to Redis, falling back to in-memory blacklist: %v", err)
			tokenBlacklist = blacklist.NewInMemoryBlacklist()
		} else {
			fmt.Println("Redis blacklist connected successfully")
			tokenBlacklist = redisBlacklist
		}
	} else {
		fmt.Println("Using in-memory token blacklist (set REDIS_ADDR for distributed blacklist)")
		tokenBlacklist = blacklist.NewInMemoryBlacklist()
	}

	tokenService := service.NewTokenService()

	loginUseCase := usecase.NewLoginUseCase(userRepo, tokenService)
	verifyTokenUseCase := usecase.NewVerifyTokenUseCase(userRepo, tokenService, tokenBlacklist)
	revokeTokenUseCase := usecase.NewRevokeTokenUseCase(tokenBlacklist, tokenService)

	idmServer := server.NewIdentityServer(
		loginUseCase,
		verifyTokenUseCase,
		revokeTokenUseCase,
	)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterIdentityServiceServer(grpcServer, idmServer)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
