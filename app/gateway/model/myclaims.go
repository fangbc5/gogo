package model

import (
	"github.com/golang-jwt/jwt/v4"
)

type MyClaims struct {
	Username string
	jwt.RegisteredClaims
}
