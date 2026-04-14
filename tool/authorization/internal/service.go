package internal

import (
	"context"
	"github.com/zzy-rabbit/xtools/xerror"
)

func (s *service) GenerateToken(ctx context.Context, plaintext []byte) ([]byte, xerror.IError) {
	ciphertext, err := s.encoder.Process(ctx, plaintext)
	if err != nil {
		s.ILogger.Error(ctx, "encode plaintext %+v fail %v", plaintext, err)
		return nil, xerror.Extend(xerror.ErrInternalError, "encode plaintext fail: %v", err)
	}
	return ciphertext, nil
}

func (s *service) ParseToken(ctx context.Context, ciphertext []byte) ([]byte, xerror.IError) {
	plaintext, err := s.decoder.Process(ctx, ciphertext)
	if err != nil {
		s.ILogger.Error(ctx, "decode ciphertext %+v fail %v", ciphertext, err)
		return nil, xerror.Extend(xerror.ErrInternalError, "decode plaintext fail: %v", err)
	}
	return plaintext, nil
}
