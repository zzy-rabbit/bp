package internal

import (
	"context"
	"encoding/json"
	logApi "github.com/zzy-rabbit/bp/tool/log/api"
	"github.com/zzy-rabbit/bp/tool/mail/api"
	"github.com/zzy-rabbit/xtools/xerror"
)

type service struct {
	config  api.Config
	ILogger logApi.IPlugin `xplugin:"bp.tool.log"`
}

func New(ctx context.Context) api.IPlugin {
	return &service{}
}

func (s *service) GetName(ctx context.Context) string {
	return api.PluginName
}

func (s *service) Init(ctx context.Context, initParam string) xerror.IError {
	err := json.Unmarshal([]byte(initParam), &s.config)
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "plugin %s init fail %v", s.GetName(ctx), err)
		return xerror.Extend(xerror.ErrInvalidParam, err.Error())
	}
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
