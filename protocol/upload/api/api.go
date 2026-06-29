package api

import (
	"context"
	"github.com/tus/tusd/pkg/handler"
	"github.com/zzy-rabbit/xtools/xerror"
	"github.com/zzy-rabbit/xtools/xplugin"
	"io"
)

const (
	PluginName = "bp.protocol.upload"
)

type Config struct {
	RootPath string `json:"root_path"`
	BaseURL  string `json:"base_url"`
	Expire   int    `json:"expire"`
	Interval int    `json:"interval"`
	MaxSize  int    `json:"max_size"`
}

type FileInfo struct {
	ID   string `json:"id"`
	Size int    `json:"size"`
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
}

// NotifyCreatedCallback 事件通知 创建成功
type NotifyCreatedCallback func(ctx context.Context, event handler.HookEvent)

// NotifyCompletedCallback 事件通知 完成
type NotifyCompletedCallback func(ctx context.Context, event handler.HookEvent)

// NotifyTerminatedCallback 事件通知 删除
type NotifyTerminatedCallback func(ctx context.Context, event handler.HookEvent)

// NotifyProgressChangedCallback 事件通知 进度变更
type NotifyProgressChangedCallback func(ctx context.Context, event handler.HookEvent)

// PreCreateCallback 钩子函数 可拦截 创建前
type PreCreateCallback func(ctx context.Context, event handler.HookEvent) error

// PreCompleteCallback 钩子函数 可拦截 删除前
type PreCompleteCallback func(ctx context.Context, event handler.HookEvent) error

type IPlugin interface {
	xplugin.IPlugin

	SetNotifyCreatedCallback(ctx context.Context, callback NotifyCreatedCallback)
	SetNotifyCompletedCallback(ctx context.Context, callback NotifyCompletedCallback)
	SetNotifyTerminatedCallback(ctx context.Context, callback NotifyTerminatedCallback)
	SetNotifyProgressChangedCallback(ctx context.Context, callback NotifyProgressChangedCallback)
	SetPreCreateCallback(ctx context.Context, callback PreCreateCallback)
	SetPreCompleteCallback(ctx context.Context, callback PreCompleteCallback)

	FileLock(ctx context.Context, id string)
	FileUnlock(ctx context.Context, id string)
	FileRLock(ctx context.Context, id string)
	FileRUnlock(ctx context.Context, id string)
	IsFileLocked(ctx context.Context, id string) bool

	GetFileInfo(ctx context.Context, id string) (FileInfo, xerror.IError)
	MoveFile(ctx context.Context, id string, path string) xerror.IError
	CopyFile(ctx context.Context, id string, w io.Writer) (FileInfo, xerror.IError)
	DeleteFile(ctx context.Context, id string) xerror.IError
}
