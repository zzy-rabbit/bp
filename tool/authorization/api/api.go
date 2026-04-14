package api

import (
	"context"
	"github.com/zzy-rabbit/xtools/xerror"
	"github.com/zzy-rabbit/xtools/xplugin"
)

const (
	PluginName = "bp.tool.authorization"
)

type Config struct {
	SecretKeyLength int `json:"secret_key_length"`
}

type IPlugin interface {
	xplugin.IPlugin
	GenerateToken(ctx context.Context, plaintext []byte) ([]byte, xerror.IError)
	ParseToken(ctx context.Context, ciphertext []byte) ([]byte, xerror.IError)
}
