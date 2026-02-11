package handlers

import (
	"clash-manager/internal/model"
	"clash-manager/internal/repository"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	Repo *repository.SettingsRepository
}

func NewSettingsHandler() *SettingsHandler {
	return &SettingsHandler{Repo: &repository.SettingsRepository{}}
}

const KeyDNS = "dns_config"

func (h *SettingsHandler) GetDNS(c *gin.Context) {
	val, err := h.Repo.Get(KeyDNS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return default if empty
	if val == "" {
		defaultDNS := model.DNSConfig{
			Enable:       true,
			Listen:       "0.0.0.0:53",
			EnhancedMode: "fake-ip",
			Nameserver:   []string{"223.5.5.5", "119.29.29.29"},
			Fallback:     []string{"8.8.8.8", "1.1.1.1"},
		}
		c.JSON(http.StatusOK, defaultDNS)
		return
	}

	// Unmarshal to ensure validity or just return raw JSON?
	// Let's unmarshal to struct to verify and return guaranteed structure
	var dnsConfig model.DNSConfig
	json.Unmarshal([]byte(val), &dnsConfig)
	c.JSON(http.StatusOK, dnsConfig)
}

func (h *SettingsHandler) SaveDNS(c *gin.Context) {
	var dnsConfig model.DNSConfig
	if err := c.ShouldBindJSON(&dnsConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bytes, err := json.Marshal(dnsConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.Set(KeyDNS, string(bytes)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dnsConfig)
}
