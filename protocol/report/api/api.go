package api

import (
	"context"
	"github.com/zzy-rabbit/xtools/xerror"
	"github.com/zzy-rabbit/xtools/xplugin"
)

const (
	PluginName = "bp.protocol.report"
)

type IPlugin interface {
	xplugin.IPlugin
	Broadcast(ctx context.Context, tag uint32, data interface{}) xerror.IError
}
