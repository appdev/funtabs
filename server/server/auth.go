package server

import (
	"funtabs-server/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type loginReq struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// Login POST /api/login
// 接收 FormData：username, password
// 返回 data: { token: "..." }
func Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBind(&req); err != nil {
		fail(c, 400, "用户名或密码不能为空")
		return
	}

	var user model.User
	if err := model.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, 400, "用户名或密码错误")
		} else {
			fail(c, 500, "服务器错误")
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		fail(c, 400, "用户名或密码错误")
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		fail(c, 500, "生成 token 失败")
		return
	}

	ok(c, gin.H{"token": token})
}

type registerReq struct {
	Username string `form:"username" binding:"required,min=2,max=32"`
	Password string `form:"password" binding:"required,min=6"`
}

// Register POST /api/register
// 接收 FormData：username, password
func Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBind(&req); err != nil {
		fail(c, 400, "用户名至少 2 位，密码至少 6 位")
		return
	}

	// 检查用户名是否已存在
	var count int64
	model.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		fail(c, 400, "用户名已被占用")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		fail(c, 500, "服务器错误")
		return
	}

	user := model.User{
		Username: req.Username,
		Password: string(hash),
	}
	if err = model.DB.Create(&user).Error; err != nil {
		fail(c, 500, "注册失败，请重试")
		return
	}

	ok(c, "注册成功")
}
