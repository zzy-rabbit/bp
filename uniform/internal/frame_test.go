package internal

import (
	"context"
	"github.com/zzy-rabbit/bp/uniform/api"
	"log"
	"reflect"
	"testing"
)

func TestFrame(t *testing.T) {
	ctx := context.Background()
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

	s := New(ctx)
	originFrame := s.NewFrame(ctx)
	originFrame.Format = api.FormatTypeBinary
	originFrame.Data = []byte{6, 7, 8, 9}
	originFrame.Encryption = api.EncryptTypeAES
	t.Logf("originFrame: %+v", originFrame)
	bytes, err := s.Marshal(ctx, &originFrame)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("bytes: %v", bytes)
	t.Logf("bytes: %s", bytes)

	parseFrame, err := s.Unmarshal(ctx, bytes)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("parseFrame: %+v", parseFrame)

	originFrame.FrameHead.CheckSum = parseFrame.FrameHead.CheckSum
	if !reflect.DeepEqual(originFrame, parseFrame) {
		t.Fatal("parseFrame != originFrame")
	}
}
