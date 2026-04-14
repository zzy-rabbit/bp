package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/zzy-rabbit/xtools/xerror"
	"github.com/zzy-rabbit/xtools/xplugin"
)

const (
	PluginName = "bp.protocol.http"
)

type IPlugin interface {
	xplugin.IPlugin
	Register(ctx context.Context, r func(ctx context.Context, fiberApp *fiber.App))
	ParseQueryParams(ctx *fiber.Ctx, header, query any) xerror.IError
	ParseBodyParams(ctx *fiber.Ctx, header, body any) xerror.IError
	CORSMiddleware(ignores ...func(ctx *fiber.Ctx) bool) fiber.Handler
	TimingMiddleware(ignores ...func(ctx *fiber.Ctx) bool) fiber.Handler
}
