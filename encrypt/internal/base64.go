package internal

import (
	"context"
	"encoding/base64"
	"github.com/zzy-rabbit/bp/encrypt/api"
	"github.com/zzy-rabbit/xtools/xerror"
)

type Base64Encoder struct {
}

func (s *service) NewBase64Encoder(ctx context.Context) api.IProcessor {
	return &Base64Encoder{}
}

func (enc *Base64Encoder) Process(ctx context.Context, plaintext []byte) ([]byte, xerror.IError) {
	cipherBuffer := make([]byte, base64.StdEncoding.EncodedLen(len(plaintext)))
	base64.StdEncoding.Encode(cipherBuffer, plaintext)
	return cipherBuffer, nil
}

type Base64Decoder struct {
}

func (s *service) NewBase64Decoder(ctx context.Context) api.IProcessor {
	return &Base64Decoder{}
}

func (enc *Base64Decoder) Process(ctx context.Context, ciphertext []byte) ([]byte, xerror.IError) {
	plainBuffer := make([]byte, base64.StdEncoding.DecodedLen(len(ciphertext)))
	n, err := base64.StdEncoding.Decode(plainBuffer, ciphertext)
	if err != nil {
		return nil, xerror.Extend(xerror.ErrInvalidParam, err.Error())
	}
	return plainBuffer[:n], nil
}
