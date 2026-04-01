package internal

import (
	"context"
	"github.com/tus/tusd/pkg/filestore"
	"github.com/tus/tusd/pkg/handler"
	"github.com/zzy-rabbit/bp/protocol/file/api"
	logApi "github.com/zzy-rabbit/bp/tool/log/api"
	"github.com/zzy-rabbit/xtools/xerror"
	"sync"
)

type Tus struct {
	ILogger logApi.IPlugin `xplugin:"bp.tool.log"`
	*handler.Handler
	handler.Config
	cancel context.CancelFunc
	mutex  sync.RWMutex

	NotifyCreatedCallback         api.NotifyCreatedCallback
	NotifyCompletedCallback       api.NotifyCompletedCallback
	NotifyTerminatedCallback      api.NotifyTerminatedCallback
	NotifyProgressChangedCallback api.NotifyProgressChangedCallback
	PreCreateCallback             api.PreCreateCallback
	PreCompleteCallback           api.PreCompleteCallback
}

func (s *service) NewTusHandler(ctx context.Context) (*Tus, xerror.IError) {
	ctx, cancel := context.WithCancel(ctx)

	store := filestore.FileStore{
		Path: s.config.RootPath,
	}

	composer := handler.NewStoreComposer()
	store.UseIn(composer)

	tusHandler := &Tus{
		ILogger: s.ILogger,
		cancel:  cancel,
	}
	tusHandler.Config = handler.Config{
		BasePath:                s.config.BaseURL,
		StoreComposer:           composer,
		NotifyCompleteUploads:   true,
		NotifyUploadProgress:    true,
		NotifyTerminatedUploads: true,
		NotifyCreatedUploads:    true,
		PreUploadCreateCallback: func(hook handler.HookEvent) error {
			s.ILogger.Info(ctx, "tusd handler pre create uploads %+v", hook)
			tusHandler.mutex.RLock()
			defer tusHandler.mutex.RUnlock()
			if tusHandler.PreCreateCallback != nil {
				return tusHandler.PreCreateCallback(ctx, hook)
			}
			return nil
		},
		PreFinishResponseCallback: func(hook handler.HookEvent) error {
			s.ILogger.Info(ctx, "tusd handler pre complete uploads %+v", hook)
			tusHandler.mutex.RLock()
			defer tusHandler.mutex.RUnlock()
			if tusHandler.PreCompleteCallback != nil {
				return tusHandler.PreCompleteCallback(ctx, hook)
			}
			return nil
		},
	}

	tus, err := handler.NewHandler(tusHandler.Config)
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "tusd handler base path %s init fail %v", s.config.BaseURL, err)
		return nil, xerror.Extend(xerror.ErrInternalError, err.Error())
	}
	tusHandler.Handler = tus

	return tusHandler, nil
}

func (t *Tus) startMonitor(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				t.ILogger.Info(ctx, "tusd handler stop")
				return
			case event := <-t.UploadProgress:
				t.ILogger.Info(ctx, "tusd handler progress changed %+v", event)
				t.mutex.RLock()
				if t.NotifyProgressChangedCallback != nil {
					go t.NotifyProgressChangedCallback(ctx, event)
				}
				t.mutex.RUnlock()
			case event := <-t.CompleteUploads:
				t.ILogger.Info(ctx, "tusd handler complete uploads %+v", event)
				t.mutex.RLock()
				if t.NotifyCompletedCallback != nil {
					t.NotifyCompletedCallback(ctx, event)
				}
				t.mutex.RUnlock()
			case event := <-t.TerminatedUploads:
				t.ILogger.Info(ctx, "tusd handler terminated uploads %+v", event)
				t.mutex.RLock()
				if t.NotifyTerminatedCallback != nil {
					t.NotifyTerminatedCallback(ctx, event)
				}
				t.mutex.RUnlock()
			case event := <-t.CreatedUploads:
				t.ILogger.Info(ctx, "tusd handler created uploads %+v", event)
				t.mutex.RLock()
				if t.NotifyCreatedCallback != nil {
					t.NotifyCreatedCallback(ctx, event)
				}
				t.mutex.RUnlock()
			}
		}
	}()
}

func (t *Tus) Stop(ctx context.Context) {
	t.cancel()
}
