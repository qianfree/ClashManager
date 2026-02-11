package main

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed dist
var distFS embed.FS

func setupStaticRoutes(r *gin.Engine) {
	// 获取 dist 子文件系统
	subFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		return
	}

	// 创建文件服务器
	fileServer := http.FileServer(http.FS(subFS))

	// 使用 NoRoute 处理静态文件请求
	// 必须在 API 路由之后设置，这样 API 路由会优先匹配
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// API 路径不处理
		if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/sub") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		// 使用文件服务器处理
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
