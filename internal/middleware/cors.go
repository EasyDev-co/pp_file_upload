package middleware

import (
	"github.com/valyala/fasthttp"
)

func (m *Middleware) CORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Добавляем CORS-заголовки
		origin := string(ctx.Request.Header.Peek("Origin"))

		if isValidOrigin(origin, m.allowedOrigins) {
			ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
			ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
			ctx.Response.Header.Set("Access-Control-Allow-Headers", "authorization, content-type")
			ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		}

		// Если это preflight-запрос (OPTIONS), сразу возвращаем 204
		if string(ctx.Method()) == fasthttp.MethodOptions {
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}

		next(ctx)
	}
}

func isValidOrigin(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}
