package middleware

import (
	"EasyDev-co/pp_file_upload/internal/consts"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func JWT(next fasthttp.RequestHandler, signingKey string) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		authHeader := string(ctx.Request.Header.Peek("Authorization"))

		if !isHeaderValid(authHeader) {
			log.Errorf("Header isn't valid.")
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		log.Infof("Token: %s", tokenString)

		// Парсим токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte(signingKey), nil
		})

		log.Infof("Token is valid: %v", token.Valid)

		if err != nil || !token.Valid {
			log.Errorf("Error: %v", err)
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Infof("Claims: %v", claims)
			ctx.SetUserValue(consts.UserIdKey, claims[consts.UserIdKey])
		}
		log.Infof("Authorized!")
		next(ctx)
	}
}

func isHeaderValid(authHeader string) bool {
	return authHeader != "" && strings.HasPrefix(authHeader, "Bearer ")
}
