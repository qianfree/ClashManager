package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"clash-manager/internal/model"
	"clash-manager/internal/repository"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type RuleHandler struct {
	Repo *repository.RuleRepository
}

func NewRuleHandler() *RuleHandler {
	return &RuleHandler{Repo: &repository.RuleRepository{}}
}

func (h *RuleHandler) ListRules(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	ruleType := c.Query("type")
	keyword := c.Query("keyword")

	// Validate parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}

	params := &repository.RuleListParams{
		Page:     page,
		PageSize: pageSize,
		Type:     ruleType,
		Keyword:  keyword,
	}

	result, err := h.Repo.FindWithPagination(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *RuleHandler) CreateRule(c *gin.Context) {
	var rule model.Rule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Repo.Create(&rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rule)
}

func (h *RuleHandler) DeleteRule(c *gin.Context) {
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

func (h *RuleHandler) UpdateRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var rule model.Rule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing rule
	existingRule, err := h.Repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rule not found"})
		return
	}

	// Update fields (preserve ID)
	rule.ID = existingRule.ID
	if err := h.Repo.Update(&rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rule)
}

// ImportRequest represents the import request structure
type ImportRequest struct {
	Content string `json:"content"`
}

func (h *RuleHandler) ImportRules(c *gin.Context) {
	var req ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse YAML content
	var yamlConfig map[string]interface{}
	if err := yaml.Unmarshal([]byte(req.Content), &yamlConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid YAML format: " + err.Error()})
		return
	}

	// Extract rules section
	rulesSection, ok := yamlConfig["rules"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No rules section found in YAML"})
		return
	}

	rulesList, ok := rulesSection.([]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rules format"})
		return
	}

	var rulesToImport []model.Rule
	priority := 0

	for _, ruleItem := range rulesList {
		ruleStr, ok := ruleItem.(string)
		if !ok {
			continue
		}

		// Parse rule string: TYPE,Payload,Target[,no-resolve]
		parts := strings.Split(ruleStr, ",")
		if len(parts) < 3 {
			continue
		}

		rule := model.Rule{
			Type:     strings.TrimSpace(parts[0]),
			Payload:  strings.TrimSpace(parts[1]),
			Target:   strings.TrimSpace(parts[2]),
			Priority: priority,
		}

		// Check for no-resolve option
		if len(parts) >= 4 && strings.TrimSpace(parts[3]) == "no-resolve" {
			rule.NoResolve = true
		}

		rulesToImport = append(rulesToImport, rule)
		priority++
	}

	if len(rulesToImport) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid rules found"})
		return
	}

	// Batch create rules
	if err := h.Repo.BatchCreate(&rulesToImport); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rules imported successfully",
		"count":   len(rulesToImport),
	})
}
