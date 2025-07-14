package main

import (
	"log"
	"net"
	"time"

	"github.com/VibeTeam/fitness-tracker-backend/shared/config"
	"github.com/VibeTeam/fitness-tracker-backend/user/auth"
	userhandler "github.com/VibeTeam/fitness-tracker-backend/user/handler"
	"github.com/VibeTeam/fitness-tracker-backend/user/models"
	"github.com/VibeTeam/fitness-tracker-backend/user/repository/gormrepository"
	usecase "github.com/VibeTeam/fitness-tracker-backend/user/use_case"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"google.golang.org/grpc"

	userv1 "github.com/VibeTeam/fitness-tracker-backend/proto/gen/go/user/v1"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to Postgres via GORM
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Migrate schema
	if !db.Migrator().HasTable(&models.User{}) {
		if err := db.AutoMigrate(&models.User{}); err != nil {
			log.Fatalf("auto-migrate: %v", err)
		}
	}

	// Build domain services
	accessTTL, err := time.ParseDuration(cfg.JWT.AccessTTL)
	if err != nil {
		log.Fatalf("parse access ttl: %v", err)
	}
	refreshTTL, err := time.ParseDuration(cfg.JWT.RefreshTTL)
	if err != nil {
		log.Fatalf("parse refresh ttl: %v", err)
	}

	tokenMgr := auth.NewManager(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, accessTTL, refreshTTL)
	userRepo := gormrepository.NewUserRepository(db)
	authSvc := usecase.NewAuthService(userRepo, tokenMgr)

	// Start gRPC server
	lis, err := net.Listen("tcp", cfg.Server.GRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userv1.RegisterUserServiceServer(grpcServer, userhandler.NewGRPCServer(authSvc, userRepo))

	log.Printf("user service gRPC listening on %s", cfg.Server.GRPCAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc serve: %v", err)
	}
}
