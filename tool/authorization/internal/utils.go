package internal

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func generateRandomKey(ctx context.Context, length int) []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	key := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		key = append(key, byte(r.Intn(256)))
	}
	fmt.Printf("key: %s\n", key)
	return key
}
