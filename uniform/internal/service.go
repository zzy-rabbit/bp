package internal

import (
	"context"
	"github.com/zzy-rabbit/bp/uniform/api"
	"github.com/zzy-rabbit/xtools/xerror"
	"time"
)

func (s *service) NextSequence() uint64 {
	return s.seq.Add(1)
}

func (s *service) NewFrame(ctx context.Context) api.Frame {
	return api.Frame{
		FrameHead: api.FrameHead{
			Delimiter:  Delimiter,
			Priority:   0,
			Type:       0,
			Version:    0x01,
			Format:     api.FormatTypeJson,
			Encryption: api.EncryptTypePlainText,
			Reserve:    [8]byte{},
			CheckSum:   0,
		},
		Header: api.Header{
			Sequence:  s.NextSequence(),
			Timestamp: uint64(time.Now().UnixMilli()),
		},
	}
}

func (s *service) Marshal(ctx context.Context, frame *api.Frame) ([]byte, xerror.IError) {
	return s.MarshalFrame(ctx, frame)
}

func (s *service) Unmarshal(ctx context.Context, data []byte) (api.Frame, xerror.IError) {
	frame := api.Frame{}
	err := s.UnmarshalFrame(ctx, data, &frame)
	if err != nil {
		return api.Frame{}, err
	}
	return frame, err
}
