package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/zzy-rabbit/bp/protocol/http/api"
	logApi "github.com/zzy-rabbit/bp/tool/log/api"
	"github.com/zzy-rabbit/xtools/xerror"
)

type service struct {
	config   api.Config
	rootApp  *fiber.App
	proxyApp *fiber.App
	ILogger  logApi.IPlugin `xplugin:"bp.tool.log"`

	mutex            sync.RWMutex
	configCallbacks  []func(ctx context.Context, fiberConfig *fiber.Config)
	handlerCallbacks []func(ctx context.Context, fiberApp *fiber.App)
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
		return xerror.Extend(xerror.ErrInvalidParam, "init param invalid")
	}
	s.ILogger.Info(ctx, "plugin %s init success", s.GetName(ctx))
	return nil
}

func (s *service) Run(ctx context.Context, runParam string) xerror.IError {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	config := fiber.Config{}
	for _, configFunc := range s.configCallbacks {
		configFunc(ctx, &config)
	}

	s.rootApp = fiber.New(config)

	for _, handlerFunc := range s.handlerCallbacks {
		handlerFunc(ctx, s.rootApp)
	}

	s.registerMiddlewares()

	// ==================== HTTPS 配置 ====================
	if s.config.Https.Enable && s.config.Https.Domain != "" {
		s.ILogger.Info(ctx, "plugin %s starting with https", s.GetName(ctx))

		// 使用 Certbot 生成的证书路径
		certPath := fmt.Sprintf("%s/live/%s/fullchain.pem",
			s.config.Https.Cert,
			s.config.Https.Domain,
		)
		keyPath := fmt.Sprintf("%s/live/%s/privkey.pem",
			s.config.Https.Cert,
			s.config.Https.Domain,
		)

		// 启动 HTTPS
		go func() {
			httpsAddr := fmt.Sprintf("%s:%d", s.config.Https.Host, s.config.Https.Port)
			err := s.rootApp.ListenTLS(httpsAddr, certPath, keyPath)
			if err != nil {
				s.ILogger.Error(ctx, "plugin %s https run at %s fail %v", s.GetName(ctx), httpsAddr, err)
			}
		}()

		// HTTP 跳转到 HTTPS
		if s.proxyApp == nil {
			s.proxyApp = fiber.New()
		}
		s.proxyApp.All("*", func(c *fiber.Ctx) error {
			host := c.Hostname()
			if host == "" {
				host = c.Get("Host")
			}
			httpsURL := "https://" + host + c.OriginalURL()
			return c.Redirect(httpsURL, fiber.StatusMovedPermanently)
		})

		go func() {
			httpAddr := fmt.Sprintf("%s:%d", s.config.Http.Host, s.config.Http.Port)
			if err := s.proxyApp.Listen(httpAddr); err != nil {
				s.ILogger.Error(ctx, "http redirect server failed at %s: %v", httpAddr, err)
			}
		}()
	} else {
		// 普通 HTTP
		go func() {
			addr := fmt.Sprintf("%s:%d", s.config.Http.Host, s.config.Http.Port)
			if err := s.rootApp.Listen(addr); err != nil {
				s.ILogger.Error(ctx, "plugin %s run %s at addr %s fail %v",
					s.GetName(ctx), runParam, addr, err)
			}
		}()
	}

	s.ILogger.Info(ctx, "plugin %s run success", s.GetName(ctx))
	return nil
}

func (s *service) Stop(ctx context.Context, stopParam string) xerror.IError {
	// 关闭 rootApp（HTTPS 或 普通 HTTP）
	if s.rootApp != nil {
		if err := s.rootApp.Shutdown(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.ILogger.Error(ctx, "root app shutdown fail %v", err)
			return xerror.Extend(xerror.ErrInternalError, "fiber shutdown root app fail %v", err)
		}
	}
	// 如果开启了 HTTPS，则额外关闭 proxyApp（HTTP 重定向）
	if s.config.Https.Enable && s.proxyApp != nil {
		if err := s.proxyApp.Shutdown(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.ILogger.Error(ctx, "proxy app shutdown fail %v", err)
			return xerror.Extend(xerror.ErrInternalError, "fiber shutdown proxy app fail %v", err)
		}
	}
	s.ILogger.Info(ctx, "plugin %s stop success", s.GetName(ctx))
	return nil
}
