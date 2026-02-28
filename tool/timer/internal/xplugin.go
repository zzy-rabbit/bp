package internal

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"

	logApi "github.com/zzy-rabbit/bp/tool/log/api"
	"github.com/zzy-rabbit/bp/tool/timer/api"
)

type service struct {
	ILogger logApi.IPlugin `xplugin:"bp.tool.log"`
	cron    *cron.Cron
	mutex   sync.RWMutex
	taskMap map[string]*api.Task
}

func New(ctx context.Context) api.IPlugin {
	return &service{
		taskMap: make(map[string]*api.Task),
	}
}

func (s *service) GetName(ctx context.Context) string {
	return api.PluginName
}

func (s *service) Init(ctx context.Context, initParam string) error {
	c := cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		),
	)
	s.cron = c
	s.ILogger.Info(ctx, "plugin %s init success", s.GetName(ctx))
	return nil
}

func (s *service) Run(ctx context.Context, runParam string) error {
	s.cron.Start()
	s.ILogger.Info(ctx, "plugin %s run success", s.GetName(ctx))
	return nil
}

func (s *service) Stop(ctx context.Context, stopParam string) error {
	stopCtx := s.cron.Stop()
	select {
	case <-stopCtx.Done():
		s.ILogger.Info(ctx, "plugin %s stop success", s.GetName(ctx))
		return nil
	case <-ctx.Done():
		s.ILogger.Info(ctx, "plugin %s stop fail %v", s.GetName(ctx), ctx.Err())
		return ctx.Err()
	}
}
