package internal

import (
	"context"
	"testing"
)

func TestBase64(t *testing.T) {
	ctx := context.Background()
	s := New(ctx)
	encoder := s.NewBase64Encoder(ctx)
	decoder := s.NewBase64Decoder(ctx)
	plaintext := []byte("hello world")
	ciphertext, err := encoder.Process(ctx, plaintext)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ciphertext: %s", ciphertext)
	plaintext, err = decoder.Process(ctx, ciphertext)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("plaintext: %s", plaintext)
	if string(plaintext) != "hello world" {
		t.Fatal("plaintext not equal")
	}
}
