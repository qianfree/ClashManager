package handlers

import (
	"clash-manager/internal/model"
	"clash-manager/internal/repository"
	"clash-manager/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SubscriptionSourceHandler struct {
	Repo      *repository.SubscriptionSourceRepository
	NodeRepo  *repository.NodeRepository
	GroupRepo *repository.GroupRepository
}

func NewSubscriptionSourceHandler() *SubscriptionSourceHandler {
	return &SubscriptionSourceHandler{
		Repo:      repository.NewSubscriptionSourceRepository(),
		NodeRepo:  &repository.NodeRepository{},
		GroupRepo: &repository.GroupRepository{},
	}
}

func (h *SubscriptionSourceHandler) ListSources(c *gin.Context) {
	sources, err := h.Repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sources)
}

func (h *SubscriptionSourceHandler) GetSource(c *gin.Context) {
	id := c.Param("id")
	var idUint uint
	if _, err := fmt.Sscanf(id, "%d", &idUint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	source, err := h.Repo.FindByID(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source not found"})
		return
	}
	c.JSON(http.StatusOK, source)
}

func (h *SubscriptionSourceHandler) CreateSource(c *gin.Context) {
	var req struct {
		Name           string `json:"name" binding:"required"`
		URL            string `json:"url" binding:"required,url"`
		Enabled        bool   `json:"enabled"`
		UpdateInterval int    `json:"updateInterval"`
		NodeTag        string `json:"nodeTag"`
		SyncMode       string `json:"syncMode"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	source := &model.SubscriptionSource{
		Name:           req.Name,
		URL:            req.URL,
		Enabled:        req.Enabled,
		UpdateInterval: req.UpdateInterval,
		NodeTag:        req.NodeTag,
		SyncMode:       req.SyncMode,
	}

	if source.UpdateInterval == 0 {
		source.UpdateInterval = 24
	}
	if source.SyncMode == "" {
		source.SyncMode = "append"
	}

	if err := h.Repo.Create(source); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create source: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, source)
}

func (h *SubscriptionSourceHandler) UpdateSource(c *gin.Context) {
	id := c.Param("id")
	var idUint uint
	if _, err := fmt.Sscanf(id, "%d", &idUint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req struct {
		Name           string `json:"name"`
		URL            string `json:"url" binding:"required,url"`
		Enabled        bool   `json:"enabled"`
		UpdateInterval int    `json:"updateInterval"`
		NodeTag        string `json:"nodeTag"`
		SyncMode       string `json:"syncMode"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	source, err := h.Repo.FindByID(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source not found"})
		return
	}

	source.Name = req.Name
	source.URL = req.URL
	source.Enabled = req.Enabled
	source.UpdateInterval = req.UpdateInterval
	source.NodeTag = req.NodeTag
	source.SyncMode = req.SyncMode

	if err := h.Repo.Update(source); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update source: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, source)
}

func (h *SubscriptionSourceHandler) DeleteSource(c *gin.Context) {
	id := c.Param("id")
	var idUint uint
	if _, err := fmt.Sscanf(id, "%d", &idUint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.Repo.Delete(idUint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete source: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SubscriptionSourceHandler) SyncSource(c *gin.Context) {
	id := c.Param("id")
	var idUint uint
	if _, err := fmt.Sscanf(id, "%d", &idUint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	source, err := h.Repo.FindByID(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source not found"})
		return
	}

	// Parse subscription
	nodes, err := service.ParseSubscription(source.URL)
	if err != nil {
		// Update error info
		now := time.Now()
		source.LastSync = &now
		source.Error = err.Error()
		h.Repo.Update(source)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse subscription: " + err.Error()})
		return
	}

	// Get existing nodes
	existingNodes, err := h.NodeRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get existing nodes"})
		return
	}

	// Merge nodes based on sync mode
	mergedNodes := service.MergeNodes(existingNodes, nodes, source.SyncMode, source.Name)

	// Determine node tag
	nodeTag := source.NodeTag
	if nodeTag == "" {
		nodeTag = source.Name
	}

	// Apply tag to nodes from this subscription
	for i := range mergedNodes {
		// Check if this node is from the current subscription (by checking if it was in newNodes)
		isFromSubscription := false
		for _, n := range nodes {
			if n.Name == mergedNodes[i].Name {
				isFromSubscription = true
				break
			}
		}

		if isFromSubscription {
			// Add or update tag for nodes from this subscription
			existingConfig := mergedNodes[i].ExtraConfig
			// Remove old tag from this source if exists
			newConfig := ""
			if existingConfig != "" {
				lines := strings.Split(existingConfig, "\n")
				for _, line := range lines {
					if !strings.HasPrefix(line, "tag: ") {
						newConfig += line
						if line != "" {
							newConfig += "\n"
						}
					}
				}
			}
			mergedNodes[i].ExtraConfig = newConfig + "tag: " + nodeTag
		}
	}

	// Build a map of merged node names
	mergedNodeNames := make(map[string]bool)
	for _, n := range mergedNodes {
		mergedNodeNames[n.Name] = true
	}

	// For replace mode: delete nodes that are not in the merged list
	if source.SyncMode == "replace" {
		for _, existingNode := range existingNodes {
			if !mergedNodeNames[existingNode.Name] {
				// This node is not in the merged list, delete it
				h.NodeRepo.Delete(existingNode.ID)
			}
		}
	}

	// Track synced node IDs for proxy group
	syncedNodeIDs := make([]uint, 0)
	nodeNameToID := make(map[string]uint)

	// Save nodes by name - update if exists, create if not
	// This preserves node IDs so existing group references remain valid
	for _, node := range mergedNodes {
		// Check if this node is from the current subscription
		isFromSubscription := false
		for _, n := range nodes {
			if n.Name == node.Name {
				isFromSubscription = true
				break
			}
		}

		existingNode, err := h.NodeRepo.FindByName(node.Name)
		if err == nil && existingNode != nil {
			// Node exists, update it while preserving ID and CreatedAt
			node.ID = existingNode.ID
			node.CreatedAt = existingNode.CreatedAt
			// Only update source info if this node is from current subscription
			if isFromSubscription {
				node.Source = "subscription"
				node.SourceID = source.ID
				node.SourceName = source.Name
			}
			h.NodeRepo.Update(&node)
			nodeNameToID[node.Name] = node.ID
			// Only add to group if it's from this subscription
			if isFromSubscription {
				syncedNodeIDs = append(syncedNodeIDs, node.ID)
			}
		} else {
			// New node, create it
			if isFromSubscription {
				node.Source = "subscription"
				node.SourceID = source.ID
				node.SourceName = source.Name
			}
			h.NodeRepo.Create(&node)
			nodeNameToID[node.Name] = node.ID
			// Only add to group if it's from this subscription
			if isFromSubscription {
				syncedNodeIDs = append(syncedNodeIDs, node.ID)
			}
		}
	}

	// Create or update proxy group with subscription source name
	groupName := source.Name
	var group *model.ProxyGroupModel

	// Try to find existing group with the same name
	existingGroups, _ := h.GroupRepo.FindAll()
	var existingGroup *model.ProxyGroupModel
	for i := range existingGroups {
		if existingGroups[i].Name == groupName {
			existingGroup = &existingGroups[i]
			break
		}
	}

	// Prepare proxy IDs JSON (only contains nodes from this subscription)
	proxyIDsJSON, _ := json.Marshal(syncedNodeIDs)

	if existingGroup != nil {
		// Update existing group
		existingGroup.ProxyIDs = string(proxyIDsJSON)
		existingGroup.Type = "select"
		existingGroup.Source = "subscription"
		h.GroupRepo.Update(existingGroup)
		group = existingGroup
	} else {
		// Create new group with subscription source
		group = &model.ProxyGroupModel{
			Name:     groupName,
			Type:     "select",
			ProxyIDs: string(proxyIDsJSON),
			URL:      "http://www.gstatic.com/generate_204",
			Interval: 300,
			Source:   "subscription",
		}
		h.GroupRepo.Create(group)
	}

	// Update last sync time
	now := time.Now()
	source.LastSync = &now
	source.Error = ""
	h.Repo.Update(source)

	c.JSON(http.StatusOK, gin.H{
		"message":    "Sync completed successfully",
		"nodesCount": len(mergedNodes),
	})
}

func (h *SubscriptionSourceHandler) TestSource(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required,url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nodes, err := service.ParseSubscription(req.URL)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"nodesCount": len(nodes),
		"preview":    nodesPreview(nodes),
	})
}

func nodesPreview(nodes []model.Node) []map[string]interface{} {
	preview := make([]map[string]interface{}, 0, len(nodes))
	maxPreview := 10

	for i, node := range nodes {
		if i >= maxPreview {
			break
		}
		preview = append(preview, map[string]interface{}{
			"name":   node.Name,
			"type":   node.Type,
			"server": node.Server,
			"port":   node.Port,
		})
	}

	return preview
}
