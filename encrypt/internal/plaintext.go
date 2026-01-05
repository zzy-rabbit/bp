package internal

import (
	"context"
	"github.com/zzy-rabbit/bp/encrypt/api"
	"github.com/zzy-rabbit/xtools/xerror"
)

type PlaintextEncoder struct {
}

func (s *service) NewPlainTextEncoder(ctx context.Context) api.IProcessor {
	return &PlaintextEncoder{}
}

func (enc *PlaintextEncoder) Process(ctx context.Context, plaintext []byte) ([]byte, xerror.IError) {
	return plaintext, nil
}

type PlaintextDecoder struct {
}

func (s *service) NewPlainTextDecoder(ctx context.Context) api.IProcessor {
	return &PlaintextDecoder{}
}

func (enc *PlaintextDecoder) Process(ctx context.Context, ciphertext []byte) ([]byte, xerror.IError) {
	return ciphertext, nil
}
