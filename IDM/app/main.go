package main

import (
	"fmt"
	"log"
	"net"

	"idmService/infrastructure/blacklist"
	"idmService/infrastructure/persistence"
	pb "idmService/proto"
	"idmService/server"
	"idmService/service"
	"idmService/usecase"

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

	// Initialize blacklist
	tokenBlacklist := blacklist.NewInMemoryBlacklist()

	// Initialize services
	tokenService := service.NewTokenService()

	// Initialize use cases
	loginUseCase := usecase.NewLoginUseCase(userRepo, tokenService)
	verifyTokenUseCase := usecase.NewVerifyTokenUseCase(userRepo, tokenService, tokenBlacklist)
	revokeTokenUseCase := usecase.NewRevokeTokenUseCase(tokenBlacklist)

	// Initialize gRPC server with use cases
	idmServer := server.NewIdentityServer(loginUseCase, verifyTokenUseCase, revokeTokenUseCase)

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
