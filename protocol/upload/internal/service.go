package internal

import (
	"context"
	"github.com/tus/tusd/pkg/handler"
	"github.com/zzy-rabbit/bp/protocol/upload/api"
	"github.com/zzy-rabbit/xtools/xerror"
)

func (t *Tus) SetNotifyCreatedCallback(ctx context.Context, callback api.NotifyCreatedCallback) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.NotifyCreatedCallback = callback
}

func (t *Tus) SetNotifyCompletedCallback(ctx context.Context, callback api.NotifyCompletedCallback) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.NotifyCompletedCallback = callback
}

func (t *Tus) SetNotifyTerminatedCallback(ctx context.Context, callback api.NotifyTerminatedCallback) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.NotifyTerminatedCallback = callback
}

func (t *Tus) SetNotifyProgressChangedCallback(ctx context.Context, callback api.NotifyProgressChangedCallback) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.NotifyProgressChangedCallback = callback
}

func eventCallback(callbacks ...func(context.Context, handler.HookEvent) error) func(context.Context, handler.HookEvent) error {
	return func(ctx context.Context, event handler.HookEvent) error {
		for _, cb := range callbacks {
			if cb != nil {
				err := cb(ctx, event)
				if xerror.Error(err) {
					return err
				}
			}
		}
		return nil
	}
}

func (t *Tus) SetPreCreateCallback(ctx context.Context, callback api.PreCreateCallback) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.PreCreateCallback = eventCallback(t.PreCreateCallback, callback)
}

func (t *Tus) SetPreCompleteCallback(ctx context.Context, callback api.PreCompleteCallback) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.PreCompleteCallback = eventCallback(t.PreCompleteCallback, callback)
}
