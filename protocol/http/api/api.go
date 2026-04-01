package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zzy-rabbit/xtools/xplugin"
)

const (
	PluginName = "bp.protocol.http"
)

type IPlugin interface {
	xplugin.IPlugin
	Register(func(fiberApp *fiber.App))
	ParseQueryParams(ctx *fiber.Ctx, query any) error
	ParseBodyParams(ctx *fiber.Ctx, body any) error
	CORSMiddleware(ignores ...func(ctx *fiber.Ctx) bool) fiber.Handler
	TimingMiddleware(ignores ...func(ctx *fiber.Ctx) bool) fiber.Handler
}
