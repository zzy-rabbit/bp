package api

import (
	"context"
	"github.com/zzy-rabbit/xtools/xerror"
	"github.com/zzy-rabbit/xtools/xplugin"
)

const (
	PluginName = "bp.tool.mail"
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type IPlugin interface {
	xplugin.IPlugin
	Send(ctx context.Context, message Message) xerror.IError
}
