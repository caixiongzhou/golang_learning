package utils

import (
	"blog-system/internal/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUtil struct {
	secretKey string
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWTUtil(cfg *config.Config) *JWTUtil {
	return &JWTUtil{
		secretKey: cfg.JWTSecret,
	}
}

// GenerateToken 生成JWT token
func (j *JWTUtil) GenerateToken(userID uint, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "blog-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken 验证JWT token
func (j *JWTUtil) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		// 🔥 只使用密钥验证签名，不查询数据库或缓存
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// 🔥 验证通过后直接从Token解析用户信息，不查询存储
		return claims, nil
	}

	return nil, errors.New("无效的token")
}

// 说明：JWT 工具封装了 token 的生成与验证逻辑。
// - 使用 jwt.RegisteredClaims 设定标准字段（过期时间、签发者等）。
// - 生成的 token 使用 HMAC-SHA256 签名，业务侧应安全保存 secretKey 并定期轮换。
