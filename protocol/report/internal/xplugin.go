package internal

import (
	"context"
	"github.com/zzy-rabbit/bp/protocol/report/api"
	websocketApi "github.com/zzy-rabbit/bp/protocol/websocket/api"
	logApi "github.com/zzy-rabbit/bp/tool/log/api"
	uniformApi "github.com/zzy-rabbit/bp/tool/uniform/api"
	"sync"
)

type service struct {
	ILogger    logApi.IPlugin       `xplugin:"bp.tool.log"`
	IUniform   uniformApi.IPlugin   `xplugin:"bp.tool.uniform"`
	IWebsocket websocketApi.IPlugin `xplugin:"bp.protocol.websocket"`
	IServer    websocketApi.IServer

	mutex   sync.RWMutex
	connMap map[string]websocketApi.IConn
}

func New(ctx context.Context) api.IPlugin {
	return &service{}
}

func (s *service) GetName(ctx context.Context) string {
	return api.PluginName
}

func (s *service) Init(ctx context.Context, initParam string) error {
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
