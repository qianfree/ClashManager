package main

import (
	"log"

	"flag"
	"fmt"
	"os"

	"clash-manager/internal/api"
	"clash-manager/internal/config"
	"clash-manager/internal/repository"

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
		handleResetInfo(*resetAdminPtr)
		return
	}

	// 2. Setup Router
	r := gin.Default()
	r.RedirectTrailingSlash = false // 禁用 trailing slash 重定向
	r.RedirectFixedPath = false     // 禁用路径修正重定向
	api.SetupRoutes(r)
	setupStaticRoutes(r) // 必须在 API 路由之后注册

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

func handleResetInfo(newPassword string) {
	fmt.Printf("Resetting admin password...\n")

	repo := &repository.UserRepository{}
	if _, err := repo.FindByUsername("admin"); err != nil {
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
}
