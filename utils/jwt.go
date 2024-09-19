package utils

import (
	user_models "ginchat/models/user_basic"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

type CustomClaims struct {
	user_models.UserBasic
	jwt.StandardClaims
}

func GenerateJWT(user user_models.UserBasic) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute) // 设置过期时间

	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := viper.GetString("jwt.key")
	tokenString, err := token.SignedString([]byte(key))
	return tokenString, err
}

// 验证 JWT
func ValidateJWT(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return viper.GetString("jwt.key"), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
