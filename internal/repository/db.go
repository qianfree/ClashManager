package repository

import (
	"clash-manager/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(dsn string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	// Migrate the schema
	return DB.AutoMigrate(
		&model.Node{},
		&model.Rule{},
		&model.ProxyGroupModel{},
		&model.GlobalSetting{},
		&model.User{},
		&model.SubscriptionLog{},
		&model.SubscriptionSource{},
	)
}

// GetDB 返回数据库实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
