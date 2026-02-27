package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/zzy-rabbit/bp/model"
	"github.com/zzy-rabbit/bp/protocol/http/api"
	logApi "github.com/zzy-rabbit/bp/tool/log/api"
	"github.com/zzy-rabbit/xtools/xerror"
)

type service struct {
	network  model.Network
	fiberApp *fiber.App
	ILogger  logApi.IPlugin `xplugin:"bp.tool.log"`
}

func New(ctx context.Context) api.IPlugin {
	return &service{}
}

func (s *service) GetName(ctx context.Context) string {
	return api.PluginName
}

func (s *service) Init(ctx context.Context, initParam string) error {
	var network model.Network
	err := json.Unmarshal([]byte(initParam), &network)
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "plugin %s init fail %v", s.GetName(ctx), err)
		return err
	}
	s.fiberApp = fiber.New()
	s.registerMiddlewares()

	s.network = network
	s.ILogger.Info(ctx, "plugin %s init success", s.GetName(ctx))
	return nil
}

func (s *service) Run(ctx context.Context, runParam string) error {
	go func() {
		addr := fmt.Sprintf("%s:%d", s.network.Host, s.network.Port)
		err := s.fiberApp.Listen(addr)
		if xerror.Error(err) {
			s.ILogger.Error(ctx, "plugin %s run %s at addr %s fail %v", s.GetName(ctx), runParam, addr, err)
			return
		}
	}()
	s.ILogger.Info(ctx, "plugin %s run success", s.GetName(ctx))
	return nil
}

func (s *service) Stop(ctx context.Context, stopParam string) error {
	err := s.fiberApp.Shutdown()
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "plugin %s stop %s fail %v", s.GetName(ctx), stopParam, err)
		return err
	}
	s.ILogger.Info(ctx, "plugin %s stop success", s.GetName(ctx))
	return nil
}
