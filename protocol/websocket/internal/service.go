package internal

import (
	"context"
	"github.com/zzy-rabbit/bp/protocol/websocket/api"
	"github.com/zzy-rabbit/xtools/xerror"
)

func (s *service) ConnTo(ctx context.Context, url string) (api.IClient, xerror.IError) {
	return s.NewClient(ctx, url)
}

func (s *service) ListenAt(ctx context.Context, addr string, callback api.OnConnCallbackFunc) (api.IServer, xerror.IError) {
	return s.NewServer(ctx, addr, callback), nil
}
