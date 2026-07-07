package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/zzy-rabbit/bp/model"
	"github.com/zzy-rabbit/xtools/xerror"
	"github.com/zzy-rabbit/xtools/xplugin"
)

const (
	PluginName = "bp.protocol.http"
)

type HttpsConfig struct {
	model.Network
	Enable  bool     `json:"enable"`
	Domains []string `json:"domains"`
}

type Config struct {
	Http  model.Network `json:"http"`
	Https HttpsConfig   `json:"https"`
}

type IPlugin interface {
	xplugin.IPlugin
	SetConfig(ctx context.Context, r func(ctx context.Context, fiberConfig *fiber.Config))
	Register(ctx context.Context, r func(ctx context.Context, fiberApp *fiber.App))
	ParseQueryParams(ctx *fiber.Ctx, header, query any) xerror.IError
	ParseBodyParams(ctx *fiber.Ctx, header, body any) xerror.IError
	CORSMiddleware(ignores ...func(ctx *fiber.Ctx) bool) fiber.Handler
	TimingMiddleware(ignores ...func(ctx *fiber.Ctx) bool) fiber.Handler
}
