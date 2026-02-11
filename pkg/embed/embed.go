package embedfs

import (
	"embed"
	"io/fs"
	"net/http"

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

	// 使用 StaticFS 提供静态文件服务
	// 注意：这会处理所有静态文件，包括 index.html
	r.StaticFS("/", http.FS(subFS))
}
