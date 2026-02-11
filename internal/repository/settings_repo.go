package repository

import (
	"clash-manager/internal/model"
	"errors"

	"gorm.io/gorm"
)

type SettingsRepository struct{}

func (r *SettingsRepository) Get(key string) (string, error) {
	var setting model.GlobalSetting
	err := DB.Where("key = ?", key).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return setting.Value, nil
}

func (r *SettingsRepository) Set(key string, value string) error {
	var setting model.GlobalSetting
	err := DB.Where("key = ?", key).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			setting = model.GlobalSetting{
				Key:   key,
				Value: value,
			}
			return DB.Create(&setting).Error
		}
		return err
	}

	setting.Value = value
	return DB.Save(&setting).Error
}
