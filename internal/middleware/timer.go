package middleware

import (
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"time"
)

func (m *Middleware) Timer(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()

		next(ctx)

		duration := time.Since(start)
		log.Printf("Method: %s, URL: %s, Time: %v", ctx.Method(), ctx.URI(), duration)
	}
}
