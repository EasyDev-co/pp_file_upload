package response

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

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

// RespondSuccessJSON отправляет успешный JSON-ответ
func RespondSuccessJSON(ctx *fasthttp.RequestCtx, data interface{}) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")

	jsonData, err := json.Marshal(data)
	if err != nil {
		RespondError(ctx, fasthttp.StatusInternalServerError, "Failed to serialize response")
		return
	}

	ctx.SetBody(jsonData)
}
