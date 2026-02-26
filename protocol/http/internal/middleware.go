package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/zzy-rabbit/xtools/xcontext"
	"log"
	"net/http"
	"strings"
)

func (s *service) corsMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods: strings.Join([]string{
			http.MethodGet,
			http.MethodPost,
			http.MethodHead,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		}, ","),
		AllowHeaders:  "*",
		ExposeHeaders: "",
		MaxAge:        0,
	})
}

func (s *service) timingMiddleware() fiber.Handler {
	ignore := func() bool {

		return false
	}
	const format = "%s %s %s\nrequestBody: %s\nresponseBody %s"

	return func(ctx *fiber.Ctx) error {
		userCtx := xcontext.Background()
		ctx.SetUserContext(userCtx)

		var reqBody []byte
		if ignore() {
			reqBody = []byte("[ignore request body]")
		} else if ctx.Method() == http.MethodGet || ctx.Method() == http.MethodHead || ctx.Method() == http.MethodOptions {
			reqBody = ctx.Request().URI().QueryString()
		} else {
			reqBody = ctx.Body()
		}

		log.Printf("HTTP_REQUEST "+format, ctx.IP(), ctx.Method(), ctx.Path(), reqBody, "")

		_ = ctx.Next()

		var respBody []byte
		if ignore() {
			respBody = []byte("[ignore response body]")
		} else {
			respBody = ctx.Response().Body()
		}
		log.Printf("HTTP_RESPONSE "+format, ctx.IP(), ctx.Method(), ctx.Path(), reqBody, respBody)
		log.Printf("HTTP_COST %s %s %v", ctx.IP(), ctx.Method(), xcontext.Since(userCtx))
		return nil
	}
}

func (s *service) registerMiddlewares() {
	middlewares := []any{"/", s.corsMiddleware, s.timingMiddleware}
	s.fiberApp.Use(middlewares...)
}
