package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/miles/rotki-demo/internal/api/handler"
)

// SetupRouter sets up the HTTP router
func SetupRouter(
	walletHandler *handler.WalletHandler,
	addressHandler *handler.AddressHandler,
	chainHandler *handler.ChainHandler,
	rpcNodeHandler *handler.RPCNodeHandler,
) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Wallet routes
		wallets := v1.Group("/wallets")
		{
			wallets.POST("", walletHandler.CreateWallet)
			wallets.GET("", walletHandler.ListWallets)
			wallets.GET("/:id", walletHandler.GetWallet)
			wallets.PUT("/:id", walletHandler.UpdateWallet)
			wallets.DELETE("/:id", walletHandler.DeleteWallet)
			wallets.POST("/:id/refresh", addressHandler.RefreshWallet)
		}

		// Address routes
		addresses := v1.Group("/addresses")
		{
			addresses.POST("", addressHandler.CreateAddress)
			addresses.GET("", addressHandler.ListAddresses)
			addresses.GET("/:id", addressHandler.GetAddress)
			addresses.PUT("/:id", addressHandler.UpdateAddress)
			addresses.DELETE("/:id", addressHandler.DeleteAddress)
			addresses.POST("/:id/refresh", addressHandler.RefreshAddress)
		}

		// Chain routes
		chains := v1.Group("/chains")
		{
			chains.GET("", chainHandler.ListChains)
		}

		// RPC Node routes
		rpcNodes := v1.Group("/rpc-nodes")
		{
			rpcNodes.POST("", rpcNodeHandler.CreateRPCNode)
			rpcNodes.GET("", rpcNodeHandler.ListRPCNodes)
			rpcNodes.GET("/grouped", rpcNodeHandler.GetRPCNodesByChain)
			rpcNodes.GET("/:id", rpcNodeHandler.GetRPCNode)
			rpcNodes.PUT("/:id", rpcNodeHandler.UpdateRPCNode)
			rpcNodes.DELETE("/:id", rpcNodeHandler.DeleteRPCNode)
			rpcNodes.POST("/:id/check", rpcNodeHandler.CheckRPCNodeConnection)
			rpcNodes.POST("/check-all", rpcNodeHandler.CheckAllRPCNodeConnections)
		}
	}

	return router
}
