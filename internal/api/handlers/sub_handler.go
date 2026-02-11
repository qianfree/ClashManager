package handlers

import (
	"net/http"

	"clash-manager/internal/repository"
	"clash-manager/internal/service"

	"github.com/gin-gonic/gin"
)

type SubHandler struct {
	Service          *service.ConfigService
	UserRepo         *repository.UserRepository
	SubHandlerHelper *SubscriptionHandler
}

func NewSubHandler() *SubHandler {
	return &SubHandler{
		Service:          service.NewConfigService(),
		UserRepo:         &repository.UserRepository{},
		SubHandlerHelper: NewSubscriptionHandler(),
	}
}

func (h *SubHandler) GetConfig(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.String(http.StatusBadRequest, "Missing token")
		return
	}

	// Validate token
	user, err := h.UserRepo.FindByToken(token)
	if err != nil {
		h.SubHandlerHelper.LogSubscription(0, token, c.ClientIP(), c.GetHeader("User-Agent"), false, "Invalid token")
		c.String(http.StatusUnauthorized, "Invalid token")
		return
	}

	// Generate config
	configBytes, err := h.Service.GenerateConfig()
	if err != nil {
		h.SubHandlerHelper.LogSubscription(user.ID, token, c.ClientIP(), c.GetHeader("User-Agent"), false, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log successful subscription
	h.SubHandlerHelper.LogSubscription(user.ID, token, c.ClientIP(), c.GetHeader("User-Agent"), true, "")

	// Set headers for Clash
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.Header("Content-Disposition", "inline; filename=clash_config.yaml")
	// Set Subscription-Userinfo header for Clash clients to display traffic info
	// Format: upload=123; download=456; total=789; expire=123456789
	// Values are in bytes. expire is user's subscription expiration unix timestamp
	// Here we use dummy values as an example. You should fetch these from your database.
	c.Header("Subscription-Userinfo", "upload=0; download=0; total=10737418240; expire=0")
	// Set profile name for Clash
	c.Header("Profile-Title", "Clash-"+user.Username)

	c.String(http.StatusOK, string(configBytes))
}
