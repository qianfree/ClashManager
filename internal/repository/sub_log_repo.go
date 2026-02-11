package repository

import (
	"clash-manager/internal/model"
	"time"
)

type SubscriptionLogRepository struct{}

type SubLogListParams struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	UserID   uint  `json:"userId"`
	Success  *bool `json:"success"`
}

type SubLogListResult struct {
	Logs       []model.SubscriptionLog `json:"logs"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	PageSize   int                     `json:"pageSize"`
	TotalPages int                     `json:"totalPages"`
}

func (r *SubscriptionLogRepository) Create(log *model.SubscriptionLog) error {
	return DB.Create(log).Error
}

func (r *SubscriptionLogRepository) FindWithPagination(params *SubLogListParams) (*SubLogListResult, error) {
	var logs []model.SubscriptionLog
	var total int64

	query := DB.Model(&model.SubscriptionLog{})

	// Filter by user
	if params.UserID > 0 {
		query = query.Where("user_id = ?", params.UserID)
	}

	// Filter by success status
	if params.Success != nil {
		query = query.Where("success = ?", *params.Success)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Fetch data with user info
	if err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(params.PageSize).Find(&logs).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &SubLogListResult{
		Logs:       logs,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetStats returns subscription statistics
func (r *SubscriptionLogRepository) GetStats(days int) (map[string]interface{}, error) {
	var stats struct {
		TotalSubscriptions int64 `json:"total_subscriptions"`
		SuccessCount       int64 `json:"success_count"`
		FailCount          int64 `json:"fail_count"`
		TopUsers           []struct {
			Username string `json:"username"`
			Count    int64  `json:"count"`
		} `json:"top_users"`
	}

	since := time.Now().AddDate(0, 0, -days)

	// Total subscriptions in the period
	DB.Model(&model.SubscriptionLog{}).Where("created_at >= ?", since).Count(&stats.TotalSubscriptions)

	// Success count
	DB.Model(&model.SubscriptionLog{}).Where("created_at >= ? AND success = ?", since, true).Count(&stats.SuccessCount)

	// Fail count
	DB.Model(&model.SubscriptionLog{}).Where("created_at >= ? AND success = ?", since, false).Count(&stats.FailCount)

	// Top users
	DB.Table("subscription_logs").
		Select("users.username, COUNT(*) as count").
		Joins("JOIN users ON users.id = subscription_logs.user_id").
		Where("subscription_logs.created_at >= ?", since).
		Group("users.id").
		Order("count DESC").
		Limit(5).
		Scan(&stats.TopUsers)

	result := map[string]interface{}{
		"total_subscriptions": stats.TotalSubscriptions,
		"success_count":       stats.SuccessCount,
		"fail_count":          stats.FailCount,
		"top_users":           stats.TopUsers,
	}

	return result, nil
}

// DeleteOldLogs deletes logs older than specified days
func (r *SubscriptionLogRepository) DeleteOldLogs(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	return DB.Where("created_at < ?", cutoff).Delete(&model.SubscriptionLog{}).Error
}
