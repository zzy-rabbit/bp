package api

import (
	"github.com/zzy-rabbit/xtools/xlog"
	"github.com/zzy-rabbit/xtools/xplugin"
)

const (
	PluginName = "xtools.plugins.log"
)

type IPlugin interface {
	xplugin.IPlugin
	xlog.ILogger
}
