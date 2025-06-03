package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/taninchot-work/backend-challenge/internal/core/config"
	"log"
	"time"
)

type JwtInterface interface {
	GenerateJwt(userId string) (string, error)
	ValidateJwt(tokenString string) (*JwtClaim, error)
}

type JwtClaim struct {
	UserId string `json:"UserId"`
	jwt.RegisteredClaims
}

func GenerateJwt(userId string) (string, error) {
	var jwtSecret = []byte(config.GetConfig().RestServer.Jwt.Secret)
	var jwtExpireIn = config.GetConfig().RestServer.Jwt.ExpireIn
	claim := JwtClaim{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			Issuer:    config.GetConfig().RestServer.Jwt.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtExpireIn) * time.Millisecond)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}
	return tokenString, nil
}

func ValidateJwt(tokenString string) (*JwtClaim, error) {
	jwtSecret := []byte(config.GetConfig().RestServer.Jwt.Secret)

	token, err := jwt.ParseWithClaims(tokenString, &JwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		log.Println("Error parsing token:", err)
		return nil, err
	}

	claim, ok := token.Claims.(*JwtClaim)
	if !ok || !token.Valid {
		log.Println("Invalid token claims")
		return nil, jwt.ErrInvalidKey
	}

	now := time.Now()
	if claim.ExpiresAt.Time.Before(now) {
		log.Println("Token has expired")
		return nil, jwt.ErrTokenExpired
	}
	if claim.IssuedAt.Time.After(now) {
		log.Println("Token is not yet valid (issued in future)")
		return nil, jwt.ErrTokenNotValidYet
	}
	if claim.NotBefore.Time.After(now) {
		log.Println("Token is not yet valid (nbf claim)")
		return nil, jwt.ErrTokenNotValidYet
	}

	return claim, nil
}
