package internal

import (
	"context"
	"encoding/json"
	"github.com/zzy-rabbit/bp/log/api"
	"github.com/zzy-rabbit/xtools/xerror"
	"github.com/zzy-rabbit/xtools/xlog"
)

type service struct {
	config xlog.Config
	xlog.ILogger
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
		return err
	}
	logger, err := xlog.New(ctx, s.config)
	if xerror.Error(err) {
		return err
	}
	s.ILogger = logger
	s.ILogger.Info(ctx, "plugin %s init success", s.GetName(ctx))
	return nil
}

func (s *service) Run(ctx context.Context, runParam string) error {
	s.ILogger.Info(ctx, "plugin %s run success", s.GetName(ctx))
	return nil
}

func (s *service) Stop(ctx context.Context, stopParam string) error {
	s.ILogger.Info(ctx, "plugin %s stop success", s.GetName(ctx))
	return nil
}
