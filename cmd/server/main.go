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
	// 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.InitLogger(&cfg.Log); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Rotki Demo application")

	// 初始化数据库
	if err := database.InitDatabase(&cfg.Database); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// 运行迁移
	// 注意：迁移由 migrations/001_initial_schema.sql 中的 SQL 文件处理
	// 为避免与手动模式冲突，禁用 AutoMigrate
	// if err := database.AutoMigrate(); err != nil {
	// 	logger.Fatal("Failed to run migrations", zap.Error(err))
	// }

	// 初始化仓储层
	db := database.GetDB()
	walletRepo := repository.NewWalletRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	chainRepo := repository.NewChainRepository(db)
	rpcNodeRepo := repository.NewRPCNodeRepository(db)

	// 从 chains.json 初始化所有链
	chainInitializer := service.NewChainInitializer(chainRepo)
	if err := chainInitializer.InitializeAllChainsFromDefault(); err != nil {
		logger.Warn("Failed to initialize chains from file", zap.Error(err))
	}

	// 初始化数据提供者
	dataProvider := debank.NewDeBankProvider(&cfg.DeBank)

	// 初始化同步服务
	syncService := service.NewSyncService(
		dataProvider,
		walletRepo,
		addressRepo,
		tokenRepo,
		chainRepo,
		cfg.Sync.GetSyncInterval(),
		cfg.Sync.BatchSize,
	)

	// 如果启用则启动同步服务
	if cfg.Sync.Enabled {
		syncService.Start()
		defer syncService.Stop()
	}

	// 初始化 RPC 节点服务
	rpcNodeService := service.NewRPCNodeService(rpcNodeRepo, logger.GetLogger())

	// 初始化处理器
	walletHandler := handler.NewWalletHandler(walletRepo)
	addressHandler := handler.NewAddressHandler(addressRepo, tokenRepo, syncService)
	chainHandler := handler.NewChainHandler(chainRepo)
	rpcNodeHandler := handler.NewRPCNodeHandler(rpcNodeService, logger.GetLogger())

	// 设置路由
	r := router.SetupRouter(walletHandler, addressHandler, chainHandler, rpcNodeHandler)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("Starting HTTP server", zap.String("address", serverAddr))

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Run(serverAddr); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号
	<-quit
	logger.Info("Shutting down server...")

	// 给未完成的请求 5 秒时间完成
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 等待上下文超时
	<-ctx.Done()
	logger.Info("Server stopped")
}
