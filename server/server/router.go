package server

import (
	"funtabs-server/model"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewRouter 创建并配置 Gin 路由
func NewRouter(origins string) *gin.Engine {
	r := gin.Default()

	// CORS 配置
	corsConfig := cors.DefaultConfig()
	if origins == "*" {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = strings.Split(origins, ",")
	}
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "token")
	r.Use(cors.New(corsConfig))

	// 静态文件服务（本地存储时使用）
	r.Static("/uploads", "./uploads")

	// 公开接口
	api := r.Group("/api")
	{
		api.POST("/login", Login)
		api.POST("/register", Register)
	}

	// 需要登录的接口
	auth := r.Group("/api", AuthRequired())
	{
		// 用户信息
		auth.GET("/getUserInfo", GetUserInfo)
		auth.POST("/changeUserInfo", ChangeUserInfo)
		auth.POST("/changePassword", ChangePassword)
		auth.POST("/deleteUserAccount", DeleteUserAccount)

		// 数据同步
		auth.POST("/saveData", SaveData)
		auth.GET("/getData", GetData)

		// 文件上传
		auth.POST("/uploadWallpaper", UploadWallpaper)
		auth.POST("/uploadFavicon", UploadFavicon)
		auth.POST("/uploadAvatar", UploadAvatar)
	}

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		ok(c, gin.H{"status": "ok", "db": model.DB != nil})
	})

	return r
}
