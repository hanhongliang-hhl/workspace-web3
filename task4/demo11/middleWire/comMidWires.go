package middleWire

import (
	"net/http"
	"strings"
	"web3Demo/task4/demo11/jwtToken"
	"github.com/gin-gonic/gin"
)

//jwt 验证中间件
func AuthMiddlewir() gin.HandlerFunc{  
	return func(c *gin.Context) { 
		authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证信息"})
            c.Abort()
            return
        }
        
        // Bearer token 格式: "Bearer <token>"
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式错误"})
            c.Abort()
            return
        }
        
        claims, err := jwtToken.ValidateJWT(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
            c.Abort()
            return
        }
        
        // 将用户信息存储在上下文中供后续使用
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Next()
	}
}