package api

import (
	"context"
	"github.com/zzy-rabbit/xtools/xerror"
	"github.com/zzy-rabbit/xtools/xplugin"
)

const PluginName = "xtools.plugins.uniform"

type Config struct {
}

type Header struct {
	Sequence      uint64 `json:"sequence"`
	Authorization string `json:"authorization"`
	Timestamp     uint64 `json:"timestamp"`
	Tag           uint32 `json:"tag"`
}

type FrameHead struct {
	Delimiter  uint32
	Priority   byte
	Type       byte
	Version    byte
	Format     byte
	Encryption byte
	Reserve    [8]byte
	CheckSum   uint32
}

type Frame struct {
	FrameHead
	Header
	Data []byte
}

const (
	FormatTypeBinary = iota + 1
	FormatTypeJson
)

const (
	EncryptTypePlainText = iota + 1
	EncryptTypeAES
)

type IPlugin interface {
	xplugin.IPlugin
	NewFrame(ctx context.Context) Frame
	Marshal(ctx context.Context, frame *Frame) ([]byte, xerror.IError)
	Unmarshal(ctx context.Context, data []byte) (Frame, xerror.IError)
}
