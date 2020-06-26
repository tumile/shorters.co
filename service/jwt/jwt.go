package jwt

import (
	"fmt"
	"os"
	"shorters/domain"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService interface {
	Sign(user domain.User) (string, error)
	Parse(token string) (domain.User, error)
}

type jwtService struct {
	signingKey []byte
}

func NewJWTService() JWTService {
	return &jwtService{[]byte(os.Getenv("key"))}
}

func (j *jwtService) Sign(user domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Audience:  "shorters.co",
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   user.Email,
	})
	return token.SignedString(j.signingKey)
}

func (j *jwtService) Parse(tokenString string) (domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.signingKey, nil
	})
	var user domain.User
	if err != nil {
		return user, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid && claims.VerifyAudience("shorters.co", true) {
		user.Email = claims["sub"].(string)
	}
	return user, nil
}
