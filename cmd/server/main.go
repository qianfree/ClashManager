package main

import (
	"clash-manager/internal/api"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"clash-manager/internal/config"
	"clash-manager/internal/repository"

	"clash-manager/web"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Parse Flags
	resetAdminPtr := flag.String("reset-admin", "", "Reset admin password to the specified value")
	portPtr := flag.String("port", "8090", "Server port (default: 8090)")
	flag.Parse()

	// 格式化端口号（确保有 : 前缀）
	serverPort := formatPort(*portPtr)

	// 1. Initialize Database
	if err := repository.InitDB(config.GetDBPath()); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Handle Reset Password Flag
	if *resetAdminPtr != "" {
		handleResetAdmin(*resetAdminPtr)
		return
	}

	// 2. Setup Router
	r := gin.Default()
	// API接口
	r.RedirectTrailingSlash = false // 禁用 trailing slash 重定向
	r.RedirectFixedPath = false     // 禁用路径修正重定向
	api.SetupRoutes(r)

	// 静态文件处理
	subFS, _ := fs.Sub(web.StaticFiles, "dist")
	staticHandler := http.FS(subFS)
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 尝试在嵌入的文件系统中打开该文件
		f, err := subFS.Open(path[1:]) // 移除开头的 "/"
		if err != nil {
			// 文件不存在，返回 index.html (交给前端路由处理)
			indexFile, _ := subFS.Open("index.html")
			defer indexFile.Close()

			// 获取文件信息以读取内容
			stat, _ := indexFile.Stat()
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html; charset=utf-8", indexFile, nil)
			return
		}
		defer f.Close()

		// 如果文件存在，直接提供服务
		http.FileServer(staticHandler).ServeHTTP(c.Writer, c.Request)
	})

	// 3. Start Server
	log.Printf("Server starting on %s...", serverPort)
	if err := r.Run(serverPort); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

// formatPort 格式化端口号，确保有 : 前缀
func formatPort(port string) string {
	if port == "" {
		return ":8090"
	}
	// 如果已经以 : 开头，直接返回
	if port[0] == ':' {
		return port
	}
	// 否则添加 : 前缀
	return ":" + port
}

func handleResetAdmin(newPassword string) {
	fmt.Printf("Database: %s\n", config.GetDBPath())
	fmt.Printf("Resetting admin password...\n")

	repo := &repository.UserRepository{}
	user, err := repo.FindByUsername("admin")
	if err != nil {
		fmt.Printf("Error: Admin user not found. Please ensure the system is initialized.\n")
		os.Exit(1)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error: Failed to encrypt password: %v\n", err)
		os.Exit(1)
	}

	if err := repo.UpdatePassword("admin", string(hashedPassword)); err != nil {
		fmt.Printf("Error: Failed to update password in DB: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success: Admin password reset to '%s'\n", newPassword)
	fmt.Printf("User ID: %d\n", user.ID)
}
