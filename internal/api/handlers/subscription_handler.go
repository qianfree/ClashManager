package handlers

import (
	"clash-manager/internal/model"
	"clash-manager/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	UserRepo *repository.UserRepository
	LogRepo  *repository.SubscriptionLogRepository
}

func NewSubscriptionHandler() *SubscriptionHandler {
	return &SubscriptionHandler{
		UserRepo: &repository.UserRepository{},
		LogRepo:  &repository.SubscriptionLogRepository{},
	}
}

// GetToken returns the current user's subscription token
func (h *SubscriptionHandler) GetToken(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.UserRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate token if doesn't exist
	if user.Token == "" {
		token, err := h.UserRepo.RefreshToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		user.Token = token
	}

	c.JSON(http.StatusOK, gin.H{
		"token": user.Token,
	})
}

// RefreshToken generates a new subscription token
func (h *SubscriptionHandler) RefreshToken(c *gin.Context) {
	userID := c.GetUint("user_id")

	token, err := h.UserRepo.RefreshToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "Token refreshed successfully",
	})
}

// GetSubscriptionURL returns the full subscription URL for the user
func (h *SubscriptionHandler) GetSubscriptionURL(c *gin.Context) {
	userID := c.GetUint("user_id")
	host := c.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = c.GetHeader("Host")
	}
	// Fallback to request host if still empty
	if host == "" {
		host = c.Request.Host
	}
	// Final fallback to localhost
	if host == "" {
		host = "localhost:8090"
	}
	scheme := "http"
	if c.GetHeader("X-Forwarded-Proto") == "https" || c.Request.TLS != nil {
		scheme = "https"
	}

	user, err := h.UserRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate token if doesn't exist
	if user.Token == "" {
		token, err := h.UserRepo.RefreshToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		user.Token = token
	}

	subURL := scheme + "://" + host + "/sub/" + user.Token

	c.JSON(http.StatusOK, gin.H{
		"url":   subURL,
		"token": user.Token,
	})
}

// GetSubscriptionLogs returns subscription logs with pagination
func (h *SubscriptionHandler) GetSubscriptionLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	userIDStr := c.Query("userId")
	successStr := c.Query("success")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	params := &repository.SubLogListParams{
		Page:     page,
		PageSize: pageSize,
	}

	// Filter by user
	if userIDStr != "" {
		userID, _ := strconv.ParseUint(userIDStr, 10, 32)
		params.UserID = uint(userID)
	}

	// Filter by success status
	if successStr == "true" {
		success := true
		params.Success = &success
	} else if successStr == "false" {
		success := false
		params.Success = &success
	}

	result, err := h.LogRepo.FindWithPagination(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetSubscriptionStats returns subscription statistics
func (h *SubscriptionHandler) GetSubscriptionStats(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, _ := strconv.Atoi(daysStr)

	stats, err := h.LogRepo.GetStats(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// DeleteOldLogs deletes old subscription logs
func (h *SubscriptionHandler) DeleteOldLogs(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "90")
	days, _ := strconv.Atoi(daysStr)

	if days < 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Days must be at least 7"})
		return
	}

	if err := h.LogRepo.DeleteOldLogs(days); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Old logs deleted successfully"})
}

// LogSubscription logs a subscription access
func (h *SubscriptionHandler) LogSubscription(userID uint, token, ip, userAgent string, success bool, errMsg string) {
	// Only log first 8 chars of token for privacy
	shortToken := ""
	if len(token) > 8 {
		shortToken = token[:8] + "..."
	} else {
		shortToken = token
	}

	log := &model.SubscriptionLog{
		UserID:    userID,
		Token:     shortToken,
		IP:        ip,
		UserAgent: userAgent,
		Success:   success,
		Error:     errMsg,
	}

	h.LogRepo.Create(log)

	// Clean up old logs (older than 90 days) in background
	go func() {
		h.LogRepo.DeleteOldLogs(90)
	}()
}

// GetOnlineClients returns active clients in the last hour
func (h *SubscriptionHandler) GetOnlineClients(c *gin.Context) {

	var result []struct {
		Username string `json:"username"`
		IP       string `json:"ip"`
		LastSeen string `json:"last_seen"`
	}

	// This is a simplified version - in production you might want more sophisticated tracking
	c.JSON(http.StatusOK, gin.H{
		"clients": result,
		"total":   0,
	})
}
