package repository

import (
	"clash-manager/internal/model"

	"crypto/rand"
	"encoding/hex"
)

type UserRepository struct{}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByToken(token string) (*model.User, error) {
	var user model.User
	err := DB.Where("token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *model.User) error {
	return DB.Create(user).Error
}

func (r *UserRepository) Count() (int64, error) {
	var count int64
	if err := DB.Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepository) UpdatePassword(username string, hashedPassword string) error {
	var user model.User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}
	user.Password = hashedPassword
	return DB.Save(&user).Error
}

// GenerateToken generates a random 32-byte token
func (r *UserRepository) GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// UpdateToken updates the user's subscription token
func (r *UserRepository) UpdateToken(userID uint, token string) error {
	var user model.User
	if err := DB.First(&user, userID).Error; err != nil {
		return err
	}
	user.Token = token
	return DB.Save(&user).Error
}

// RefreshToken generates and updates a new token for the user
func (r *UserRepository) RefreshToken(userID uint) (string, error) {
	token := r.GenerateToken()
	if err := r.UpdateToken(userID, token); err != nil {
		return "", err
	}
	return token, nil
}
