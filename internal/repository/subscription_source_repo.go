package repository

import (
	"clash-manager/internal/model"
	"gorm.io/gorm"
)

type SubscriptionSourceRepository struct {
	DB *gorm.DB
}

func NewSubscriptionSourceRepository() *SubscriptionSourceRepository {
	return &SubscriptionSourceRepository{DB: GetDB()}
}

func (r *SubscriptionSourceRepository) FindAll() ([]model.SubscriptionSource, error) {
	var sources []model.SubscriptionSource
	err := r.DB.Order("created_at DESC").Find(&sources).Error
	return sources, err
}

func (r *SubscriptionSourceRepository) FindByID(id uint) (*model.SubscriptionSource, error) {
	var source model.SubscriptionSource
	err := r.DB.First(&source, id).Error
	if err != nil {
		return nil, err
	}
	return &source, nil
}

func (r *SubscriptionSourceRepository) Create(source *model.SubscriptionSource) error {
	return r.DB.Create(source).Error
}

func (r *SubscriptionSourceRepository) Update(source *model.SubscriptionSource) error {
	return r.DB.Save(source).Error
}

func (r *SubscriptionSourceRepository) Delete(id uint) error {
	return r.DB.Delete(&model.SubscriptionSource{}, id).Error
}

func (r *SubscriptionSourceRepository) UpdateLastSync(id uint, lastSync *model.SubscriptionSource) error {
	return r.DB.Model(&model.SubscriptionSource{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_sync": lastSync.LastSync,
		"error":     lastSync.Error,
	}).Error
}

func (r *SubscriptionSourceRepository) FindEnabled() ([]model.SubscriptionSource, error) {
	var sources []model.SubscriptionSource
	err := r.DB.Where("enabled = ?", true).Find(&sources).Error
	return sources, err
}
