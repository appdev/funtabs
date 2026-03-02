package server

import "github.com/gin-gonic/gin"

// ok 返回成功响应，data 会放在 data 字段
func ok(c *gin.Context, data any) {
	c.JSON(200, gin.H{
		"code": 200,
		"data": data,
		"msg":  "",
	})
}

// fail 返回业务错误响应，HTTP 状态码始终为 200（前端通过 code 字段判断）
func fail(c *gin.Context, code int, msg string) {
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
	})
}

// authFail 返回 token 失效响应，前端收到 code=-3 后会清除本地 token 并刷新页面
func authFail(c *gin.Context) {
	fail(c, -3, "身份验证失败，请重新登录")
}
