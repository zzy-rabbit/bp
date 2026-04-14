package internal

import (
	"context"
	"encoding/json"
	httpApi "github.com/zzy-rabbit/bp/protocol/http/api"
	"github.com/zzy-rabbit/bp/protocol/upload/api"
	logApi "github.com/zzy-rabbit/bp/tool/log/api"
	"github.com/zzy-rabbit/xtools/xerror"
	"github.com/zzy-rabbit/xtools/xfile"
	"os"
	"sync"
)

type service struct {
	ILogger logApi.IPlugin  `xplugin:"bp.tool.log"`
	IHttp   httpApi.IPlugin `xplugin:"bp.protocol.http"`
	config  api.Config
	*Tus
	cancel    context.CancelFunc
	busyMutex sync.RWMutex
	busyFiles map[string]*fileSync
}

func New(ctx context.Context) api.IPlugin {
	return &service{busyFiles: make(map[string]*fileSync)}
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

	if !xfile.IsExist(ctx, s.config.RootPath) {
		s.ILogger.Info(ctx, "plugin %s tus root path %s not exist, try to create", s.GetName(ctx), s.config.RootPath)
		err = os.MkdirAll(s.config.RootPath, os.ModePerm)
		if xerror.Error(err) {
			s.ILogger.Error(ctx, "plugin %s tus create root path %s fail %v", s.GetName(ctx), s.config.RootPath, err)
			return err
		}
	}

	tusHandler, xerr := s.NewTusHandler(ctx)
	if xerror.Error(xerr) {
		s.ILogger.Error(ctx, "plugin %s init tus handler fail %v", s.GetName(ctx), err)
		return xerr
	}
	s.Tus = tusHandler
	s.ILogger.Info(ctx, "plugin %s init tus handler by config %+v success", s.GetName(ctx), s.config)

	s.IHttp.Register(ctx, s.registerRouter)

	s.ILogger.Info(ctx, "plugin %s init success", s.GetName(ctx))
	return nil
}

func (s *service) Run(ctx context.Context, runParam string) error {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.Tus.startEventMonitor(ctx)
	s.startExpireMonitor(ctx)
	s.ILogger.Info(ctx, "plugin %s run success", s.GetName(ctx))
	return nil
}

func (s *service) Stop(ctx context.Context, stopParam string) error {
	s.cancel()
	s.ILogger.Info(ctx, "plugin %s stop success", s.GetName(ctx))
	return nil
}
