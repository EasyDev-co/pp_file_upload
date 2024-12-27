package middleware

import "EasyDev-co/pp_file_upload/internal/services"

type Middleware struct {
	jwtService     services.JWTService
	allowedOrigins []string
}

func NewMiddleware(jwtService services.JWTService, allowedOrigins []string) *Middleware {
	return &Middleware{
		jwtService:     jwtService,
		allowedOrigins: allowedOrigins,
	}
}
