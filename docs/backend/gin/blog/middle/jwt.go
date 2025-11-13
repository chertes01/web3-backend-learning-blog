package middle

import (
	"net/http"
	"strings"
	"time"

	"blog/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type JWTclaims struct {
	ID uint
	jwt.RegisteredClaims
}

// 定义用于签名和验证 JWT 的密钥
var jwtKey = []byte("insert_your_own_secret_key")

// GenerateToken 根据用户 ID 生成 JWT Token
func GenerateToken(userID uint, c *gin.Context) (string, error) {
	claims := JWTclaims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "your_app_name",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		config.Log.Error("系统错误：JWT 签名失败 ", err)
		return "", err
	}
	return tokenString, nil
}

// JWTAuthMiddleware 是一个 Gin 中间件函数，用于验证请求中的 JWT Token
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取 Authorization Header，标准格式是"Authorization: Bearer <token>"
		AuthHeader := c.GetHeader("Authorization")

		if AuthHeader == "" {
			// [日志] 记录未提供 token 的信息
			config.Log.WithFields(logrus.Fields{
				"ip": c.ClientIP(),
			}).Warn("请求失败：未提供 token")
			// 如果没有提供 Authorization Header，则返回 401 未授权错误
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供 token"})
			// 中止请求，后续的中间件和目标路由处理函数（Handler）都不会被执行
			c.Abort()
			return
		}
		//提取 Token 字符串
		tokenStr := strings.TrimPrefix(AuthHeader, "Bearer ")

		// jwt.MapClaims 是 map[string]interface{} 的别名，用于存储键值对数据
		claims := jwt.MapClaims{}
		// 解析和验证 Token
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		//检查验证结果
		// !token.Valid 确保 Token 最终被标记为“有效”
		if err != nil || !token.Valid {
			// 必须处理 err 为 nil 的情况，防止 panic
			var errMsg string
			if err != nil {
				// 当 err 不为 nil 时，使用其错误信息
				errMsg = err.Error()
			} else {
				// 当 err 为 nil 时，说明 Token 无效但没有具体错误信息
				errMsg = "Token 无效或已过期"
			}
			// [日志] 记录 Token 验证失败的信息
			config.Log.WithFields(logrus.Fields{
				"ip":    c.ClientIP(),
				"error": errMsg, // 使用安全的变量
			}).Warn("Token 验证失败")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 token"})
			// 中止请求，后续的中间件和目标路由处理函数（Handler）都不会被执行
			c.Abort()
			return
		}
		// 从 claims 中提取用户 ID 并设置到 Gin 的上下文中，供后续处理函数使用,并转换为 uint 类型
		if id, ok := claims["ID"].(float64); ok {
			c.Set("user_id", uint(id))
		} else {
			// 防御性编程：如果 Token 里没有 ID 字段
			config.Log.Error("Token 解析异常：缺少 ID 字段")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 格式错误"})
			c.Abort()
			return
		}

		c.Next()
	}
}
