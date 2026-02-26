package api

import (
	"github.com/zzy-rabbit/xtools/xlog"
	"github.com/zzy-rabbit/xtools/xplugin"
)

const (
	PluginName = "bp.tool.log"
)

type IPlugin interface {
	xplugin.IPlugin
	xlog.ILogger
}
