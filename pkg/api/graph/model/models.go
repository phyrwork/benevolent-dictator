package model

import "github.com/golang-jwt/jwt"

type UserClaims struct {
	UserID int `json:"userId"`
	jwt.StandardClaims
}
