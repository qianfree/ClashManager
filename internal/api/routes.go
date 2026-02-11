package api

import (
	"clash-manager/internal/api/handlers"
	"clash-manager/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Initialize handlers
	nodeHandler := handlers.NewNodeHandler()
	ruleHandler := handlers.NewRuleHandler()
	groupHandler := handlers.NewGroupHandler()
	subHandler := handlers.NewSubHandler()
	settingsHandler := handlers.NewSettingsHandler()
	authHandler := handlers.NewAuthHandler() // Add auth handler
	subscriptionHandler := handlers.NewSubscriptionHandler()

	// Public Auth Routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/setup", authHandler.Setup)
	}

	// Protected API Group
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Auth management (create new user)
		api.POST("/auth/register", authHandler.CreateUser)
		api.POST("/auth/password", authHandler.ChangePassword) // Change password

		// Node routes
		api.GET("/nodes", nodeHandler.ListNodes)
		api.POST("/nodes", nodeHandler.CreateNode)
		api.POST("/nodes/import", nodeHandler.ImportNode)
		api.PUT("/nodes/:id", nodeHandler.UpdateNode)
		api.DELETE("/nodes/:id", nodeHandler.DeleteNode)

		// Rule routes
		api.GET("/rules", ruleHandler.ListRules)
		api.POST("/rules", ruleHandler.CreateRule)
		api.POST("/rules/import", ruleHandler.ImportRules)
		api.PUT("/rules/:id", ruleHandler.UpdateRule)
		api.DELETE("/rules/:id", ruleHandler.DeleteRule)

		// Group routes
		api.GET("/groups", groupHandler.ListGroups)
		api.POST("/groups", groupHandler.CreateGroup)
		api.PUT("/groups/:id", groupHandler.UpdateGroup)
		api.DELETE("/groups/:id", groupHandler.DeleteGroup)

		// Settings routes
		api.GET("/settings/dns", settingsHandler.GetDNS)
		api.POST("/settings/dns", settingsHandler.SaveDNS)

		// Subscription management routes
		api.GET("/subscription/token", subscriptionHandler.GetToken)
		api.POST("/subscription/token/refresh", subscriptionHandler.RefreshToken)
		api.GET("/subscription/url", subscriptionHandler.GetSubscriptionURL)
		api.GET("/subscription/logs", subscriptionHandler.GetSubscriptionLogs)
		api.GET("/subscription/stats", subscriptionHandler.GetSubscriptionStats)
		api.DELETE("/subscription/logs/old", subscriptionHandler.DeleteOldLogs)
		api.GET("/subscription/online", subscriptionHandler.GetOnlineClients)
	}

	// Subscription route (public, requires token)
	r.GET("/sub/:token", subHandler.GetConfig)
}
