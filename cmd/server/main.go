package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/miles/rotki-demo/internal/api/handler"
	"github.com/miles/rotki-demo/internal/api/router"
	"github.com/miles/rotki-demo/internal/config"
	"github.com/miles/rotki-demo/internal/database"
	"github.com/miles/rotki-demo/internal/logger"
	"github.com/miles/rotki-demo/internal/provider/debank"
	"github.com/miles/rotki-demo/internal/repository"
	"github.com/miles/rotki-demo/internal/service"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.InitLogger(&cfg.Log); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Rotki Demo application")

	// Initialize database
	if err := database.InitDatabase(&cfg.Database); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Run migrations
	// Note: Migrations are handled by SQL file in migrations/001_initial_schema.sql
	// AutoMigrate is disabled to avoid conflicts with manual schema
	// if err := database.AutoMigrate(); err != nil {
	// 	logger.Fatal("Failed to run migrations", zap.Error(err))
	// }

	// Initialize repositories
	db := database.GetDB()
	walletRepo := repository.NewWalletRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	chainRepo := repository.NewChainRepository(db)
	rpcNodeRepo := repository.NewRPCNodeRepository(db)

	// Initialize all chains from chains.json
	chainInitializer := service.NewChainInitializer(chainRepo)
	if err := chainInitializer.InitializeAllChainsFromDefault(); err != nil {
		logger.Warn("Failed to initialize chains from file", zap.Error(err))
	}

	// Initialize data provider
	dataProvider := debank.NewDeBankProvider(&cfg.DeBank)

	// Initialize sync service
	syncService := service.NewSyncService(
		dataProvider,
		walletRepo,
		addressRepo,
		tokenRepo,
		chainRepo,
		cfg.Sync.GetSyncInterval(),
		cfg.Sync.BatchSize,
	)

	// Start sync service if enabled
	if cfg.Sync.Enabled {
		syncService.Start()
		defer syncService.Stop()
	}

	// Initialize RPC node service
	rpcNodeService := service.NewRPCNodeService(rpcNodeRepo, logger.GetLogger())

	// Initialize handlers
	walletHandler := handler.NewWalletHandler(walletRepo)
	addressHandler := handler.NewAddressHandler(addressRepo, tokenRepo, syncService)
	chainHandler := handler.NewChainHandler(chainRepo)
	rpcNodeHandler := handler.NewRPCNodeHandler(rpcNodeService, logger.GetLogger())

	// Setup router
	r := router.SetupRouter(walletHandler, addressHandler, chainHandler, rpcNodeHandler)

	// Start server
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("Starting HTTP server", zap.String("address", serverAddr))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Run(serverAddr); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-quit
	logger.Info("Shutting down server...")

	// Give outstanding requests 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Wait for context timeout
	<-ctx.Done()
	logger.Info("Server stopped")
}
