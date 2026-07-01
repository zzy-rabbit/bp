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
	"sync"
)

type service struct {
	network  model.Network
	fiberApp *fiber.App
	ILogger  logApi.IPlugin `xplugin:"bp.tool.log"`

	mutex    sync.RWMutex
	configs  []func(ctx context.Context, fiberConfig *fiber.Config)
	handlers []func(ctx context.Context, fiberApp *fiber.App)
}

func New(ctx context.Context) api.IPlugin {
	return &service{}
}

func (s *service) GetName(ctx context.Context) string {
	return api.PluginName
}

func (s *service) Init(ctx context.Context, initParam string) xerror.IError {
	var network model.Network
	err := json.Unmarshal([]byte(initParam), &network)
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "plugin %s init fail %v", s.GetName(ctx), err)
		return xerror.Extend(xerror.ErrInvalidParam, "init param invalid")
	}

	s.network = network
	s.ILogger.Info(ctx, "plugin %s init success", s.GetName(ctx))
	return nil
}

func (s *service) Run(ctx context.Context, runParam string) xerror.IError {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	config := fiber.Config{}
	for _, configFunc := range s.configs {
		configFunc(ctx, &config)
	}

	s.fiberApp = fiber.New(config)

	for _, handlerFunc := range s.handlers {
		handlerFunc(ctx, s.fiberApp)
	}

	s.registerMiddlewares()

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

func (s *service) Stop(ctx context.Context, stopParam string) xerror.IError {
	err := s.fiberApp.Shutdown()
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "plugin %s stop %s fail %v", s.GetName(ctx), stopParam, err)
		return xerror.Extend(xerror.ErrInternalError, "fiber shutdown fail")
	}
	s.ILogger.Info(ctx, "plugin %s stop success", s.GetName(ctx))
	return nil
}
