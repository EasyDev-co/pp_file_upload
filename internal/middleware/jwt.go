package middleware

import (
	"EasyDev-co/pp_file_upload/internal/consts"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (m *Middleware) JWT(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		authHeader := string(ctx.Request.Header.Peek("Authorization"))

		if !isHeaderValid(authHeader) {
			log.Errorf("Header isn't valid.")
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := m.jwtService.ParseJWT(tokenString)
		if err != nil {
			log.Errorf("Token parse error: %v", err)
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		}
		if token == nil {
			log.Errorf("Token is nil.")
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx.SetUserValue(consts.UserIdKey, claims[consts.UserIdKey])
		}
		next(ctx)
	}
}

func isHeaderValid(authHeader string) bool {
	return authHeader != "" && strings.HasPrefix(authHeader, "Bearer ")
}
