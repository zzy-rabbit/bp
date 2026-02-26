package api

import (
	"context"
	uniformApi "github.com/zzy-rabbit/bp/tool/uniform/api"
	"github.com/zzy-rabbit/xtools/xerror"
	"github.com/zzy-rabbit/xtools/xplugin"
	"net"
	"net/http"
	"time"
)

const (
	PluginName = "bp.protocol.websocket"
)

type IPlugin interface {
	xplugin.IPlugin
}

type Config struct {
}

type OnReceiveCallbackFunc func(ctx context.Context, frame uniformApi.Frame)

type IConn interface {
	RemoteAddr(ctx context.Context) net.Addr
	Send(ctx context.Context, frame uniformApi.Frame, timeout time.Duration) (uniformApi.Frame, xerror.IError)
	Post(ctx context.Context, frame uniformApi.Frame) xerror.IError
	SetCallback(ctx context.Context, callback OnReceiveCallbackFunc)
}

type IClient interface {
	IConn
}

type OnConnCallbackFunc func(ctx context.Context, conn IConn, req Request)

type IServer interface {
}

type Request struct {
	Headers http.Header
}

type ITransport interface {
	xplugin.IPlugin
	ConnTo(ctx context.Context, url string) (IClient, xerror.IError)
	ListenAt(ctx context.Context, addr string, callback OnConnCallbackFunc) (IServer, xerror.IError)
}
