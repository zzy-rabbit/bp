package internal

import (
	"context"
	"github.com/zzy-rabbit/bp/tool/encrypt/api"
	logApi "github.com/zzy-rabbit/bp/tool/log/api"
	"github.com/zzy-rabbit/xtools/xerror"
	"sync/atomic"
)

type service struct {
	seq     atomic.Uint64
	ILogger logApi.IPlugin `xplugin:"bp.tool.log"`
}

func New(ctx context.Context) api.IPlugin {
	return &service{}
}

func (s *service) GetName(ctx context.Context) string {
	return api.PluginName
}

func (s *service) Init(ctx context.Context, initParam string) xerror.IError {
	s.ILogger.Info(ctx, "plugin %s init success", s.GetName(ctx))
	return nil
}

func (s *service) Run(ctx context.Context, runParam string) xerror.IError {
	s.ILogger.Info(ctx, "plugin %s run success", s.GetName(ctx))
	return nil
}

func (s *service) Stop(ctx context.Context, stopParam string) xerror.IError {
	s.ILogger.Info(ctx, "plugin %s stop success", s.GetName(ctx))
	return nil
}
