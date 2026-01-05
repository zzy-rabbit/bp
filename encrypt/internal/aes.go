package internal

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/zzy-rabbit/bp/encrypt/api"
	"github.com/zzy-rabbit/xtools/xerror"
	"io"
)

type AESEncoder struct {
	key []byte
}

func (s *service) NewAESEncoder(ctx context.Context, key []byte) api.IProcessor {
	return &AESEncoder{key: []byte(key)}
}

func (enc *AESEncoder) Process(ctx context.Context, plaintext []byte) ([]byte, xerror.IError) {
	block, err := aes.NewCipher(enc.key)
	if err != nil {
		return nil, xerror.Extend(xerror.ErrInternalError, err.Error())
	}

	// 填充明文以满足块大小
	plaintext = pkcs7Pad(plaintext, aes.BlockSize)

	// IV需要是唯一的，但不一定是保密的
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, xerror.Extend(xerror.ErrInternalError, err.Error())
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

type AESDecoder struct {
	key []byte
}

func (s *service) NewAESDecoder(ctx context.Context, key []byte) api.IProcessor {
	return &AESDecoder{key: []byte(key)}
}

func (enc *AESDecoder) Process(ctx context.Context, ciphertext []byte) ([]byte, xerror.IError) {
	block, err := aes.NewCipher(enc.key)
	if err != nil {
		return nil, xerror.Extend(xerror.ErrInternalError, err.Error())
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, xerror.Extend(xerror.ErrInvalidParam, "ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, xerror.Extend(xerror.ErrInvalidParam, "ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// 去除填充
	ciphertext = pkcs7UnPad(ciphertext)
	return ciphertext, nil
}

// PKCS7填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7去除填充
func pkcs7UnPad(data []byte) []byte {
	length := len(data)
	unPadding := int(data[length-1])
	return data[:(length - unPadding)]
}
