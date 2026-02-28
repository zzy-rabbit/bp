package api

import (
	"time"

	"github.com/zzy-rabbit/xtools/xplugin"
)

const (
	PluginName = "bp.tool.timer"
)

type Job func()

type Task struct {
	Name       string
	Spec       string
	Plugin     string
	EntryID    int
	NextRun    time.Time
	CreateTime time.Time
}

type IPlugin interface {
	xplugin.IPlugin
}
