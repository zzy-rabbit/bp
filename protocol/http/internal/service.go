package internal

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/zzy-rabbit/bp/model"
	"github.com/zzy-rabbit/xtools/xerror"
)

func (s *service) ParseQueryParams(ctx *fiber.Ctx, query any) error {
	header := model.Header{}
	err := ctx.ReqHeaderParser(&header)
	if xerror.Error(err) {
		return ctx.JSON(&model.HttpResponse{
			IError: xerror.Extend(xerror.ErrInvalidParam, "parse request header fail"),
			Data:   json.RawMessage{},
		})
	}
	err = ctx.QueryParser(query)
	if xerror.Error(err) {
		return ctx.JSON(&model.HttpResponse{
			IError: xerror.Extend(xerror.ErrInvalidParam, "parse request query fail"),
		})
	}
	return nil
}

func (s *service) ParseBodyParams(ctx *fiber.Ctx, body any) error {
	header := model.Header{}
	err := ctx.ReqHeaderParser(&header)
	if xerror.Error(err) {
		return ctx.JSON(&model.HttpResponse{
			IError: xerror.Extend(xerror.ErrInvalidParam, "parse request header fail"),
			Data:   json.RawMessage{},
		})
	}
	err = ctx.BodyParser(body)
	if xerror.Error(err) {
		return ctx.JSON(&model.HttpResponse{
			IError: xerror.Extend(xerror.ErrInvalidParam, "parse request body fail"),
		})
	}
	return nil
}

func (s *service) Register(r func(fiberApp *fiber.App)) {
	r(s.fiberApp)
}
