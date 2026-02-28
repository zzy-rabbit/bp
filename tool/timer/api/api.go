package api

import (
	"context"
	"github.com/zzy-rabbit/xtools/xerror"
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
	Register(ctx context.Context, name string, spec string, job Job) xerror.IError
	Unregister(ctx context.Context, name string)
}
