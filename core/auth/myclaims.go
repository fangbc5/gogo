package auth

import (
	"errors"
	"github.com/fangbc5/gogo/constant"
	"github.com/fangbc5/gogo/core/config"
	"github.com/fangbc5/gogo/core/db"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type MyClaims struct {
	Username string
	jwt.RegisteredClaims
}

func VerifyToken(tokenString string) (*MyClaims, error) {
	tokenString = strings.ReplaceAll(tokenString, constant.Bearer, "")
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Get().Auth.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		//校验成功应该刷新token过期时间
		db.RedisCache("expire", constant.TokenKey+tokenString, config.Get().Auth.TokenLife)
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
