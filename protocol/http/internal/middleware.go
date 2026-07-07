package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/zzy-rabbit/xtools/xcontext"
	"net/http"
	"strconv"
	"strings"
)

func (s *service) CORSMiddleware(ignores ...func(ctx *fiber.Ctx) bool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		for _, ignore := range ignores {
			if ignore != nil && ignore(ctx) {
				return ctx.Next()
			}
		}
		return cors.New(cors.Config{
			AllowOrigins:     "*",
			AllowCredentials: false,
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
		})(ctx)
	}
}

func (s *service) TimingMiddleware(ignores ...func(ctx *fiber.Ctx) bool) fiber.Handler {
	const format = "%s %s %s\nrequestBody: %s\nresponseBody %s"
	return func(ctx *fiber.Ctx) error {
		userCtx := xcontext.Background()
		ctx.SetUserContext(userCtx)

		for _, ignore := range ignores {
			if ignore != nil && ignore(ctx) {
				_ = ctx.Next()
				s.ILogger.Info(userCtx, "HTTP_COST %s %s %s %v", ctx.IP(), ctx.Method(), ctx.Path(), xcontext.Since(userCtx))
				return nil
			}
		}

		var reqBody []byte
		if ctx.Method() == http.MethodGet || ctx.Method() == http.MethodHead || ctx.Method() == http.MethodOptions {
			reqBody = ctx.Request().URI().QueryString()
		} else {
			reqBody = ctx.Body()
		}

		s.ILogger.Info(userCtx, "HTTP_REQUEST "+format, ctx.IP(), ctx.Method(), ctx.Path(), reqBody, "")

		_ = ctx.Next()

		var respBody []byte
		respBody = ctx.Response().Body()
		s.ILogger.Info(userCtx, "HTTP_RESPONSE "+format+" "+strconv.Itoa(ctx.Response().StatusCode()), ctx.IP(), ctx.Method(), ctx.Path(), reqBody, respBody)
		s.ILogger.Info(userCtx, "HTTP_COST %s %s %s %v", ctx.IP(), ctx.Method(), ctx.Path(), xcontext.Since(userCtx))
		return nil
	}
}

func (s *service) registerMiddlewares() {
	//middlewares := []any{"/", s.corsMiddleware(), s.timingMiddleware()}
	//s.rootApp.Use(middlewares...)
}
