package response

import "github.com/valyala/fasthttp"

// RespondError отправляет ошибку с заданным статус-кодом и сообщением.
func RespondError(ctx *fasthttp.RequestCtx, statusCode int, message string) {
	ctx.SetStatusCode(statusCode)
	ctx.SetBodyString(message)
}

// RespondSuccess отправляет успешный ответ с заданным сообщением.
func RespondSuccess(ctx *fasthttp.RequestCtx, message string) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(message)
}
