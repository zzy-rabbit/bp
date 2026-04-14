package internal

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func (s *service) registerRouter(ctx context.Context, fiberApp *fiber.App) {
	group := fiberApp.Group(
		s.config.BaseURL,
		s.IHttp.TimingMiddleware(func(ctx *fiber.Ctx) bool {
			return true
		}),
		func(c *fiber.Ctx) error {
			err := c.Next() // 先执行 tus
			c.Response().Header.Del("Access-Control-Allow-Origin")
			c.Set("Access-Control-Allow-Origin", "*")
			c.Set("Access-Control-Allow-Methods", "GET,POST,HEAD,PATCH,DELETE,OPTIONS")
			c.Set("Access-Control-Allow-Headers", "*")
			c.Set("Access-Control-Expose-Headers", "Location,Upload-Offset,Upload-Length,Tus-Resumable")
			return err
		},
		adaptor.HTTPMiddleware(s.Tus.Middleware),
	)
	group.Post("", adaptor.HTTPHandlerFunc(s.Tus.PostFile))
	group.Head(":id", adaptor.HTTPHandlerFunc(s.Tus.HeadFile))
	group.Patch(":id", adaptor.HTTPHandlerFunc(s.Tus.PatchFile))
	if !s.Tus.Config.DisableDownload {
		group.Get(":id", adaptor.HTTPHandlerFunc(s.Tus.GetFile))
	}
	if s.Tus.Config.StoreComposer.UsesTerminater && !s.Tus.Config.DisableTermination {
		group.Delete(":id", adaptor.HTTPHandlerFunc(s.Tus.DelFile))
	}
}
