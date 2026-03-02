package server

import (
	"funtabs-server/model"
	"funtabs-server/storage"

	"github.com/gin-gonic/gin"
)

// UploadWallpaper POST /api/uploadWallpaper
func UploadWallpaper(c *gin.Context) {
	uploadFile(c, "wallpapers")
}

// UploadFavicon POST /api/uploadFavicon
func UploadFavicon(c *gin.Context) {
	uploadFile(c, "favicons")
}

// UploadAvatar POST /api/uploadAvatar
// 上传成功后同时更新数据库中的用户头像字段
func UploadAvatar(c *gin.Context) {
	userID := c.GetUint("userID")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		fail(c, 400, "请选择要上传的文件")
		return
	}
	defer file.Close()

	url, err := storage.S.Save(file, header, "avatars")
	if err != nil {
		fail(c, 500, "上传失败："+err.Error())
		return
	}

	// 删除旧头像（忽略错误）
	var user model.User
	if model.DB.First(&user, userID).Error == nil && user.Avatar != "" {
		_ = storage.S.Delete(user.Avatar)
	}

	// 更新头像 URL
	model.DB.Model(&model.User{}).Where("id = ?", userID).Update("avatar", url)

	ok(c, gin.H{"url": url})
}

// uploadFile 是通用文件上传逻辑，返回 data: { url: "..." }
func uploadFile(c *gin.Context, dir string) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		fail(c, 400, "请选择要上传的文件")
		return
	}
	defer file.Close()

	url, err := storage.S.Save(file, header, dir)
	if err != nil {
		fail(c, 500, "上传失败："+err.Error())
		return
	}

	ok(c, gin.H{"url": url})
}
