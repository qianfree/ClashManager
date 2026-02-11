package handlers

import (
	"net/http"
	"strconv"

	"clash-manager/internal/model"
	"clash-manager/internal/repository"
	"clash-manager/internal/service" // Added this import

	"github.com/gin-gonic/gin"
)

type NodeHandler struct {
	Repo *repository.NodeRepository
}

func NewNodeHandler() *NodeHandler {
	return &NodeHandler{Repo: &repository.NodeRepository{}}
}

func (h *NodeHandler) ListNodes(c *gin.Context) {
	nodes, err := h.Repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nodes)
}

func (h *NodeHandler) CreateNode(c *gin.Context) {
	var node model.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Repo.Create(&node); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, node)
}

func (h *NodeHandler) UpdateNode(c *gin.Context) {
	// For simplicity, we just bind and save. Real app might check ID match.
	var node model.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Ensure ID from param matches or is handled
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	node.ID = uint(id)

	if err := h.Repo.Update(&node); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, node)
}

func (h *NodeHandler) DeleteNode(c *gin.Context) {
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

func (h *NodeHandler) ImportNode(c *gin.Context) {
	var req struct {
		Link string `json:"link"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	node, err := service.ParseLink(req.Link)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link format: " + err.Error()})
		return
	}

	if err := h.Repo.Create(node); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save node: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, node)
}
