package main

import (
	"embed"
	"io/fs"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed dist
var distFS embed.FS

// SetupStaticRoutes 设置静态文件服务路由
func SetupStaticRoutes(r *gin.Engine) {
	// 获取 dist 子文件系统
	subFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		return
	}

	// 使用 NoRoute 处理静态文件请求
	r.NoRoute(func(c *gin.Context) {
		reqPath := c.Request.URL.Path

		// API 路径返回 404
		if strings.HasPrefix(reqPath, "/api") || strings.HasPrefix(reqPath, "/sub") {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}

		// 去掉开头的 /
		filePath := strings.TrimPrefix(reqPath, "/")
		if filePath == "" {
			filePath = "index.html"
		}

		// 尝试读取请求的文件
		content, err := fs.ReadFile(subFS, filePath)
		if err != nil {
			// 文件不存在，返回 index.html (SPA 路由回退)
			content, err = fs.ReadFile(subFS, "index.html")
			if err != nil {
				c.String(500, "Failed to load index.html")
				return
			}
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.Data(200, "text/html; charset=utf-8", content)
			return
		}

		// 设置 Content-Type
		contentType := "text/html; charset=utf-8"
		switch path.Ext(filePath) {
		case ".js":
			contentType = "application/javascript; charset=utf-8"
		case ".css":
			contentType = "text/css; charset=utf-8"
		case ".json":
			contentType = "application/json; charset=utf-8"
		case ".png":
			contentType = "image/png"
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".svg":
			contentType = "image/svg+xml"
		case ".ico":
			contentType = "image/x-icon"
		case ".woff":
			contentType = "font/woff"
		case ".woff2":
			contentType = "font/woff2"
		case ".ttf":
			contentType = "font/ttf"
		}
		c.Header("Content-Type", contentType)
		c.Data(200, contentType, content)
	})
}
