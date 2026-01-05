package api

import (
	"context"
	"github.com/zzy-rabbit/xtools/xerror"
	"github.com/zzy-rabbit/xtools/xplugin"
)

const PluginName = "bp.encrypt"

const (
	EncryptTypePlainText = iota + 1
	EncryptTypeAES
	EncryptTypeBase64
)

type IProcessor interface {
	Process(ctx context.Context, content []byte) ([]byte, xerror.IError)
}

type IWorkflow interface {
	IProcessor
}

type Config struct {
}

type IPlugin interface {
	xplugin.IPlugin
	Workflow(ctx context.Context, processors ...IProcessor) IProcessor
	NewPlainTextEncoder(ctx context.Context) IProcessor
	NewPlainTextDecoder(ctx context.Context) IProcessor
	NewAESEncoder(ctx context.Context, key []byte) IProcessor
	NewAESDecoder(ctx context.Context, key []byte) IProcessor
	NewBase64Encoder(ctx context.Context) IProcessor
	NewBase64Decoder(ctx context.Context) IProcessor
}
