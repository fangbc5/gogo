package auth

import (
	"errors"
	"strings"

	"github.com/fangbc5/gogo/constant"
	// cache "github.com/fangbc5/gogo/core/cache/redis"
	"github.com/fangbc5/gogo/core/config/consul"

	"github.com/golang-jwt/jwt/v4"
)

type MyClaims struct {
	Username string
	jwt.RegisteredClaims
}

func VerifyToken(tokenString string) (*MyClaims, error) {
	tokenString = strings.ReplaceAll(tokenString, constant.Bearer, "")
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(consul.Get().Auth.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		//校验成功应该刷新token过期时间
		// cache.RedisCache("expire", constant.TokenKey+tokenString, consul.Get().Auth.TokenLife)
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
