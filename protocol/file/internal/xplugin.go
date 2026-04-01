package internal

import (
	"context"
	"encoding/json"
	"github.com/zzy-rabbit/bp/protocol/file/api"
	httpApi "github.com/zzy-rabbit/bp/protocol/http/api"
	logApi "github.com/zzy-rabbit/bp/tool/log/api"
	"github.com/zzy-rabbit/xtools/xerror"
)

type service struct {
	ILogger logApi.IPlugin  `xplugin:"bp.tool.log"`
	IHttp   httpApi.IPlugin `xplugin:"bp.protocol.http"`
	config  api.Config
	*Tus
}

func New(ctx context.Context) api.IPlugin {
	return &service{}
}

func (s *service) GetName(ctx context.Context) string {
	return api.PluginName
}

func (s *service) Init(ctx context.Context, initParam string) error {
	err := json.Unmarshal([]byte(initParam), &s.config)
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "plugin %s init fail %v", s.GetName(ctx), err)
		return err
	}

	tusHandler, xerr := s.NewTusHandler(ctx)
	if xerror.Error(xerr) {
		s.ILogger.Error(ctx, "plugin %s init tus handler fail %v", s.GetName(ctx), err)
		return xerr
	}
	s.Tus = tusHandler
	s.ILogger.Info(ctx, "plugin %s init tus handler by config %+v success", s.GetName(ctx), s.config)

	s.IHttp.Register(s.registerRouter)

	s.ILogger.Info(ctx, "plugin %s init success", s.GetName(ctx))
	return nil
}

func (s *service) Run(ctx context.Context, runParam string) error {
	s.Tus.startMonitor(ctx)
	s.ILogger.Info(ctx, "plugin %s run success", s.GetName(ctx))
	return nil
}

func (s *service) Stop(ctx context.Context, stopParam string) error {
	s.ILogger.Info(ctx, "plugin %s stop success", s.GetName(ctx))
	return nil
}
