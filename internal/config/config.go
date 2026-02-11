package config

import (
	"os"
	"path/filepath"
)

const (
	ServerPort = ":8090"
)

// GetDBPath 获取数据库文件的绝对路径
// 基于 exe 文件所在目录，确保在任何位置运行都能找到数据库
func GetDBPath() string {
	// 获取可执行文件所在目录
	exePath, err := os.Executable()
	if err != nil {
		// 如果获取失败，使用当前工作目录
		return "data/clash.db"
	}

	exeDir := filepath.Dir(exePath)
	dbPath := filepath.Join(exeDir, "data", "clash.db")

	// 确保数据目录存在
	dataDir := filepath.Join(exeDir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		// 如果创建目录失败，返回相对路径
		return "data/clash.db"
	}

	return dbPath
}
