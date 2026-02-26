package internal

import (
	"context"
	"encoding/json"
	uniformApi "github.com/zzy-rabbit/bp/tool/uniform/api"
	"github.com/zzy-rabbit/xtools/xerror"
)

func (s *service) broadcast(ctx context.Context, frame uniformApi.Frame) xerror.IError {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, conn := range s.connMap {
		err := conn.Post(ctx, frame)
		if xerror.Error(err) {
			s.ILogger.Error(ctx, "conn %s post error: %v", conn.RemoteAddr(ctx), err)
			continue
		}
	}
	return nil
}

func (s *service) Broadcast(ctx context.Context, tag uint32, data interface{}) xerror.IError {
	bytes, err := json.Marshal(data)
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "json marshal error: %v", err)
		return xerror.Extend(xerror.ErrInternalError, err.Error())
	}

	frame := s.IUniform.NewFrame(ctx)
	frame.Format = uniformApi.FormatTypeJson
	frame.Data = bytes
	frame.Encryption = uniformApi.EncryptTypePlainText
	frame.Header.Tag = tag
	return s.broadcast(ctx, frame)
}
