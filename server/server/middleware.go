package server

import (
	"errors"
	"funtabs-server/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	UserID uint `json:"uid"`
	jwt.RegisteredClaims
}

// generateToken 为指定用户 ID 签发 JWT
func generateToken(userID uint) (string, error) {
	expire := time.Duration(config.Cfg.JWT.Expire) * 24 * time.Hour
	claims := jwtClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWT.Secret))
}

// parseToken 验证并解析 JWT，返回 claims
func parseToken(tokenStr string) (*jwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("签名算法不匹配")
		}
		return []byte(config.Cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token 无效")
	}
	return claims, nil
}

// AuthRequired 是 JWT 鉴权中间件，从 token 请求头读取并验证 JWT
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("token")
		if tokenStr == "" || tokenStr == "undefined" {
			authFail(c)
			c.Abort()
			return
		}
		claims, err := parseToken(tokenStr)
		if err != nil {
			authFail(c)
			c.Abort()
			return
		}
		// 将用户 ID 注入到上下文，后续 handler 通过 c.GetUint("userID") 获取
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
