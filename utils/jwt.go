package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Claims represents the JWT claims structure.
type Claims struct {
	UserID uint   `json:"id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func VerifyJWT(tokenString string, secretKey []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid JWT token")
	}

	// Check if the token has expired
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("JWT token has expired")
	}

	return claims, nil
}
