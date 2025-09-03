package jwtToken

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


// JWT 配置
var jwtKey = []byte("task4") // 实际项目中应使用环境变量

// Claims 定义 JWT 中包含的信息
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}


// GenerateJWT 生成 JWT token
func GenerateJWT(userID uint, username string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour) // token 24小时过期
    
    claims := &Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "blog-app",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

// ValidateJWT 验证 JWT token
func ValidateJWT(tokenString string) (*Claims, error) {
    claims := &Claims{}
    
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    
    return claims, nil
}