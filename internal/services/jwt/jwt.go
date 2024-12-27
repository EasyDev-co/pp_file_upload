package jwt

import (
	"EasyDev-co/pp_file_upload/internal/consts"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type JWTService struct {
	signingKey string
}

func NewJWTService(signingKey string) *JWTService {
	return &JWTService{signingKey: signingKey}
}

func (j *JWTService) GenerateJWT(userId string) (string, error) {
	claims := jwt.MapClaims{
		consts.UserIdKey: userId,
		"exp":            time.Now().Add(time.Minute * 30).Unix(),
		"token_type":     "access",
		"jti":            uuid.New().String(),
		"iat":            time.Now().Unix(),
	}

	// Создаём токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретом
	tokenString, err := token.SignedString([]byte(j.signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTService) ParseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(j.signingKey), nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}
