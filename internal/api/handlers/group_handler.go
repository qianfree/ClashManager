package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"clash-manager/internal/model"
	"clash-manager/internal/repository"

	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	Repo *repository.GroupRepository
}

func NewGroupHandler() *GroupHandler {
	return &GroupHandler{Repo: &repository.GroupRepository{}}
}

func (h *GroupHandler) ListGroups(c *gin.Context) {
	groups, err := h.Repo.FindAllWithNodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, groups)
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req struct {
		Name     string `json:"Name"`
		Type     string `json:"Type"`
		ProxyIDs []uint `json:"ProxyIDs"`
		URL      string `json:"URL"`
		Interval int    `json:"Interval"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert ProxyIDs array to JSON string
	proxyIDsJSON := "[]"
	if len(req.ProxyIDs) > 0 {
		data, err := json.Marshal(req.ProxyIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		proxyIDsJSON = string(data)
	}

	group := model.ProxyGroupModel{
		Name:     req.Name,
		Type:     req.Type,
		ProxyIDs: proxyIDsJSON,
		Use:      "",
		URL:      req.URL,
		Interval: req.Interval,
	}

	if err := h.Repo.Create(&group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, group)
}

func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.Repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req struct {
		Name     string `json:"Name"`
		Type     string `json:"Type"`
		ProxyIDs []uint `json:"ProxyIDs"`
		URL      string `json:"URL"`
		Interval int    `json:"Interval"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing group
	group, err := h.Repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Convert ProxyIDs array to JSON string
	proxyIDsJSON := "[]"
	if len(req.ProxyIDs) > 0 {
		data, err := json.Marshal(req.ProxyIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		proxyIDsJSON = string(data)
	}

	// Update fields
	group.Name = req.Name
	group.Type = req.Type
	group.ProxyIDs = proxyIDsJSON
	group.URL = req.URL
	group.Interval = req.Interval

	if err := h.Repo.Update(group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, group)
}
