package internal

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/zzy-rabbit/bp/protocol/websocket/api"
	"github.com/zzy-rabbit/xtools/xerror"
)

type client struct {
	api.IConn
	*service
}

func (s *service) NewClient(ctx context.Context, url string) (api.IClient, xerror.IError) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		s.ILogger.Error(ctx, "websocket dial at url %s error: %v", url, err)
		return nil, xerror.Extend(xerror.ErrInternalError, err.Error())
	}
	return &client{
		IConn:   s.NewConnection(ctx, conn),
		service: s,
	}, nil
}
