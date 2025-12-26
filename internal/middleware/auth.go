package middleware

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	Username string `json:"username"`
	Role	 string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateAuthToken(username,role string ) (	string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &claims{
		Username: username,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //non hashed token

	return token.SignedString([]byte(jwtSecret)) //hashed token


}

func ValidateAuthToken(tokenString string) (*claims, error) {
	claims := &claims{}

	 token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil

}
