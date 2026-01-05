package api

import (
	"github.com/zzy-rabbit/xtools/xlog"
	"github.com/zzy-rabbit/xtools/xplugin"
)

const (
	PluginName = "bp.log"
)

type IPlugin interface {
	xplugin.IPlugin
	xlog.ILogger
}
