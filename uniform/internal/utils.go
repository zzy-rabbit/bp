package internal

import (
	"bytes"
	"context"
	"encoding/binary"
	encryptApi "github.com/zzy-rabbit/bp/encrypt/api"
	"github.com/zzy-rabbit/xtools/xerror"
	"hash/crc32"
)

func CheckSum(data []byte) uint32 {
	return crc32.Checksum(data, crc32.MakeTable(crc32.Koopman))
}

func calculateCheckSum(data []byte) ([]byte, xerror.IError) {
	sum := CheckSum(data)
	buf := bytes.NewBuffer(make([]byte, 0, 4))
	err := binary.Write(buf, binary.BigEndian, sum)
	if err != nil {
		return nil, xerror.Extend(xerror.ErrInternalError, "write check sum error")
	}
	return buf.Bytes(), nil
}

func (s *service) Encode(ctx context.Context, typ int, content []byte) ([]byte, xerror.IError) {
	switch typ {
	case encryptApi.EncryptTypePlainText:
		return s.IEncrypt.NewPlainTextEncoder(ctx).Process(ctx, content)
	case encryptApi.EncryptTypeAES:
		return s.IEncrypt.NewAESEncoder(ctx, []byte("1234567890123456")).Process(ctx, content)
	case encryptApi.EncryptTypeBase64:
		return s.IEncrypt.NewBase64Encoder(ctx).Process(ctx, content)
	}
	return nil, xerror.Extend(xerror.ErrInvalidParam, "invalid encrypt type")
}

func (s *service) Decode(ctx context.Context, typ int, content []byte) ([]byte, xerror.IError) {
	switch typ {
	case encryptApi.EncryptTypePlainText:
		return s.IEncrypt.NewPlainTextDecoder(ctx).Process(ctx, content)
	case encryptApi.EncryptTypeAES:
		return s.IEncrypt.NewAESDecoder(ctx, []byte("1234567890123456")).Process(ctx, content)
	case encryptApi.EncryptTypeBase64:
		return s.IEncrypt.NewBase64Decoder(ctx).Process(ctx, content)
	}
	return nil, xerror.Extend(xerror.ErrInvalidParam, "invalid encrypt type")
}
