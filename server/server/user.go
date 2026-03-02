package server

import (
	"funtabs-server/model"
	"funtabs-server/storage"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetUserInfo GET /api/getUserInfo
// 返回 data: { id, username, avatar }
func GetUserInfo(c *gin.Context) {
	userID := c.GetUint("userID")

	var user model.User
	if err := model.DB.First(&user, userID).Error; err != nil {
		authFail(c)
		return
	}

	ok(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"avatar":   user.Avatar,
	})
}

type changeUserInfoReq struct {
	Username string `json:"username" binding:"required,min=2,max=32"`
}

// ChangeUserInfo POST /api/changeUserInfo
func ChangeUserInfo(c *gin.Context) {
	userID := c.GetUint("userID")

	var req changeUserInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "用户名不合法")
		return
	}

	// 检查新用户名是否被其他人占用
	var count int64
	model.DB.Model(&model.User{}).
		Where("username = ? AND id != ?", req.Username, userID).
		Count(&count)
	if count > 0 {
		fail(c, 400, "用户名已被占用")
		return
	}

	if err := model.DB.Model(&model.User{}).Where("id = ?", userID).
		Update("username", req.Username).Error; err != nil {
		fail(c, 500, "更新失败")
		return
	}

	ok(c, "更新成功")
}

type changePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword POST /api/changePassword
func ChangePassword(c *gin.Context) {
	userID := c.GetUint("userID")

	var req changePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误，新密码至少 6 位")
		return
	}

	var user model.User
	if err := model.DB.First(&user, userID).Error; err != nil {
		authFail(c)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		fail(c, 400, "原密码错误")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		fail(c, 500, "服务器错误")
		return
	}

	model.DB.Model(&model.User{}).Where("id = ?", userID).Update("password", string(hash))
	ok(c, "密码修改成功")
}

// DeleteUserAccount POST /api/deleteUserAccount
func DeleteUserAccount(c *gin.Context) {
	userID := c.GetUint("userID")

	// 删除用户数据
	model.DB.Where("user_id = ?", userID).Delete(&model.UserData{})

	// 删除头像文件（如果有）
	var user model.User
	if model.DB.First(&user, userID).Error == nil && user.Avatar != "" {
		_ = storage.S.Delete(user.Avatar)
	}

	// 删除用户记录
	model.DB.Delete(&model.User{}, userID)

	ok(c, "账号已注销")
}
