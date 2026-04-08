package internal

import (
	"context"
	"github.com/zzy-rabbit/xtools/xsync"
	"github.com/zzy-rabbit/xtools/xtrace"
)

type fileSync struct {
	xsync.Mutex
}

func (s *service) getFileSync(ctx context.Context, id string) *fileSync {
	s.busyMutex.Lock()
	fs, ok := s.busyFiles[id]
	if !ok {
		fs = &fileSync{}
		s.busyFiles[id] = fs
	}
	s.busyMutex.Unlock()
	return fs
}

func (s *service) deleteFileSync(ctx context.Context, id string) {
	defer xtrace.Trace(ctx)(id)
	s.busyMutex.Lock()
	defer s.busyMutex.Unlock()
	delete(s.busyFiles, id)
}

func (s *service) FileLock(ctx context.Context, id string) {
	defer xtrace.Trace(ctx)(id)
	fs := s.getFileSync(ctx, id)
	fs.Lock(ctx)
}

func (s *service) FileUnlock(ctx context.Context, id string) {
	defer xtrace.Trace(ctx)(id)
	fs := s.getFileSync(ctx, id)
	fs.Unlock(ctx)
}

func (s *service) FileRLock(ctx context.Context, id string) {
	defer xtrace.Trace(ctx)(id)
	fs := s.getFileSync(ctx, id)
	fs.RLock(ctx)
}

func (s *service) FileRUnlock(ctx context.Context, id string) {
	defer xtrace.Trace(ctx)(id)
	fs := s.getFileSync(ctx, id)
	fs.RUnlock(ctx)
}

func (s *service) IsFileLocked(ctx context.Context, id string) bool {
	s.busyMutex.Lock()
	fs, ok := s.busyFiles[id]
	if !ok {
		s.busyMutex.Unlock()
		return false
	}
	s.busyMutex.Unlock()
	return fs.Locked(ctx)
}
