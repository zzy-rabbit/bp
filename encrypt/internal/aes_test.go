package internal

import (
	"context"
	"testing"
)

func TestAES(t *testing.T) {
	ctx := context.Background()
	s := New(ctx)
	encoder := s.NewAESEncoder(ctx, []byte("1234567812345678"))
	decoder := s.NewAESDecoder(ctx, []byte("1234567812345678"))
	plaintext := []byte("hello world")
	ciphertext, err := encoder.Process(ctx, plaintext)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ciphertext: %+v", ciphertext)
	plaintext2, err := decoder.Process(ctx, ciphertext)
	if err != nil {
		t.Fatal(err)
	}
	if string(plaintext) != string(plaintext2) {
		t.Fatal("plaintext not equal")
	}
}
