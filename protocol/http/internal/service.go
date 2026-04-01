package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zzy-rabbit/xtools/xerror"
)

func (s *service) ParseQueryParams(ctx *fiber.Ctx, header, query any) xerror.IError {
	if header != nil {
		err := ctx.ReqHeaderParser(header)
		if xerror.Error(err) {
			return xerror.Extend(xerror.ErrInvalidParam, "parse request header fail "+err.Error())
		}
	}
	if query != nil {
		err := ctx.QueryParser(query)
		if xerror.Error(err) {
			return xerror.Extend(xerror.ErrInvalidParam, "parse request query fail "+err.Error())
		}
	}
	return nil
}

func (s *service) ParseBodyParams(ctx *fiber.Ctx, header, body any) xerror.IError {
	if header != nil {
		err := ctx.ReqHeaderParser(header)
		if xerror.Error(err) {
			return xerror.Extend(xerror.ErrInvalidParam, "parse request header fail "+err.Error())
		}
	}
	if body != nil {
		err := ctx.BodyParser(body)
		if xerror.Error(err) {
			return xerror.Extend(xerror.ErrInvalidParam, "parse request body fail "+err.Error())
		}
	}
	return nil
}

func (s *service) Register(r func(fiberApp *fiber.App)) {
	r(s.fiberApp)
}
